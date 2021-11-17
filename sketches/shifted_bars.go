package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()
	col := rand.Color(rand.Palette())

	m := doc.NewLayer()
	gfx.Clear(m, RGBA{0, 0, 0, 1})
	gfx.Fill{
		Shape: c,
		Color: RGBA{1, 0, 0, 1},
	}.Draw(m)

	vc := VoronoiCells{
		Rect: Rect{
			A: XY{.1, .1},
			B: XY{.9, .9},
		},
		Spacing: math.Interp(0.003, 0.03, .1),
	}

	s := doc.NewLayer()
	gfx.Clear(s, RGBA{1, 1, 1, 1})
	gfx.Fill{
		Shape: vc.Mesh(),
		Color: col,
	}.Draw(s)

	N := 7

	for i := 0; i < 5; i++ {
		y := .2 + float32(i)/5 + rand.Range(-.01, .01)
		l := NewLine(XY{0, y}, XY{1, y})
		sub := Subdivide(l, N)

		for i, _ := range sub {
			dy := rand.Range(-.1, .1)
			sub[i].Y += dy
		}

		modified := XYsToLines(sub...)
		path := Path{}

		for _, line := range modified {
			path = append(path, line)
		}

		gfx.Stroke{
			Shape: path,
			Width: .02,
			Color: RGBA{0, 0, 0, 1},
		}.Draw(m)

		gfx.Stroke{
			Shape: path,
			Width: -.002,
			Color: col,
		}.Draw(s)

		gfx.Stroke{
			Shape: path,
			Width: .022,
			Color: col,
		}.Draw(s)
	}

	gfx.Stroke{
		Shape: c,
		Width: .003,
		Color: col,
	}.Draw(s)

	gfx.Mask{
		Rect:   gfx.Fullscreen,
		Source: s,
		Mask:   m,
	}.Draw(doc)
}

type VoronoiCells struct {
	Rect
	Spacing float32
}

func (vc VoronoiCells) Mesh() Mesh {
	bn := rand.BlueNoise{
		Rect:    vc.Rect,
		Spacing: vc.Spacing,
	}
	noise := bn.Generate()
	v := voronoi.New(noise, vc.Rect)

	var meshes []Mesh
	for _, e := range v.Edges() {
		meshes = append(meshes, e.Stroke(StrokeOpt{}))
	}
	return Merge(meshes...)
}
