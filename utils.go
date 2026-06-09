package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image/color"
	"math"
)

func generateRandomSeed() string {
	b := make([]byte, 8)
	crand.Read(b)
	return fmt.Sprintf("%x", b)
}

func seedFromString(s string) int64 {
	h := sha256.Sum256([]byte(s))
	return int64(binary.BigEndian.Uint64(h[:8]))
}

func HSLToRGB(h, s, l float64) color.RGBA {
	// Achromatic case
	if s == 0 {
		v := uint8(math.Round(l * 255))
		return color.RGBA{v, v, v, 255}
	}

	c := (1 - math.Abs(2*l-1)) * s
	hPrime := h / 60
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))

	var r1, g1, b1 float64
	switch {
	case hPrime < 1:
		r1, g1, b1 = c, x, 0
	case hPrime < 2:
		r1, g1, b1 = x, c, 0
	case hPrime < 3:
		r1, g1, b1 = 0, c, x
	case hPrime < 4:
		r1, g1, b1 = 0, x, c
	case hPrime < 5:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}

	m := l - c/2
	r := uint8(math.Round((r1 + m) * 255))
	g := uint8(math.Round((g1 + m) * 255))
	b := uint8(math.Round((b1 + m) * 255))
	return color.RGBA{r, g, b, 255}
}
