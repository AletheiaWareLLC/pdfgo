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
	"github.com/AletheiaWareLLC/pdfgo"
)

type MaxLayout struct {
	Boxes []Box
}

func (l *MaxLayout) Add(box Box) {
	l.Boxes = append(l.Boxes, box)
}

func (l *MaxLayout) SetBounds(bounds *Rectangle) (*Rectangle, error) {
	for _, b := range l.Boxes {
		if _, err := b.SetBounds(bounds); err != nil {
			return nil, err
		}
	}
	return bounds, nil
}

func (l *MaxLayout) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	for _, b := range l.Boxes {
		if err := b.Write(p, buffer); err != nil {
			return err
		}
	}
	return nil
}
