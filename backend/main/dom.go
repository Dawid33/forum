package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
)

const templateThread = `<li class="post">
<a href="https://google.com"><span>%s</span></a>
<a href="https://google.com">(google.com)</a><br/>
<span class="post-subtext">400 points by jim 2 hours age | hide | <a href="thread?id=%d">198
comments</a></span>
</li>`

const templateComment = `
<div class="comment">
  <div class="comment-heading">
	<span>
	  <a class="commenter-username" href="http://user">%s</a>
	  <span class="time-commented">43 minutes ago</span>
	  <input
		type="checkbox"
		id="collapsible"
		class="toggle collabsible-input"
	  />
	  <label for="collapsible" class="lbl-toggle"></label>
	  <div class="collapsible-content">
		<div class="content-inner">
		  <p>%s</p>
		</div>
	  </div>
	</span>
  </div>
</div>`

const templateCommentForm = `
<form method="post" action="/comment">
	<input type="hidden" name="thread" value="%d" />
	<input type="hidden" name="parentId" value="%d" />
	<input type="hidden" name="goto" value="thread?id=%d" />
	<textarea name="text" rows="6" cols="60">Default Comment</textarea>
	<br /><br />
	<input type="submit" value="add comment" />
</form>
`

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