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

type Object interface {
	SetName(int)
	GetName() int
	SetAddress(int)
	GetAddress() int
	GetGeneration() int
	Write(io.Writer) (int, error)
}

type Metadata struct {
	Name       int
	Address    int
	Generation int
}

func (m *Metadata) SetName(name int) {
	m.Name = name
}

func (m *Metadata) GetName() int {
	return m.Name
}

func (m *Metadata) SetAddress(address int) {
	m.Address = address
}

func (m *Metadata) GetAddress() int {
	return m.Address
}

func (m *Metadata) GetGeneration() int {
	return m.Generation
}
