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

package graphics

import (
	"math"
	"strconv"
)

type Rectangle struct {
	Left, Right, Top, Bottom float64
}

func (r *Rectangle) DX() float64 {
	return r.Right - r.Left
}

func (r *Rectangle) DY() float64 {
	return r.Top - r.Bottom
}

func (r *Rectangle) Max(o *Rectangle) *Rectangle {
	return &Rectangle{
		Left:   math.Min(r.Left, o.Left),
		Right:  math.Max(r.Right, o.Right),
		Top:    math.Max(r.Top, o.Top),
		Bottom: math.Min(r.Bottom, o.Bottom),
	}
}

func NegativeRectangle() *Rectangle {
	return &Rectangle{
		Left:   math.MaxFloat64,
		Right:  -math.MaxFloat64,
		Top:    -math.MaxFloat64,
		Bottom: math.MaxFloat64,
	}
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
