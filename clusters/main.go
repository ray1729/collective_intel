package main

import (
  "bufio"
  "fmt"
  //"math"
  "os"
  "strconv"
  "strings"

  "github.com/BenLubar/memoize"
  "github.com/montanaflynn/stats"
)

type Dataset struct {
  Rownames []string
  Colnames []string
  Data     []stats.Float64Data
}

type BiCluster struct {
  Left, Right *BiCluster
  Vec stats.Float64Data
  Id int
  Distance float64
}

func main() {
  filename := os.Args[1]
  ds, _ := ReadDataset(filename)
  clust := hcluster(ds.Data, pearson)
  fmt.Println(len(ds.Rownames))
  printclust(clust, ds.Rownames, 0)
}

func printclust(clust *BiCluster, labels []string, n int) {
  for i := 0; i < n; i++ {
    fmt.Print(" ")
  }
  if clust.Id < 0 {
    fmt.Println("-")
  } else {
    fmt.Println(labels[clust.Id])
  }
  if clust.Left != nil {
    printclust(clust.Left, labels, n+1)
  }
  if clust.Right != nil {
    printclust(clust.Right, labels, n+1)
  }
}

func pearson (x, y stats.Float64Data) float64 {
  c, _ := x.Correlation(y)
  return 1.0 - c
}

func hcluster(rows []stats.Float64Data, distance func(stats.Float64Data, stats.Float64Data) float64) *BiCluster {
  clust := make(map[int]*BiCluster)
  for clusterid, vec := range rows {
    clust[clusterid] = &BiCluster{Id: clusterid, Vec: vec}
  }

  n_elems := rows[0].Len()

  dist := memoize.Memoize(func (i,j int) float64 {
    return distance(clust[i].Vec, clust[j].Vec)
  }).(func(i,j int) float64)

  d := func(i,j int) float64 {
    if i < j {
      return dist(i, j)
    }
    return dist(j, i)
  }

  clusterid := 0
  for len(clust) > 1 {
    // Find the closest clusters
    clusterids := make([]int, 0, len(clust))
    for id := range(clust) {
      clusterids = append(clusterids, id)
    }
    lx := clusterids[0]
    ly := clusterids[1]
    closest := d(lx, ly)
    for i := 0; i < len(clusterids); i++ {
      for j := i+1; j < len(clusterids); j++ {
        this_d := d(clusterids[i], clusterids[j])
        if this_d < closest {
          closest = this_d
          lx = clusterids[i]
          ly = clusterids[j]
        }
      }
    }
    // Calculate the average of the closest clusters
    vec := make([]float64, n_elems, n_elems)
    for i := 0; i < n_elems; i++ {
      vec[i] = (clust[lx].Vec.Get(i) + clust[ly].Vec.Get(i))/2.0
    }
    // Create new cluster as this average
    clusterid--
    newcluster := BiCluster{Id: clusterid, Vec: vec, Left: clust[lx], Right: clust[ly], Distance: closest}
    clust[clusterid] = &newcluster
    // Delete the clusters we merged
    delete(clust, lx)
    delete(clust, ly)
  }

  // There is only one cluster in the clust map
  var result *BiCluster
  for _, v := range clust {
    result = v
  }
  return result
}

func ReadDataset(filename string) (*Dataset, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  colnames := ParseHeader(scanner)
  rownames, data, err := ParseRows(scanner)
  if err != nil {
    return nil, err
  }
  res := &Dataset{Colnames: colnames, Rownames: rownames, Data: data}
  return res, nil
}

func ParseRows(scanner *bufio.Scanner) ([]string, []stats.Float64Data, error) {
  var rownames []string
  var data []stats.Float64Data
  for scanner.Scan() {
    xs := strings.Split(scanner.Text(), "\t")
    name := xs[0]
    var row []float64
    for i := 1; i < len(xs); i++ {
      v, err := strconv.ParseFloat(xs[i], 64)
      if err != nil {
        return nil, nil, fmt.Errorf("Error parsing %s row %d: %v", name, i, err)
      }
      row = append(row, v)
    }
    rownames = append(rownames, name)
    data = append(data, stats.Float64Data(row))
  }
  return rownames, data, nil
}

func ParseHeader(scanner *bufio.Scanner) []string {
  scanner.Scan()
  xs := strings.Split(scanner.Text(), "\t")
  return xs[1:]
}
