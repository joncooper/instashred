package main

import (
	"image"
	"testing"
)

type pixelSimilarityTestCase struct {
	pixel1, pixel2 image.Color
	expected       float64
}

var pixelSimilarityTests = []pixelSimilarityTestCase{
	pixelSimilarityTestCase{
		image.RGBAColor{255, 255, 255, 255},
		image.RGBAColor{255, 255, 255, 255},
		1.0,
	},
	pixelSimilarityTestCase{
		image.RGBAColor{0, 0, 0, 0},
		image.RGBAColor{0, 0, 0, 0},
		1.0,
	},
	pixelSimilarityTestCase{
		image.RGBAColor{255, 255, 255, 255},
		image.RGBAColor{255, 0, 0, 255},
		1.0 / 3.0,
	},
	pixelSimilarityTestCase{
		image.RGBAColor{0, 0, 0, 128},
		image.RGBAColor{255, 255, 255, 128},
		0.0,
	},
}

func TestShredSimilarity(t *testing.T) {
	for _, testCase := range pixelSimilarityTests {
		shred1 := image.NewRGBA(256, 256)
		shred2 := image.NewRGBA(256, 256)
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				shred1.Set(i, j, testCase.pixel1)
				shred2.Set(i, j, testCase.pixel2)
			}
		}
	}
}

func TestPixelSimilarity(t *testing.T) {
	for _, testCase := range pixelSimilarityTests {
		actual := PixelSimilarity(testCase.pixel1, testCase.pixel2)
		if actual != testCase.expected {
			t.Errorf("PixelSimilarity: %+v\nActual: %v\n", testCase, actual)
		}
	}
}

func TestCopyShredToImage(t *testing.T) {
	src_image := image.NewRGBA(256, 256)
	dest_image := image.NewRGBA(256, 256)

	for i := 0; i < 8; i++ {
		CopyShredToImage(dest_image, src_image, i, i, 32)
	}

	for i := range src_image.Pix {
		if src_image.Pix[i] != dest_image.Pix[i] {
			t.Errorf("CopyShredToImage failed at pixel %v (%v vs %v)\n", i, dest_image.Pix[i], src_image.Pix[i])
		}
	}
}
