package gram

import (
  "image"
  "testing"
)

type pixelSimilarityTestCase struct {
  pixel1, pixel2 image.Color
  expected float32
}

var pixelSimilarityTests = []pixelSimilarityTestCase {
  pixelSimilarityTestCase {
    image.NRGBAColor { 255, 255, 255, 255 },
    image.NRGBAColor { 255, 255, 255, 255 },
    1.0,
  },
  pixelSimilarityTestCase{
    image.NRGBAColor { 0, 0, 0, 255 },
    image.NRGBAColor { 255, 255, 255, 0},
    0.0,
  },
}

func TestPixelSimilarity(t *testing.T) {
  for _, testCase := range pixelSimilarityTests {
    actual := PixelSimilarity(testCase.pixel1, testCase.pixel2)
    if (actual != testCase.expected) {
      t.Errorf("PixelSimilarity: %+v\n", testCase)
    }
  }
}
