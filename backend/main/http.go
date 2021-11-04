package main

import (
	"fmt"
	"log"
	"net/http"
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
	publicFiles := http.FileServer(http.Dir("./public/forum-templates"))
	mux.Handle("/", publicFiles)

	consoleFiles := http.FileServer(http.Dir("./public/console"))
	mux.Handle("/console/", http.StripPrefix("/console/", consoleFiles))

	mux.HandleFunc(fmt.Sprintf("/%s/", apiPath), restHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

// Function that handles all requests to /console/**
func consoleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Frame-Options", "DENY")
	fmt.Println("Accepted console request")
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



