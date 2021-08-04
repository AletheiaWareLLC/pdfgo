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

type Direction int

const (
	LeftRight Direction = iota
	RightLeft
	TopBottom
	BottomTop
)

type ListLayout struct {
	Rectangle
	Direction      Direction
	Padding        float64
	Boxes          []Box
	MinimumVisible int
	Visible        int
}

func (l *ListLayout) Add(box Box) {
	l.Boxes = append(l.Boxes, box)
}

func (l *ListLayout) SetBounds(bounds *Rectangle) (*Rectangle, error) {
	left := bounds.Left
	top := bounds.Top
	right := bounds.Right
	bottom := bounds.Bottom
	used := NegativeRectangle()
	l.Visible = 0
	switch l.Direction {
	case TopBottom:
		for _, b := range l.Boxes {
			u, err := b.SetBounds(&Rectangle{
				Left:   left,
				Top:    top,
				Right:  right,
				Bottom: bottom,
			})
			if err != nil {
				return nil, err
			}
			top -= u.DY()
			if top < bottom {
				break
			}
			l.Visible += 1
			used = used.Max(u)
			top -= l.Padding
		}
	default:
		return nil, fmt.Errorf("Direction not implemented: %d", l.Direction)
	}
	if l.Visible < l.MinimumVisible {
		l.Visible = 0
		return bounds, nil
	}
	l.Left = used.Left
	l.Top = used.Top
	l.Right = used.Right
	l.Bottom = used.Bottom
	return used, nil
}

func (l *ListLayout) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	for i, b := range l.Boxes {
		if i >= l.Visible {
			break
		}
		if err := b.Write(p, buffer); err != nil {
			return err
		}
	}
	return nil
}
