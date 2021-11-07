package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
)

const (
	Undefined = 0
	SubFileQuery = 1
	DatabaseQuery = 2
)

var validQueries = []string {
	"id",
}

type QueryParameter struct {
	key string
	value string
}

type Query struct {
	queryType int
	params []QueryParameter
}

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
	return nil, errors.New(fmt.Sprintf("Missing id %s in the node tree", id))
}
func getNodeByTag(doc *html.Node, tag string) (*html.Node, error) {
	var contentNode *html.Node
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if node.Data == tag {
				contentNode = node
				return
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
	return nil, errors.New(fmt.Sprintf("Missing tag %s in the node tree", tag))
}
func getNodesWithAttributes(doc *html.Node, attributes []string) []*html.Node {
	var contentNode []*html.Node
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, x := range node.Attr {
				for _, y := range attributes {
					if x.Key == y {
						contentNode = append(contentNode, node)
						return
					}
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	return contentNode
}

func parseAttributes(node *html.Node, req *http.Request) (Query, error) {
	var output Query
	for _, attr := range node.Attr {
		switch attr.Key {
		case "file":
			output.queryType = SubFileQuery
		case "query":
			output.queryType = DatabaseQuery
		default:
			// If the requested data starts with ?, it must come from the url.
			if strings.HasPrefix(attr.Key, "?") {
				for _, x := range validQueries {
					if "?" + x == attr.Key {
						parameter := QueryParameter{key: attr.Key, value: req.URL.Query().Get(x)}
						output.params = append(output.params, parameter)
						break
					}
				}
			} else {
				output.params = append(output.params, QueryParameter{key: attr.Key, value: attr.Val})
			}
		}
	}


	if output.queryType == Undefined {
		return Query{}, errors.New("cannot parse attributes")
	}
	return output, nil
}

func fulfillQuery(node *html.Node, query Query) error {
	switch query.queryType {
	case SubFileQuery:
		for _, x := range query.params {
			fmt.Println("Getting file : ", x.key + ".html")
			content := getSubTemplateFile(x.key + ".html")
			addContentToNode(node, content)
		}
	case DatabaseQuery:
		var requestedDataType = ""
		var category = ""
		var template = ""
		var requestedData []QueryParameter
		for _, x := range query.params {
			switch x.key {
			case "type":
				requestedDataType = x.value
			case "category":
				category = x.value
			case "orderby":
				//TODO order requested info
			case "template":
				template = x.value
			default:
				// If there is no key / value pair, must be the requested information
				// Or if it starts with ?, the information must come from url
				if x.value == "" {
					requestedData = append(requestedData, QueryParameter{
						key: x.key,
						value: "",
					})
				}
				if strings.HasPrefix(x.value, "?") {
					requestedData = append(requestedData, QueryParameter{
						key: x.key,
						value: x.value,
					})
				}
			}
		}

		// Must fulfill basic requirements for db query
		if requestedDataType != "" && requestedData != nil {
			var file = ""
			if template != "" {
				newFile := getSubTemplateFile(template + ".html")
				file = newFile
			} else if node.FirstChild != nil {
				start := node.FirstChild
				var output string
				for start != nil {
					output += htmlToString(start)
					start = start.NextSibling
					if start == node.FirstChild {
						break
					}
				}
				// Get data and orphan child as we will be replacing it with data.
				node.FirstChild = nil
				node.LastChild = nil
				file = output
			}

			var newContent string
			switch requestedDataType {
			case "threads":
				var threads []Thread
				// DB Query here.
				if category != "" {
					newThreads, err := getThreadsWithCategory(category)
					threads = newThreads
					PrintError(err)
				} else {
					//TODO: Get all threads in the database
					threads = []Thread{}
				}

				for _, x := range threads {
					var dbInfo []string
					for _, y := range requestedData {
						switch y.key {
						case "title":
							dbInfo = append(dbInfo, x.title)
						case "content":
							dbInfo = append(dbInfo, x.content)
						case "userid":
							dbInfo = append(dbInfo, x.userId)
						case "threadid":
							dbInfo = append(dbInfo, strconv.FormatUint(x.threadId, 10))
						default:
							if strings.HasPrefix(y.key, "?") {
								dbInfo = append(dbInfo, y.value)
							}
						}
					}
					dbInfoInterface := make([]interface{}, len(dbInfo))
					for i, v := range dbInfo {
						dbInfoInterface[i] = v
					}

					newItem := fmt.Sprintf(file, dbInfoInterface...)
					newContent += newItem
				}
			case "comments":

			}

			addContentToNode(node, newContent)
		}
	default:
		return errors.New("cannot recognise query")
	}
	return nil
}

func addContentToNode(input *html.Node, newContent string) {
	var newNode = &html.Node{
		Type:        html.TextNode,
		Data:        newContent,
	}
	input.AppendChild(newNode)
}

func htmlToString(input *html.Node) string {
	buffer := bytes.NewBufferString("")
	err := html.Render(buffer, input)
	CheckError(err)
	gohtml.Condense = true
	return gohtml.Format(html.UnescapeString(buffer.String()))
}