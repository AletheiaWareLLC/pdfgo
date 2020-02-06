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

import (
	"fmt"
	"io"
	"log"
)

type PDF struct {
	Version        string
	Catalog        *DictionaryObject
	Pages          *ArrayObject
	PagesReference *ObjectReference
	PageCount      *NumberObject
	Annotations    *ArrayObject
	Objects        []Object
}

func NewPDF() *PDF {
	p := &PDF{
		Version:     "1.7",
		Pages:       &ArrayObject{},
		PageCount:   &NumberObject{},
		Annotations: &ArrayObject{},
	}
	p.Catalog = p.NewDictionaryObject()
	p.Catalog.AddNameNameEntry("Type", "Catalog")
	ps := p.NewDictionaryObject()
	ps.AddNameNameEntry("Type", "Pages")
	ps.AddNameObjectEntry("Kids", p.Pages)
	ps.AddNameObjectEntry("Count", p.PageCount) // Leaf node count
	p.PagesReference = NewObjectReference(ps)
	p.Catalog.AddNameObjectEntry("Pages", p.PagesReference)
	return p
}

func (p *PDF) AddPage(width, height float64, resources, contents Object) {
	page := p.NewDictionaryObject()
	page.AddNameNameEntry("Type", "Page")
	page.AddNameObjectEntry("Parent", p.PagesReference)
	page.AddNameObjectEntry("MediaBox", &ArrayObject{
		Array: []Object{
			&NumberObject{Number: 0},
			&NumberObject{Number: 0},
			&NumberObject{Number: width},
			&NumberObject{Number: height},
		},
	})
	page.AddNameObjectEntry("Annots", p.Annotations)
	if resources != nil {
		page.AddNameObjectEntry("Resources", resources)
	}
	if contents != nil {
		page.AddNameObjectEntry("Contents", contents)
	}
	p.Pages.Array = append(p.Pages.Array, NewObjectReference(page))
	p.PageCount.Number = float64(len(p.Pages.Array))
}

func NewObjectReference(object Object) *ObjectReference {
	return &ObjectReference{
		Object: object,
	}
}

func (p *PDF) NewArrayObject(array []Object) *ArrayObject {
	o := &ArrayObject{
		Array: array,
	}
	p.add(o)
	return o
}

func (p *PDF) NewBooleanObject(boolean bool) *BooleanObject {
	o := &BooleanObject{
		Boolean: boolean,
	}
	p.add(o)
	return o
}

func (p *PDF) NewDictionaryObject() *DictionaryObject {
	o := &DictionaryObject{
		Dictionary: make(map[*NameObject]Object),
	}
	p.add(o)
	return o
}

func (p *PDF) NewNameObject(name string) *NameObject {
	o := &NameObject{
		Name: name,
	}
	p.add(o)
	return o
}

func (p *PDF) NewNumberObject(number float64) *NumberObject {
	o := &NumberObject{
		Number: number,
	}
	p.add(o)
	return o
}

func (p *PDF) NewStreamObject() *StreamObject {
	o := &StreamObject{}
	p.add(o)
	return o
}

func (p *PDF) NewStringObject(str string) *StringObject {
	o := &StringObject{
		String: str,
	}
	p.add(o)
	return o
}

func (p *PDF) add(object Object) {
	p.Objects = append(p.Objects, object)
	object.SetName(len(p.Objects))
}

func (p *PDF) AddAnnotation(annotation Annotation) {
	a := p.NewDictionaryObject()
	a.AddNameNameEntry("Type", "Annot")
	a.AddNameNameEntry("Subtype", annotation.GetSubtype())
	a.AddNameObjectEntry("Rect", annotation.GetRectangle())
	a.AddNameObjectEntry("Contents", annotation.GetContents())
	a.AddNameObjectEntry("Border", annotation.GetBorder())
	action := annotation.GetAction()
	if action != nil {
		a.AddNameObjectEntry("A", action)
	} else {
		a.AddNameObjectEntry("Dest", annotation.GetDestination())
	}
	p.Annotations.Array = append(p.Annotations.Array, NewObjectReference(a))
}

func (p *PDF) Write(out io.Writer) error {
	// Write Header
	var count int
	n, err := WriteF(out, "%%PDF-%s\n", p.Version)
	if err != nil {
		return err
	}
	count += n
	log.Println("Wrote Header", count)

	// Write Body
	for _, o := range p.Objects {
		o.SetAddress(count)
		n, err = WriteF(out, "%d 0 obj ", o.GetName())
		if err != nil {
			return err
		}
		count += n
		n, err = o.Write(out)
		if err != nil {
			return err
		}
		count += n
		n, err = WriteS(out, " endobj\n")
		if err != nil {
			return err
		}
		count += n
	}
	log.Println("Wrote Body", count)

	// Write Cross Reference
	xrefOffset := count
	n, err = WriteS(out, "xref\n")
	if err != nil {
		return err
	}
	count += n
	n, err = WriteF(out, "0 %d\n0000000000 65535 f\n", len(p.Objects)+1)
	if err != nil {
		return err
	}
	count += n
	for _, o := range p.Objects {
		n, err = WriteF(out, "%010d %05d n\n", o.GetAddress(), o.GetGeneration())
		if err != nil {
			return err
		}
		count += n
	}
	log.Println("Wrote Cross Reference", count)

	// Write Trailer
	n, err = WriteF(out, "trailer <</Size %d /Root %d 0 R>>\n", len(p.Objects)+1, p.Catalog.GetName())
	if err != nil {
		return err
	}
	count += n
	n, err = WriteF(out, "startxref\n%d\n", xrefOffset)
	if err != nil {
		return err
	}
	count += n
	n, err = WriteS(out, "%%EOF\n")
	if err != nil {
		return err
	}
	count += n
	log.Println("Wrote Trailer", count)
	return nil
}

func WriteF(out io.Writer, format string, args ...interface{}) (int, error) {
	return WriteS(out, fmt.Sprintf(format, args...))
}

func WriteS(out io.Writer, data string) (int, error) {
	return io.WriteString(out, data)
}
