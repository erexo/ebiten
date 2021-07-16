// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package vector provides functions for vector graphics rendering.
//
// This package is under experiments and the API might be changed with breaking backward compatibility.
package vector

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type point struct {
	x float32
	y float32
}

// Path represents a collection of path segments.
type Path struct {
	segs [][]point
	cur  point
}

// MoveTo skips the current position of the path to the given position (x, y) without adding any strokes.
func (p *Path) MoveTo(x, y float32) {
	p.cur = point{x: x, y: y}
	p.segs = append(p.segs, []point{p.cur})
}

// LineTo adds a line segument to the path, which starts from the current position and ends to the given position (x, y).
//
// LineTo updates the current position to (x, y).
func (p *Path) LineTo(x, y float32) {
	if len(p.segs) == 0 {
		p.segs = append(p.segs, []point{p.cur})
	}
	seg := p.segs[len(p.segs)-1]
	if seg[len(seg)-1].x != x || seg[len(seg)-1].y != y {
		p.segs[len(p.segs)-1] = append(seg, point{x: x, y: y})
	}
	p.cur = point{x: x, y: y}
}

// QuadTo adds a quadratic Bézier curve to the path.
// (x1, y1) is the control point, and (x2, y2) is the destination.
//
// QuadTo updates the current position to (x2, y2).
func (p *Path) QuadTo(x1, y1, x2, y2 float32) {
	p.quadTo(x1, y1, x2, y2, 0)
}

// isPointCloseToSegment detects the distance between a segment (x0, y0)-(x1, y1) and a point (x, y) is less than allow.
func isPointCloseToSegment(x, y, x0, y0, x1, y1 float32, allow float32) bool {
	// Line passing through (x0, y0) and (x1, y1) in the form of ax + by + c = 0
	a := y1 - y0
	b := -(x1 - x0)
	c := (x1-x0)*y0 - (y1-y0)*x0

	// The distance between a line ax+by+c=0 and (x0, y0) is
	//     |ax0 + by0 + c| / √(a² + b²)
	return allow*allow*(a*a+b*b) > (a*x+b*y+c)*(a*x+b*y+c)
}

func (p *Path) quadTo(x1, y1, x2, y2 float32, level int) {
	if level > 10 {
		return
	}

	x0 := p.cur.x
	y0 := p.cur.y
	if isPointCloseToSegment(x1, y1, x0, y0, x2, y2, 0.5) {
		p.LineTo(x2, y2)
		return
	}

	x01 := (x0 + x1) / 2
	y01 := (y0 + y1) / 2
	x12 := (x1 + x2) / 2
	y12 := (y1 + y2) / 2
	x012 := (x01 + x12) / 2
	y012 := (y01 + y12) / 2
	p.quadTo(x01, y01, x012, y012, level+1)
	p.quadTo(x12, y12, x2, y2, level+1)
}

// CubicTo adds a cubic Bézier curve to the path.
// (x1, y1) and (x2, y2) are the control points, and (x3, y3) is the destination.
//
// CubicTo updates the current position to (x3, y3).
func (p *Path) CubicTo(x1, y1, x2, y2, x3, y3 float32) {
	p.cubicTo(x1, y1, x2, y2, x3, y3, 0)
}

func (p *Path) cubicTo(x1, y1, x2, y2, x3, y3 float32, level int) {
	if level > 10 {
		return
	}

	x0 := p.cur.x
	y0 := p.cur.y
	if isPointCloseToSegment(x1, y1, x0, y0, x3, y3, 0.5) && isPointCloseToSegment(x2, y2, x0, y0, x3, y3, 0.5) {
		p.LineTo(x3, y3)
		return
	}

	x01 := (x0 + x1) / 2
	y01 := (y0 + y1) / 2
	x12 := (x1 + x2) / 2
	y12 := (y1 + y2) / 2
	x23 := (x2 + x3) / 2
	y23 := (y2 + y3) / 2
	x012 := (x01 + x12) / 2
	y012 := (y01 + y12) / 2
	x123 := (x12 + x23) / 2
	y123 := (y12 + y23) / 2
	x0123 := (x012 + x123) / 2
	y0123 := (y012 + y123) / 2
	p.cubicTo(x01, y01, x012, y012, x0123, y0123, level+1)
	p.cubicTo(x123, y123, x23, y23, x3, y3, level+1)
}

func normalize(x, y float32) (float32, float32) {
	len := float32(math.Hypot(float64(x), float64(y)))
	return x / len, y / len
}

// ArcTo adds an arc curve to the path. (x1, y1) is the control point, and (x2, y2) is the destination.
//
// ArcTo updates the current position to (x2, y2).
func (p *Path) ArcTo(x1, y1, x2, y2, radius float32) {
	x0 := p.cur.x
	y0 := p.cur.y
	dx0 := x0 - x1
	dy0 := y0 - y1
	dx1 := x2 - x1
	dy1 := y2 - y1
	dx0, dy0 = normalize(dx0, dy0)
	dx1, dy1 = normalize(dx1, dy1)

	// theta is the angle between two vectors (dx0, dy0) and (dx1, dy1).
	theta := math.Acos(float64(dx0*dx1 + dy0*dy1))
	// TODO: When theta is bigger than π/2, the arc should be split into two.

	// dist is the distance between the control point and the arc's begenning and ending points.
	dist := radius / float32(math.Tan(theta/2))

	// TODO: What if dist is too big?

	// (ax0, ay0) is the start of the arc.
	ax0 := x1 + dx0*dist
	ay0 := y1 + dy0*dist

	// (ax1, ay1) is the end of the arc.
	ax1 := x1 + dx1*dist
	ay1 := y1 + dy1*dist

	p.LineTo(ax0, ay0)

	// Calculate the control points for an approximated Bézier curve.
	// See https://docs.microsoft.com/en-us/xamarin/xamarin-forms/user-interface/graphics/skiasharp/curves/beziers.
	alpha := math.Pi - theta
	l := radius * float32(math.Tan(alpha/4)*4/3)
	cx0 := ax0 + l*(-dx0)
	cy0 := ay0 + l*(-dy0)
	cx1 := ax1 + l*(-dx1)
	cy1 := ay1 + l*(-dy1)
	p.CubicTo(cx0, cy0, cx1, cy1, ax1, ay1)

	p.LineTo(x2, y2)
}

// AppendVerticesAndIndicesForFilling appends vertices and indices to fill this path and returns them.
// AppendVerticesAndIndicesForFilling works in a similar way to the built-in append function.
// If the arguments are nils, AppendVerticesAndIndices returns new slices.
//
// The returned vertice's SrcX and SrcY are 0, and ColorR, ColorG, ColorB, and ColorA are 1.
//
// The returned values are intended to be passed to DrawTriangles or DrawTrianglesShader with EvenOdd option
// in order to render a complex polygon like a concave polygon, a polygon with holes, or a self-intersecting polygon.
func (p *Path) AppendVerticesAndIndicesForFilling(vertices []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16) {
	// TODO: Add tests.

	var base uint16
	for _, seg := range p.segs {
		if len(seg) < 3 {
			continue
		}
		for i, pt := range seg {
			vertices = append(vertices, ebiten.Vertex{
				DstX:   pt.x,
				DstY:   pt.y,
				SrcX:   0,
				SrcY:   0,
				ColorR: 1,
				ColorG: 1,
				ColorB: 1,
				ColorA: 1,
			})
			if i < 2 {
				continue
			}
			indices = append(indices, base, base+uint16(i-1), base+uint16(i))
		}
		base += uint16(len(seg))
	}
	return vertices, indices
}
