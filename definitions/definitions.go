package definitions

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WordDB interface {
	AddWord(s string)
}

func LoadWords(filename string, db WordDB) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		db.AddWord(scanner.Text())
	}

	return scanner.Err()
}

func defineWord(word string) (string, error) {
	form := url.Values{
		"dictWord": {word},
	}
	resp, err := http.PostForm("https://scrabble.hasbro.com/en-us/tools", form)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	def := doc.Find(".word-definition")
	definition := def.Text()
	matcher := regexp.MustCompile("(?is)" + word + "(.*)")
	core := matcher.FindStringSubmatch(definition)
	if core == nil {
		return "", errors.New("No definition")
	}
	return strings.Join(strings.Fields(core[1]), " "), nil
}

// func checkWords(filename string) {
// 	words, err := LoadWords(filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	start := 0

// 	output, err := os.OpenFile("fast-words.csv", os.O_RDWR|os.O_APPEND, 0660)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer output.Close()

// 	type Result struct {
// 		Word, Definition, Error string
// 	}

// 	jobs := make(chan string, 100)
// 	results := make(chan Result, 100)
// 	for i := 0; i < 100; i++ {
// 		go func() {
// 			for j := range jobs {
// 				definition, err := defineWord(j)
// 				errString := ""
// 				if err != nil {
// 					errString = err.Error()
// 				}
// 				results <- Result{
// 					Word:       j,
// 					Error:      errString,
// 					Definition: definition,
// 				}
// 			}
// 		}()
// 	}

// 	go func() {
// 		for _, word := range words[start:] {
// 			jobs <- word
// 		}
// 		close(jobs)
// 	}()

// 	wr := csv.NewWriter(output)
// 	for res := range results {
// 		err = wr.Write([]string{
// 			res.Word,
// 			res.Definition,
// 			res.Error,
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(res.Word, err == nil)
// 		wr.Flush()
// 	}
// 	wr.Flush()

// }
