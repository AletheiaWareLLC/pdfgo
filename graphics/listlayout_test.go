package graphics_test

import (
	"bytes"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/graphics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListLayout_nesting(t *testing.T) {
	// TODO
}

func TestListLayout_overflow(t *testing.T) {
	b1 := &fixedSizeBox{}
	b1.width = 100
	b1.height = 100

	b2 := &fixedSizeBox{}
	b2.width = 100
	b2.height = 100

	l := &graphics.ListLayout{
		Direction: graphics.TopBottom,
	}
	l.Add(b1)
	l.Add(b2)
	u, err := l.SetBounds(&graphics.Rectangle{
		Left:   10,
		Right:  150,
		Top:    150,
		Bottom: 10,
	})
	assert.Nil(t, err)
	assert.Equal(t, 10., u.Left)
	assert.Equal(t, 110., u.Right)
	assert.Equal(t, 150., u.Top)
	assert.Equal(t, 50., u.Bottom)
	assert.Equal(t, 1, l.Visible)
}

type fixedSizeBox struct {
	width, height float64
}

func (b *fixedSizeBox) SetBounds(bounds *graphics.Rectangle) (*graphics.Rectangle, error) {
	return &graphics.Rectangle{
		Left:   bounds.Left,
		Right:  bounds.Left + b.width,
		Top:    bounds.Top,
		Bottom: bounds.Top - b.height,
	}, nil
}

func (b *fixedSizeBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	// Nothing
	return nil
}
