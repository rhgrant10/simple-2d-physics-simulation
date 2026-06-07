package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/scene"
)

var (
	avgBallSize     float64  = .25
	cameraZoom               = 1.0
	cpuProf         *os.File = nil
	enableProfiling          = false
	maxBallVariance          = 0.1
	numBalls                 = 8000
	timeScale                = 0.50
)

func generateRandomSeed() string {
	b := make([]byte, 8)
	crand.Read(b)
	return fmt.Sprintf("%x", b)
}

func seedFromString(s string) int64 {
	h := sha256.Sum256([]byte(s))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

// func randomBall(minR, maxR float64) scene.Ball {

// }

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
	// camera.Pos.Y += (float64(-cameraHeight) / 2) / camera.Scale
	camera.Scale /= 2
	camera.Scale *= cameraZoom

	// create the objects in the scene
	containerColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	floor := scene.NewBox(
		geo.Vec2{X: -20, Y: -5},
		100,
		10,
		containerColor,
		0.65,
		physics.FrictionIce,
	)
	leftWall := scene.NewBox(
		geo.Vec2{X: -65, Y: 10},
		10,
		50,
		containerColor,
		0.65,
		physics.FrictionIce,
	)
	entities := []scene.Entity{&floor, &leftWall}

	maxBallRadius := 0.0
	for range numBalls {
		radiusJitter := (rng.Float64() - 0.5) * maxBallVariance
		radius := avgBallSize + radiusJitter
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
			rng.Float64()*0.1+20,
			scene.HSLToRGB(
				rng.Float64()*360,
				0.2*float64(rng.Intn(4)+1),
				rng.Float64()*0.5+0.5,
			),
			rng.Float64()*0.1+0.55,
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
