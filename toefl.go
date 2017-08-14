package main

import (
	"os"
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"strconv"
	"time"
	"unicode"

	"gopkg.in/iconv.v1"
	"gopkg.in/toast.v1"
)

func main() {
	var minutes time.Duration = 5
	
	if len(os.Args) > 1 {
		arg := os.Args[1]
		i, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("please input correct number for time interval")
			fmt.Println("use default time interval: 5 min")
		} else {
			minutes = time.Duration(i);
			fmt.Println("time interval: " + arg + " min")
		}
	} else {
		fmt.Println("use default time interval: 5 min")
	}
	for {
		word := random()
		fmt.Println(word)

		cd, err := iconv.Open("gbk", "utf-8")
		if err != nil {
			fmt.Println("iconv.Open failed!")
			return
		}
		defer cd.Close()

		meaning_gbk := cd.ConvString(word.meaningString())

		notification := toast.Notification{
    	    AppID: "TOEFL words",
    	    Title: word.spelling,
			Message: meaning_gbk,
			Audio: toast.Silent,
		}
		notification.Push()

		time.Sleep(minutes * 60000 * time.Millisecond)
	}
}

type word struct {
	spelling string
	meanings []string
}

func (w word) String() string {
	var buf bytes.Buffer

	buf.WriteString(w.spelling)
	buf.WriteRune('\n')
	for _, meaning := range w.meanings {
		buf.WriteRune('\n')
		buf.WriteString(meaning)
	}

	return buf.String()
}

func (w word) meaningString() string {
	var buf bytes.Buffer

	for _, meaning := range w.meanings {
		buf.WriteRune('\n')
		buf.WriteString(meaning)
	}

	return buf.String()
}

func random() word {
	words := load("words/wangyumei-toefl-words.txt")
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(words))
	return words[idx]
}

func load(path string) (words []word) {
	asset := MustAsset(path)
	for _, line := range strings.Split(string(asset), "\n") {
		line = strings.TrimFunc(line, unicode.IsSpace)
		if line == "" {
			continue
		}

		words = append(words, parse(line))
	}
	return
}

func parse(raw string) word {
	parts := strings.Split(raw, "#")
	spelling := parts[0]

	var meanings []string
	for _, meaning := range strings.Split(parts[1], ";") {
		meaning := strings.TrimFunc(meaning, unicode.IsSpace)
		meanings = append(meanings, meaning)
	}
	return word{spelling, meanings}
}
