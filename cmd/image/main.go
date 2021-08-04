package main

import (
	"bytes"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/graphics"
	"io/ioutil"
	"log"
	"flag"
	"os"
)

var jpeg = flag.String("jpeg", "", "the jpeg to include")

func main() {
	flag.ParseFlags()

	p := pdfgo.NewPDF()

	// Create Resources
	resources := p.NewDictionaryObject()
	xs := p.NewDictionaryObject()
	resources.AddNameObjectEntry("XObject", pdfgo.NewObjectReference(xs))

	data, err := ioutil.ReadFile(*jpeg)
	if err != nil {
		log.Fatal(err)
	}
	ir, w, h, err := p.AddImage("image/jpg", data)
	if err != nil {
		log.Fatal(err)
	}
	id := "/img"
	xs.AddNameObjectEntry(id, ir)

	// Create Contents
	box := &graphics.ImageBox{
		ImageID: id,
		Width:   w,
		Height:  h,
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

	// Write
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
