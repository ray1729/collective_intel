package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "strings"
)

type Dataset struct {
  Rownames []string
  Colnames []string
  Data     [][]float64
}

func main() {
  filename := os.Args[1]
  ds, _ := ReadDataset(filename)
  fmt.Println(ds)
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
