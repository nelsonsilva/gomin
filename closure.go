package main

// Closure compiler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const CLOSURE_URI = "http://closure-compiler.appspot.com/compile"

type ClosureJSCompiler struct {
	compilationLevel string
}

func (c *ClosureJSCompiler) Compile(files []string) (code []byte, err error) {
	fmt.Println("[Closure] Compiling JS ...")

	jscode := ConcatFiles(files)

	resp, err := http.PostForm(CLOSURE_URI, url.Values{
		"compilation_level": {c.compilationLevel},
		"output_info":       {"compiled_code", "statistics", "errors", "warnings"},
		"output_format":     {"json"},
		"js_code":           {jscode}})

	if err != nil {
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var t closure_response
	err = decoder.Decode(&t)
	if err != nil {
		return
	}

	fmt.Printf("[Closure]\n%s\n... done!\n", t.Statistics)

	code = []byte(t.CompiledCode)

	return
}

type closure_compile_message struct {
	Charno int
	Error  string
	Lineno int
	File   string
	Type   string
	Line   string
}

type closure_server_error struct {
	Code  int
	Error string
}

type closure_statistics struct {
	OriginalSize   int
	CompressedSize int
	CompileTime    int
}

func (s closure_statistics) String() string {
	return fmt.Sprintf("OriginalSize:\t%d\nCompressedSize:\t%d\nCompileTime:\t%d\n",
		s.OriginalSize, s.CompressedSize, s.CompileTime)
}

type closure_response struct {
	CompiledCode string
	Errors       []closure_compile_message
	Warnings     []closure_compile_message
	ServerErrors []closure_server_error
	Statistics   closure_statistics
}
