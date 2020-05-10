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
	"github.com/AletheiaWareLLC/pdfgo"
)

type Type1Font struct {
	Name      string
	Reference *pdfgo.ObjectReference
}

func NewType1Font(p *pdfgo.PDF, name string) (*Type1Font, error) {
	font := p.NewDictionaryObject()
	font.AddNameNameEntry("Type", "Font")
	font.AddNameNameEntry("Subtype", "Type1")
	font.AddNameNameEntry("BaseFont", name)
	font.AddNameObjectEntry("FirstChar", &pdfgo.NumberObject{
		Number: 32,
	})
	font.AddNameObjectEntry("LastChar", &pdfgo.NumberObject{
		Number: 255,
	})
	font.AddNameNameEntry("Encoding", "") // "MacRomanEncoding" or "WinAnsiEncoding"
	var widths []pdfgo.Object
	/* TODO
	   for i, w := range f.Widths {
	       widths = append(widths, &pdfgo.NumberObject{Number: 0})
	   }
	*/
	width := p.NewArrayObject(widths)
	font.AddNameObjectEntry("Widths", pdfgo.NewObjectReference(width))
	descriptor := p.NewDictionaryObject()
	descriptor.AddNameNameEntry("Type", "FontDescriptor")
	descriptor.AddNameNameEntry("FontName", name)
	/*
	   descriptor.AddNameNameEntry("Ascent",)
	   descriptor.AddNameNameEntry("Descent",)
	   descriptor.AddNameNameEntry("CapHeight",)
	   descriptor.AddNameNameEntry("Flags",)
	   descriptor.AddNameNameEntry("FontBBox",)
	   descriptor.AddNameNameEntry("ItalicAngle",)
	   descriptor.AddNameNameEntry("Leading",)
	   descriptor.AddNameNameEntry("StemV",)
	   descriptor.AddNameNameEntry("StemH",)
	   descriptor.AddNameNameEntry("AvgWidth",)
	   descriptor.AddNameNameEntry("MaxWidth",)
	   descriptor.AddNameNameEntry("MissingWidth",)
	   descriptor.AddNameNameEntry("FontFile",)
	*/
	font.AddNameObjectEntry("FontDescriptor", pdfgo.NewObjectReference(descriptor))
	return &Type1Font{
		Name:      name,
		Reference: pdfgo.NewObjectReference(font),
	}, nil
}

func (f *Type1Font) GetReference() *pdfgo.ObjectReference {
	return f.Reference
}

func (f *Type1Font) MeasureText(text []rune, fontSize float64) float64 {
	// TODO
	return 0.5 * float64(len(text)) * fontSize
}
