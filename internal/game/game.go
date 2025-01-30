package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Drawable interface {
	Draw(screen *ebiten.Image)
}

type Updatable interface {
	Update()
}

type Game struct {
	AudioContext *audio.Context
	screenWidth  int
	screenHeight int
	Components   []interface{}
}

func New(screenWidth, screenHeight int) *Game {
	return &Game{
		AudioContext: audio.NewContext(44100),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (g *Game) AddComponent(c interface{}) {
	g.Components = append(g.Components, c)
}

func (g *Game) RemoveComponent(c interface{}) {
	for i, component := range g.Components {
		if component == c {
			g.Components[i] = nil
			break
		}
	}
}

func (g *Game) cleanupComponents() {
	components := []interface{}{}
	for _, c := range g.Components {
		if c != nil {
			components = append(components, c)
		}
	}

	g.Components = components
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, c := range g.Components {
		if updatable, ok := c.(Drawable); ok {
			updatable.Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) Update() error {
	for _, c := range g.Components {
		if c == nil {
			continue
		}
		if updatable, ok := c.(Updatable); ok {
			updatable.Update()
		}

	}

	g.cleanupComponents()

	return nil
}
