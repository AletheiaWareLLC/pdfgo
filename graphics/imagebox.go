/*
 * Copyright 2021 Aletheia Ware LLC
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
	"math"
)

type ImageBox struct {
	Rectangle
	ImageID                     string
	Width, Height               float64
	MinimumWidth, MinimumHeight float64
}

func (b *ImageBox) SetBounds(bounds *Rectangle) (*Rectangle, error) {
	dx := bounds.Right - bounds.Left
	dy := bounds.Top - bounds.Bottom
	scaledWidth, scaledHeight := b.scales(dx, dy)
	b.Left = bounds.Left
	b.Top = bounds.Top
	b.Right = bounds.Left + scaledWidth
	b.Bottom = bounds.Top - scaledHeight
	return &b.Rectangle, nil
}

func (b *ImageBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	dx := b.DX()
	dy := b.DY()
	scaledWidth, scaledHeight := b.scales(dx, dy)
	translateX := b.Left + ((dx - scaledWidth) / 2)
	translateY := b.Bottom + ((dy - scaledHeight) / 2)
	buffer.WriteString(fmt.Sprintf("q %f 0 0 %f %f %f cm /%s Do Q\n", scaledWidth, scaledHeight, translateX, translateY, b.ImageID))
	return nil
}

func (b *ImageBox) scales(dx, dy float64) (float64, float64) {
	scale := math.Max(b.Width/dx, b.Height/dy)
	maximumScale := math.Max(b.Width/b.MinimumWidth, b.Height/b.MinimumHeight)
	if scale > maximumScale {
		scale = maximumScale
	}
	scaledWidth := b.Width / scale
	scaledHeight := b.Height / scale
	return scaledWidth, scaledHeight
}
