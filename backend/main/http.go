package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
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
	if req.URL.Query().Has("file"){
		serveFiles(w,req)
		return
	}

	switch req.Method {
	case "POST":
		// Redirect after POST
		defer http.Redirect(w, req, "/", http.StatusMovedPermanently)

		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Cannot parse form: %v", err)
			return
		}
		text := req.FormValue("text")
		db, err := ConnectToDB()
		CheckError(err)
		_, err = db.Exec("INSERT INTO forum.posts (userid, post) VALUES ($1, $2);", "no_user_id", text)
		CheckError(err)
		err = db.Close()
		CheckError(err)
	case "GET":
		sendIndexFileWithPosts(w, req)
	default:
		fmt.Fprintf(w, "Cannot handle method %s", req.Method)
	}
}

func sendIndexFileWithPosts(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		db, err := ConnectToDB()
		CheckError(err)
		rows, err := db.Query("SELECT post FROM forum.posts")
		var newContent string
		for rows.Next() {
			var post string
			err = rows.Scan(&post)
			CheckError(err)
			newContent += fmt.Sprintf(templatePost, post)
		}
		fmt.Fprintf(w, addContentToFile("index.html", newContent))
		db.Close()
	}
}

func addContentToFile(file string, newContent string) string {
	content, err := ioutil.ReadFile("./forum-templates/" + file)
	CheckError(err)
	doc, err := html.Parse(bytes.NewReader(content))
	CheckError(err)
	contentNode, err := getNodeById(doc, "posts")
	CheckError(err)
	var newNode = &html.Node{
		Type:        html.TextNode,
		Data:        newContent,
	}
	contentNode.AppendChild(newNode)

	buffer := bytes.NewBufferString("")
	err = html.Render(buffer, doc)
	output := html.UnescapeString(buffer.String())
	return output
}


func serveFiles(w http.ResponseWriter, req *http.Request) {
	// TODO: This needs a lot of work, but its not neccesary right now.
	fsPath := ""
	if strings.HasSuffix(req.URL.Path, ".css")  {
		w.Header().Set("Content-Type", "text/css")
		fsPath = "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".ico") {
		fsPath = "./forum-templates" + req.URL.Path
	}
	http.ServeFile(w, req, fsPath)
}