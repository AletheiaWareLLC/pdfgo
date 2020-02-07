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

package graphics_test

import (
	"github.com/AletheiaWareLLC/pdfgo/graphics"
	"testing"
)

func TestTextBox_SplitLines(t *testing.T) {
	given := []rune("This\r\nis\na\r\ntest\n")
	expected := [][]rune{
		[]rune("This"),
		[]rune("is"),
		[]rune("a"),
		[]rune("test"),
	}
	actual := graphics.SplitLines(given)
	if len(actual) != len(expected) {
		t.Errorf("Incorrect line count; expected '%d', got '%d'", len(expected), len(actual))
	} else {
		for i := 0; i < len(expected); i++ {
			a := actual[i]
			e := expected[i]
			if len(a) != len(e) {
				t.Errorf("Incorrect line length; expected '%d', got '%d'", len(e), len(a))
			} else {
				for j := range e {
					if a[j] != e[j] {
						t.Errorf("Incorrect splitting; expected '%c', got '%c'", e[j], a[j])
					}
				}
			}
		}
	}
}

func TestTextBox_PDFEscapeString(t *testing.T) {
	given := []rune("ATTACK) Tj\n\nET")
	expected := "ATTACK\\) Tj\n\nET"
	actual := graphics.PDFEscapeString(given)
	if actual != expected {
		t.Errorf("Incorrect escaping; expected '%s', got '%s'", expected, actual)
	}
}
