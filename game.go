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
	bvh := geo.NewBVH()
	// noBroadphase := collision.NewNoBroadphase()
	// spatialHash := geo.NewSpatialHash(dynamicObjectRadius*40, len(dynamicEntities), 10 * len(dynamicEntities))
	return &Game{
		StaticEntities:  staticEntities,
		DynamicEntities: dynamicEntities,
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
	dt := now.Sub(g.lastUpdatedAt).Seconds() * timeScale
	g.lastUpdatedAt = now

	// collisions with dynamic objects - build
	bvs := make([]geo.AABB, 0, len(g.DynamicEntities))
	for i := range g.DynamicEntities {
		if !g.DynamicEntities[i].GetRigidBody().IsStatic() {
			transform := *g.DynamicEntities[i].GetTransform()
			aabb := g.DynamicEntities[i].GetCollider().AABB(transform)
			bvs = append(bvs, aabb)
		}
	}
	for range 1 {
		g.Broadphaser.Clear()
		g.Broadphaser.Build(&bvs)

		// collisions with dynamic objects - find & process
		g.closePairs = g.Broadphaser.CandidatePairs(g.closePairs)
		collisionCount := 0
		for _, p := range g.closePairs {
			a, b := p.A(), p.B()
			if info, ok := collision.Detect(g.DynamicEntities[a], g.DynamicEntities[b]); ok {
				collisionCount++
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
	}

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
	g.Renderer.Begin()
	for i := range g.DynamicEntities {
		g.DynamicEntities[i].GetBrush(&g.Cam).AppendVertices(&g.Renderer.Vertices, &g.Renderer.Indices)
	}
	for i := range g.StaticEntities {
		g.StaticEntities[i].GetBrush(&g.Cam).AppendVertices(&g.Renderer.Vertices, &g.Renderer.Indices)
	}
	g.Renderer.Flush(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	return outsideWidth * scale, outsideHeight * scale
}
