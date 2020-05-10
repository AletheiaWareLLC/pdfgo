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
	descriptor.AddNameNameEntry("FontName", basename)
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
	   descriptor.AddNameNameEntry("FontFile2", /Subtype "TrueType")
	*/
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
