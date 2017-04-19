package main

import (
  "fmt"
  "log"
  "math/rand"
  "os"
  "strconv"

  "github.com/montanaflynn/stats"
  "github.com/ray1729/collective_intel/dataset"
)

type Distance func(stats.Float64Data, stats.Float64Data) float64

func main() {
  infile := os.Args[1]
  k, err := strconv.Atoi(os.Args[2])
  if err != nil {
    log.Fatal("Failed to parse number of clusters")
  }
  ds, _ := dataset.ReadDataset(infile)
  clusters, err := kcluster(ds, pearson, k)
  if err != nil {
    log.Fatal(err)
  }
  for i := 0; i < k; i++ {
    fmt.Println("======================")
    for j := 0; j < ds.NRows; j++ {
      if clusters[j] == i {
        fmt.Println(ds.RowLabels[j])
      }
    }
  }
}

func pearson (x, y stats.Float64Data) float64 {
  c, _ := x.Correlation(y)
  return 1.0 - c
}

func kcluster(ds *dataset.Dataset, distance Distance, k int) ([]int, error) {
  // Start with k random centroids
  maxima := make([]float64, ds.NCols, ds.NCols)
  minima := make([]float64, ds.NCols, ds.NCols)
  for i := 0; i < ds.NCols; i++ {
    col := ds.Col(i)
    max, err := stats.Max(col)
    if err != nil {
      return nil, err
    }
    min, err := stats.Min(col)
    if err != nil {
      return nil, err
    }
    maxima[i] = max
    minima[i] = min
  }
  var centroids []stats.Float64Data
  for i := 0; i < k; i++ {
    centroids = append(centroids, randomCentroid(minima, maxima))
  }

  var lastmatches []int
  bestmatches := make([]int, ds.NRows, ds.NRows)

  closestCentroid := func(row stats.Float64Data) int {
    var bestmatch int
    for c, centroid := range centroids {
      if distance(row, centroid) < distance(row, centroids[bestmatch]) {
        bestmatch = c
      }
    }
    return bestmatch
  }

  for t := 0; t < 100; t++ {
    fmt.Printf("Iteration %d\n", t)
    // For each row, find its closest centroid
    for i := 0; i < ds.NRows; i++ {
      bestmatches[i] = closestCentroid(ds.Row(i))
    }
    if match_eq(bestmatches, lastmatches) {
      break
    }
    lastmatches = bestmatches
    // Move the centroids to the means of their members
    for c := 0; c < k; c++ {
      var m float64
      avgs := make([]float64, ds.NCols, ds.NCols)
      for j := 0; j < ds.NRows; j++ {
        if bestmatches[j] == c {
          m++
          for i, v := range ds.Row(j) {
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
        centroids[c] = randomCentroid(minima, maxima)
      }
    }
  }

  return bestmatches, nil
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

func randomCentroid(minima, maxima stats.Float64Data) stats.Float64Data {
  var centroid []float64
  for i := 0; i < len(minima); i++ {
    min := minima[i]
    max := maxima[i]
    v := rand.Float64()*(max - min) + min
    centroid = append(centroid, v)
  }
  return centroid
}
