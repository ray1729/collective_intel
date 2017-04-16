package dataset

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"

  "github.com/montanaflynn/stats"
)

type Dataset struct {
  RowLabels []string
  ColLabels []string
  Data      []stats.Float64Data
}

func NewDataset(rownames, colnames []string, data []stats.Float64Data) (*Dataset, error) {
  var n int
  if colnames != nil {
    n = len(colnames)
  } else {
    n = len(data[0])
  }
  for i, row := range data {
    if len(row) != n {
      return nil, fmt.Errorf("Error constructing dataset: row %d has unexpected length", i)
    }
  }
  return &Dataset{rownames, colnames, data}, nil
}

func (ds *Dataset) NRows() int {
  return len(ds.Data)
}

func (ds *Dataset) NCols() int {
  return len(ds.Data[0])
}

func (ds *Dataset) Get(row, col int) float64 {
  return ds.Data[row][col]
}

func (ds *Dataset) Row(row int) stats.Float64Data {
  return ds.Data[row]
}

func (ds *Dataset) Col(col int) stats.Float64Data {
  res := make([]float64, 0, len(ds.Data))
  for row := 0; row < len(ds.Data); row++ {
    res = append(res, ds.Data[row][col])
  }
  return res
}

func (ds *Dataset) Transpose() *Dataset {
  data := make([]stats.Float64Data, 0, len(ds.ColLabels))
  for col := 0; col < len(ds.ColLabels); col++ {
    data = append(data, ds.Col(col))
  }
  return &Dataset{ds.ColLabels, ds.RowLabels, data}
}

func (ds *Dataset) MapRows(f func (stats.Float64Data) float64) stats.Float64Data {
  nr := ds.NRows()
  res := make([]float64, 0, nr)
  for i := 0; i < nr; i++ {
    res = append(res, f(ds.Row(i)))
  }
  return res
}

func (ds *Dataset) MapCols(f func (stats.Float64Data) float64) stats.Float64Data {
  nc := ds.NCols()
  res := make([]float64, 0, nc)
  for i := 0; i < nc; i++ {
    res = append(res, f(ds.Col(i)))
  }
  return res
}

func ReadDataset(filename string) (*Dataset, error) {
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  colnames, err := ParseHeader(scanner)
  if err != nil {
    return nil, fmt.Errorf("Error reading dataset %s: %v", filename, err)
  }
  rownames, data, err := ParseRows(scanner)
  if err != nil {
    return nil, fmt.Errorf("Error reading dataset %s: %v", filename, err)
  }
  return NewDataset(rownames, colnames, data)
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
  if err := scanner.Err(); err != nil {
    return nil, nil, err
  }
  return rownames, data, nil
}

func ParseHeader(scanner *bufio.Scanner) ([]string, error) {
  if ! scanner.Scan() {
    if err := scanner.Err(); err != nil {
      return nil, fmt.Errorf("Could not parse header: %v", err)
    }
    return nil, fmt.Errorf("Could not parse header: scan failed")
  }
  xs := strings.Split(scanner.Text(), "\t")
  return xs[1:], nil
}
