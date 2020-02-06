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

type ArrayObject struct {
	Metadata
	Array []Object
}

func (o *ArrayObject) Write(out io.Writer) (int, error) {
	var count int
	n, err := WriteS(out, "[")
	if err != nil {
		return 0, err
	}
	count += n
	for i, a := range o.Array {
		if i != 0 {
			n, err = WriteS(out, " ")
			if err != nil {
				return 0, err
			}
			count += n
		}
		n, err = a.Write(out)
		if err != nil {
			return 0, err
		}
		count += n
	}
	n, err = WriteS(out, "]")
	if err != nil {
		return 0, err
	}
	count += n
	return count, nil
}
