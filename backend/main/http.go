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
<span class="post-subtext">400 points by jim 2 hours age | hide | <a href="item?id=%d">198
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
	fsPath := ""
	if strings.HasSuffix(req.URL.Path, ".css")  {
		w.Header().Set("Content-Type", "text/css")
		fsPath = "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".js")  {
		w.Header().Set("Content-Type", "application/javascript")
		fsPath = "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".ico") {
		fsPath = "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".html") {
		fsPath = "./forum-templates" + req.URL.Path
	}
	if fsPath != "" {
		http.ServeFile(w, req, fsPath)
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

		sendPopulatedHtmlFile(w, req)

	default:
		fmt.Fprintf(w, "Cannot handle method %s", req.Method)
	}
}

func sendPopulatedHtmlFile(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/index.html" {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
		return
	}
	if req.URL.Path == "/"{
		db, err := ConnectToDB()
		CheckError(err)
		rows, err := db.Query("SELECT * FROM forum.posts")
		var newPosts string
		var prefetchLinks string
		for rows.Next() {
			var postid uint64
			var userid string
			var post string
			err = rows.Scan(&postid, &userid, &post)
			CheckError(err)
			newPosts += fmt.Sprintf(templatePost, post, postid)
			prefetchLinks += fmt.Sprintf("	<link rel=\"prefetch\" href=\"item?%d\">\n", postid)
		}
		doc := getTemplateFile("index.html")
		addContentToTagByIdInDoc(doc, "posts", newPosts)
		addContentToTagInDoc(doc, "head", prefetchLinks)
		fmt.Fprintf(w, htmlNodeToString(doc))
		db.Close()
	}
	// Query specific comment / post.
	if req.URL.Path == "/item" {
		// Search posts to see if I need to load an entire post.
		if req.URL.Query().Has("id") {
			post, err := GetPost(req.URL.Query().Get("id"))
			CheckError(err)
			doc := getTemplateFile("post.html")
			addContentToTagByIdInDoc(doc, "post", post.post)
			fmt.Fprintf(w, htmlNodeToString(doc))
		} else {

		}
	}
}

func addContentToTagByIdInDoc(input *html.Node, id string, newContent string) {
	contentNode, err := getNodeById(input, id)
	CheckError(err)
	var newNode = &html.Node{
		Type:        html.TextNode,
		Data:        newContent,
	}
	contentNode.AppendChild(newNode)
}
func addContentToTagInDoc(input *html.Node, id string, newContent string) {
	contentNode, err := getNodeByTag(input, id)
	CheckError(err)
	var newNode = &html.Node{
		Type:        html.TextNode,
		Data:        newContent,
	}
	contentNode.AppendChild(newNode)
}

func getTemplateFile(file string) *html.Node {
	content, err := ioutil.ReadFile("./forum-templates/" + file)
	CheckError(err)
	doc, err := html.Parse(bytes.NewReader(content))
	CheckError(err)
	return doc
}

func htmlNodeToString(input *html.Node) string {
	buffer := bytes.NewBufferString("")
	err := html.Render(buffer, input)
	CheckError(err)
	return html.UnescapeString(buffer.String())
}