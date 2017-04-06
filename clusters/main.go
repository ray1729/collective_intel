package main

import (
  "bufio"
  "fmt"
  "math"
  "os"
  "strconv"
  "strings"
)

type Dataset struct {
  Rownames []string
  Colnames []string
  Data     [][]float64
}

type BiCluster struct {
  Left, Right *BiCluster
  Vec []float64
  Id int
  Distance float64
}

func main() {
  filename := os.Args[1]
  ds, _ := ReadDataset(filename)
  fmt.Println(ds)
}

func hcluster(rows [][]float64, distance func([]float64, []float64) float64) *BiCluster {

}

func pearson(v1, v2 []float64) float64 {
  n := float64(len(v1))
  sum1 := sum(v1)
  sum2 := sum(v2)
  sum1sq := sum_squares(v1)
  sum2sq := sum_squares(v2)
  psum := sum_products(v1, v2)
  num := psum - (sum1*sum2)/n
  den := math.Sqrt((sum1sq-square(sum1)/n)*((sum2sq-square(sum2))/n))
  if den == 0 {
    return 0.0
  }
  return 1.0-(num/den)
}

func sum_products(xs, ys []float64) float64 {
  s := 0.0
  for i := range xs {
    s += xs[i]*ys[i]
  }
  return s
}

func sum_squares(xs []float64) float64 {
  s := 0.0
  for _, v := range xs {
    s += square(v)
  }
  return s
}

func square(x float64) float64 {
  return x*x
}

func sum(xs []float64) float64 {
  s := 0.0
  for _, v := range xs {
    s += v
  }
  return s
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

func ParseRows(scanner *bufio.Scanner) ([]string, [][]float64, error) {
  var rownames []string
  var data [][]float64
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
    data = append(data, row)
  }
  return rownames, data, nil
}

func ParseHeader(scanner *bufio.Scanner) []string {
  scanner.Scan()
  xs := strings.Split(scanner.Text(), "\t")
  return xs[1:]
}
