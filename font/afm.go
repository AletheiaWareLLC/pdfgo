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

package font

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
)

type Direction struct {
	Direction          int
	UnderlinePosition  float64
	UnderlineThickness float64
	ItalicAngle        float64
	CharWidth          [2]float64
	IsFixedPitch       bool
}

type CharacterMetric struct {
	Code     int
	WidthX   [2]float64
	WidthY   [2]float64
	VVector  [2]float64
	BBox     [4]float64
	Ligature map[string]string
}

type TrackKern struct {
	MinSize, MaxSize float64
	MinKern, MaxKern float64
}

type KernPair struct {
	X float64
	Y float64
}

type KernData struct {
	TrackKern map[int]*TrackKern
	KernPairs [2]map[string]map[string]*KernPair
}

type CompositePart struct {
	Name string
	X    float64
	Y    float64
}

type Composite struct {
	Parts []*CompositePart
}

type AFM struct {
	FileVersion    float64
	FontName       string
	FullName       string
	FamilyName     string
	Weight         string
	FontBBox       [4]float64
	FontVersion    string
	Notice         string
	EncodingScheme string
	MappingScheme  int
	EscChar        int
	CharacterSet   string
	Characters     int
	IsBaseFont     bool
	VVector        [2]float64
	IsFixedV       bool
	Directions     [2]*Direction
	CapHeight      float64
	XHeight        float64
	Ascender       float64
	Descender      float64
	StdHW          float64
	StdVW          float64
	Metrics        map[string]*CharacterMetric
	KernData       *KernData
	Composites     map[string]*Composite
}

func ReadAFM(reader io.Reader) (*AFM, error) {
	var afm *AFM
	var direction *Direction
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "Comment":
			// Ignored
		case "StartFontMetrics":
			version, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			direction = &Direction{
				Direction: 0,
			}
			afm = &AFM{
				FileVersion: version,
				Directions:  [2]*Direction{direction, nil},
				Metrics:     make(map[string]*CharacterMetric),
				Composites:  make(map[string]*Composite),
			}
		case "EndFontMetrics":
			break
		case "FontName":
			afm.FontName = fields[1]
		case "FullName":
			afm.FullName = fields[1]
		case "FamilyName":
			afm.FamilyName = fields[1]
		case "Weight":
			afm.Weight = fields[1]
		case "FontBBox":
			llx, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			lly, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, err
			}
			urx, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			ury, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				return nil, err
			}
			afm.FontBBox = [4]float64{llx, lly, urx, ury}
		case "Version":
			afm.FontVersion = fields[1]
		case "Notice":
			afm.Notice = fields[1]
		case "EncodingScheme":
			afm.EncodingScheme = fields[1]
		case "MappingScheme":
			scheme, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			afm.MappingScheme = scheme
		case "EscChar":
			esc, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			afm.EscChar = esc
		case "CharacterSet":
			afm.CharacterSet = fields[1]
		case "Characters":
			characters, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			afm.Characters = characters
		case "IsBaseFont":
			b, err := strconv.ParseBool(fields[1])
			if err != nil {
				return nil, err
			}
			afm.IsBaseFont = b
		case "VVector":
			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, err
			}
			afm.VVector = [2]float64{x, y}
		case "IsFixedV":
			b, err := strconv.ParseBool(fields[1])
			if err != nil {
				return nil, err
			}
			afm.IsFixedV = b
		case "StartDirection":
			d, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			if direction == nil || direction.Direction != d {
				direction = &Direction{
					Direction: d,
				}
				afm.Directions[d] = direction
			}
		case "EndDirection":
			direction = nil
		case "UnderlinePosition":
			underlineposition, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			direction.UnderlinePosition = underlineposition
		case "UnderlineThickness":
			underlinethickness, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			direction.UnderlineThickness = underlinethickness
		case "ItalicAngle":
			italicangle, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			direction.ItalicAngle = italicangle
		case "CharWidth":
			x, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				return nil, err
			}
			direction.CharWidth = [2]float64{x, y}
		case "IsFixedPitch":
			b, err := strconv.ParseBool(fields[1])
			if err != nil {
				return nil, err
			}
			direction.IsFixedPitch = b
		case "CapHeight":
			capheight, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.CapHeight = capheight
		case "XHeight":
			xheight, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.XHeight = xheight
		case "Ascender":
			ascender, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.Ascender = ascender
		case "Descender":
			descender, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.Descender = descender
		case "StdHW":
			stdhw, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.StdHW = stdhw
		case "StdVW":
			stdvw, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return nil, err
			}
			afm.StdVW = stdvw
		case "StartCharMetrics":
			count, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			for i := 0; i < count && scanner.Scan(); i++ {
				line := scanner.Text()
				if len(line) == 0 {
					continue
				}
				metric := &CharacterMetric{}
				for _, part := range strings.Split(line, ";") {
					fields := strings.Fields(part)
					if len(fields) == 0 {
						continue
					}
					switch fields[0] {
					case "C":
						code, err := strconv.Atoi(fields[1])
						if err != nil {
							return nil, err
						}
						metric.Code = code
					case "CH":
						log.Println("TODO parse Hex", line)
					case "WX":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX = [2]float64{w, w}
					case "W0X":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX[0] = w
					case "W1X":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX[1] = w
					case "WY":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthY = [2]float64{w, w}
					case "W0Y":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthY[0] = w
					case "W1Y":
						w, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthY[1] = w
					case "W":
						x, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(fields[2], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX = [2]float64{x, x}
						metric.WidthY = [2]float64{y, y}
					case "W0":
						x, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(fields[2], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX[0] = x
						metric.WidthY[0] = y
					case "W1":
						x, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(fields[2], 64)
						if err != nil {
							return nil, err
						}
						metric.WidthX[1] = x
						metric.WidthY[1] = y
					case "VV":
						x, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(fields[2], 64)
						if err != nil {
							return nil, err
						}
						metric.VVector = [2]float64{x, y}
					case "N":
						afm.Metrics[fields[1]] = metric
					case "B":
						llx, err := strconv.ParseFloat(fields[1], 64)
						if err != nil {
							return nil, err
						}
						lly, err := strconv.ParseFloat(fields[2], 64)
						if err != nil {
							return nil, err
						}
						urx, err := strconv.ParseFloat(fields[3], 64)
						if err != nil {
							return nil, err
						}
						ury, err := strconv.ParseFloat(fields[4], 64)
						if err != nil {
							return nil, err
						}
						metric.BBox = [4]float64{llx, lly, urx, ury}
					case "L":
						if metric.Ligature == nil {
							metric.Ligature = make(map[string]string)
						}
						metric.Ligature[fields[1]] = fields[2]
					default:
						return nil, errors.New("Unrecognized Character Metric: " + line)
					}
				}
			}
			if scanner.Scan() && scanner.Text() != "EndCharMetrics" {
				return nil, errors.New("Expected EndCharMetrics")
			}
		case "StartKernData":
			afm.KernData = &KernData{}
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) == 0 {
					continue
				}
				if line == "EndKernData" {
					break
				}
				fields := strings.Fields(line)
				if len(fields) == 0 {
					continue
				}
				switch fields[0] {
				case "StartTrackKern":
					count, err := strconv.Atoi(fields[1])
					if err != nil {
						return nil, err
					}
					for i := 0; i < count && scanner.Scan(); i++ {
						line := scanner.Text()
						if len(line) == 0 {
							continue
						}
						fields := strings.Fields(line)
						if len(fields) == 0 {
							continue
						}
						if fields[0] == "TrackKern" {
							degree, err := strconv.Atoi(fields[1])
							if err != nil {
								return nil, err
							}
							minSize, err := strconv.ParseFloat(fields[2], 64)
							if err != nil {
								return nil, err
							}
							minKern, err := strconv.ParseFloat(fields[3], 64)
							if err != nil {
								return nil, err
							}
							maxSize, err := strconv.ParseFloat(fields[4], 64)
							if err != nil {
								return nil, err
							}
							maxKern, err := strconv.ParseFloat(fields[5], 64)
							if err != nil {
								return nil, err
							}
							afm.KernData.TrackKern[degree] = &TrackKern{
								MinSize: minSize,
								MinKern: minKern,
								MaxSize: maxSize,
								MaxKern: maxKern,
							}
						} else {
							return nil, errors.New("Unrecognized Track Kern " + line)
						}
					}
					if scanner.Scan() && scanner.Text() != "EndTrackKern" {
						return nil, errors.New("Expected EndTrackKern")
					}
				case "StartKernPairs":
					count, err := strconv.Atoi(fields[1])
					if err != nil {
						return nil, err
					}
					pairs, err := ReadKernPairs(scanner, count)
					if err != nil {
						return nil, err
					}
					afm.KernData.KernPairs[0] = pairs
					if scanner.Scan() && scanner.Text() != "EndKernPairs" {
						return nil, errors.New("Expected EndKernPairs")
					}
				case "StartKernPairs0":
					count, err := strconv.Atoi(fields[1])
					if err != nil {
						return nil, err
					}
					pairs, err := ReadKernPairs(scanner, count)
					if err != nil {
						return nil, err
					}
					afm.KernData.KernPairs[0] = pairs
					if scanner.Scan() && scanner.Text() != "EndKernPairs0" {
						return nil, errors.New("Expected EndKernPairs0")
					}
				case "StartKernPairs1":
					count, err := strconv.Atoi(fields[1])
					if err != nil {
						return nil, err
					}
					pairs, err := ReadKernPairs(scanner, count)
					if err != nil {
						return nil, err
					}
					afm.KernData.KernPairs[1] = pairs
					if scanner.Scan() && scanner.Text() != "EndKernPairs1" {
						return nil, errors.New("Expected EndKernPairs1")
					}
				default:
					return nil, errors.New("Unrecognized Kern Data: " + line)
				}
			}
		case "StartComposites":
			count, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			for i := 0; i < count && scanner.Scan(); i++ {
				line := scanner.Text()
				if len(line) == 0 {
					continue
				}
				for _, part := range strings.Split(line, ";") {
					fields := strings.Fields(part)
					if len(fields) == 0 {
						continue
					}
					var composite *Composite
					switch fields[0] {
					case "CC":
						composite = &Composite{}
						name := fields[1]
						afm.Composites[name] = composite
					case "PCC":
						name := fields[1]
						x, err := strconv.ParseFloat(fields[3], 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(fields[4], 64)
						if err != nil {
							return nil, err
						}
						composite.Parts = append(composite.Parts, &CompositePart{
							Name: name,
							X:    x,
							Y:    y,
						})
					default:
						return nil, errors.New("Unrecognized Composite: " + line)
					}
				}
			}
			if scanner.Scan() && scanner.Text() != "EndComposites" {
				return nil, errors.New("Expected EndComposites")
			}
		default:
			return nil, errors.New("Unrecognized Line: " + line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return afm, nil
}

func ReadKernPairs(scanner *bufio.Scanner, count int) (map[string]map[string]*KernPair, error) {
	pairs := make(map[string]map[string]*KernPair)
	for i := 0; i < count && scanner.Scan(); i++ {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "KP":
			first := fields[1]
			second := fields[2]
			x, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				return nil, err
			}
			p, ok := pairs[first]
			if !ok {
				p = make(map[string]*KernPair)
				pairs[first] = p
			}
			p[second] = &KernPair{
				X: x,
				Y: y,
			}
		case "KPH":
			log.Println("TODO parse Hex", line)
		case "KPX":
			first := fields[1]
			second := fields[2]
			x, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			p, ok := pairs[first]
			if !ok {
				p = make(map[string]*KernPair)
				pairs[first] = p
			}
			p[second] = &KernPair{
				X: x,
			}
		case "KPY":
			first := fields[1]
			second := fields[2]
			y, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				return nil, err
			}
			p, ok := pairs[first]
			if !ok {
				p = make(map[string]*KernPair)
				pairs[first] = p
			}
			p[second] = &KernPair{
				Y: y,
			}
		default:
			return nil, errors.New("Unrecognized Kern Pair: " + line)
		}
	}
	return pairs, nil
}

func (a *AFM) GetName(r rune) (string, error) {
	c := uint32(r)
	if a.Metrics != nil {
		for n, m := range a.Metrics {
			if uint32(m.Code) == c {
				return n, nil
			}
		}
	}
	return "", errors.New("Cannot get name for " + string(r) + " (" + strconv.Itoa(int(c)) + ")")
}

func (a *AFM) GetKernPair(direction int, first, last string) *KernPair {
	kd := a.KernData
	if kd != nil {
		kp := kd.KernPairs[direction]
		if kp != nil {
			k := kp[first]
			if k != nil {
				return k[last]
			}
		}
	}
	return nil
}

func (a *AFM) GetWidth(direction int, name string) float64 {
	m, ok := a.Metrics[name]
	if ok {
		return m.WidthX[direction]
	}
	return 500
}
