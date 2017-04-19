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
  RowLabels, ColLabels []string
  NRows, NCols int
  Data []stats.Float64Data
}

func NewDataset(rownames, colnames []string, data []stats.Float64Data) (*Dataset, error) {
  var nrows, ncols int
  nrows = len(data)
  if rownames != nil && len(rownames) != nrows {
    return nil, fmt.Errorf("Error constructing dataset: %d row labels for %d rows", len(rownames), nrows)
  }
  if colnames != nil {
    ncols = len(colnames)
  } else {
    ncols = len(data[0])
  }
  for i, row := range data {
    if len(row) != ncols {
      return nil, fmt.Errorf("Error constructing dataset: row %d has unexpected length", i)
    }
  }

  return &Dataset{RowLabels: rownames, ColLabels: colnames, NRows: nrows, NCols: ncols, Data: data}, nil
}

func (ds *Dataset) Get(row, col int) float64 {
  return ds.Data[row][col]
}

func (ds *Dataset) Row(row int) stats.Float64Data {
  return ds.Data[row]
}

func (ds *Dataset) Col(col int) stats.Float64Data {
  res := make([]float64, 0, ds.NRows)
  for row := 0; row < ds.NRows; row++ {
    res = append(res, ds.Data[row][col])
  }
  return res
}

func (ds *Dataset) Transpose() *Dataset {
  data := make([]stats.Float64Data, 0, ds.NCols)
  for col := 0; col < ds.NCols; col++ {
    data = append(data, ds.Col(col))
  }
  return &Dataset{RowLabels: ds.ColLabels, ColLabels: ds.RowLabels, NRows: ds.NCols, NCols: ds.NRows, Data: data}
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
