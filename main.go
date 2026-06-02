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
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/physics"
	"github.com/rhgrant10/sim/scene"
)

type Game struct {
	Cam           camera.Camera
	Entities      []scene.Entity
	lastUpdatedAt time.Time
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Until next time")
	}

	if g.lastUpdatedAt.IsZero() {
		g.lastUpdatedAt = time.Now()
		return nil
	}

	now := time.Now()
	dt := now.Sub(g.lastUpdatedAt).Seconds()
	g.lastUpdatedAt = now

	for i := range g.Entities {
		for j := range g.Entities {
			if i == j {
				continue
			}
			if info, ok := physics.DetectCollision(g.Entities[i], g.Entities[j]); ok {
				physics.ResolveCollision(g.Entities[i], g.Entities[j], info)
			}
		}
	}

	for i := range g.Entities {
		g.Entities[i].Update(dt)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	for i := range g.Entities {
		g.Entities[i].Draw(screen, &g.Cam)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
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
	containerColor := color.White
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
	for range 1000 {
		ball := scene.NewBall(
			geo.Vec2{
				X: rng.Float64()*80 - 40 - 20,
				Y: rng.Float64()*50 + 10,
			},
			geo.Vec2FromPolar(
				rng.Float64()*20,
				rng.Float64()*360,
			),
			rng.Float64()*.2+0.5,
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
	game := Game{
		Cam:      camera,
		Entities: entities,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
