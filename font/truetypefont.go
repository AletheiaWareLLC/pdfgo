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
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"io/ioutil"
	"log"
)

type TrueTypeFont struct {
	Reference *pdfgo.ObjectReference
	Font      *truetype.Font
}

func NewTrueTypeFont(p *pdfgo.PDF, file string) (*TrueTypeFont, error) {
	log.Println("Loading Font:", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	/**/
	for i, n := range []truetype.NameID{
		truetype.NameIDCopyright,
		truetype.NameIDFontFamily,
		truetype.NameIDFontSubfamily,
		truetype.NameIDUniqueSubfamilyID,
		truetype.NameIDFontFullName,
		truetype.NameIDNameTableVersion,
		truetype.NameIDPostscriptName,
		truetype.NameIDTrademarkNotice,
		truetype.NameIDManufacturerName,
		truetype.NameIDDesignerName,
		truetype.NameIDFontDescription,
		truetype.NameIDFontVendorURL,
		truetype.NameIDFontDesignerURL,
		truetype.NameIDFontLicense,
		truetype.NameIDFontLicenseURL,
		truetype.NameIDPreferredFamily,
		truetype.NameIDPreferredSubfamily,
		truetype.NameIDCompatibleName,
		truetype.NameIDSampleText,
	} {
		name := f.Name(n)
		if name != "" {
			log.Println("Font Name:", i, name)
		}
	}
	/**/
	basename := f.Name(truetype.NameIDPostscriptName)
	if basename == "" {
		basename = f.Name(truetype.NameIDFontFamily)
	}
	font := p.NewDictionaryObject()
	font.AddNameNameEntry("Type", "Font")
	font.AddNameNameEntry("Subtype", "TrueType")
	font.AddNameNameEntry("BaseFont", basename)
	// TODO font.AddNameNameEntry("Encoding", "") // "MacRomanEncoding" or "WinAnsiEncoding"
	var (
		end      int
		widths   []pdfgo.Object
		sum, max float64
	)
	scale := fixed.Int26_6(f.FUnitsPerEm())
	for i := 32; i < 65536; i++ {
		index := f.Index(rune(i))
		if index == 0 {
			continue
		}
		end = i
		width := float64(f.HMetric(scale, index).AdvanceWidth)
		sum += width
		if width > max {
			max = width
		}
		//log.Println(i, string(rune(i)), index, width)
		widths = append(widths, &pdfgo.NumberObject{
			Number: width,
		})
	}
	font.AddNameObjectEntry("FirstChar", &pdfgo.NumberObject{
		Number: 32,
	})
	font.AddNameObjectEntry("LastChar", &pdfgo.NumberObject{
		Number: float64(end),
	})
	width := p.NewArrayObject(widths)
	font.AddNameObjectEntry("Widths", pdfgo.NewObjectReference(width))
	descriptor := p.NewDictionaryObject()
	descriptor.AddNameNameEntry("Type", "FontDescriptor")
	descriptor.AddNameNameEntry("FontName", basename)

	// TODO descriptor.AddNameNameEntry("Ascent",)
	// TODO descriptor.AddNameNameEntry("Descent",)
	// TODO descriptor.AddNameNameEntry("CapHeight",)
	// TODO descriptor.AddNameNameEntry("Flags", integer)
	bounds := f.Bounds(scale)
	descriptor.AddNameObjectEntry("FontBBox", &pdfgo.ArrayObject{
		Array: []pdfgo.Object{
			&pdfgo.NumberObject{Number: float64(bounds.Min.X)},
			&pdfgo.NumberObject{Number: float64(bounds.Min.Y)},
			&pdfgo.NumberObject{Number: float64(bounds.Max.X)},
			&pdfgo.NumberObject{Number: float64(bounds.Max.Y)},
		},
	})
	// TODO descriptor.AddNameNameEntry("ItalicAngle", integer)
	// TODO descriptor.AddNameNameEntry("Leading",)
	// TODO descriptor.AddNameNameEntry("StemV",)
	// TODO descriptor.AddNameNameEntry("StemH",)
	average := sum / (224)
	descriptor.AddNameObjectEntry("AvgWidth", &pdfgo.NumberObject{
		Number: average,
	})
	descriptor.AddNameObjectEntry("MaxWidth", &pdfgo.NumberObject{
		Number: max,
	})
	descriptor.AddNameObjectEntry("MissingWidth", &pdfgo.NumberObject{
		Number: average,
	})
	stream := p.NewStreamObject()
	stream.Data = b
	descriptor.AddNameObjectEntry("FontFile2", pdfgo.NewObjectReference(stream))
	font.AddNameObjectEntry("FontDescriptor", pdfgo.NewObjectReference(descriptor))
	return &TrueTypeFont{
		Reference: pdfgo.NewObjectReference(font),
		Font:      f,
	}, nil
}

func (f *TrueTypeFont) GetReference() *pdfgo.ObjectReference {
	return f.Reference
}

func (f *TrueTypeFont) MeasureText(text []rune, fontSize float64) float64 {
	scale := fixed.Int26_6(f.Font.FUnitsPerEm())
	var previous truetype.Index
	var width fixed.Int26_6
	for i, c := range text {
		index := f.Font.Index(c)
		if i > 0 {
			kern := f.Font.Kern(scale, previous, index)
			width += kern
		}
		awidth := f.Font.HMetric(scale, index).AdvanceWidth
		width += awidth
		previous = index
	}
	return float64(width) * fontSize / float64(f.Font.FUnitsPerEm())
}
