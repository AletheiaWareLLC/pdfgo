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

package pdfgo_test

import (
	"bytes"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/font"
	"github.com/AletheiaWareLLC/pdfgo/graphics"
	"log"
	"os"
)

func ExamplePDF_empty_page() {
	p := pdfgo.NewPDF()

	// Create Empty Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, nil, nil)

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj
	// 3 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots []>> endobj
	// xref
	// 0 4
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// trailer <</Size 4 /Root 1 0 R>>
	// startxref
	// 191
	// %%EOF
}

func ExamplePDF_hyperlink() {
	p := pdfgo.NewPDF()

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, nil, nil)

	p.AddAnnotation(pdfgo.NewHyperlink(10, 10, 100, 100, "https://example.com"))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [3 0 R] /Count 1>> endobj
	// 3 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [4 0 R]>> endobj
	// 4 0 obj <</Type /Annot /Subtype /Link /Rect [10 10 100 100] /Contents (https://example.com) /Border [0 0 0] /A <</Type /Action /S /URI /URI (https://example.com)>>>> endobj
	// xref
	// 0 5
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000196 00000 n
	// trailer <</Size 5 /Root 1 0 R>>
	// startxref
	// 369
	// %%EOF
}

func ExamplePDF_stream_colourbox() {
	p := pdfgo.NewPDF()

	box := &graphics.ColourBox{
		BorderColour: []float64{1.0, 0.0, 0.0},
		FillColour:   []float64{0.0, 0.0, 1.0},
	}

	if _, err := box.SetBounds(&graphics.Rectangle{
		Left:   50,
		Right:  350,
		Top:    550,
		Bottom: 50,
	}); err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := box.Write(p, &buffer); err != nil {
		log.Fatal(err)
	}
	contents := p.NewStreamObject()
	contents.Data = buffer.Bytes()

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, nil, pdfgo.NewObjectReference(contents))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [4 0 R] /Count 1>> endobj
	// 3 0 obj <</Length 59>>
	// stream
	// q
	// 0 0 1 rg
	// 50 50 300 500 re f
	// 1 0 0 RG
	// 50 50 300 500 re S
	// Q
	// endstream endobj
	// 4 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [] /Contents 3 0 R>> endobj
	// xref
	// 0 5
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000218 00000 n
	// trailer <</Size 5 /Root 1 0 R>>
	// startxref
	// 314
	// %%EOF
}

func ExamplePDF_stream_line() {
	p := pdfgo.NewPDF()

	// Create Stream to Draw a Line
	contents := p.NewStreamObject()
	contents.Data = []byte("BT\n0 0 m\n400 600 l h S\nET")

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, nil, pdfgo.NewObjectReference(contents))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [4 0 R] /Count 1>> endobj
	// 3 0 obj <</Length 25>>
	// stream
	// BT
	// 0 0 m
	// 400 600 l h S
	// ET
	// endstream endobj
	// 4 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [] /Contents 3 0 R>> endobj
	// xref
	// 0 5
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000184 00000 n
	// trailer <</Size 5 /Root 1 0 R>>
	// startxref
	// 280
	// %%EOF
}

func ExamplePDF_stream_rectangle() {
	p := pdfgo.NewPDF()

	// Create Stream to Draw a Rectangle
	contents := p.NewStreamObject()
	contents.Data = []byte("BT\n0.1 0.1 0.9 rg\n100 100 200 400 re f\nET")

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, nil, pdfgo.NewObjectReference(contents))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [4 0 R] /Count 1>> endobj
	// 3 0 obj <</Length 41>>
	// stream
	// BT
	// 0.1 0.1 0.9 rg
	// 100 100 200 400 re f
	// ET
	// endstream endobj
	// 4 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [] /Contents 3 0 R>> endobj
	// xref
	// 0 5
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000200 00000 n
	// trailer <</Size 5 /Root 1 0 R>>
	// startxref
	// 296
	// %%EOF
}

func ExamplePDF_stream_text() {
	p := pdfgo.NewPDF()

	// Create Font
	f1, err := font.NewCoreFont(p, "Helvetica")
	if err != nil {
		log.Fatal(err)
	}
	fonts := p.NewDictionaryObject()
	fonts.AddNameObjectEntry("F1", f1.GetReference())

	// Create Resources
	resources := p.NewDictionaryObject()
	resources.AddNameObjectEntry("Font", pdfgo.NewObjectReference(fonts))

	// Create Stream to Display "Hello World!"
	contents := p.NewStreamObject()
	contents.Data = []byte("BT\n/F1 24 Tf\n200 300 Td\n(Hello World!) Tj\nET")

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, pdfgo.NewObjectReference(resources), pdfgo.NewObjectReference(contents))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [7 0 R] /Count 1>> endobj
	// 3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding>> endobj
	// 4 0 obj <</F1 3 0 R>> endobj
	// 5 0 obj <</Font 4 0 R>> endobj
	// 6 0 obj <</Length 44>>
	// stream
	// BT
	// /F1 24 Tf
	// 200 300 Td
	// (Hello World!) Tj
	// ET
	// endstream endobj
	// 7 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [] /Resources 5 0 R /Contents 6 0 R>> endobj
	// xref
	// 0 8
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000206 00000 n
	// 0000000235 00000 n
	// 0000000266 00000 n
	// 0000000358 00000 n
	// trailer <</Size 8 /Root 1 0 R>>
	// startxref
	// 471
	// %%EOF
}

func ExamplePDF_stream_textbox() {
	p := pdfgo.NewPDF()

	// Create Font
	f1, err := font.NewCoreFont(p, "Helvetica")
	if err != nil {
		log.Fatal(err)
	}
	fonts := p.NewDictionaryObject()
	fonts.AddNameObjectEntry("F1", f1.GetReference())

	// Create Resources
	resources := p.NewDictionaryObject()
	resources.AddNameObjectEntry("Font", pdfgo.NewObjectReference(fonts))

	// Create Contents
	box := &graphics.TextBox{
		Text:       []rune("Hello World!"),
		FontID:     "F1",
		Font:       f1,
		FontSize:   32,
		FontColour: []float64{1.0, 0.0, 0.0},
		Align:      graphics.Center,
	}

	if _, err := box.SetBounds(&graphics.Rectangle{
		Left:   50,
		Right:  350,
		Top:    550,
		Bottom: 50,
	}); err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := box.Write(p, &buffer); err != nil {
		log.Fatal(err)
	}
	contents := p.NewStreamObject()
	contents.Data = buffer.Bytes()

	// Create Page
	width := 400.0
	height := 600.0
	p.AddPage(width, height, pdfgo.NewObjectReference(resources), pdfgo.NewObjectReference(contents))

	// Write to Standard Out
	p.Write(os.Stdout)

	// Output:
	// %PDF-1.7
	// 1 0 obj <</Type /Catalog /Pages 2 0 R>> endobj
	// 2 0 obj <</Type /Pages /Kids [7 0 R] /Count 1>> endobj
	// 3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding>> endobj
	// 4 0 obj <</F1 3 0 R>> endobj
	// 5 0 obj <</Font 4 0 R>> endobj
	// 6 0 obj <</Length 121>>
	// stream
	// q
	// BT
	// /F1 32 Tf
	// 1 0 0 rg
	// 1 0 0 1 50 518 Tm
	// 32 TL
	// 0 Tc
	// 0 Tw
	// 100 Tz
	// 63.760000000000005 0 Td
	// 0 Ts
	// 0 Tr
	// (Hello World!) Tj
	// ET
	// Q
	// endstream endobj
	// 7 0 obj <</Type /Page /Parent 2 0 R /MediaBox [0 0 400 600] /Annots [] /Resources 5 0 R /Contents 6 0 R>> endobj
	// xref
	// 0 8
	// 0000000000 65535 f
	// 0000000009 00000 n
	// 0000000056 00000 n
	// 0000000111 00000 n
	// 0000000206 00000 n
	// 0000000235 00000 n
	// 0000000266 00000 n
	// 0000000436 00000 n
	// trailer <</Size 8 /Root 1 0 R>>
	// startxref
	// 549
	// %%EOF
}
