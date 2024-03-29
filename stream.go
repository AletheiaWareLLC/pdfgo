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

import "io"

type StreamObject struct {
	DictionaryObject
	// TODO support encoding
	Data []byte
}

func (o *StreamObject) Write(out io.Writer) (int, error) {
	addKey := true
	for k := range o.Dictionary {
		if k.Name == "Length" {
			addKey = false
			break
		}
	}
	key := &NameObject{
		Name: "Length",
	}
	if addKey {
		o.Keys = append(o.Keys, key)
	}
	o.Dictionary[key] = &NumberObject{
		Number: float64(len(o.Data)),
	}
	var count int
	n, err := o.DictionaryObject.Write(out)
	if err != nil {
		return 0, err
	}
	count += n
	n, err = WriteF(out, "\nstream\n")
	if err != nil {
		return 0, err
	}
	count += n
	n, err = out.Write(o.Data)
	if err != nil {
		return 0, err
	}
	count += n
	n, err = WriteS(out, "\nendstream")
	if err != nil {
		return 0, err
	}
	count += n
	return count, nil
}
