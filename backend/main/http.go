package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const templatePost = `<li class="post">
<a href="https://google.com"><span>%s</span></a>
<a href="https://google.com">(google.com)</a><br/>
<span class="post-subtext">400 points by jim 2 hours age | hide | <a href="post.html">198
comments</a></span>
</li>`

var PORT = 3000
var isDebug = true
const apiPath = "api"

type commentRequest struct {
	Post string
}

type Page struct {
	body []byte
}

func startHttpServer() {
	fmt.Println("Starting HTTP Server!")

	mux := http.NewServeMux()
	mux.HandleFunc("/", fileSendHandler)

	mux.HandleFunc(fmt.Sprintf("/%s/", apiPath), restHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

// Function that handles all regular requests
func fileSendHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("path :", req.URL.Path)
	fsPath := ""
	if req.URL.Path == "/" {
		fsPath = "./forum-templates/index.html"
	} else {
		fsPath = "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}

	content, err := ioutil.ReadFile(fsPath)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(content)
}