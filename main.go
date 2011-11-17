package main

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
  shred_ordering[0] = 0
  for i := 1; i < shred_count; i++ {
    shred_ordering[i], _ = MaximumSimilarityShredIndex(shred_ordering[i-1])
  }
  unshredded_image := image.NewRGBA(image_size.X, image_size.Y)
  for i := 0; i < shred_count; i++ {
    CopyShredToImage(unshredded_image, shredded_image, i, shred_ordering[i], SHRED_WIDTH)
  }
  return unshredded_image
}

type Similarity struct {
  left_i, middle_i, right_i int
  left_middle, middle_right float64
  similarity float64
  delta float64
}

func Reindex() {
  similarities := make([]*Similarity, 20)
  for i := 0; i < 20; i++ {

    similarity := new(Similarity)
    similarity.left_i = i
    similarity.middle_i = (similarity.left_i + 1) % 20
    similarity.right_i = (similarity.left_i + 2) % 20
    similarity.left_middle = ShredSimilarity(GetShred(similarity.left_i), GetShred(similarity.middle_i))
    similarity.middle_right = ShredSimilarity(GetShred(similarity.middle_i), GetShred(similarity.right_i))
    similarity.similarity = (similarity.left_middle + similarity.middle_right) / 2.0

    similarities[i] = similarity
  }
  for i := 0; i < 20; i++ {
    s1 := similarities[i]
    s2 := similarities[(i+1)%20]
    s1.delta = math.Fabs(s1.similarity - s2.similarity)
  }

  for _, s := range similarities {
    fmt.Printf("%3d%3d%3d%6.2f%6.2f%6.2f%6.2f%6.2f\n",
      s.left_i, s.middle_i, s.right_i,
      s.left_middle*100.0, s.middle_right*100.0, s.similarity*100.0, s.delta*100.0)
  }


  // rotate until the goodness-of-fit at the beginning is > that at the end
  // (this is because the end wraps around; we don't want good goodness of fit
  // while we are wrapping around the end!)
  // i.e. (0,1,2)'s similarity is > than that of (19,0,1)
  // i.e. max goodness-of-fit of position 0 while minimizing that of position last
  // rather max(gof[i]-gof[i-1 in ring])
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


