package main

/* Produce a script that reads in a shredded image (like the one below)
   and produces the original image. For this image, you can assume shreds
   are 32 pixels wide and uniformly spaced across the image horizontally.
   These shreds are scattered at random and if rearranged, will yield the
   original image. */

import (
  "fmt"
  "log"
  "os"
  "image"
  "image/draw"
  "image/png"
  "math"
)

const INPUT_FILENAME string = "TokyoPanoramaShredded.png"
const OUTPUT_FILENAME string = "TokyoPanorama.png"
const SHRED_WIDTH int = 32

var shredded_image image.Image
var image_size image.Point
var shred_count int

const DEBUG bool = false

func DbgPrintln(to_print ...interface{}) (n int, err os.Error) {
  if (DEBUG) {
    return fmt.Println(to_print)
  }
  return
}

func PixelChannelSimilarity(channel1, channel2 uint32) float64 {
  c1 := float64(channel1) / float64(0xFFFF)
  c2 := float64(channel2) / float64(0xFFFF)
  return (1.0 - math.Fabs(c1 - c2)) / 1.0
}

func PixelSimilarity(pixel1, pixel2 image.Color) float64 {
  p1r, p1g, p1b, _ := pixel1.RGBA()
  p2r, p2g, p2b, _ := pixel2.RGBA()
  similarity := 0.0
  similarity += (PixelChannelSimilarity(p1r, p2r) / 3.0)
  similarity += (PixelChannelSimilarity(p1g, p2g) / 3.0)
  similarity += (PixelChannelSimilarity(p1b, p2b) / 3.0)
  return similarity
}

// Compare the rightmost column of the left shred to the leftmost column of the right shred
//
func ShredSimilarity(left_shred, right_shred image.Image) float64 {
  left_shred_rightmost_column_index := left_shred.Bounds().Max.X - 1
  left_shred_height := left_shred.Bounds().Max.Y

  right_shred_leftmost_column_index := right_shred.Bounds().Min.X
  right_shred_height := right_shred.Bounds().Max.Y

  if (left_shred_height != right_shred_height) {
    log.Fatal("Shreds have different Y heights. (%v vs %v).\n", left_shred_height, right_shred_height)
  }

  similarity := 0.0
  for i := 0; i < left_shred_height; i++ {
    left_pixel := left_shred.At(left_shred_rightmost_column_index, i)
    right_pixel := right_shred.At(right_shred_leftmost_column_index, i)
    pixel_similarity := PixelSimilarity(left_pixel, right_pixel)

    DbgPrintln(i, left_pixel, right_pixel, pixel_similarity)

    similarity += pixel_similarity
 }
  similarity /= float64(left_shred_height)
  return similarity
}

func CopyShredToImage(dest_image draw.Image, src_image image.Image, dest_shred_index, src_shred_index, shred_width int) {
  // TODO: handle the case where we get a bad index
  src_point := image.ZP // Zero point, i.e. (0,0)
  src_point.X = shred_width * src_shred_index

  dest_rect := image.Rectangle{
    image.Point{
      shred_width * dest_shred_index,
      0,
    },
    image.Point{
      (shred_width * (dest_shred_index+1)),
      dest_image.Bounds().Max.Y,
    },
  }
  // The second coordinate of the Rectangle (Rectangle.Max) looks counterintuitive here.
  // image.Bounds() defines a Rectangle including Bounds.Min (i.e. (x0,y0)) but excluding Bounds.Max (i.e. (x1,y1)).
  // That's why.
  draw.Draw(dest_image, dest_rect, src_image, src_point, draw.Src)
} 

func ReadPNGFile(filename string) image.Image {
  f, err := os.Open(INPUT_FILENAME)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  image, err := png.Decode(f)
  if err != nil {
    log.Fatal(err)
  }
  return image
}

func WritePNGFile(filename string, image image.Image) {
  f, err := os.Create(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  err = png.Encode(f, image)
  if err != nil {
    log.Fatal(err)
  }
}




func GetShred(shred_index int) *image.RGBA {
  bounds := shredded_image.Bounds()
  shred := image.NewRGBA(SHRED_WIDTH, bounds.Dy())
  draw_rect := shred.Bounds()
  src_point := bounds.Min.Add(image.Pt(SHRED_WIDTH*shred_index, 0))

  draw.Draw(shred, draw_rect, shredded_image, src_point, draw.Src)
  return shred
}


func PrintSimilarityMatrix() {
  fmt.Printf("%6v", "")
  for i := 0; i < shred_count; i++ {
    fmt.Printf("%6d", i)
  }
  fmt.Printf("\n")
  for i := 0; i < shred_count; i++ {
    fmt.Printf("%6v", i)
    for j := 0; j < shred_count; j++ {
      shred1 := GetShred(i)
      shred2 := GetShred(j)
      similarity := ShredSimilarity(shred1, shred2)
      fmt.Printf("%6.2f", similarity*100)
    }
    fmt.Printf("\n")
  }
}

func MaximumSimilarityShredIndex(left_shred_index int) (int, float64) {
  left_shred := GetShred(left_shred_index)
  maximum_similarity := 0.0
  maximum_similarity_shred_index := -1
  for i := 0; i < shred_count; i++ {
    if i == left_shred_index {
      continue
    }
    right_shred := GetShred(i)
    similarity := ShredSimilarity(left_shred, right_shred)
    if (similarity > maximum_similarity) {
      maximum_similarity = similarity
      maximum_similarity_shred_index = i
    }
  }
  return maximum_similarity_shred_index, maximum_similarity
}

func Unshred() image.Image {
  shred_ordering := make([]int, shred_count)
  shred_similarity := make([]float64, shred_count)
  shred_ordering[0] = 0
  for i := 1; i < shred_count; i++ {
    shred_ordering[i], shred_similarity[i-1] = MaximumSimilarityShredIndex(shred_ordering[i-1])
  }
  last_shred := GetShred(shred_ordering[shred_count-1])
  first_shred := GetShred(shred_ordering[0])
  shred_similarity[shred_count-1] = ShredSimilarity(last_shred, first_shred)

  // rotate until the goodness-of-fit at the beginning is > that at the end
  // (this is because the end wraps around; we don't want good goodness of fit
  // while we are wrapping around the end!)
  // i.e. max goodness-of-fit of position 0 while minimizing that of position last
  // rather max(gof[i]-gof[i-1 in ring])
  // put another way, we want the best fit possible between shreds (0,1)
  // while minimizing the fit between shreds(19,0)

  max_delta := 0.0
  rightmost_shred := shred_ordering[shred_count-1]
  for i, similarity := range shred_similarity {
    left_goodness_of_fit := shred_similarity[((i-1)+shred_count)%shred_count]
    right_goodness_of_fit := similarity
    goodness_of_fit_delta := left_goodness_of_fit - right_goodness_of_fit
    if goodness_of_fit_delta > max_delta {
      rightmost_shred = i
      max_delta = goodness_of_fit_delta
    }
    fmt.Println(i, similarity, left_goodness_of_fit, right_goodness_of_fit, goodness_of_fit_delta)
    fmt.Println(rightmost_shred, max_delta)
  }

  unshredded_image := image.NewRGBA(image_size.X, image_size.Y)
  shred_index := (rightmost_shred+1)%20
  for i := 0; i < shred_count; i++ {
    CopyShredToImage(unshredded_image, shredded_image, i, shred_ordering[shred_index], SHRED_WIDTH)
    shred_index = (shred_index+1)%20
  }
  fmt.Println(shred_ordering)
  return unshredded_image
}

func main() {
  shredded_image = ReadPNGFile(INPUT_FILENAME)
  image_size = shredded_image.Bounds().Size()
  shred_count = image_size.X / SHRED_WIDTH

  if (DEBUG) {
    PrintSimilarityMatrix()
  }

  unshredded_image := Unshred()
  WritePNGFile(OUTPUT_FILENAME, unshredded_image)
}


