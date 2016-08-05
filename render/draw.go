package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
	"container/heap"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
)

var (
	rh                *RenderableHeap
	toPushRenderables []Renderable
	preDrawBind       event.Binding
	resetHeap         bool
)

type RenderableHeap []Renderable

// Satisfying the Heap interface
func (h RenderableHeap) Len() int           { return len(h) }
func (h RenderableHeap) Less(i, j int) bool { return h[i].GetLayer() < h[j].GetLayer() }
func (h RenderableHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *RenderableHeap) Push(x interface{}) {
	*h = append(*h, x.(Renderable))
}

func (h_p *RenderableHeap) Pop() interface{} {
	h := *h_p
	n := len(h)
	x := h[n-1]
	*h_p = h[0 : n-1]
	return x
}

// ResetDrawHeap sets a flag to clear the drawheap
// at the next predraw phase
func ResetDrawHeap() {
	resetHeap = true
}

func InitDrawHeap() {
	rh = &RenderableHeap{}
	heap.Init(rh)
	preDrawBind, _ = event.GlobalBind(PreDraw, "PreDraw")
}

// Drawing does not actually immediately draw a renderable,
// instead the renderable is added to a list of elements to
// be drawn next frame. This avoids issues where elements
// are added to the heap while it is being drawn.
func Draw(r Renderable, l int) Renderable {
	r.SetLayer(l)
	toPushRenderables = append(toPushRenderables, r)
	return r
}

// PreDraw parses through renderables to be pushed
// and adds them to the drawheap.
func PreDraw(no int, nothing interface{}) int {
	if resetHeap == true {
		InitDrawHeap()
		resetHeap = false
	} else {
		for _, r := range toPushRenderables {
			heap.Push(rh, r)
		}
	}
	toPushRenderables = []Renderable{}
	return 0
}

// LoadSpriteAndDraw is shorthand for LoadSprite
// followed by Draw.
func LoadSpriteAndDraw(filename string, l int) *Sprite {
	s := LoadSprite(filename)
	return Draw(s, l).(*Sprite)
}

// DrawColor is equivalent to LoadSpriteAndDraw,
// but with colorboxes.
func DrawColor(c color.Color, x1, y1, x2, y2 float64, l int) {
	cb := NewColorBox(int(x2), int(y2), c)
	cb.ShiftX(x1)
	cb.ShiftY(y1)
	Draw(cb, l)
}

// DrawHeap takes every element in the heap
// and draws them as it removes them. It
// filters out elements who have the layer
// -1, reserved for elements to be undrawn.
func DrawHeap(b screen.Buffer) {
	newRh := &RenderableHeap{}
	for rh.Len() > 0 {
		r := heap.Pop(rh).(Renderable)
		if r.GetLayer() != -1 {
			r.Draw(b.RGBA())
			heap.Push(newRh, r)
		}
	}
	rh = newRh
}

// ShinyDraw performs a draw operation at -x, -y, because
// shiny/screen represents quadrant 4 as negative in both axes.
// draw.Over will merge two pixels at a given position based on their
// alpha channel.
func ShinyDraw(buff draw.Image, img image.Image, x, y int) {
	draw.Draw(buff, buff.Bounds(),
		img, image.Point{-x, -y}, draw.Over)
}

// draw.Src will overwrite pixels beneath the given image regardless of
// the new image's alpha.
func ShinyOverwrite(buff screen.Buffer, img image.Image, x, y int) {
	draw.Draw(buff.RGBA(), buff.Bounds(),
		img, image.Point{-x, -y}, draw.Src)
}
