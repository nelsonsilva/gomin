package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	blockProcessors = make(map[string]BlockProcessor)
	// JS Block Processor
	JSBlockProcessor = RegexpBlockProcessor{
		blockType:   "js",
		regexp:      `<script(?:.*)src="([^"]*)`,
		replacement: `<script src="$FILE" type="text/javascript"></script>`,
		inline:      `<script type="text/javascript">$CONTENT</script>`,
		compiler:    &ClosureJSCompiler{compilationLevel: "SIMPLE_OPTIMIZATIONS"},
	}
	// CSS Block Processor
	CSSBlockProcessor = RegexpBlockProcessor{
		blockType:   "css",
		regexp:      `<link(?:.*)href="([^"]*)`,
		replacement: `<link href="$FILE" rel="stylesheet">`,
		inline:      `<style>$CONTENT</style>`,
		compiler:    &YUICompiler{fileType: "CSS"},
	}
)

type BlockProcessor interface {
	Init()
	GetType() string
	Process(block []byte) (output []byte, files []string, err error)
	GetReplacement(file string) string
	GetInlineReplacement(content string) string
}

// Regexp Block Processor
type RegexpBlockProcessor struct {
	blockType   string
	regexp      string
	replacement string
	inline      string
	compiler    Compiler
	_re         *regexp.Regexp
}

func (p *RegexpBlockProcessor) Init() {
	re, err := regexp.Compile(p.regexp)
	if err != nil {
		return // there was a problem with the regular expression.
	}
	p._re = re
}

func (p *RegexpBlockProcessor) GetType() string {
	return p.blockType
}

func (p *RegexpBlockProcessor) GetInlineReplacement(content string) string {
	return strings.Replace(p.inline, "$CONTENT", content, -1)
}

func (p *RegexpBlockProcessor) GetReplacement(filename string) string {
	return strings.Replace(p.replacement, "$FILE", filename, -1)
}

func (p *RegexpBlockProcessor) Process(block []byte) (output []byte, files []string, err error) {
	matches := p._re.FindAllSubmatch(block, -1)
	files = make([]string, len(matches))
	for i, match := range matches {
		fmt.Printf("Processing %s ...\n", match[1])
		files[i] = string(match[1])
	}

	output, err = p.compiler.Compile(files)
	return
}

/// Compiler
type Compiler interface {
	Compile(files []string) (code []byte, err error)
}

func RegisterBlockProcessor(b BlockProcessor) {
	blockProcessors[b.GetType()] = b
	b.Init()
}

func GetBlockProcessor(t string) BlockProcessor {
	return blockProcessors[t]
}
