package linksSearcher

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const minLengthURL = 11

const (
	HTTP       = "http://"
	SPACEHTTP  = " http://"
	HTTPS      = "https://"
	SPACEHTTPS = " https://"
)

var reURL = regexp.MustCompile(`http[s]?://[\w\[\]\-./?=&:#@!$'()*+,;_]*`)

func FindLinks(procNumber int, inputFilename, outputFilename string) error {

	in := input(inputFilename)
	out := output(outputFilename)

	wp := NewWP(procNumber)
	go wp.Run()

	done := make(chan struct{})
	resMap := make(map[int][]string)
	go func() {
		resMap = writeResultInMap(resMap, wp.Result())
		close(done)
	}()

	if err := addJobs(in, wp); err != nil {
		return err
	}

	<-done

	return write(out, resMap)
}

func addJobs(r io.Reader, wp *WorkerPool) error {
	reader := bufio.NewReader(r)

	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		wp.AddJob(Job{Id: i, Arg: string(line), ExecFn: findLinks})
	}
	wp.EndJob()
	return nil
}

func input(inputFilename string) *os.File {
	in, err := os.Open(inputFilename)
	if err != nil {
		in = os.Stdin
	}
	return in
}

func output(outputFilename string) *os.File {
	out, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		out = os.Stdout
	}
	return out
}

func writeResultInMap(m map[int][]string, res <-chan Result) map[int][]string {
	for r := range res {
		m[r.Id] = r.Value
	}
	return m
}

func write(w io.StringWriter, m map[int][]string) error {
	for i := 0; i < len(m); i++ {
		for _, link := range m[i] {
			if _, err := w.WriteString(link + "\n"); err != nil {
				return fmt.Errorf("coulndn't write to input file: %v", err)
			}
		}
	}
	return nil
}

func findLinks(line string) []string {

	links := splitLineIntoLinks(line)

	var result []string
	for _, link := range links {
		link, err := checkURL(link)
		if err != nil {
			continue
		}
		result = append(result, link)
	}
	return result
}

func splitLineIntoLinks(line string) []string {
	line = strings.ReplaceAll(line, HTTP, SPACEHTTP)
	line = strings.ReplaceAll(line, HTTPS, SPACEHTTPS)

	return reURL.FindAllString(line, -1)
}

func checkURL(link string) (string, error) {

	if err := checkBaseURL(link); err != nil {
		return "", err
	}

	n := utf8.RuneCountInString(link)

	for n >= minLengthURL {

		res, err := http.Get(link[:n])
		if err != nil || res.StatusCode > 399 && res.StatusCode < 500 {
			n--
			continue
		}

		return link[:n], nil
	}
	return "", errors.New("error: invalid link")
}

func checkBaseURL(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return err
	}
	baseURL := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
	}
	_, err = http.Get(baseURL.String())
	return err
}
