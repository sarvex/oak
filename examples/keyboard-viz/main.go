package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/debugtools/inputviz"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	oak.AddScene("keyviz", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("start")
			fnt, _ := render.DefFontGenerator.RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
				fg.Color = image.NewUniform(color.RGBA{0, 0, 0, 255})
				fg.Size = 13
				return fg
			})
			m := inputviz.Keyboard{
				Rect:             floatgeom.NewRect2(0, 0, float64(ctx.Window.Width()), float64(ctx.Window.Height())),
				BaseLayer:        -1,
				RenderCharacters: true,
				Font:             fnt,
			}
			m.RenderAndListen(ctx, 0)
		},
	})
	err := oak.Init("keyviz", func(c oak.Config) (oak.Config, error) {
		c.Debug.Level = dlog.VERBOSE.String()
		c.Screen.Width = 800
		c.Screen.Height = 300
		return c, nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
