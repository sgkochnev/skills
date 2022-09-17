package linksSearcher

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	input             string
	inputFilename     string
	outputFilename    string
	result            string
	numberOfProcesses int
}{
	{
		input:             "https://www.digitalocean.com/community/tutorials*8;9joasd\nasasfkpml;mlhttps://ru.wikipedia.org lpmhttps://stackoverflow.com/questions|l;q\nhttps://www.digitalocean.com/community/tutorials*8;",
		inputFilename:     "link_searcher_test1_input.txt",
		outputFilename:    "link_searcher_test1_output.txt",
		result:            "https://www.digitalocean.com/community/tutorials\nhttps://ru.wikipedia.org\nhttps://stackoverflow.com/questions\nhttps://www.digitalocean.com/community/tutorials\n",
		numberOfProcesses: 5,
	},
	{
		input:             "asahttp://sfkp0i_iwlhttps://ru.wikipedia.org lpmhttps://nonexistentlink.com/forum|l;q\nhttps://www.digitalocean.com/community/tutorials{!imklsdf",
		inputFilename:     "link_searcher_test2_input.txt",
		outputFilename:    "link_searcher_test2_output.txt",
		result:            "https://ru.wikipedia.org\nhttps://www.digitalocean.com/community/tutorials\n",
		numberOfProcesses: 2,
	},
}

var testCases2 = []struct {
	input             string
	inputFilename     string
	outputFilename    string
	result            []string
	numberOfProcesses int
}{
	{
		input:             strings.Repeat("https://www.digitalocean.com/community/tutorials*8;9joasd\nasasfkpml;mlhttps://ru.wikipedia.org lpmhttps://stackoverflow.com/questions|l;q\nhttps://www.digitalocean.com/community/tutorials*8;", 12),
		inputFilename:     "link_searcher_test1_input.txt",
		outputFilename:    "link_searcher_test1_output.txt",
		result:            strings.Split(strings.TrimSuffix(strings.Repeat("https://www.digitalocean.com/community/tutorials\nhttps://ru.wikipedia.org\nhttps://stackoverflow.com/questions\nhttps://www.digitalocean.com/community/tutorials\n", 12), "\n"), "\n"),
		numberOfProcesses: 5,
	},
	{
		input:          "asahttp://sfkp0i_iwlhttps://ru.wikipedia.org lpmhttps://nonexistentlink.com/forum|l;q\nhttps://www.digitalocean.com/community/tutorials{!imklsdf",
		inputFilename:  "link_searcher_test2_input.txt",
		outputFilename: "link_searcher_test2_output.txt",
		result: []string{
			"https://ru.wikipedia.org",
			"https://www.digitalocean.com/community/tutorials",
		},
		numberOfProcesses: 2,
	},
}

// Functional test
func TestLinksSearcher(t *testing.T) {
	for _, testCase := range testCases {
		inputFile, err := createInputFile(testCase.inputFilename, testCase.input)
		if err != nil {
			t.Errorf("Error during test setup: %v", err)
			return
		}
		err = FindLinks(testCase.numberOfProcesses, inputFile.Name(), testCase.outputFilename)

		if err != nil {
			t.Errorf("Error occured: %v", err)
			return
		}

		var outputString string
		outputString, err = readOutputFile(testCase.outputFilename)
		if err != nil {
			t.Errorf("Error during test setup: %v", err)
			return
		}

		if outputString != testCase.result {
			t.Errorf("Result doesn't match expected. Got: %s.Expected: %s.", outputString, testCase.result)
			return
		}
	}
}

func TestLinksSearcher2(t *testing.T) {
	for _, testCase := range testCases2 {
		inputFile, err := createInputFile(testCase.inputFilename, testCase.input)
		if err != nil {
			t.Errorf("Error during test setup: %v", err)
			return
		}
		err = FindLinks(testCase.numberOfProcesses, inputFile.Name(), testCase.outputFilename)

		if err != nil {
			t.Errorf("Error occured: %v", err)
			return
		}

		var outputString []string
		outputString, err = readOutputFile2(testCase.outputFilename)
		if err != nil {
			t.Errorf("Error during test setup: %v", err)
			return
		}
		assert.ElementsMatch(t, outputString, testCase.result, "Result doesn't match expected. Got: %s.Expected: %s.", outputString, testCase.result)

		err = deleteFiles(testCase.inputFilename, testCase.outputFilename)
		if err != nil {
			t.Errorf("Error occured: %v", err)
			return
		}
	}
}
