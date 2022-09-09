package linksSearcher

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"
)

const minLengthURL = 11

const (
	_http       = "http://"
	_spacehttp  = " http://"
	_https      = "https://"
	_spacehttps = " https://"
)

var reURL = regexp.MustCompile(`http[s]?://[\w\[\]\-./?=&:#@!$'()*+,;_]*`)

func FindLinks(procNumber int, inputFilename, outputFilename string) error {

	runtime.GOMAXPROCS(procNumber)

	in := input(inputFilename)
	out := output(outputFilename)

	chanLinks := make(chan string)

	go func() {
		if err := writeLinks(chanLinks, out); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}()

	reader := bufio.NewReader(in)
	wg := sync.WaitGroup{}

loop:
	for {
		for i := 0; i < procNumber; i++ {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break loop
			}
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(line string) {
				findLinks(line, chanLinks)
				wg.Done()
			}(string(line))
		}
		wg.Wait()
	}
	wg.Wait()
	close(chanLinks)

	return nil
}

func input(inputFilename string) *os.File {
	in := &os.File{}
	in, err := os.Open(inputFilename)
	if err != nil {
		in = os.Stdin
	}
	return in
}

func output(outputFilename string) *os.File {
	out := &os.File{}
	out, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		out = os.Stdout
	}
	return out
}

func writeLinks(links <-chan string, w io.StringWriter) error {
	for l := range links {
		if _, err := w.WriteString(l + "\n"); err != nil {
			return fmt.Errorf("coulndn't write to input file: %v", err)
		}
	}
	return nil
}

func findLinks(line string, chOut chan string) {

	links := splitLineIntoLinks(line)

	linkChan := make(chan string)

	for _, link := range links {
		go checkURL(link, linkChan)
	}

	for i := 0; i <= len(links); i++ {
		if link := <-linkChan; link != "" {
			chOut <- link
		}
	}
}

func splitLineIntoLinks(line string) []string {
	line = strings.ReplaceAll(line, _http, _spacehttp)
	line = strings.ReplaceAll(line, _https, _spacehttps)

	return reURL.FindAllString(line, -1)
}

func checkURL(link string, chanLinks chan<- string) {

	if err := checkBaseURL(link); err != nil {
		chanLinks <- ""
		return
	}

	n := utf8.RuneCountInString(link)

	for n >= minLengthURL {

		res, err := http.Get(link[:n])
		if err != nil || res.StatusCode > 399 {
			n--
			continue
		}

		chanLinks <- link[:n]
		break
	}
	chanLinks <- ""
}

func checkBaseURL(link string) error {
	u, err := url.Parse(link)
	if err == nil {
		baseURL := &url.URL{
			Scheme: u.Scheme,
			Host:   u.Hostname(),
		}
		_, err = http.Get(baseURL.String())
		return err
	}
	return err
}
