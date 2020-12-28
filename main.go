package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const EARRING_URL = "https://electriccowboy.bigcartel.com"

type ProductInfo struct {
	Name    string
	Url     string
	InStock bool
}

/**
 * functions for extracting all the necessary info from a given card
 */
func getName(n html.Node) *string {
	var name *string
	name = nil
	traverse(n, func(n html.Node) {
		if hasClass(n, "product-list-thumb-name") {
			child := n.FirstChild
			name = &((*child).Data)
		}
	})
	return name
}

func getURL(n html.Node) *string {
	var url *string
	url = nil
	traverse(n, func(n html.Node) {
		if hasClass(n, "product-list-link") {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					val := attr.Val
					url = &val
				}
			}
		}
	})
	return url
}

func isInStock(n html.Node) bool {
	inStock := false
	traverse(n, func(n html.Node) {
		if hasClass(n, "product-list-thumb-status") {
			child := n.FirstChild
			if (*child).Data == "Sold out" {
				inStock = false
			} else {
				inStock = true
			}
		}
	})
	return inStock
}

/**
 * Utility function to check if a given node has the specified class name
 */
func hasClass(n html.Node, classname string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, s := range strings.Split(attr.Val, " ") {
				if s == classname {
					return true
				}
			}
		}
	}
	return false
}

/**
 * Utility function to find the HTML elements for a product card
 */
func isProductCard(n html.Node) bool {
	return hasClass(n, "product-list-thumb")
}

/**
 * Utility function to traverse all nodes in a given tree and run a
 * callback on each one
 */
func traverse(n html.Node, f func(n html.Node)) {
	// run the callback on the current node
	f(n)
	// get all children of the current node
	if n.FirstChild != nil {
		traverse(*n.FirstChild, f)
	}
	// get all the siblings of the node
	if n.NextSibling != nil {
		traverse(*n.NextSibling, f)
	}
}

/**
 * Format Output
 */
func format(products []ProductInfo) string {
	productMkdownItems := []string{}
	for _, p := range products {
		if p.InStock {
			str := fmt.Sprintf("  - [%s](%s%s)", p.Name, EARRING_URL, p.Url)
			productMkdownItems = append(productMkdownItems, str)
		}
	}
	if len(productMkdownItems) > 0 {
		title := "@reccanti Hoi! Some earrings are in stock!"
		return fmt.Sprintf("# Update\n%s\n\n%s", title, strings.Join(productMkdownItems, "\n"))
	}
	return ""
}

func main() {
	// Create an HTML tree from our given URL
	res, err := http.Get(EARRING_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	tree, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// first, create a list of all the Product Card nodes
	cards := []html.Node{}
	traverse(*tree, func(n html.Node) {
		if isProductCard(n) {
			cards = append(cards, n)
		}
	})

	// get product info from each of the cards
	products := []ProductInfo{}
	for _, c := range cards {

		name := *(getName(*c.FirstChild))
		url := *(getURL(*c.FirstChild))
		inStock := isInStock(*c.FirstChild)

		p := ProductInfo{
			Name:    name,
			Url:     url,
			InStock: inStock,
		}

		products = append(products, p)
	}

	os.Stdout.Write([]byte(format(products)))
}
