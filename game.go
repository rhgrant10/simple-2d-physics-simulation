package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rhgrant10/sim/camera"
	"github.com/rhgrant10/sim/collision"
	"github.com/rhgrant10/sim/geo"
	"github.com/rhgrant10/sim/scene"
)

type Game struct {
	Renderer        scene.Renderer
	Cam             camera.Camera
	CollisionSolver collision.Solver
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
	broad := geo.NewBVH()
	narrow := collision.GJK{}
	detector := collision.NewDefaultDetector(&broad, &narrow)
	resolver := collision.NewImpulseResolver()
	solver := collision.NewSolver(&detector, resolver)

	return &Game{
		StaticEntities:  staticEntities,
		DynamicEntities: dynamicEntities,
		CollisionSolver: solver,
		Renderer:        *scene.NewRenderer(),
		Cam:             camera,
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

	if enableProfiling && time.Since(g.startTime) >= profilingDuration {
		pprof.StopCPUProfile()
		cpuProf.Close()
		g.statefulMessage = "Profiling... Done!"
	}

	now := time.Now()
	dt := now.Sub(g.lastUpdatedAt).Seconds() * timeScale
	g.lastUpdatedAt = now

	collidables := make([]collision.Collidable, 0, len(g.DynamicEntities)+len(g.StaticEntities))
	for i := range g.DynamicEntities {
		collidables = append(collidables, g.DynamicEntities[i])
	}
	for i := range g.StaticEntities {
		collidables = append(collidables, g.StaticEntities[i])
	}
	g.CollisionSolver.Solve(collidables)

	physicsHz := 120.0
	maxDt := 1 / physicsHz
	for dt > 0 {
		stepDt := min(maxDt, dt)
		for i := range g.DynamicEntities {
			g.DynamicEntities[i].Update(stepDt)
		}
		dt -= maxDt
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f\nFPS: %.1f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), g.statefulMessage))
	g.Renderer.Render(screen, &g.Cam, append(g.DynamicEntities, g.StaticEntities...))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return outsideWidth * scale, outsideHeight * scale
}
