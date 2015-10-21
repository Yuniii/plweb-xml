package main

import (
	"io/ioutil"
	"strings"
	"fmt"
	"sync"
	"strconv"
	"encoding/xml"
	"encoding/base64"
	"github.com/qiniu/iconv"
)

type Property struct {
	Key string `xml:"key"`
	Value string `xml:"value"`
}

type File struct {
	Path string `xml:"path"`
	Content string `xml:"content"`
}

type Lesson struct {
	XMLName xml.Name `xml:"project"`
	Id int `xml:"id"`
	Title string `xml:"title"`
	Properties []Property `xml:"property"`
	Files []File `xml:"file"`
}

var wg sync.WaitGroup

func main() {
	for i := 1; i <= 24; i++ {
		wg.Add(1)
		go processFile("xml/247/lesson" + strconv.Itoa(i) + ".xml", "lesson" + strconv.Itoa(i) + ".xml")
	}
	wg.Wait()
}

func processFile(filename, outputFilename string) {
	lesson := Lesson{}

	xmlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("read file err: %v", err)
		return
	}

	err = xml.Unmarshal(xmlFile, &lesson)
	if err != nil {
		fmt.Printf("parse xml err: %v", err)
		return
	}

	var path, content string
	newXmlContent := &Lesson{}

	cd, _ := iconv.Open("UTF-8", "big5")
	defer cd.Close()

	for i := 0; i < len(lesson.Files); i++ {
		path = lesson.Files[i].Path
		
		if  ! isUsefulFileType(path) {
			continue
		}

		content = DecodeStr(lesson.Files[i].Content)
		if strings.Contains(path, ".cond") {
			content = cd.ConvString(content)
		}
		newXmlContent.Files = append(newXmlContent.Files, File{path, content})
	}
	newXmlOutput, _ := xml.MarshalIndent(newXmlContent, "", "  ")
	_ = ioutil.WriteFile("output/" + outputFilename, newXmlOutput, 0644)
	defer wg.Done()
}

func isUsefulFileType(path string) bool {
	if strings.Contains(path, ".html") || strings.Contains(path, ".cond") || strings.Contains(path, ".java") || strings.Contains(path, ".part") {
		return true
	}
	return false
}

func Decode(base64str string) []byte {
	result, _ := base64.StdEncoding.DecodeString(base64str)
	return result
}

func DecodeStr(base64str string) string {
	result := Decode(base64str)
	return string(result)
}
