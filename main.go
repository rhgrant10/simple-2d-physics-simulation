package main

import (
	"errors"
	"image/color"
	"log"
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

func main() {
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
	// and align the bottom of the screen with world Y=0
	camera.Pos.Y += (float64(-cameraHeight) / 2) / camera.Scale
	camera.Scale /= 2

	ball1 := scene.NewBall(
		geo.Vec2{X: -20, Y: 9},
		geo.Vec2{X: 7, Y: 8},
		3,
		2,
		color.White,
		0.75,
	)
	ball2 := scene.NewBall(
		geo.Vec2{X: 5, Y: 11},
		geo.Vec2{X: -8, Y: 7},
		3,
		2,
		color.RGBA{R: 0, G: 0, B: 255, A: 255},
		0.75,
	)
	floor := scene.NewBox(
		geo.Vec2{X: -20, Y: -1},
		80,
		2,
		color.RGBA{R: 255, G: 128, B: 128, A: 255},
		0.65,
	)

	game := Game{
		Cam:      camera,
		Entities: []scene.Entity{&ball1, &ball2, &floor},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
