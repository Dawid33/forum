package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func restHandler(w http.ResponseWriter, req *http.Request) {

	// Return
	/*
	if req.URL.Query().Has("threadID") {

	}

	 */

	requestItem := strings.TrimPrefix(req.URL.Path, fmt.Sprintf("/%s/", apiPath))
	fsPath := strings.Join([]string{"./public/forum/", requestItem},"")
	content, err := ioutil.ReadFile(fsPath)

	if err != nil {
		log.Println(err)
		return
	}

	var contentToAdd string = ""
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
		fmt.Fprintf(w, "Request method not supported.")
		return
	}

	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	contentNode, err := getNodeById(doc)
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

/*
func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "http://localhost:3001" + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	fmt.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
*/
