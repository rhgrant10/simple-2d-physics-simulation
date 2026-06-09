package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	_ "net/http/pprof"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/mesh"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/world"
)

var (
	avgBallSize     float64 = .25
	cameraZoom              = 1.0
	maxBallVariance         = 0.1
	numBalls                = 1000
)

func main() {
	// establish repeatable randomness
	var seedStr string
	if len(os.Args) < 2 {
		seedStr = generateRandomSeed()
	} else {
		seedStr = os.Args[1]
	}
	fmt.Printf("Randomness seed is: %s\n", seedStr)
	seed := seedFromString(seedStr)
	rng := rand.New(rand.NewSource(seed))

	// set up a full screen situation
	ebiten.SetFullscreen(true)
	pixelWidth, pixelHeight := ebiten.Monitor().Size()
	scale := ebiten.Monitor().DeviceScaleFactor()
	cameraWidth := float64(pixelWidth) * scale
	cameraHeight := float64(pixelHeight) * scale
	camera := camera.Camera{
		Pos:       geo.Vec2{X: 0, Y: 0},
		Width:     cameraWidth,
		Height:    cameraHeight,
		Scale:     50,
		Antialias: true,
	}
	// align the bottom of the screen with world Y=0, then zoom out
	// camera.Pos.Y += (float64(-cameraHeight) / 2) / camera.Scale
	camera.Scale /= 2
	camera.Scale *= cameraZoom

	// create the objects in the scene
	containerColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	w := world.NewWorld()
	floor := world.Spawn(w)
	var width, height float64

	width, height = 100.0, 10.0
	world.AddComponent(w, floor, geo.Transform{Pos: geo.Vec2{X: -20, Y: -5}})
	world.AddComponent(w, floor, mesh.Rectangle{Width: 100, Height: 10, Color: containerColor})
	world.AddComponent(w, floor, collision.Collider(&collision.BoxCollider{HalfW: width / 2, HalfH: height / 2}))
	world.AddComponent(w, floor, physics.Material{Restitution: 0.65, Friction: physics.FrictionIce})

	leftWall := world.Spawn(w)
	width, height = 10.0, 50.0
	world.AddComponent(w, leftWall, geo.Transform{Pos: geo.Vec2{X: -65, Y: 10}})
	world.AddComponent(w, leftWall, mesh.Rectangle{Width: width, Height: height, Color: containerColor})
	world.AddComponent(w, leftWall, collision.Collider(&collision.BoxCollider{HalfW: width / 2, HalfH: height / 2}))
	world.AddComponent(w, leftWall, physics.Material{Restitution: 0.65, Friction: physics.FrictionIce})

	maxBallRadius := 0.0
	for range numBalls {
		radiusJitter := (rng.Float64() - 0.5) * maxBallVariance
		radius := avgBallSize + radiusJitter
		if radius > maxBallRadius {
			maxBallRadius = radius
		}
		ball := world.Spawn(w)
		density := rng.Float64()*0.1 + 20
		world.AddComponent(w, ball, geo.Transform{
			Pos: geo.Vec2{
				X: rng.Float64()*80 - 40 - 20,
				Y: rng.Float64()*50 + 10,
			}})
		world.AddComponent(w, ball, physics.RigidBody{
			Velocity: geo.Vec2FromPolar(
				rng.Float64()*20,
				rng.Float64()*360,
			),
			Mass: density * math.Pi * radius * radius,
		})
		world.AddComponent(w, ball, mesh.Circle{Radius: radius, Color: HSLToRGB(
			rng.Float64()*360,
			0.2*float64(rng.Intn(4)+1),
			rng.Float64()*0.5+0.5,
		)})
		world.AddComponent(w, ball, collision.Collider(&collision.CircleCollider{Radius: radius}))
		world.AddComponent(w, ball, physics.Material{
			Restitution: rng.Float64()*0.1 + 0.55,
			Friction:    physics.KineticFriction(rng.Float64()*0.2 + 0.1)})
	}

	// run the simulation
	game := DefaultGame(*w, camera)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
