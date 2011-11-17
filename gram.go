/* Produce a script that reads in a shredded image (like the one below)
   and produces the original image. For this image, you can assume shreds
   are 32 pixels wide and uniformly spaced across the image horizontally.
   These shreds are scattered at random and if rearranged, will yield the
   original image. */

package main

import (
  "fmt"
  "image"
  "image/draw"
  "log"
  "math"
  "os"
)

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
