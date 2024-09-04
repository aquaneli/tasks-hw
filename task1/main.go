package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

var dirName string = "files/"

func main() {
	dir := ParseDir()

	ch := make(chan *int)
	linesCount := 0
	go func() {
		ch <- &linesCount
	}()

	m := Processing(ch, dir)
	time.Sleep(time.Second * 2)

	fmt.Println(*<-ch)
	for _, v := range *m {
		for i := 0; i < len(v); i++ {
			fmt.Println(i+1, v[i])
		}
	}

}

func ParseDir() []fs.DirEntry {
	dir, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func Processing(ch chan *int, dir []fs.DirEntry) *[]map[int]string {
	m := make([]map[int]string, len(dir))
	for i, v := range dir {
		info, err := v.Info()
		if err != nil {
			log.Fatal(err)
		}
		m[i] = make(map[int]string)
		go ParseFile(ch, &m[i], info)
	}
	return &m
}

func ParseFile(ch chan *int, m *map[int]string, info os.FileInfo) {
	file, err := os.Open(dirName + info.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linesCount := <-ch
	for i := 0; scanner.Scan(); i++ {
		(*m)[i] = scanner.Text()
	}
	*linesCount += len(*m)

	ch <- linesCount
}
