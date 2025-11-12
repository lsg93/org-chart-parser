package parser

import (
	"bufio"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/lsg93/org-chart-parser/internal/model"
)

// Not going to be used now, but if doing this properly, you'd probably want to not couple too tightly to files.
// Probably better to use an interface, and change implementations as necessary.
// In a real project, to aid OCP it would probably be better to split the parsing logic out from this contract into a separate file - e.g. file_parser.go
type OrganisationChartParser interface {
	Parse() (model.OrganisationChart, error)
}

var (
	errParserScanError         = errors.New("The input data could not be scanned line by line.")
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

func (parser *orgChartFileParser) Parse() (model.OrganisationChart, error) {
	chart := model.OrganisationChart{}

	scanner := bufio.NewScanner(parser.input)

	i := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines - this accommodates leading whitespace.
		if len(line) == 0 {
			continue
		}

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

	// If no scanner iterations took place, then the input was likely empty.
	if i == 0 {
		return chart, errParserEmptyInput
	}

	if err := scanner.Err(); err != nil {
		return chart, errParserScanError
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
	name := s[1]
	managerId := s[2]

	// Edge case for empty rows
	if employeeId == "" && name == "" && managerId == "" {
		return nil, nil
	}

	if employeeId == managerId || employeeId == "" {
		return nil, errParserInvalidIdField
	}

	// Check employee ID is numeric.
	if _, err := strconv.Atoi(s[0]); err != nil {
		return nil, errParserInvalidIdField
	}

	// Check manager ID is numeric.
	if _, err := strconv.Atoi(s[2]); err != nil {
		// Only error if the error occurs when the manager ID is not blank.
		if s[2] != "" {
			return nil, errParserInvalidIdField
		}
	}

	return s, nil
}

func (parser *orgChartFileParser) marshalLine(s []string) model.Employee {
	// Fairly confident the errors can be ignored, as input should have been validated @ this point.
	// This could be better though I think.
	employeeId, _ := strconv.Atoi(s[0])
	name := s[1]
	managerId, _ := strconv.Atoi(s[2])

	employee := model.Employee{
		Id:        employeeId,
		Name:      name,
		ManagerId: managerId,
	}

	return employee
}

// Might need to move these helpers later down the line.

func lowercaseSlice(s []string) []string {
	ls := []string{}

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
