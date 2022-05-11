package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// OpenFile open a file and return its handle.
func OpenFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file error: %s", err)
	}
	return file
}

// ReadLine read a line of data from a file, unaffected by the buffer.
func ReadLine(b *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		var lineEnd bool
		lineTemp, lineEnd, err := b.ReadLine()
		if err != nil {
			return nil, err
		}
		line = append(line, lineTemp...)
		if !lineEnd {
			break
		}
	}
	return line, nil
}

// GetWordBlock used to get a complete word block from a file.
func GetWordBlock(b *bufio.Reader) ([]byte, error) {
	//content := []byte{60, 47, 62, 10}
	content := []byte{10, 60, 47, 62}
	for {
		bytes, err := ReadLine(b)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(string(bytes)) == "</>" {
			if content == nil {
				continue
			}
			break
		}
		content = append(content, 10)       // 追加 \n 换行符
		content = append(content, bytes...) // 追加内容
	}
	content = append(content, 10, 60, 47, 62) // 追加 \n 换行符
	//content = append(content, 10) // 追加 \n 换行符
	return content, nil
}
