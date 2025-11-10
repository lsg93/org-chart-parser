package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type OrgChartParserInput struct {
	filepath           string
	firstEmployeeName  string
	secondEmployeeName string
}

var errArgValidationBlankArgumentProvided = errors.New("One, or many of the arguments provided are blank.")
var errArgValidationIncorrectArgumentAmount = errors.New("One or more of the expected arguments have not been provided.")

func Run() {
	_, err := parseArguments()

	if err != nil {
		fmt.Println(err)
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
	for _, item := range args {
		if strings.TrimSpace(item) == "" {
			return errArgValidationBlankArgumentProvided
		}
	}

	if len(args) != 3 {
		return errArgValidationIncorrectArgumentAmount
	}

	return nil
}
