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

func TestRectangle_Width(t *testing.T) {
	r := &graphics.Rectangle{
		Left:   0,
		Right:  10,
		Top:    10,
		Bottom: 0,
	}
	expected := 10.0
	actual := r.DX()
	if expected != actual {
		t.Errorf("Incorrect width; expected '%f', got '%f'", expected, actual)
	}
}

func TestRectangle_Height(t *testing.T) {
	r := &graphics.Rectangle{
		Left:   0,
		Right:  10,
		Top:    10,
		Bottom: 0,
	}
	expected := 10.0
	actual := r.DY()
	if expected != actual {
		t.Errorf("Incorrect height; expected '%f', got '%f'", expected, actual)
	}
}

func TestRectangle_FloatToString(t *testing.T) {
	expected := "3.14"
	actual := graphics.FloatToString(3.14)
	if expected != actual {
		t.Errorf("Incorrect result; expected '%s', got '%s'", expected, actual)
	}
}
