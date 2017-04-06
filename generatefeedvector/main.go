package main

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "regexp"
  "strings"

  "github.com/mmcdole/gofeed"
  "golang.org/x/net/html"
)

var wordRE = regexp.MustCompile("\\w+")

func main() {
  n_blogs := 0
  apcount := make(map[string]int)
  wordcounts := make(map[string]map[string]int)
  input := bufio.NewScanner(os.Stdin)
  for input.Scan() {
    url := input.Text()
    title, wc, err := getFeedWordCount(url)
    if err != nil {
      log.Printf("Error parsing feed %s: %v", url, err)
      continue
    }
    n_blogs++
    wordcounts[title] = wc
    for word, _ := range wc {
      apcount[word]++
    }
  }
  var wordlist []string
  for word, count := range apcount {
    freq := float64(count)/float64(n_blogs)
    if freq > 0.1 && freq < 0.5 {
      wordlist = append(wordlist, word)
    }
  }
  fmt.Printf("Blog")
  for _, word := range wordlist {
    fmt.Printf("\t%s", word)
  }
  fmt.Printf("\n")
  for blog, wc := range wordcounts {
    fmt.Printf(blog)
    for _, word := range wordlist {
      fmt.Printf("\t%d", wc[word])
    }
    fmt.Printf("\n")
  }
}

func getFeedWordCount(url string) (string, map[string]int, error) {
  log.Println("Parsing feed " + url)
  fp := gofeed.NewParser()
  feed, err := fp.ParseURL(url)
  if err != nil {
    return "", nil, err
  }
  wc := make(map[string]int)
  for _, item := range feed.Items {
    item_count, err := getItemWordCount(item)
    if err != nil {
      return "", nil, fmt.Errorf("Failed to parse item %s: %v", item.Title, err)
    }
    for w, c := range item_count {
      wc[w] += c
    }
  }
  return feed.Title, wc, nil
}

func getItemWordCount(item *gofeed.Item) (map[string]int, error) {
  log.Printf("Parsing item: %s", item.Title)
  wc := make(map[string]int)
  content := getItemContent(item)
  if len(content) == 0 {
    return wc, nil
  }
  r := strings.NewReader(content)
  doc, err := html.Parse(r)
  if err != nil {
    return nil, err
  }
  visit(wc, doc)
  return wc, nil
}

func visit(words map[string]int, n *html.Node) {
  if n.Type == html.TextNode {
    for _, w := range wordRE.FindAllString(n.Data, -1) {
      words[strings.ToLower(w)]++
    }
  }
  for c := n.FirstChild; c != nil; c = c.NextSibling {
    visit(words, c)
  }
}

func getItemContent(item *gofeed.Item) string {
  if len(item.Content) > 0 {
    return item.Content
  }
  if item.Extensions["content"]["encoded"] != nil {
    v :=  item.Extensions["content"]["encoded"][0].Value
    if len(v) > 0 {
      return v
    }
  }
  return item.Description
}
