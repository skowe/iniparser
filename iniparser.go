package iniparser

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)
const (
	LineEnd = '\n'
	CommentStart = ';'
	KeyValSep = "="	
)

// INI type contains a slice of blocks 
// Blocks are held inside the Blocks []Block variable
// INI type also contains the Raw data
// RawTrimmed contains all the data with comments trimmed
// inside the Raw []byte variable
type INI struct {
	Blocks []Block
	Raw []byte
	RawTrimmed []byte
}


// The Block type contains data all the data that describes the block
// With comments ommited
// Content variable represents the raw data
// while the Data variable represents the key value pairs
type Block struct{
	Name string
	Content []byte
	Data map[string]string
}

// Lines Method is used to return a slice of all lines in the config file
// if trimmed flag is false comments will be included
// otherwise the comments are excluded from the result 
func (c *INI) Lines(trimmed bool) [][]byte{
	lines := make([][]byte, 0)
	var b *bytes.Buffer
	if !trimmed{
		b = bytes.NewBuffer(c.Raw)
	} else {
		b = bytes.NewBuffer(c.RawTrimmed)
	}
	loop := true
	for loop{
		line, err := b.ReadBytes(LineEnd)
		loop = err != io.EOF
		if err != nil && err != io.EOF {
				log.Fatal("eror reading a line:", err)
		}
		if len(line) == 0{
			continue
		} else {
			lines = append(lines, line)
		}
	}

	return lines
}


// Returns a pointer to unparsed INI struct
// Raw attribute should contain all the file contents as a byte slice
// once the INI is loaded call its Parse method to parse all the data in Raw attribute
func NewINI(pathToIni string) *INI {
	f, err := os.Open(pathToIni)
	if err != nil {
		log.Fatal("Failed to open INI file:", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	
	if err != nil {
		log.Fatal("Failed to open INI file:", err)
	}
	return &INI{
		Raw: b,
		Blocks: make([]Block, 0),
	}
}


// Parses all the data in the INI file which includes
// Removing comments
// Check if all keys have valid names
// Creation of blocks and
// parsing of block data
func (c *INI) Parse() {
	trimComments(c)
	extractBlocks(c)

	log.Println(c.Blocks)
}

func extractBlocks(c *INI){
	lines := c.Lines(true)
	regexpCh := regexp.MustCompile(`\[([\w\$-]+)\]`)
	var blk Block
	cnt := 0
	for i, line := range lines {
		if regexpCh.Match(line){
			if cnt != 0 {
				c.Blocks = append(c.Blocks, blk)
			}
			blk = Block{
				Name: regexpCh.FindStringSubmatch(string(line))[1],
				Content: make([]byte, 0),
				Data: make(map[string]string),
			}
			blk.AppendContent(line)
			cnt++
		} else if cnt != 0 {
			kv := strings.Split(string(line), KeyValSep)
			if validKey(kv[0]) && len(kv) == 2{
				blk.AppendContent(line)
				
				blk.AddData(kv)
			} else if len(kv) == 2 {
				log.Fatal("A key can not contain the ';' character in its name")
			} else if len(kv) <2 {
				log.Panicf("Line %d is invalid, is not a valid key value pair of a comment", i)
			}
		}
	}
	c.Blocks = append(c.Blocks, blk)
	
}


func validKey(key string) bool{
	return !strings.Contains(key, ";")
}

func (b *Block) AppendContent(c []byte){
	b.Content = append(b.Content, c...)
}

func (b *Block) AddData(kv []string){
	b.Data[kv[0]] = kv[1]
}

// Trims out all lines starting with ;
// and all content starting with ; preceeded by a blank character
func trimComments(c *INI){
	lines := c.Lines(false)
	res := make([]byte, 0)
	for _, line := range lines {
		
		if line[0] == CommentStart {
			continue
		} else {
			for i, x := range line{
				if x == CommentStart && (line[i-1] == ' ' || line[i-1] == '\t' || line[i-1] == ']'){
					line = line[0:i]
					line = append(line, '\n')
					break
				}
			}
			res = append(res, line...)
		}
	}
	c.RawTrimmed = res
}