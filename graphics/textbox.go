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
	"github.com/AletheiaWareLLC/pdfgo/font"
	"log"
	"unicode"
)

const TEXT_PADDING = 8

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
	JustifiedLeft
	JustifiedCenter
	JustifiedRight
)

type Line struct {
	CharacterSpacing  float64
	WordSpacing       float64
	HorizontalScaling float64
	Indent            float64
	Rise              float64
	Render            int
	Text              string
}

type TextBox struct {
	Text             []rune
	FontID           string
	Font             font.Font
	FontSize         float64
	FontColour       []float64
	Align            Alignment
	ShrinkToFit      bool
	OriginX, OriginY float64
	Lines            []*Line
}

func (b *TextBox) AddMeasuredLine(text []rune, width, delta float64) {
	var indent float64
	for _, l := range b.Lines {
		indent += l.Indent
	}
	var line *Line
	switch b.Align {
	case Center:
		line = newCenteredLine(text, width, delta, indent)
	case Right:
		line = newRightLine(text, width, delta, indent)
	case JustifiedLeft:
		line = newJustifiedLine(text, width, delta, indent, Left)
	case JustifiedRight:
		line = newJustifiedLine(text, width, delta, indent, Right)
	case JustifiedCenter:
		line = newJustifiedLine(text, width, delta, indent, Center)
	default:
		b.Align = Left
		fallthrough
	case Left:
		line = newLeftLine(text, width, delta, indent)
	}
	b.Lines = append(b.Lines, line)
}

func newLeftLine(text []rune, width, delta, indent float64) *Line {
	return &Line{
		HorizontalScaling: 100,
		Indent:            indent,
		Text:              PDFEscapeString(text),
	}
}

func newCenteredLine(text []rune, width, delta, indent float64) *Line {
	return &Line{
		HorizontalScaling: 100,
		Indent:            (delta / 2) - indent,
		Text:              PDFEscapeString(text),
	}
}

func newRightLine(text []rune, width, delta, indent float64) *Line {
	return &Line{
		HorizontalScaling: 100,
		Indent:            delta - indent,
		Text:              PDFEscapeString(text),
	}
}

func newJustifiedLine(text []rune, width, delta, indent float64, secondary Alignment) *Line {
	if delta > (width / 4) {
		// Use secondary alignment for short lines
		switch secondary {
		case Center:
			return newCenteredLine(text, width, delta, indent)
		case Right:
			return newRightLine(text, width, delta, indent)
		default:
			fallthrough
		case Left:
			return newLeftLine(text, width, delta, indent)
		}
	}

	line := &Line{
		HorizontalScaling: 100,
		Indent:            -indent,
		Text:              PDFEscapeString(text),
	}
	characterSpaces := float64(len(text) - 1)
	var wordSpaces float64
	for _, c := range text {
		if unicode.IsSpace(c) {
			wordSpaces++
		}
	}
	if wordSpaces > 1 {
		line.WordSpacing = (delta / 2) / wordSpaces
		line.CharacterSpacing = (delta / 2) / characterSpaces
		//line.HorizontalScaling += (((delta / 3) / (width + delta)) / 15) * 100
	} else {
		line.CharacterSpacing = (delta / 1) / characterSpaces
		//line.HorizontalScaling += (delta / width) * 100
	}
	return line
}

func (b *TextBox) SetBounds(bounds *Rectangle) (*Rectangle, error) {
	b.OriginX = bounds.Left
	b.OriginY = bounds.Top - b.FontSize
	maxWidth := bounds.DX()
	maxHeight := bounds.DY()

	var height float64
	b.Lines = nil
	for _, line := range SplitLines(b.Text) {
		textWidth := b.Font.MeasureText(line, b.FontSize)
		delta := maxWidth - textWidth
		if delta < 0 && !b.ShrinkToFit {
			// Split line
			wrappoint := -1
			start := 0
			end := 0
			for end < len(line) {
				c := line[end]
				if unicode.IsSpace(c) {
					wrappoint = end
				}
				substring := line[start:end]
				textWidth := b.Font.MeasureText(line[start:end+1], b.FontSize)
				delta := maxWidth - textWidth
				if delta < 0 {
					if wrappoint == -1 {
						if len(b.Lines) == 0 || height+b.FontSize <= maxHeight {
							b.AddMeasuredLine(substring, textWidth, delta)
							height += b.FontSize
						} else {
							log.Println("Cannot fit line:", string(substring))
							break
						}
						end++
					} else {
						substring = line[start:wrappoint]
						if len(b.Lines) == 0 || height+b.FontSize <= maxHeight {
							b.AddMeasuredLine(substring, textWidth, delta)
							height += b.FontSize
						} else {
							log.Println("Cannot fit line:", string(substring))
							break
						}
						end = wrappoint + 1
					}
					start = end
					wrappoint = -1
				} else {
					end++
				}
			}
			if end-start > 0 && height+b.FontSize <= maxHeight {
				substring := line[start:end]
				textWidth := b.Font.MeasureText(substring, b.FontSize)
				delta := maxWidth - textWidth
				b.AddMeasuredLine(substring, textWidth, delta)
				height += b.FontSize
			}
		} else {
			// Shrink To Fit
			for delta < 0 {
				b.FontSize--
				b.OriginY = bounds.Top - b.FontSize
				textWidth = b.Font.MeasureText(line, b.FontSize)
				delta = maxWidth - textWidth
			}
			if len(b.Lines) == 0 || height+b.FontSize <= maxHeight {
				b.AddMeasuredLine(line, textWidth, delta)
				height += b.FontSize
			} else {
				log.Println("Cannot fit line", string(line))
			}
		}
	}
	return &Rectangle{
		Left:   bounds.Left,
		Right:  bounds.Right,
		Top:    bounds.Top,
		Bottom: bounds.Top - height,
	}, nil
}

func (b *TextBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	buffer.WriteString("q\nBT\n")
	buffer.WriteString(fmt.Sprintf("/%s %s Tf\n", b.FontID, FloatToString(b.FontSize)))
	buffer.WriteString(fmt.Sprintf("%s %s %s rg\n", FloatToString(b.FontColour[0]), FloatToString(b.FontColour[1]), FloatToString(b.FontColour[2])))
	buffer.WriteString(fmt.Sprintf("1 0 0 1 %s %s Tm\n", FloatToString(b.OriginX), FloatToString(b.OriginY)))
	buffer.WriteString(fmt.Sprintf("%s TL\n", FloatToString(b.FontSize)))
	for i, l := range b.Lines {
		buffer.WriteString(fmt.Sprintf("%s Tc\n", FloatToString(l.CharacterSpacing)))
		buffer.WriteString(fmt.Sprintf("%s Tw\n", FloatToString(l.WordSpacing)))
		buffer.WriteString(fmt.Sprintf("%s Tz\n", FloatToString(l.HorizontalScaling)))
		buffer.WriteString(fmt.Sprintf("%s 0 Td\n", FloatToString(l.Indent)))
		buffer.WriteString(fmt.Sprintf("%s Ts\n", FloatToString(l.Rise)))
		buffer.WriteString(fmt.Sprintf("%d Tr\n", l.Render))
		if i == 0 { // TODO try with UTF-16BE
			buffer.WriteString(fmt.Sprintf("(%s) Tj\n", l.Text))
		} else {
			buffer.WriteString(fmt.Sprintf("(%s) '\n", l.Text))
		}
	}
	buffer.WriteString("ET\nQ")
	return nil
}

func SplitLines(text []rune) [][]rune {
	var lines [][]rune
	var line []rune
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c == '\r' {
			// Ignore Carriage Return unless followed by New Line
			i++
			if i < len(text) {
				c = text[i]
			} else {
				break
			}
		}
		if c == '\n' {
			lines = append(lines, line)
			line = nil
		} else {
			line = append(line, c)
		}
	}
	if line != nil {
		lines = append(lines, line)
	}
	return lines
}

func PDFEscapeString(text []rune) string {
	var output []rune
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch c {
		case '\\':
			if i < len(text)-1 {
				i++
				output = append(output, c, text[i])
			} else {
				output = append(output, '\\', c)
			}
		case '(':
			fallthrough
		case ')':
			output = append(output, '\\', c)
		default:
			output = append(output, c)
		}

	}
	return string(output)
}
