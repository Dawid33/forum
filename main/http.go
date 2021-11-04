package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const templatePost = `<li class="post">
<a href="https://google.com"><span>%s</span></a>
<a href="https://google.com">(google.com)</a><br/>
<span class="post-subtext">400 points by jim 2 hours age | hide | <a href="post.html">198
comments</a></span>
</li>`

var PORT = 3000

type commentRequest struct {
	Post string
}

type Page struct {
	body []byte
}

func startHttpServer() {
	mux := http.NewServeMux()
	publicFiles := http.FileServer(http.Dir("./public/forum"))
	consoleFiles := http.FileServer(http.Dir("./public/console-dev"))

	mux.Handle("/", publicFiles)
	mux.Handle("/console/", http.StripPrefix("/console/", consoleFiles))
	mux.HandleFunc("/api/", restHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("MY-HEADER", "COOL")
	fmt.Fprintf(w, "Hi there, this will be the management forum!")
}

func restHandler(w http.ResponseWriter, req *http.Request) {
	requestItem := strings.TrimPrefix(req.URL.Path, "/api/")
	fsPath := strings.Join([]string{"./public/forum/", requestItem},"")
	content, err := ioutil.ReadFile(fsPath)
	if err != nil {
		log.Println(err)
		return
	}

	var contentToAdd string = ""
	strings.NewReplacer()
	switch req.Method {
	case "POST":
		decoder := json.NewDecoder(req.Body)

		var comment commentRequest
		err := decoder.Decode(&comment)
		if err != nil && err != io.EOF {
			fmt.Fprintf(w, err.Error())
			return
		}
		contentToAdd += comment.Post
	default:
		fmt.Fprintf(w, "Request method supported.")
		return
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	contentNode, err := getContent(doc)
	if err != nil {
		log.Fatal(err)
	}
	var newNode = &html.Node{
		Parent:      nil,
		FirstChild:  nil,
		LastChild:   nil,
		PrevSibling: nil,
		NextSibling: nil,
		Type:        html.TextNode,
		DataAtom:    0,
		Data:        fmt.Sprintf(templatePost, contentToAdd),
		Namespace:   "",
		Attr:        nil,
	}

	contentNode.AppendChild(newNode)
	buffer := bytes.NewBufferString("")
	if err := html.Render(buffer, doc); err != nil {
		log.Fatal(err)
	}
	output := html.UnescapeString(buffer.String())

	err = os.Remove(fsPath)
	if err != nil {
		log.Fatal(err)
	}
	fileHandle, err := os.Create(fsPath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fileHandle.WriteString(output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Success")
}

func getContent(doc *html.Node) (*html.Node, error) {
	var contentNode *html.Node
	var crawler func(*html.Node)

	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, x := range node.Attr {
				if x.Key == "id" && x.Val == "posts" {
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

