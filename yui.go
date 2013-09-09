package main

// YUI compiler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const YUI_URI = "http://refresh-sf.com/yui/"

type YUICompiler struct {
	fileType string
}

func (c *YUICompiler) Compile(files []string) (code []byte, err error) {
	fmt.Printf("[YUI] Compiling %s ...\n", c.fileType)

	cssCode := ConcatFiles(files)

	resp, err := http.PostForm(YUI_URI, url.Values{
		"type":         {c.fileType},
		"redirect":     {"1"},
		"compresstext": {cssCode}})

	if err != nil {
		return
	}
	defer resp.Body.Close()
	code, err = ioutil.ReadAll(resp.Body)

	fmt.Printf("[YUI] ... done!\n")

	return
}
