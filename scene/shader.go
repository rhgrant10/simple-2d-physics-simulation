package scene

// shaderSrc is the Kage shader source. Kage is Ebiten's shading language,
// which is a subset of Go with some GLSL-like builtins.
//
// This shader runs once per pixel for every triangle we submit.
// The GPU runs all these pixels IN PARALLEL — that's why one
// DrawTriangles call with 1000 shapes is fast.
//
// Inputs available per fragment (pixel):
//
//	dst  - position of this pixel in the destination image
//	src  - the interpolated SrcX/SrcY from our vertex data (our UVs)
//	col  - the interpolated RGBA from our vertex color fields
//
// We use the UV range as a shape-type flag:
//
//	src in [0,1]  → circle, apply SDF
//	src in [2,3]  → box, fill solid
var shaderSrc = []byte(`
//kage:unit pixels
package scene
 
func Fragment(dst vec4, src vec2, col vec4) vec4 {
	// -----------------------------------------------------------------
	// Box: UVs are in [2,3] range — just return the vertex color.
	// The quad geometry already defines the box shape exactly.
	// -----------------------------------------------------------------
	if src.x >= 2.0 {
		return col
	}
 
	// -----------------------------------------------------------------
	// Circle: UVs are in [0,1] range.
	//
	// The quad is a square bounding box around the circle.
	// We use a Signed Distance Field (SDF) to determine whether
	// each pixel is inside or outside the circle.
	//
	// SDF for a circle centered at (0.5, 0.5) with radius 0.5:
	//   distance = length(uv - center) - radius
	//   < 0 → inside
	//   > 0 → outside
	//
	// We remap src from [0,1] to [-0.5, 0.5] so the center is origin.
	// -----------------------------------------------------------------
	uv := src - vec2(0.5, 0.5)
	dist := length(uv)  // distance from center, in UV space (0..0.5 = edge)
 
	// smoothstep gives us anti-aliased edges for free.
	// The range [0.47, 0.5] controls the softness of the edge.
	// Smaller range = sharper edge.
	alpha := 1.0 - smoothstep(0.47, 0.5, dist)
 
	// Discard fully transparent pixels (outside the circle).
	// This is what "punches out" the corners of the bounding quad.
	if alpha <= 0.0 {
		return vec4(0)
	}
 
	// Premultiply alpha: Ebiten expects premultiplied RGBA.
	// col.a * alpha is the final opacity; RGB is scaled to match.
	a := col.a * alpha
	return vec4(col.rgb * a, a)
}
`)
