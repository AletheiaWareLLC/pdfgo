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

package main

import (
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/font"
	"log"
	"os"
)

func main() {
	p := pdfgo.NewPDF()
	width := 595.28
	height := 841.89
	f1, err := font.NewTrueTypeFont(p, "NotoSerif-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	fonts := p.NewDictionaryObject()
	fonts.AddNameObjectEntry("F1", f1.GetReference())
	resources := p.NewDictionaryObject()
	resources.AddNameObjectEntry("Font", pdfgo.NewObjectReference(fonts))
	contents := p.NewStreamObject()
	contents.Data = []byte("BT\n0.1 0.1 0.9 rg\n100 100 395.28 641.89 re f\n0.9 0.1 0.1 RG\n0 0 m 595.28 841.89 l h S\n0.1 0.9 0.1 rg\n/F1 24 Tf\n297.64 420.945 Td\n(Hello World!) Tj\nET")
	p.AddPage(width, height, pdfgo.NewObjectReference(resources), pdfgo.NewObjectReference(contents))
	writer := os.Stdout
	if len(os.Args) > 1 {
		log.Println("Writing:", os.Args[1])
		file, err := os.OpenFile(os.Args[1], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer = file
	}
	p.Write(writer)
}
