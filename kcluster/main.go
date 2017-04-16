package main

import (
  "fmt"
  //"math"
  "math/rand"
  "os"

  "github.com/montanaflynn/stats"
  "github.com/ray1729/collective_intel/dataset"
)

type Distance func(stats.Float64Data, stats.Float64Data) float64

func main() {
  infile := os.Args[1]
  ds, _ := dataset.ReadDataset(infile)
  kcluster(ds.Data, pearson, 4)
}

func pearson (x, y stats.Float64Data) float64 {
  c, _ := x.Correlation(y)
  return 1.0 - c
}

type Range struct {
  min, max float64
}

func min (xs stats.Float64Data) float64 {
  min, _ := stats.Min(xs)
  return min
}

func max (xs stats.Float64Data) float64 {
  max, _ := stats.Max(xs)
  return max
}

func kcluster(rows []stats.Float64Data, distance Distance, k int) {
  var ranges []Range
  n := len(rows)
  for i := 0; i < n; i++ {
    col := getCol(rows, i)
    min, _ := col.Min()
    max, _ := col.Max()
    ranges = append(ranges, Range{min, max})
  }

  // Start with k random centroids
  var centroids []stats.Float64Data
  for i := 0; i < k; i++ {
    centroids = append(centroids, randomCentroid(ranges))
  }

  var lastmatches []int
  bestmatches := make([]int, n, n)

  for t := 0; t < 100; t++ {
    fmt.Printf("Iteration %d\n", t)
    // For each row, find its closest centroid
    for i, row := range rows {
      var bestmatch int
      for c, centroid := range centroids {
        if distance(row, centroid) < distance(row, centroids[bestmatch]) {
          bestmatch = c
        }
      }
      bestmatches[i] = bestmatch
    }
    if match_eq(bestmatches, lastmatches) {
      break
    }
    lastmatches = bestmatches
    // Move the centroids to the means of their members
    for c := 0; c < k; c++ {
      var m float64
      avgs := make([]float64, len(rows[0]), len(rows[0]))
      for j := 0; j < n; j++ {
        if bestmatches[j] == c {
          m++
          for i, v := range rows[j] {
            avgs[i] += v
          }
        }
      }
      if m > 0 {
        for i := range avgs {
          avgs[i] /= m
        }
        centroids[c] = avgs
      } else {
        centroids[c] = randomCentroid(ranges)
      }
    }
  }
}

func match_eq(xs, ys []int) bool {
  if len(xs) != len(ys) {
    return false
  }
  for i := 0; i < len(xs); i++ {
    if xs[i] != ys[i] {
      return false
    }
  }
  return true
}

func randomCentroid(ranges []Range) stats.Float64Data {
  var centroid []float64
  for _, r := range ranges {
    centroid = append(centroid, randInRange(r))
  }
  return centroid
}

func randInRange(x Range) float64 {
  return rand.Float64()*(x.max - x.min) + x.min
}

func getCol(rows []stats.Float64Data, i int) stats.Float64Data {
  n := len(rows)
  col := make([]float64, n, n)
  for j := 0; j < len(rows); j++ {
    col[j] = rows[j].Get(i)
  }
  return col
}
