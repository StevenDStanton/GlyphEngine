package main

import (
	_ "image/png"
	"log"

	"github.com/StevenDStanton/the-social-shift/internal/game"
	"github.com/StevenDStanton/the-social-shift/internal/intro"
	"github.com/StevenDStanton/the-social-shift/internal/level"
	"github.com/StevenDStanton/the-social-shift/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	title        = "The Social Contract"
	screenWidth  = 1280
	screenHeight = 720
)

func main() {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := game.New(screenWidth, screenHeight)

	p := player.New(g.AudioContext, 10, 10)
	i := intro.New()

	l := level.New()
	l.Player = p

	p.Level = l

	i.Game = g
	i.Level = l

	g.AddComponent(i)
	g.AddComponent(l)
	g.AddComponent(p)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
