package main

import (
  "log"
  "os"
  "image"
  "image/png"
)

const INPUT_FILENAME string = "TokyoPanoramaShredded.png"
const OUTPUT_FILENAME string = "TokyoPanorama.png"
const SHRED_WIDTH int = 32

func ReadShreddedImageFile(filename string) image.Image {
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

func WriteImageFile(filename string, image image.Image) {
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

func main() {
  shredded_image := ReadShreddedImageFile(INPUT_FILENAME)
  image_size := shredded_image.Bounds().Size()
  unshredded_image := image.NewRGBA(image_size.X, image_size.Y)

  for i := 0; i < 32; i++ {
    CopyShredToImage(unshredded_image, shredded_image, i, i, 32)
  }

  WriteImageFile(OUTPUT_FILENAME, unshredded_image)
}


