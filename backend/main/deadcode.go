package main

/*
func serveIndexFile(w http.ResponseWriter, req *http.Request) {
	db := connectToDB()
	rows, err := db.Query("SELECT * FROM forum.posts")

	var prefetchLinks string
	for rows.Next() {
		var threadId uint64
		var userId string
		var category string
		var title string
		var content string
		err = rows.Scan(&threadId, &userId, &category, &title, &content)
		CheckError(err)
		//newThreads += fmt.Sprintf(, title, threadId)
		prefetchLinks += fmt.Sprintf("	<link rel=\"prefetch\" href=\"thread?%d\">\n", threadId)
	}
	doc := getTemplateFile("index.html")
	/*
	addContentToTagByIdInDoc(doc, "thread-form", getSubTemplateFile("thread_form.html"))
	addContentToTagByIdInDoc(doc, "posts", newThreads)
	addContentToTagInDoc(doc, "head", prefetchLinks)

	_, err = fmt.Fprintf(w, htmlToString(doc))
	CheckError(err)
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


		addContentToTagByIdInDoc(doc, "comment-form", fmt.Sprintf(getSubTemplateFile("comment_form.html"), post.threadId, 0, post.threadId))
		addContentToTagByIdInDoc(doc, "post-title", post.title)
		addContentToTagByIdInDoc(doc, "post-content", post.content)
		addContentToTagByIdInDoc(doc, "comments", comments)

		fmt.Fprintf(w, htmlToString(doc))
	} else {
		redirectToUrl(w, req, "404.html")
	}
}
*/