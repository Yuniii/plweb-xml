package main

import (
	"io/ioutil"
	"strings"
	"fmt"
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

func main() {
	lesson := Lesson{}

	xmlFile, err := ioutil.ReadFile("xml/247/lesson1.xml")
	if err != nil {
		fmt.Printf("read file err: %v", err)
		return
	}

	err = xml.Unmarshal(xmlFile, &lesson)
	if err != nil {
		fmt.Printf("parse xml err: %v", err)
		return
	}

	var filename, content string
	newXmlContent := &Lesson{}

	cd, _ := iconv.Open("UTF-8", "big5")
	defer cd.Close()

	for i := 0; i < len(lesson.Files); i++ {
		filename = lesson.Files[i].Path
		
		if  ! isUsefulFileType(filename) {
			continue
		}

		content = DecodeStr(lesson.Files[i].Content)
		if strings.Contains(filename, ".cond") {
			content = cd.ConvString(content)
		}
		newXmlContent.Files = append(newXmlContent.Files, File{filename, content})
	}
	newXmlOutput, _ := xml.MarshalIndent(newXmlContent, "", "  ")
	_ = ioutil.WriteFile("output/1.xml", newXmlOutput, 0644)
}

func isUsefulFileType(filename string) bool {
	if strings.Contains(filename, ".html") || strings.Contains(filename, ".cond") || strings.Contains(filename, ".java") || strings.Contains(filename, ".part") {
		return true;
	}
	return false;
}

func Decode(base64str string) []byte {
	result, _ := base64.StdEncoding.DecodeString(base64str)
	return result
}

func DecodeStr(base64str string) string {
	result := Decode(base64str)
	return string(result)
}
