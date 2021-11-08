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

func getTemplateFile(file string) (*html.Node, error) {
	content, err := ioutil.ReadFile("./forum-templates/" + file)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func getSubTemplateFile(file string) string {
	doc, err := getTemplateFile(file)
	Panic(err)
	// Parsing html adds stuff that doesn't exist in the file.
	doc, err = getNodeByTag(doc, "body")
	Panic(err)
	return getContentFromNode(doc)
}