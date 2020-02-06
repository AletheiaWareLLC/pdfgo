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
	"github.com/AletheiaWareLLC/pdfgo"
)

type FibonacciLayout struct {
	Sizes []float64
	Boxes []Box
}

func (l *FibonacciLayout) Add(box Box) {
	l.Boxes = append(l.Boxes, box)
}

func (l *FibonacciLayout) SetBounds(bounds *Rectangle) error {
	x := bounds.Left
	y := bounds.Top
	for i, b := range l.Boxes {
		size := l.Sizes[i]
		if err := b.SetBounds(&Rectangle{
			Left:   x,
			Top:    y,
			Right:  x + size,
			Bottom: y - size,
		}); err != nil {
			return err
		}
		switch i % 2 {
		case 0:
			// Move Down
			y -= size
		case 1:
			// Move Right
			x += size
		}
	}
	return nil
}

func (l *FibonacciLayout) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	for _, b := range l.Boxes {
		if err := b.Write(p, buffer); err != nil {
			return err
		}
	}
	return nil
}
