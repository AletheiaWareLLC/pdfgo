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
	"archive/zip"
	"errors"
	"github.com/AletheiaWareLLC/pdfgo"
	"log"
)

const (
	CORE_FONT_AFM_ZIP_URL = "ftp://ftp.adobe.com/pub/adobe/devnet/font/pdfs/Core14_AFMs.zip"
	CORE_FONT_AFM_ZIP     = "Core14_AFMs.zip"
)

type CoreFont struct {
	Name      string
	Metrics   *AFM
	Reference *pdfgo.ObjectReference
}

func NewCoreFont(p *pdfgo.PDF, name string) (*CoreFont, error) {
	// Open Core Font AFM ZIP
	r, err := zip.OpenReader(CORE_FONT_AFM_ZIP)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	// Load Font Metrics
	var metrics *AFM
	for _, f := range r.File {
		if f.Name == name+".afm" {
			file, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()
			metrics, err = ReadAFM(file)
			if err != nil {
				return nil, err
			}
		}
	}
	if metrics == nil {
		return nil, errors.New("Could not get font: " + name)
	}
	font := p.NewDictionaryObject()
	font.AddNameNameEntry("Type", "Font")
	font.AddNameNameEntry("Subtype", "Type1")
	font.AddNameNameEntry("BaseFont", name)
	if name != "Symbol" && name != "ZapfDingbats" {
		font.AddNameNameEntry("Encoding", "WinAnsiEncoding")
	}
	return &CoreFont{
		Name:      name,
		Metrics:   metrics,
		Reference: pdfgo.NewObjectReference(font),
	}, nil
}

func (f *CoreFont) GetReference() *pdfgo.ObjectReference {
	return f.Reference
}

func (f *CoreFont) MeasureText(text []rune, fontSize float64) float64 {
	var width float64
	var previous string
	for i, c := range text {
		name, err := f.Metrics.GetName(c)
		if err != nil {
			log.Fatal(i, err)
		}
		if i > 0 {
			kp := f.Metrics.GetKernPair(0, previous, name)
			if kp != nil {
				width += kp.X
			}
		}
		w := f.Metrics.GetWidth(0, name)
		width += w
		previous = name
	}
	return width * fontSize / 1000
}
