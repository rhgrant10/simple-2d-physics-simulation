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
	"github.com/rhgrant10/sim/systems"
	"github.com/rhgrant10/sim/world"
)

type ProfilingConfig struct {
	enabled     bool
	cpuFilename string
	cpuFile     *os.File
	memFilename string
	memFile     *os.File
	duration    time.Duration
}

func DefaultProfilingConfig() ProfilingConfig {
	return ProfilingConfig{
		enabled:     true,
		cpuFilename: "cpu.prof",
		memFilename: "mem.prof",
		duration:    10 * time.Second,
	}
}

type Renderer interface {
	Draw(screen *ebiten.Image, camera *camera.Camera, world *world.World)
}

type Simulation interface {
	Step(world *world.World, dt float64)
}

type Game struct {
	renderer        Renderer
	simulation      Simulation
	camera          camera.Camera
	world           world.World
	lastUpdatedAt   time.Time
	startTime       time.Time
	profiling       ProfilingConfig
	statefulMessage string
}

func NewGame(w world.World, r Renderer, s Simulation, c camera.Camera, p ProfilingConfig) *Game {
	return &Game{
		renderer:   r,
		simulation: s,
		camera:     c,
		world:      w,
		profiling:  p,
	}
}

func DefaultGame(w world.World, c camera.Camera) *Game {
	renderer := scene.NewBasicRenderer()
	bvh := geo.NewBVH()
	detector := collision.NewDefaultDetector(&bvh, &collision.GJK{})
	resolver := collision.NewImpulseResolver()
	simulation := systems.BasicSimulator{
		Integrator: &systems.VerletIntegrator{},
		Solver: &systems.BasicCollisionSystem{
			Detector: &detector,
			Resolver: &resolver,
		},
		// StepInterval: 1 * time.Second / 120.0,
	}
	profileConfig := DefaultProfilingConfig()
	return NewGame(w, renderer, &simulation, c, profileConfig)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("Until next time")
	}

	if g.lastUpdatedAt.IsZero() {
		g.startTime = time.Now()
		g.lastUpdatedAt = g.startTime
		if g.profiling.enabled {
			g.statefulMessage = fmt.Sprintf("Profiling for %v...", g.profiling.duration)
			f, _ := os.Create(g.profiling.cpuFilename)
			g.profiling.cpuFile = f
			pprof.StartCPUProfile(g.profiling.cpuFile)
		}
		return nil
	}

	if g.profiling.enabled && time.Since(g.startTime) >= g.profiling.duration {
		pprof.StopCPUProfile()
		g.profiling.cpuFile.Close()
		g.profiling.memFile, _ = os.Create(g.profiling.memFilename)
		pprof.WriteHeapProfile(g.profiling.memFile)
		g.profiling.memFile.Close()
		g.statefulMessage = "Profiling... Done!"
	}

	now := time.Now()
	dt := now.Sub(g.lastUpdatedAt).Seconds()
	g.lastUpdatedAt = now
	g.simulation.Step(&g.world, dt)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.1f\nFPS: %.1f\n%s", ebiten.ActualTPS(), ebiten.ActualFPS(), g.statefulMessage))
	g.renderer.Draw(screen, &g.camera, &g.world)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return outsideWidth * scale, outsideHeight * scale
}
