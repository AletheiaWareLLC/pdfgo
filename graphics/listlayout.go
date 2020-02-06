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
	"errors"
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
	Direction Direction
	Padding   float64
	Boxes     []Box
}

func (l *ListLayout) Add(box Box) {
	l.Boxes = append(l.Boxes, box)
}

func (l *ListLayout) SetBounds(bounds *Rectangle) error {
	l.Left = bounds.Left
	l.Top = bounds.Top
	l.Right = bounds.Right
	l.Bottom = bounds.Bottom
	left := bounds.Left
	top := bounds.Top
	right := bounds.Right
	bottom := bounds.Bottom
	switch l.Direction {
	case TopBottom:
		for _, b := range l.Boxes {
			if err := b.SetBounds(&Rectangle{
				Left:   left,
				Top:    top,
				Right:  right,
				Bottom: bottom,
			}); err != nil {
				return err
			}
			top -= b.GetHeight() + l.Padding
		}
	default:
		return errors.New(fmt.Sprintf("Direction not implemented: %d", l.Direction))
	}
	return nil
}

func (l *ListLayout) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	for _, b := range l.Boxes {
		if err := b.Write(p, buffer); err != nil {
			return err
		}
	}
	return nil
}
