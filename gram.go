/* Produce a script that reads in a shredded image (like the one below)
   and produces the original image. For this image, you can assume shreds
   are 32 pixels wide and uniformly spaced across the image horizontally.
   These shreds are scattered at random and if rearranged, will yield the
   original image. */

package main

import (
  "image"
  "log"
  "os"
  "image/png"
)

const INPUT_FILENAME string = "TokyoPanoramaShredded.png"
const OUTPUT_FILENAME string = "TokyoPanorama.png"

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
  WriteImageFile(OUTPUT_FILENAME, shredded_image)
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
