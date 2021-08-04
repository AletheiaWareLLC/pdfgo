/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package graphics

import (
	"bytes"
	"fmt"
	"github.com/AletheiaWareLLC/pdfgo"
)

type Gravity int

const (
	Middle Gravity = iota
	North
	East
	South
	West
)

type GravityLayout struct {
	Rectangle
	Box     Box
	Gravity Gravity
}

func (l *GravityLayout) Add(box Box) {
	l.Box = box
}

func (l *GravityLayout) SetBounds(bounds *Rectangle) (*Rectangle, error) {
	l.Left = bounds.Left
	l.Top = bounds.Top
	l.Right = bounds.Right
	l.Bottom = bounds.Bottom
	// Set bounds once to calculate width & height
	used, err := l.Box.SetBounds(bounds)
	if err != nil {
		return nil, err
	}
	var b *Rectangle
	switch l.Gravity {
	case North:
		dx := (l.DX() - used.DX()) / 2
		b = &Rectangle{
			Left:   l.Left + dx,
			Top:    l.Top,
			Right:  l.Right - dx,
			Bottom: l.Bottom,
		}
	case Middle:
		dx := (l.DX() - used.DX()) / 2
		dy := (l.DY() - used.DY()) / 2
		b = &Rectangle{
			Left:   l.Left + dx,
			Top:    l.Top - dy,
			Right:  l.Right - dx,
			Bottom: l.Bottom + dy,
		}
	case South:
		dx := (l.DX() - used.DX()) / 2
		dy := (l.DY() - used.DY())
		b = &Rectangle{
			Left:   l.Left + dx,
			Top:    l.Top - dy,
			Right:  l.Right - dx,
			Bottom: l.Bottom,
		}
	default:
		return nil, fmt.Errorf("Gravity not implemented: %d", l.Gravity)
	}
	// Set final bounds
	if _, err := l.Box.SetBounds(b); err != nil {
		return nil, err
	}
	return bounds, nil
}

func (l *GravityLayout) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	return l.Box.Write(p, buffer)
}
