package intro

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

type Intro struct {
	images     []*ebiten.Image
	timer      int
	imageIndex int
	sceneCount int
	Level      Level
	Game       Game
}

type Level interface {
	LoadLevel()
}

type Game interface {
	RemoveComponent(c interface{})
}

func New() *Intro {
	i := &Intro{}
	i.loadIntro()
	return i
}

//go:embed img/4.png
var aptImage []byte

func (i *Intro) loadIntro() {

	i.sceneCount = 1

	aptImageData, _, err := image.Decode(bytes.NewReader(aptImage))
	if err != nil {
		log.Fatal(err)
	}

	i.images = []*ebiten.Image{
		ebiten.NewImageFromImage(aptImageData),
	}

}

func (i *Intro) Update() {
	i.timer++
	if i.timer%120 == 0 {
		i.imageIndex++
		if i.imageIndex >= i.sceneCount {
			i.Level.LoadLevel()
			// //go:embed img/1.png
			// var goLangImage []byte

			// //go:embed img/2.png
			// var gopherImage []byte

			// //go:embed img/3.png
			// var thesimpledevImage []byte
			i.Game.RemoveComponent(i)
		}
	}
}

func (i *Intro) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := screen.Size()

	switch i.imageIndex {
	case 0:
		screen.Fill(color.RGBA{0x38, 0x39, 0x3a, 0xff})
		iw4, ih4 := i.images[i.imageIndex].Size()
		op4 := &ebiten.DrawImageOptions{}

		scaleW := float64(screenWidth) / float64(iw4)
		scaleH := float64(screenHeight) / float64(ih4)
		scale4 := math.Min(scaleW, scaleH)

		op4.GeoM.Scale(scale4, scale4)

		scaledW4 := float64(iw4) * scale4
		scaledH4 := float64(ih4) * scale4
		x4 := (float64(screenWidth) - scaledW4) / 2
		y4 := (float64(screenHeight) - scaledH4) / 2
		op4.GeoM.Translate(x4, y4)

		screen.DrawImage(i.images[i.imageIndex], op4)

	}
}
