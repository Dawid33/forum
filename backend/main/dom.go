package main

import (
	"errors"
	"golang.org/x/net/html"
)

func getNodeById(doc *html.Node, id string) (*html.Node, error) {
	var contentNode *html.Node
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, x := range node.Attr {
				if x.Key == "id" && x.Val == id {
					contentNode = node
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if contentNode != nil {
		return contentNode, nil
	}
	return nil, errors.New("Missing <content> in the node tree")
}
