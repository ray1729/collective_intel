package main

import (
  "fmt"
  "log"
  "os"

  "golang.org/x/net/html"
)

func main () {
  if len(os.Args) != 2 {
    log.Fatalf("Usage: parsefeedlist SOURCE")
  }
  source := os.Args[1]
  file, err := os.Open(source)
  if err != nil {
    log.Fatalf("%v", err)
  }
  defer file.Close()
  doc, err := html.Parse(file)
  for _, url := range find_feeds(nil, doc) {
    fmt.Println(url)
  }
}

func find_feeds(feeds []string, node *html.Node) []string {
  if node.Type == html.ElementNode && node.Data == "tbody" {
    feeds = find_feeds_in_table(feeds, node)
    return feeds
  }
  for c := node.FirstChild; c != nil; c = c.NextSibling {
    feeds = find_feeds(feeds, c)
  }
  return feeds
}

func find_feeds_in_table(feeds []string, node *html.Node) []string {
  if node.Type == html.ElementNode && node.Data == "a" {
    for _, attr := range node.Attr {
      if attr.Key == "href" {
        feeds = append(feeds, attr.Val)
        return feeds
      }
    }
  }
  for c := node.FirstChild; c != nil; c = c.NextSibling {
    feeds = find_feeds_in_table(feeds, c)
  }
  return feeds
}

// func find_feeds(feeds []string, node *html.Node) []string {
//   if node.Type == html.ElementNode && node.Data == "td" {
//     for c := node.FirstChild; c != nil; c = c.NextSibling {
//       if c.Type == html.ElementNode && c.Data == "a" {
//         attrs := make(map[string]string)
//         for _, a := range c.Attr {
//           attrs[a.Key] = a.Val
//         }
//         if attrs["rel"] == "nofollow" && attrs["class"] != "external" {
//           feeds = append(feeds, attrs["href"])
//           break
//         }
//       }
//     }
//     return feeds
//   }
//   for c := node.FirstChild; c != nil; c = c.NextSibling {
//     feeds = find_feeds(feeds, c)
//   }
//   return feeds
// }
