package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/scene"
)

var (
	avgBallSize              = 5.0
	cpuProf         *os.File = nil
	enableProfiling          = true
	maxBallVariance          = 0.1
	numBalls                 = 5
	gravityScale             = 1.0
)

type Game struct {
	Renderer        scene.Renderer
	Cam             camera.Camera
	Broadphaser     collision.Broadphaser
	DynamicEntities []scene.Entity
	StaticEntities  []scene.Entity
	closePairs      []geo.PairKey
	lastUpdatedAt   time.Time
	startTime       time.Time
	statefulMessage string
}

func NewGame(entities []scene.Entity, camera camera.Camera, dynamicObjectRadius float64) *Game {
	staticEntities := []scene.Entity{}
	dynamicEntities := []scene.Entity{}
	for i := range entities {
		if entities[i].GetRigidBody().IsStatic() {
			staticEntities = append(staticEntities, entities[i])
		} else {
			dynamicEntities = append(dynamicEntities, entities[i])
		}
	}

	bvh := geo.BVH{}
	// noBroadphase := collision.NewNoBroadphase()
	// spatialHash := geo.NewSpatialHash(dynamicObjectRadius*40, len(dynamicEntities), 10 * len(dynamicEntities))

	return &Game{
		DynamicEntities: dynamicEntities,
		StaticEntities:  staticEntities,
		Renderer:        *scene.NewRenderer(),
		Cam:             camera,
		Broadphaser:     &bvh,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Until next time")
	}

	if g.lastUpdatedAt.IsZero() {
		g.startTime = time.Now()
		g.lastUpdatedAt = g.startTime
		if enableProfiling {
			g.statefulMessage = "Profiling..."
			cpuProf, _ := os.Create("cpu.prof")
			pprof.StartCPUProfile(cpuProf)
		}
		return nil
	}

	if enableProfiling && time.Since(g.startTime) >= 10*time.Second {
		pprof.StopCPUProfile()
		cpuProf.Close()
		g.statefulMessage = "Profiling... Done!"
	}

	now := time.Now()
	dt := now.Sub(g.lastUpdatedAt).Seconds() * gravityScale
	g.lastUpdatedAt = now

	// collisions with dynamic objects - build
	g.Broadphaser.Clear()
	for i := range g.DynamicEntities {
		if !g.DynamicEntities[i].GetRigidBody().IsStatic() {
			transform := *g.DynamicEntities[i].GetTransform()
			aabb := g.DynamicEntities[i].GetCollider().AABB(transform)
			g.Broadphaser.Insert(i, aabb)
		}
	}

	// collisions with dynamic objects - find & process
	g.closePairs = g.Broadphaser.CandidatePairs(g.closePairs)
	for _, p := range g.closePairs {
		a, b := p.A(), p.B()
		if info, ok := collision.Detect(g.DynamicEntities[a], g.DynamicEntities[b]); ok {
			collision.Resolve(g.DynamicEntities[a], g.DynamicEntities[b], info)
		}
	}

	// collisions with static objects
	for d := range g.DynamicEntities {
		for s := range g.StaticEntities {
			if info, ok := collision.Detect(g.DynamicEntities[d], g.StaticEntities[s]); ok {
				collision.Resolve(g.DynamicEntities[d], g.StaticEntities[s], info)
			}
		}
	}

	for i := range g.DynamicEntities {
		g.DynamicEntities[i].Update(dt)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f\nFPS: %.1f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), g.statefulMessage))
	g.Renderer.Begin()
	for i := range g.StaticEntities {
		g.StaticEntities[i].GetBrush(&g.Cam).AppendVertices(&g.Renderer.Vertices, &g.Renderer.Indices)
	}
	for i := range g.DynamicEntities {
		g.DynamicEntities[i].GetBrush(&g.Cam).AppendVertices(&g.Renderer.Vertices, &g.Renderer.Indices)
	}
	g.Renderer.Flush(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth, screenHeight = ebiten.Monitor().Size()
	// screenWidth = (screenWidth * 7) / 10
	// screenHeight = (screenHeight * 7) / 10
	return
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return outsideWidth * scale, outsideHeight * scale
}

func generateRandomSeed() string {
	b := make([]byte, 8)
	crand.Read(b)
	return fmt.Sprintf("%x", b)
}

func seedFromString(s string) int64 {
	h := sha256.Sum256([]byte(s))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

func main() {
	// set up a profiler
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()

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
	camera.Pos.Y += (float64(-cameraHeight) / 2) / camera.Scale
	camera.Scale /= 2

	// create the objects in the scene
	containerColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	floor := scene.NewBox(
		geo.Vec2{X: -20, Y: -1},
		80,
		2,
		containerColor,
		0.65,
		physics.FrictionConcrete,
	)
	leftWall := scene.NewBox(
		geo.Vec2{X: -61, Y: 18},
		2,
		40,
		containerColor,
		0.65,
		physics.FrictionConcrete,
	)
	entities := []scene.Entity{&floor, &leftWall}

	maxBallRadius := 0.0
	for range numBalls {
		radius := rng.Float64()*(avgBallSize*maxBallVariance) + avgBallSize
		if radius > maxBallRadius {
			maxBallRadius = radius
		}
		ball := scene.NewBall(
			geo.Vec2{
				X: rng.Float64()*80 - 40 - 20,
				Y: rng.Float64()*50 + 10,
			},
			geo.Vec2FromPolar(
				rng.Float64()*20,
				rng.Float64()*360,
			),
			radius,
			rng.Float64()*0.1+2,
			scene.HSLToRGB(
				rng.Float64()*360,
				0.2*float64(rng.Intn(4)+1),
				rng.Float64()*0.5+0.5,
			),
			rng.Float64()*0.4+0.55,
			physics.KineticFriction(rng.Float64()*0.2+0.1),
		)
		entities = append(entities, &ball)
	}

	// run the simulation
	game := NewGame(entities, camera, maxBallRadius)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
