package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lsg93/org-chart-parser/internal/analysis"
	"github.com/lsg93/org-chart-parser/internal/parser"
)

type OrgChartParserInput struct {
	filepath           string
	firstEmployeeName  string
	secondEmployeeName string
}

var (
	errArgValidationBlankArgumentProvided   = errors.New("One, or many of the arguments provided are blank.")
	errArgValidationIncorrectArgumentAmount = errors.New("One or more of the expected arguments (filepath, start name, target name) have not been provided.")
	errCouldNotReadFile                     = errors.New("There was an error reading the file.")
)

func Run() {
	input, err := parseArguments()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err := readFile(input.filepath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser, err := parser.NewOrganisationChartParser(bytes.NewReader(data))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	chart, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	analyser := analysis.NewOrganisationChartAnalyser(os.Stdin, chart)
	err = analyser.Analyse(input.firstEmployeeName, input.secondEmployeeName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func parseArguments() (OrgChartParserInput, error) {
	flag.Parse()
	args := flag.Args()

	err := validateArguments(args)

	if err != nil {
		return OrgChartParserInput{}, err
	}

	res := OrgChartParserInput{
		filepath:           args[0],
		firstEmployeeName:  args[1],
		secondEmployeeName: args[2],
	}

	return res, nil
}

func validateArguments(args []string) error {
	if len(args) != 3 {
		return errArgValidationIncorrectArgumentAmount
	}

	for _, item := range args {
		if strings.TrimSpace(item) == "" {
			return errArgValidationBlankArgumentProvided
		}
	}

	return nil
}

func readFile(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errCouldNotReadFile
	}
	return bytes, nil
}
