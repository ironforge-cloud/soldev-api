package utils

import (
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetImageIfExist will try to find an image using the html meta tags
func GetImageIfExist(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error: %v ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("Error while parsing html: %v ", err)
		return ""
	}

	imageUrl := parseDocument(doc)

	// Fixing Figment img 404
	if imageUrl == "https://figment.io/wp-content/uploads/2019/08/figment-networks-logo.jpg" {
		return ""
	}

	return imageUrl
}

// Utility to parse html document.
func parseDocument(n *html.Node) string {
	if n.Type == html.ElementNode && n.Parent.Data == "head" && n.Data == "meta" && len(n.Attr) == 2 {
		tagName := n.Attr[0]
		tagContent := n.Attr[1]

		if tagName.Val == "og:image" {
			return tagContent.Val
		}

		if tagName.Val == "twitter:image" {
			return tagContent.Val
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		img := parseDocument(c)
		if img != "" {
			return img
		}
	}

	return ""
}
