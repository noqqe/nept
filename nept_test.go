package main

import "testing"

// Low Level Tests

// uint32 to RGB Conversion
func TestUint32ToRGB(t *testing.T) {

	m := map[uint32]uint8{
		0:     0,
		1:     0,
		32768: 127,
		60500: 235,
		65535: 255,
	}

	for c, want := range m {
		if got := uint32ToRGB(c); got != want {
			t.Errorf("%d, want %d", got, want)
		}
	}
}

// percentage calc
func TestPercentToInt(t *testing.T) {

	m := map[uint32]uint32{
		100: 65534,
		50:  32767,
		10:  6553,
		1:   655,
		0:   0,
	}

	for c, want := range m {
		if got := percentToInt(c); got != want {
			t.Errorf("%d, want %d", got, want)
		}
	}
}

// SUB
func TestSubInt(t *testing.T) {
	var want uint32 = 20
	var c uint32 = 60
	var v uint32 = 40
	if got := subInt(c, v); got != want {
		t.Errorf("%d, want %d", got, want)
	}
}

func TestSubIntUnderflow(t *testing.T) {
	var want uint32 = 0
	var c uint32 = 60
	var v uint32 = 80
	if got := subInt(c, v); got != want {
		t.Errorf("%d, want %d", got, want)
	}
}

// ADD
func TestAddIntWithinLimits(t *testing.T) {
	var want uint32 = 100
	var c uint32 = 60
	var v uint32 = 40
	if got := addInt(c, v); got != want {
		t.Errorf("%d, want %d", got, want)
	}
}

func TestAddIntOverflow(t *testing.T) {
	var want uint32 = 65535
	var c uint32 = 60000
	var v uint32 = 30000
	if got := addInt(c, v); got != want {
		t.Errorf("%d, want %d", got, want)
	}
}

// Features Test

func TestBrighten(t *testing.T) {

	m := map[Pixel]Pixel{
		Pixel{r: 10, g: 10, b: 10, a: 255}:          Pixel{r: 32777, g: 32777, b: 32777, a: 255},
		Pixel{r: 40000, g: 40000, b: 40000, a: 255}: Pixel{r: 65535, g: 65535, b: 65535, a: 255},
	}

	for in, want := range m {
		if got := brighten(in, 50); got != want {
			t.Errorf("%v, want %v", got, want)
		}
	}

}

func TestDarken(t *testing.T) {

	m := map[Pixel]Pixel{
		Pixel{r: 32777, g: 32777, b: 32777, a: 255}: Pixel{r: 10, g: 10, b: 10, a: 255},
		Pixel{r: 20000, g: 20000, b: 20000, a: 255}: Pixel{r: 0, g: 0, b: 0, a: 255},
	}

	for in, want := range m {
		if got := darken(in, 50); got != want {
			t.Errorf("%v, want %v", got, want)
		}
	}

}

func TestIso(t *testing.T) {

	m := map[Pixel]Pixel{
		Pixel{r: 32777, g: 32777, b: 32777, a: 255}: Pixel{r: 10, g: 10, b: 10, a: 255},
		Pixel{r: 20000, g: 20000, b: 20000, a: 255}: Pixel{r: 44247, g: 44247, b: 44247, a: 255},
	}

	for in, _ := range m {
		if got := isoify(in, 10); got.r < in.r {
			t.Errorf("%v, want %v>%v", got.r, got.r, in.r)
		}
	}

}

func TestFlatten(t *testing.T) {

	m := map[Pixel]Pixel{

		Pixel{r: 32777, g: 32777, b: 32777, a: 255}: Pixel{r: 32777, g: 32777, b: 32777, a: 255}, // value should not be changed
		Pixel{r: 10000, g: 10000, b: 10000, a: 255}: Pixel{r: 13107, g: 13107, b: 13107, a: 255}, // value should be changed
	}

	for in, want := range m {
		if got := flatten(in, 20); got != want {
			t.Errorf("%v, want %v", got, want)
		}
	}

}
