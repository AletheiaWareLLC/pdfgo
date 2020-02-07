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

type ColourBox struct {
	Rectangle
	BorderColour []float64
	FillColour   []float64
}

func (b *ColourBox) SetBounds(bounds *Rectangle) error {
	b.Left = bounds.Left
	b.Top = bounds.Top
	b.Right = bounds.Right
	b.Bottom = bounds.Bottom
	return nil
}

func (b *ColourBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	buffer.WriteString("q\n")
	// Fill
	if b.FillColour != nil {
		buffer.WriteString(fmt.Sprintf("%s %s %s rg\n", FloatToString(b.FillColour[0]), FloatToString(b.FillColour[1]), FloatToString(b.FillColour[2])))
		buffer.WriteString(fmt.Sprintf("%s %s %s %s re f\n", FloatToString(b.Left), FloatToString(b.Bottom), FloatToString(b.GetWidth()), FloatToString(b.GetHeight())))
	}
	// Border
	if b.BorderColour != nil {
		buffer.WriteString(fmt.Sprintf("%s %s %s RG\n", FloatToString(b.BorderColour[0]), FloatToString(b.BorderColour[1]), FloatToString(b.BorderColour[2])))
		buffer.WriteString(fmt.Sprintf("%s %s %s %s re S\n", FloatToString(b.Left), FloatToString(b.Bottom), FloatToString(b.GetWidth()), FloatToString(b.GetHeight())))
	}
	buffer.WriteString("Q")
	return nil
}
