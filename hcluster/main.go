package main

import (
  "fmt"
  "image/color"
  "math"
  "os"

  "github.com/BenLubar/memoize"
  "github.com/fogleman/gg"
  "github.com/montanaflynn/stats"
  "github.com/ray1729/collective_intel/dataset"
)

type BiCluster struct {
  Left, Right *BiCluster
  Vec stats.Float64Data
  Id int
  Distance float64
}

func main() {
  infile := os.Args[1]
  outfile := os.Args[2]
  ds, _ := dataset.ReadDataset(infile)
  clust := hcluster(ds.Data, pearson)
  drawdendrogram(clust, ds.RowLabels, outfile)
}

func drawdendrogram(clust *BiCluster, labels []string, filename string) {
  h := getheight(clust)*20.0
  w := 1500.0
  depth := getdepth(clust)
  // Width is fixed, so scale to fit
  sf := (w-250.0)/depth
  dc := gg.NewContext(int(w), int(h))
  dc.SetColor(color.White)
  dc.DrawRectangle(0.0, 0.0, w, h)
  dc.Fill()
  dc.SetColor(color.Black)
  dc.SetLineWidth(1.0)
  dc.DrawLine(0, h/2, 10.0, h/2)
  dc.Stroke()
  drawnode(dc, clust, 10.0, h/2, sf, labels)
  dc.SavePNG(filename)
}

func drawnode(dc *gg.Context, clust *BiCluster, x float64, y float64, sf float64, labels []string) {
  if clust.Id < 0 {
    h1 := getheight(clust.Left)*20.0
    h2 := getheight(clust.Right)*20.0
    top := y - (h1+h2)/2.0
    bottom := y + (h1+h2)/2.0
    ll := clust.Distance*sf
    // Vertical line from this cluster to children
    dc.DrawLine(x, top+h1/2.0, x, bottom-h2/2.0)
    // Horizontal line to left item
    dc.DrawLine(x, top+h1/2.0, x+ll, top+h1/2.0)
    // Horizontal line to right item
    dc.DrawLine(x, bottom-h2/2.0, x+ll, bottom-h2/2.0)
    dc.Stroke()
    // Left and right nodes
    drawnode(dc, clust.Left, x+ll, top+h1/2.0, sf, labels)
    drawnode(dc, clust.Right, x+ll, bottom-h2/2.0, sf, labels)
  } else {
    // Leaf node, render the label
    dc.DrawString(labels[clust.Id], x+5.0, y+2.0)
  }
}

func getheight(clust *BiCluster) float64 {
  // If this is a leaf node, the height is 1
  if clust.Left == nil && clust.Right == nil {
    return 1.0
  }
  // Otherwise the height is the sum of the heights
  // of the two branches
  return getheight(clust.Left) + getheight(clust.Right)
}

func getdepth(clust *BiCluster) float64 {
  // The depth of a leaf node is 0
  if clust.Left == nil && clust.Right == nil {
    return 0.0
  }
  // The depth of a branch is the greater of its two
  // sides plus its own distance
  return math.Max(getdepth(clust.Left), getdepth(clust.Right)) + clust.Distance
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
