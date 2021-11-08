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
	"threadid",
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
					if "?" + x == attr.Key && req.URL.Query().Has(x){
						parameter := QueryParameter{key: attr.Key, value: req.URL.Query().Get(x)}
						output.params = append(output.params, parameter)
						break
					}
				}
			} else if strings.HasPrefix(attr.Val, "?") {
				for _, x := range validQueries {
					if "?" + x == attr.Val && req.URL.Query().Has(x){
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
			replaceNodeWithContent(node, content)
		}
	case DatabaseQuery:
		var requestedDataType = ""
		var category = ""
		var template = ""
		var threadid = ""
		var requestedData []QueryParameter
		for _, x := range query.params {
			switch x.key {
			case "type":
				requestedDataType = x.value
			case "category":
				category = x.value
			case "template":
				template = x.value
			case "threadid":
				if x.value != "" {
					threadid = x.value
				} else {
					requestedData = append(requestedData, QueryParameter{
						key: x.key,
						value: x.value,
					})
				}
			case "?threadid":
				requestedData = append(requestedData, QueryParameter{
					key: x.key,
					value: x.value,
				})
			default:
				if x.value != "" && x.key == "threadid" {
					threadid = x.value
				}
				// If there is no key / value pair, must be the requested information
				// Or if it starts with ?, the information must come from url
				if x.value == "" {
					requestedData = append(requestedData, QueryParameter{
						key:   x.key,
						value: "",
					})
				}
			}
		}

		// Must fulfill basic requirements for db query
		if requestedData != nil {
			var file = ""
			if template != "" {
				newFile := getSubTemplateFile(template + ".html")
				file = newFile
			} else if node.FirstChild != nil {
				// Get data and orphan child as we will be replacing it with data.
				file = getContentFromNode(node)
				node.FirstChild = nil
				node.LastChild = nil
			}

			var newContent string
			switch requestedDataType {
			case "threads":
				var threads []Thread
				// TODO: Allow more complex queries
				var err error = nil
				if category != "" {
					categoryQuery := fmt.Sprintf("WHERE posts.category = '%s'", category)

					threads, err = getThreads(db, categoryQuery)

				} else {
					threads, err = getThreads(db, "")
				}
				if err != nil || len(threads) == 0{
					PrintError(err)
					return err
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
				if threadid == "" {
					err := errors.New("no thread ID")
					fmt.Println(err)
					return err
				}
				commentQuery := fmt.Sprintf("WHERE comments.threadid = %s", threadid)
				comments, err := getComments(db, commentQuery)
				if err != nil {
					PrintError(err)
					return err
				}

				for _, x := range comments {
					var dbInfo []string
					for _, y := range requestedData {
						switch y.key {
						case "content":
							dbInfo = append(dbInfo, x.content)
						case "userid":
							dbInfo = append(dbInfo, x.userId)
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
			default:

				var dbInfo []string
				for _, y := range requestedData {
					switch y.key {
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

				newContent = fmt.Sprintf(file, dbInfoInterface...)
			}
			replaceNodeWithContent(node, newContent)
		}
	default:
		return errors.New("cannot recognise query")
	}
	return nil
}

func getContentFromNode(node *html.Node) string {
	if node.FirstChild == nil {
		return ""
	}
	start := node.FirstChild
	var output string
	for start != nil {
		output += htmlToString(start)
		start = start.NextSibling
		if start == node.FirstChild {
			break
		}
	}
	return output
}

func replaceNodeWithContent(input *html.Node, newContent string) {
	input.FirstChild = nil
	input.LastChild = nil
	input.Type = html.TextNode
	input.Data = newContent
}

func htmlToString(input *html.Node) string {
	buffer := bytes.NewBufferString("")
	err := html.Render(buffer, input)
	CheckError(err)
	gohtml.Condense = true
	return gohtml.Format(html.UnescapeString(buffer.String()))
}