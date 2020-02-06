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
	FontId           string
	Font             font.Font
	FontSize         float64
	FontColour       []float64
	Align            Alignment
	OriginX, OriginY float64
	Width, Height    float64
	Lines            []*Line
}

func (b *TextBox) AddLine(text []rune) {
	b.AddMeasuredLine(text, b.Font.MeasureText(text, b.FontSize))
}

func (b *TextBox) AddMeasuredLine(text []rune, width float64) {
	var indent float64
	for _, l := range b.Lines {
		indent += l.Indent
	}
	log.Println("NewLine:", len(text), string(text))
	log.Println("TextWidth:", width)
	delta := b.Width - width
	log.Println("Delta:", delta)
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
	log.Println("Indent:", line.Indent)
	b.Lines = append(b.Lines, line)
	b.Height += b.FontSize
	log.Println("TextHeight:", b.Height)
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
	log.Println("Word Spaces:", wordSpaces)
	log.Println("Character Spaces:", characterSpaces)
	if wordSpaces > 1 {
		line.WordSpacing = (delta / 2) / wordSpaces
		line.CharacterSpacing = (delta / 2) / characterSpaces
		//line.HorizontalScaling += (((delta / 3) / (width + delta)) / 15) * 100
	} else {
		line.CharacterSpacing = (delta / 1) / characterSpaces
		//line.HorizontalScaling += (delta / width) * 100
	}
	log.Println("Word Spacing:", line.WordSpacing)
	log.Println("Character Spacing:", line.CharacterSpacing)
	log.Println("Horizontal Scaling:", line.HorizontalScaling)
	return line
}

func (b *TextBox) GetWidth() float64 {
	return b.Width
}

func (b *TextBox) GetHeight() float64 {
	return b.Height
}

func (b *TextBox) SetBounds(bounds *Rectangle) error {
	b.OriginX = bounds.Left
	log.Println("OriginX:", b.OriginX)
	b.OriginY = bounds.Top - b.FontSize
	log.Println("OriginY:", b.OriginY)
	maxWidth := bounds.GetWidth()
	log.Println("MaxWidth:", maxWidth)
	b.Width = maxWidth
	b.Height = 0
	b.Lines = nil
	maxHeight := bounds.GetHeight()
	log.Println("MaxHeight:", maxHeight)
	for _, line := range SplitLines(b.Text) {
		textWidth := b.Font.MeasureText(line, b.FontSize)
		log.Println("TextWidth:", textWidth)
		if textWidth > (maxWidth * 2) {
			// Multiple Lines

			// TODO consider splitting text into words and incrementally add to a list until the width is greater than maxWidth.
			//  If the last word is short, remove from list and put on a new line.
			//  If the last word is long, hyphenate.

			// TODO design a backtracking algorithm to find the most efficient layout with least "badness"

			wrappoint := -1
			start := 0
			end := 0
			for end < len(line) {
				c := line[end]
				if unicode.IsSpace(c) {
					wrappoint = end
				}
				substring := line[start:end]
				if c == '\n' {
					log.Println("New Line Character:", end)
					if b.Height+b.FontSize <= maxHeight {
						b.AddLine(substring)
					} else {
						log.Println("Cannot fit line")
						break
					}
					end++
					start = end
					wrappoint = -1
				} else {
					textWidth := b.Font.MeasureText(line[start:end+1], b.FontSize)
					if textWidth > maxWidth {
						if wrappoint == -1 {
							log.Println("Hard break:", start, end)
							if b.Height+b.FontSize <= maxHeight {
								b.AddLine(substring)
							} else {
								log.Println("Cannot fit line")
								break
							}
							end++
						} else {
							log.Println("Break at wrap point:", start, wrappoint)
							substring = line[start:wrappoint]
							if b.Height+b.FontSize <= maxHeight {
								b.AddLine(substring)
							} else {
								log.Println("Cannot fit line")
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
			}
			if end-start > 0 && b.Height+b.FontSize <= maxHeight {
				b.AddLine(line[start:end])
			}
		} else {
			// Single Line
			for textWidth > maxWidth {
				b.FontSize--
				log.Println("FontSize:", b.FontSize)
				b.OriginY = bounds.Top - b.FontSize
				log.Println("OriginY:", b.OriginY)
				textWidth = b.Font.MeasureText(line, b.FontSize)
				log.Println("TextWidth:", textWidth)
			}
			if b.Height+b.FontSize <= maxHeight {
				b.AddMeasuredLine(line, textWidth)
			} else {
				log.Println("Cannot fit line")
			}
		}
	}
	return nil
}

func (b *TextBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	buffer.WriteString("q\nBT\n")
	buffer.WriteString(fmt.Sprintf("/%s %s Tf\n", b.FontId, FloatToString(b.FontSize)))
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
		if i == 0 {
			buffer.WriteString(fmt.Sprintf("(%s) Tj\n", l.Text))
		} else {
			buffer.WriteString(fmt.Sprintf("(%s) '\n", l.Text))
		}
	}
	buffer.WriteString("ET\nQ\n")
	return nil
}

func SplitLines(text []rune) [][]rune {
	var lines [][]rune
	var line []rune
	for i := 0; i < len(text); i++ {
		c := text[i]
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
