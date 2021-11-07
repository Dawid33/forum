package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
)

const templateThread = ``

func GetSQLFile(name string) string {
	data, _ := f.ReadFile(fmt.Sprintf("sql/%s.sql",name))
	return string(data)
}

func getTemplateFile(file string) *html.Node {
	content, err := ioutil.ReadFile("./forum-templates/" + file)
	Panic(err)
	doc, err := html.Parse(bytes.NewReader(content))
	Panic(err)
	return doc
}

func getSubTemplateFile(file string) string {
	doc := getTemplateFile(file)
	// Parsing html adds stuff that doesn't exist in the file.
	// TODO: Make this more efficent
	doc, err := getNodeByTag(doc, "body")
	Panic(err)
	start := doc.FirstChild
	var output string
	for start != nil {
		output += htmlToString(start)
		start = start.NextSibling
		if start == doc.FirstChild {
			break
		}
	}
	return output
}