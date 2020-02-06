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

package pdfgo

type Annotation interface {
	GetSubtype() string
	GetRectangle() *ArrayObject
	GetContents() *StringObject
	GetBorder() *ArrayObject
	GetAction() *DictionaryObject
	GetDestination() *DictionaryObject
}

type Hyperlink struct {
	Left,
	Bottom,
	Right,
	Top float64
	URI string
}

func NewHyperlink(left, bottom, right, top float64, uri string) *Hyperlink {
	return &Hyperlink{
		Left:   left,
		Bottom: bottom,
		Right:  right,
		Top:    top,
		URI:    uri,
	}
}

func (h *Hyperlink) GetSubtype() string {
	return "Link"
}

func (h *Hyperlink) GetRectangle() *ArrayObject {
	return &ArrayObject{
		Array: []Object{
			&NumberObject{Number: h.Left},
			&NumberObject{Number: h.Bottom},
			&NumberObject{Number: h.Right},
			&NumberObject{Number: h.Top},
		},
	}
}

func (h *Hyperlink) GetContents() *StringObject {
	return &StringObject{
		String: h.URI,
	}
}

func (h *Hyperlink) GetBorder() *ArrayObject {
	return &ArrayObject{
		Array: []Object{
			&NumberObject{Number: 0},
			&NumberObject{Number: 0},
			&NumberObject{Number: 0},
		},
	}
}

func (h *Hyperlink) GetAction() *DictionaryObject {
	a := &DictionaryObject{
		Dictionary: make(map[*NameObject]Object),
	}
	a.AddNameNameEntry("Type", "Action")
	a.AddNameNameEntry("S", "URI")
	a.AddNameObjectEntry("URI", &StringObject{
		String: h.URI,
	})
	return a
}

func (h *Hyperlink) GetDestination() *DictionaryObject {
	return nil
}
