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

func Shred(shredded_image image.Image, shred_index int) *image.RGBA {
  bounds := shredded_image.Bounds()
  shred := image.NewRGBA(SHRED_WIDTH, bounds.Dy())
  draw.Draw(shred, shred.Bounds(), shredded_image, bounds.Min.Add(image.Pt(SHRED_WIDTH*shred_index, 0)), draw.Src)
  return shred
}


func PrintSimilarityMatrix(shredded_image image.Image, shred_count int) {
  fmt.Printf("%6v", "")
  for i := 0; i < shred_count; i++ {
    fmt.Printf("%6d", i)
  }
  fmt.Printf("\n")
  for i := 0; i < shred_count; i++ {
    fmt.Printf("%6v", i)
    for j := 0; j < shred_count; j++ {
      shred1 := Shred(shredded_image, i)
      shred2 := Shred(shredded_image, j)
      similarity := ShredSimilarity(shred1, shred2)
      fmt.Printf("%6.2f", similarity*100)
    }
    fmt.Printf("\n")
  }
}

func MaximumSimilarityShredIndex(shredded_image image.Image, left_shred_index, shred_count int) (int, float64) {
  left_shred := Shred(shredded_image, left_shred_index)
  maximum_similarity := 0.0
  maximum_similarity_shred_index := -1
  for i := 0; i < shred_count; i++ {
    if i == left_shred_index {
      continue
    }
    right_shred := Shred(shredded_image, i)
    similarity := ShredSimilarity(left_shred, right_shred)
    if (similarity > maximum_similarity) {
      maximum_similarity = similarity
      maximum_similarity_shred_index = i
    }
  }
  return maximum_similarity_shred_index, maximum_similarity
}

// 8 is left 9 is right

func Unshred(shredded_image image.Image, shred_count int) image.Image {
  shred_ordering := make([]int, shred_count)
  minimum_similarity := 1.0
  minimum_similarity_index := -1

  for i := 0; i < shred_count; i++ {
    index, similarity := MaximumSimilarityShredIndex(shredded_image, i, shred_count)
    if similarity < minimum_similarity {
      minimum_similarity = similarity
      minimum_similarity_index = index
    }
    shred_ordering[i] = index
  }
  fmt.Println(shred_ordering)
  unshredded_image := image.NewRGBA(shredded_image.Bounds().Dx(), shredded_image.Bounds().Dy())
  current_index := 0

  emitted_order := make([]int, shred_count)
  for i := 0; i < shred_count; i++ {
    CopyShredToImage(unshredded_image, shredded_image, i, current_index, SHRED_WIDTH)
    emitted_order[i] = current_index
    current_index = shred_ordering[current_index]
  }
  fmt.Println(minimum_similarity_index)
  fmt.Println(emitted_order)

  return unshredded_image
}

type Similarity struct {
  left_i, middle_i, right_i int
  left_middle, middle_right float64
  similarity float64
  delta float64
}

func main() {
  shredded_image := ReadShreddedImageFile(INPUT_FILENAME)
  image_size := shredded_image.Bounds().Size()
  shred_count := image_size.X / SHRED_WIDTH

  PrintSimilarityMatrix(shredded_image, shred_count)

  unshredded_image := Unshred(shredded_image, shred_count)

  PrintSimilarityMatrix(unshredded_image, shred_count)

  similarities := make([]*Similarity, 20)
  for i := 0; i < 20; i++ {

    similarity := new(Similarity)
    similarity.left_i = i
    similarity.middle_i = (similarity.left_i + 1) % 20
    similarity.right_i = (similarity.left_i + 2) % 20
    similarity.left_middle = ShredSimilarity(Shred(unshredded_image, similarity.left_i), Shred(unshredded_image, similarity.middle_i))
    similarity.middle_right = ShredSimilarity(Shred(unshredded_image, similarity.middle_i), Shred(unshredded_image, similarity.right_i))
    similarity.similarity = (similarity.left_middle + similarity.middle_right) / 2.0

    similarities[i] = similarity
  }
  for i := 0; i < 20; i++ {
    s1 := similarities[i]
    s2 := similarities[(i+1)%20]
    s1.delta = math.Fabs(s1.similarity - s2.similarity)
  }

  for _, s := range similarities {
    fmt.Printf("%3d%3d%3d%6.2f%6.2f%6.2f%6.2f\n",
      s.left_i, s.middle_i, s.right_i,
      s.left_middle*100.0, s.middle_right*100.0, s.similarity*100.0, s.delta*100.0)
  }


  // rotate until the goodness-of-fit at the beginning is > that at the end
  // (this is because the end wraps around; we don't want good goodness of fit
  // while we are wrapping around the end!)
  // i.e. (0,1,2)'s similarity is > than that of (19,0,1)
  // i.e. max goodness-of-fit of position 0 while minimizing that of position last

/* WriteImageFile(OUTPUT_FILENAME, unshredded_image)*/
  /* unshredded_image := image.NewRGBA(image_size.X, image_size.Y)*/

  /* for i := 0; i < 20; i++ {*/
  /*   shred1 := Shred(shredded_image, 0)*/
  /*   shred2 := Shred(shredded_image, i)*/
  /*   fmt.Printf("%3d: %f\n", i, ShredSimilarity(shred1, shred2))*/
  /*   WriteImageFile(fmt.Sprintf("%d.png", i), shred2)*/
  /* }*/

  /* WriteImageFile(OUTPUT_FILENAME, unshredded_image)*/
}


