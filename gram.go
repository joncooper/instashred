/* Produce a script that reads in a shredded image (like the one below)
   and produces the original image. For this image, you can assume shreds
   are 32 pixels wide and uniformly spaced across the image horizontally.
   These shreds are scattered at random and if rearranged, will yield the
   original image. */

package main

import (
  "image"
  "log"
  "math"
  "os"
  "image/draw"
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

func PixelChannelSimilarity(channel1, channel2 uint32) float64 {
  c1 := float64(channel1 / 0xFFFF)
  c2 := float64(channel2 / 0xFFFF)
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

type pixelSimilarityTestCase struct {
  pixel1, pixel2 image.Color
  expected float64
}

var pixelSimilarityTests = []pixelSimilarityTestCase {
  pixelSimilarityTestCase {
    image.RGBAColor { 255, 255, 255, 255 },
    image.RGBAColor { 255, 255, 255, 255 },
    1.0,
  },
  pixelSimilarityTestCase {
    image.RGBAColor { 0, 0, 0, 0 },
    image.RGBAColor { 0, 0, 0, 0 },
    1.0,
  },
  pixelSimilarityTestCase {
    image.RGBAColor { 255, 255, 255, 255 },
    image.RGBAColor { 255, 0, 0, 255 },
    1.0/3.0,
  },
  pixelSimilarityTestCase{
    image.RGBAColor { 0, 0, 0, 128 },
    image.RGBAColor { 255, 255, 255, 128},
    0.0,
  },
}

func TestPixelSimilarity() {
  for _, testCase := range pixelSimilarityTests {
    actual := PixelSimilarity(testCase.pixel1, testCase.pixel2)
    if (actual != testCase.expected) {
      log.Fatalf("PixelSimilarity: %+v\nActual: %v\n", testCase, actual)
    }
  }
}

func CopyShredToImage(dest_image, src_image *image.Image, strip_index int) {
  
}


func main() {
  shredded_image := ReadShreddedImageFile(INPUT_FILENAME)
  image_size := shredded_image.Bounds().Size()
  unshredded_image := image.NewRGBA(image_size.X, image_size.Y)

  dest_rect := unshredded_image.Bounds()
  src_rect := shredded_image.Bounds()
  draw.Draw(unshredded_image, dest_rect, shredded_image, src_rect.Min, draw.Src)

  WriteImageFile(OUTPUT_FILENAME, unshredded_image)
}

/*
import (
  "crypto/md5"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "path"
  "path/filepath"
)

var verbose *bool = flag.Bool("verbose", false, "Print the list of duplicate files.")
var rootDir string = "."
var fullPathsByFilename map[string][]string

type DupeChecker struct{}

func (dc DupeChecker) VisitDir(fullpath string, f *os.FileInfo) bool {
  return true
}

func (dc DupeChecker) VisitFile(fullpath string, f *os.FileInfo) {
  filename := path.Base(fullpath)
  fullPathsByFilename[filename] = append(fullPathsByFilename[filename], fullpath)
}

func MD5OfFile(fullpath string) []byte {
  if contents, err := ioutil.ReadFile(fullpath); err == nil { 
    md5sum := md5.New()
    md5sum.Write(contents)
    return md5sum.Sum()
  }
  return nil
}

func PrintResults() {
  dupes := 0
  for key, value := range fullPathsByFilename {
    if (len(value) < 2) {
      continue
    }
    dupes++
    if (*verbose) {
      println(key, ":")
      for _, filename := range value {
        println("  ", filename)
        fmt.Printf("    %x\n", MD5OfFile(filename))
      }
    }
  }
  println("Total duped files found:", dupes)
}

func FindDupes(root string) {
  fullPathsByFilename = make(map[string][]string)
  filepath.Walk(root, DupeChecker{}, nil)
}

func ParseArgs() {
  flag.Parse()
  if (len(flag.Args()) > 0) { 
    rootDir = flag.Arg(0)
  } 
}

func main() {
  ParseArgs()
  FindDupes(rootDir)
  PrintResults()
}

*/
