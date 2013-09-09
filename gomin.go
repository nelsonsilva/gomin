package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"
)

const (
	NAME    = "gomin"
	VERSION = "0.0.1"
)

var (
	basePath   = flag.String("b", "", "base path to the resources (js, css, etc)")
	outputFile = flag.String("o", "", "the output file")
	help       = flag.Bool("h", false, "show this help")
	version    = flag.Bool("v", false, "version info")
)

// Block
type Block struct {
	blockType   string
	idx         []int
	outFilename string // the ouput filename
	content     []byte
}

func (b *Block) process() (err error) {
	p := GetBlockProcessor(b.blockType)

	if p == nil {
		fmt.Printf("No block processor registered for type %s\n", b.blockType)
		return
	}

	output, err := p.Process(b.content)

	fmt.Printf("Writing %s ...\n", b.blockType)

	outputFile := path.Join(*basePath, b.outFilename)
	err = ioutil.WriteFile(outputFile, output, 0644)

	return
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gomin file.html [flags]", "")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	if *version {
		fmt.Printf("%s v%s\n", NAME, VERSION)
		return
	}

	if flag.NArg() != 1 {
		usage()
	}

	file := flag.Arg(0)
	fileContent, err := ioutil.ReadFile(file)

	if err != nil {
		return // there was a problem reading the file.
	}

	// Register the block processors
	RegisterBlockProcessor(&JSBlockProcessor)
	RegisterBlockProcessor(&CSSBlockProcessor)

	// Find with regexp (non greedy .*?)
	// <!-- build:js js/scroll.min.js --> * <!-- endbuild -->
	blockre, err := regexp.Compile(`<!--(?:\s)*build\:(\w+)(?:\s)*(\S+)(?:\s)*-->((?s).*?)<!--(?:\s)*endbuild(?:\s)*-->`)
	if err != nil {
		return // there was a problem with the regular expression.
	}

	matches := blockre.FindAllSubmatchIndex(fileContent, -1)
	var blocks []*Block = make([]*Block, len(matches))

	// Get the list of blocks
	for idx, match := range matches {

		block := &Block{
			idx:         []int{match[0], match[1]},              // indexes of the block
			blockType:   string(fileContent[match[2]:match[3]]), // the type of the block
			outFilename: string(fileContent[match[4]:match[5]]), // the ouput filename
			content:     fileContent[match[6]:match[7]],         // the block's content
		}

		fmt.Printf("Found %s block[%d:%d] (ouput:%s)\n", block.blockType, block.idx[0], block.idx[1], block.outFilename)

		blocks[idx] = block
	}

	// Process the blocks
	var wg sync.WaitGroup

	for _, block := range blocks {

		// Increment the WaitGroup counter.
		wg.Add(1)

		// Launch a goroutine to process the block.
		go func(block *Block) {

			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			err := block.process()

			if err != nil {
				fmt.Printf("Failed to process block: %s", err)
			}

		}(block)
	}

	// Wait for all block processors complete.
	wg.Wait()

	// Rewrite the file
	var (
		idx     int = 0
		content []byte
	)

	// open the output file
	fo, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(fo)

	for _, block := range blocks {

		// write everything up to the beginning of the block
		content = fileContent[idx:block.idx[0]]

		// write a chunk
		if _, err := w.Write(content); err != nil {
			panic(err)
		}

		// write the replacement
		p := GetBlockProcessor(block.blockType)

		if _, err := w.Write([]byte(p.GetReplacement(block.outFilename))); err != nil {
			panic(err)
		}

		// jump to the end of the block
		idx = block.idx[1]
	}

	// write everything up to the end of the file
	content = fileContent[idx:]

	// write a chunk
	if _, err := w.Write(content); err != nil {
		panic(err)
	}

	// Flush the buffer
	if err = w.Flush(); err != nil {
		panic(err)
	}
}
