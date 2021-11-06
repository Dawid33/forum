package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const PORT = 3000

type commentRequest struct {
	Thread string
}

type Page struct {
	body []byte
}

func startHttpServer() {
	fmt.Println("Starting HTTP Server!")

	mux := http.NewServeMux()
	mux.HandleFunc("/", fileSendHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

// Function that handles all regular requests
func fileSendHandler(w http.ResponseWriter, req *http.Request) {

	isStatic, path := isStaticFile(w, req)
	if isStatic {
		http.ServeFile(w, req, path)
	}

	switch req.Method {
	case "POST":
		err := req.ParseForm()
		redirectTo404OnError(w, req, err)

		switch req.URL.Path {
		case "/post":
			acceptNewThread(connectToDB(), w, req)
		case "/comment":
			acceptNewComment(connectToDB(), w, req)
		default:
			redirectToUrl(w, req, "404.html")
		}
	case "GET":
		sendGeneratedHtmlFile(w, req)
	default:
		fmt.Fprintf(w, "Cannot handle method %s", req.Method)
	}
}

func sendGeneratedHtmlFile(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/index.html" {
		redirectToUrl(w, req, "/")
		return
	}

	switch req.URL.Path {
	case "/":
		serveIndexFile(w, req)
	case "/thread":
		serveThreadFile(w, req)
	}
}

func serveIndexFile(w http.ResponseWriter, req *http.Request) {
	db := connectToDB()
	rows, err := db.Query("SELECT * FROM forum.posts")
	var newThreads string
	var prefetchLinks string
	for rows.Next() {
		var threadId uint64
		var userId string
		var title string
		var content string
		err = rows.Scan(&threadId, &userId, &title, &content)
		CheckError(err)
		newThreads += fmt.Sprintf(templateThread,title, threadId)
		prefetchLinks += fmt.Sprintf("	<link rel=\"prefetch\" href=\"thread?%d\">\n", threadId)
	}
	doc := getTemplateFile("index.html")
	addContentToTagByIdInDoc(doc, "posts", newThreads)
	addContentToTagInDoc(doc, "head", prefetchLinks)
	_, err = fmt.Fprintf(w, htmlNodeToString(doc))
	CheckError(err)
	_ = db.Close()
}

func serveThreadFile (w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Has("id") {
		id := req.URL.Query().Get("id")
		post, err := getThread(id)
		if err != nil {
			redirectToUrl(w, req, "404.html")
		}
		comments, err := getCommentsInThread(connectToDB(), id)

		doc := getTemplateFile("thread.html")

		addContentToTagByIdInDoc(doc, "comment-form", fmt.Sprintf(templateCommentForm, post.threadId, 0, post.threadId))
		addContentToTagByIdInDoc(doc, "post-title", post.title)
		addContentToTagByIdInDoc(doc, "post-content", post.content)
		addContentToTagByIdInDoc(doc, "comments", comments)
		fmt.Fprintf(w, htmlNodeToString(doc))
	} else {
		redirectToUrl(w, req, "404.html")
	}
}

func acceptNewThread(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	text := req.FormValue("text")
	if text == ""{
		redirectTo404OnError(w, req, errors.New("text field does not exist in post form"))
		return
	}

	_, err := db.Exec("INSERT INTO forum.posts (userid, title, content) VALUES ($1, $2, $3);", "no_user_id", "Default Title", text)
	if err != nil {redirectTo503OnError(w, req, err); return }
	_ = db.Close()

	if gotoUrl := req.FormValue("goto"); gotoUrl == "" {
		redirectTo404OnError(w, req, err)
	}
	redirectToUrl(w, req, req.FormValue("goto"))
}

func acceptNewComment(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	if req.FormValue("parentId") == "" {
		redirectTo404OnError(w, req, errors.New("parentId field does not exist in post form"))
		return
	}
	parentId, err := strconv.ParseInt(req.FormValue("parentId"), 10, 32)
	if err != nil {redirectTo503OnError(w, req, err); return }

	if req.FormValue("thread") == "" {
		redirectTo404OnError(w, req, errors.New("thread field does not exist in post form"))
		return
	}
	threadId, err := strconv.ParseInt(req.FormValue("thread"), 10, 32)
	if err != nil {redirectTo503OnError(w, req, err); return }

	content := req.FormValue("text")

	_, err = db.Exec("INSERT INTO forum.comments (threadid, parentid, kidsid, userid, content) VALUES ($1, $2, $3, $4, $5);", threadId, parentId, nil, "Anonymous", content)
	if err != nil {redirectTo503OnError(w, req, err); return }

	if gotoUrl := req.FormValue("goto"); gotoUrl == "" {
		redirectTo404OnError(w, req, err)
	}
	redirectToUrl(w, req, req.FormValue("goto"))
}

// TODO: This code needs to be rewritten and done properly.
func isStaticFile(w http.ResponseWriter, req *http.Request) (bool, string) {
	if strings.HasSuffix(req.URL.Path, ".css")  {
		w.Header().Set("Content-Type", "text/css")
		return true, "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".js")  {
		w.Header().Set("Content-Type", "application/javascript")
		return true, "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".ico") {
		return true, "./forum-templates" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".html") {
		return true, "./forum-templates" + req.URL.Path
	}
	return false, ""
}

func redirectToUrl(w http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

func redirectTo404OnError(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		PrintError(err)
		redirectToUrl(w, req, "404.html")
	}
}

func redirectTo503OnError(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		PrintError(err)
		redirectToUrl(w, req, "404.html")
	}
}