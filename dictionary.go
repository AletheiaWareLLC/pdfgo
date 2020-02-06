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

type DictionaryObject struct {
	Metadata
	Keys       []*NameObject
	Dictionary map[*NameObject]Object
}

func (o *DictionaryObject) AddNameNameEntry(key, value string) {
	o.AddNameObjectEntry(key, &NameObject{Name: value})
}

func (o *DictionaryObject) AddNameObjectEntry(key string, value Object) {
	o.AddObjectObjectEntry(&NameObject{Name: key}, value)
}

func (o *DictionaryObject) AddObjectObjectEntry(key *NameObject, value Object) {
	o.Keys = append(o.Keys, key)
	o.Dictionary[key] = value
}

func (o *DictionaryObject) Write(out io.Writer) (int, error) {
	var count int
	n, err := WriteS(out, "<<")
	if err != nil {
		return 0, err
	}
	count += n
	for i, k := range o.Keys {
		if i != 0 {
			n, err = WriteS(out, " ")
			if err != nil {
				return 0, err
			}
			count += n
		}
		n, err = k.Write(out)
		if err != nil {
			return 0, err
		}
		count += n
		n, err = WriteS(out, " ")
		if err != nil {
			return 0, err
		}
		count += n
		n, err = o.Dictionary[k].Write(out)
		if err != nil {
			return 0, err
		}
		count += n
	}
	n, err = WriteS(out, ">>")
	if err != nil {
		return 0, err
	}
	count += n
	return count, nil
}
