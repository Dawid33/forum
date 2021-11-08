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
	/*
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	*/
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

// Function that handles all regular requests
func fileSendHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		err := req.ParseForm()
		if err != nil {
			PrintError(err)
			redirectToUrl(w,req, "404.html")
			return
		}

		switch req.URL.Path {
		case "/post-thread":
			acceptNewThread(db, w, req)
		case "/post-comment":
			acceptNewComment(db, w, req)
		default:
			redirectToUrl(w, req, "404.html")
		}
	case "GET":
		isStatic, path := isStaticFile(w, req)
		if isStatic {
			http.ServeFile(w, req, path)
			return
		}
		// Need to generate html file from template.
		sendGeneratedHtmlFile(w, req)
	default:
		fmt.Fprintf(w, "Cannot handle method %s", req.Method)
	}
}

func sendGeneratedHtmlFile(w http.ResponseWriter, req *http.Request) {
	// redirect index.html to /
	if req.URL.Path == "/index.html" {
		redirectToUrl(w, req, "/")
		return
	}

	var filename string
	if req.URL.Path == "/" {
		filename = "index.html"
	} else {
		filename = strings.TrimPrefix(req.URL.Path, "/") + ".html"
	}

	doc, err := getTemplateFile(filename)
	if err != nil {
		PrintError(err)
		redirectToUrl(w, req, "404.html")
		return
	}

	var  attributes = []string {
		"query",
		"file",
	}
	nodes := getNodesWithAttributes(doc, attributes)
	if len(nodes) != 0 {
		for _, node := range nodes {
			query, err := parseAttributes(node, req)
			if err != nil {
				continue
			}
			err = fulfillQuery(node, query)
			if err != nil {
				continue
			}
		}
	}
	fmt.Fprintf(w, htmlToString(doc))
}

func acceptNewThread(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	text, err := getFieldFromPost("text", w, req)
	if err != nil {redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery); return}
	title, err := getFieldFromPost("title", w, req)
	if err != nil {redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery); return}

	_, err = db.Exec("INSERT INTO forum.posts (userid, category, title, content) VALUES ($1, $2, $3, $4);", "Anonymous", "default", title, text)
	if err != nil {
		PrintError(err)
		redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery)
		return
	}

	if gotoUrl := req.FormValue("goto"); gotoUrl == "" {
		redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery)
		return
	}
	redirectToUrl(w, req, req.FormValue("goto"))
}

func acceptNewComment(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	thread, err := getFieldFromPost("thread", w, req)
	if err != nil {redirectToUrl(w, req, "/"); return}

	content, err := getFieldFromPost("text", w, req)
	if err != nil {redirectToUrl(w, req, "/thread?threadid=" + thread); return}

	parentId, err := getFieldFromPost("parentId", w, req)
	if err != nil {redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery); return}

	threadId, err := strconv.ParseInt(thread, 10, 32)
	if err != nil {redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery); return }

	_, err = db.Exec("INSERT INTO forum.comments (threadid, parentid, kidsid, userid, content) VALUES ($1, $2, $3, $4, $5);", threadId, parentId, nil, "Anonymous", content)
	if err != nil {redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery); return }

	if gotoUrl := req.FormValue("goto"); gotoUrl == "" {
		redirectToUrl(w, req, req.URL.RawPath + req.URL.RawQuery)
		return
	}
	redirectToUrl(w, req, req.FormValue("goto"))
}

func getFieldFromPost(field string, w http.ResponseWriter, req *http.Request) (string, error){
	text := req.FormValue(field)
	if text == ""{
		err := errors.New(fmt.Sprintf("%s field does not exist in post form", field))
		return "", err
	}
	return text, nil
}

// TODO: This code needs to be rewritten and done properly.
func isStaticFile(w http.ResponseWriter, req *http.Request) (bool, string) {
	if strings.HasSuffix(req.URL.Path, ".css")  {
		w.Header().Set("Content-Type", "text/css")
		return true, "./public" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".js")  {
		w.Header().Set("Content-Type", "application/javascript")
		return true, "./public" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".ico") {
		return true, "./public" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".html") {
		return true, "./public" + req.URL.Path
	}
	if strings.HasSuffix(req.URL.Path, ".js.map") {
		return true, "./public" + req.URL.Path
	}
	return false, ""
}

func redirectToUrl(w http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

func redirectTo503OnError(w http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		PrintError(err)
		redirectToUrl(w, req, "404.html")
	}
}