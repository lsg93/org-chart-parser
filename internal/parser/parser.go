package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Employee struct {
	id        int
	name      string
	managerId int
}

type OrganisationChart = []Employee

// Not going to be used now, but if doing this properly, you'd probably want to not couple too tightly to files.
// Probably better to use an interface, and change implementations as necessary.
type OrganisationChartParser interface {
	Parse() (OrganisationChart, error)
	validateHeader(string) bool
	validateLine(string) ([]string, error) // Returned errors can be different.
	marshalLine([]string) Employee
}

var (
	errParserEmptyInput        = errors.New("Provided input to the parser was empty.")
	errParserInvalidHeader     = errors.New("No header with appropriate column names was found in given input.")
	errParserInvalidIdField    = errors.New("A problem was encountered when parsing the ID field - Check that your input has correct ID fields.")
	errParserInvalidLineLength = errors.New("One of the lines in the input has too many, or too few fields.")
)

type orgChartFileParser struct {
	input io.Reader
}

func NewOrganisationChartParser(input io.Reader) (OrganisationChartParser, error) {
	return &orgChartFileParser{input: input}, nil

}

func (parser *orgChartFileParser) Parse() (OrganisationChart, error) {
	chart := OrganisationChart{}

	scanner := bufio.NewScanner(parser.input)

	i := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if i == 0 {
			// Split header and check it has correct column names
			// If not, return error, as input is malformed.
			if !parser.validateHeader(line) {
				return chart, errParserInvalidHeader
			}
			i++
			continue
		}

		// Stopping on failure is better for something without a UI I think.
		validated, err := parser.validateLine(line)

		if err != nil {
			return chart, err
		}

		if validated == nil {
			// Empty row - continue on.
			continue
		}

		// marshal line into struct.
		employee := parser.marshalLine(validated)
		chart = append(chart, employee)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return chart, nil

}

func (parser *orgChartFileParser) validateHeader(headerLine string) bool {
	headerNames := []string{"employee id", "name", "manager id"}
	colNames := normaliseLineSlice(strings.Split(headerLine, "|"))

	if len(colNames) != 3 {
		return false
	}

	return slices.Equal(headerNames, lowercaseSlice(colNames))
}

func (parser *orgChartFileParser) validateLine(line string) ([]string, error) {

	s := normaliseLineSlice(strings.Split(line, "|"))

	if len(s) != 3 {
		return nil, errParserInvalidLineLength
	}

	employeeId := s[0]
	name := s[0]
	managerId := s[2]

	// Edge case for empty rows
	if employeeId == "" && name == "" && managerId == "" {
		return nil, nil
	}

	if employeeId == managerId || employeeId == "" {
		return nil, errParserInvalidIdField
	}

	return s, nil
}

func (parser *orgChartFileParser) marshalLine(s []string) Employee {
	// Fairly confident the errors can be ignored, as they'd be zero values in the struct if invalid.
	employeeId, _ := strconv.Atoi(s[0])
	name := s[1]
	managerId, _ := strconv.Atoi(s[2])

	employee := Employee{
		id:        employeeId,
		name:      name,
		managerId: managerId,
	}

	return employee
}

// Might need to move these helpers later down the line.

func lowercaseSlice(s []string) []string {
	ls := []string{}

	fmt.Println(strings.Join(s, ","))

	for i := range s {
		ls = append(ls, strings.ToLower(s[i]))
	}

	return ls
}

func normaliseLineSlice(s []string) []string {
	ts := []string{}

	for _, v := range s[1 : len(s)-1] {
		ts = append(ts, strings.TrimSpace(v))
	}

	return ts
}
