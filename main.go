package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func isProductCard(t html.Token) bool {
	for _, attr := range t.Attr {
		if attr.Key == "class" {
			for _, s := range strings.Split(attr.Val, " ") {
				if s == "product-list-thumb" {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	res, err := http.Get("https://electriccowboy.bigcartel.com/products")
	if err != nil {
		fmt.Println(err)
		return
	}
	tokenizer := html.NewTokenizer(res.Body)
	cards := []html.Token{}
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		// fmt.Println(tokenizer.Token().Attr)
		// Process the current token.
		t := tokenizer.Token()
		if isProductCard(t) {
			cards = append(cards, t)
		}
	}
	fmt.Println(len(cards))
}
