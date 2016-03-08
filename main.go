package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocument(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("div.article-body-inner > a[target=_blank]").Each(func(_ int, s *goquery.Selection) {
		if imgurl, ok := s.Attr("href"); ok {
			for {
				fmt.Printf("Downloading %s...", imgurl)
				err = download(imgurl)
				if err != nil {
					fmt.Println(err)
					fmt.Println("Retry...")
					continue
				}
				fmt.Println("done")
				break
			}
		}
	})
}

func download(rawurl string) error {
	filename, err := fileNameOf(rawurl)
	if err != nil {
		return err
	}
	resp, err := http.Get(rawurl)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

var reInPath = regexp.MustCompile("[^/]+$")

func fileNameOf(rawurl string) (string, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	file := reInPath.FindString(url.Path)
	if file == "" {
		return "", fmt.Errorf("Filename not found: %s", rawurl)
	}
	return file, nil
}
