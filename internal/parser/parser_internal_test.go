package parser

import (
	"slices"
	"strings"
	"testing"
)

func setupParser(input string, t *testing.T) OrganisationChartParser {
	parser, err := NewOrganisationChartParser(strings.NewReader(input))

	if err != nil {
		t.Fatalf("An error occurred initialising the parser with the given data.")
	}

	return parser
}

func TestParsesOrgChartTextSuccesfully(t *testing.T) {

	type testCase struct {
		input          string
		expectedResult OrganisationChart
	}

	testCases := map[string]testCase{
		"with example data": {
			input: `|Employee ID|Name|Manager ID|
			|1|Lawrence||
			|2|Adrian|1|
			|3|Joshua|2|`,
			expectedResult: OrganisationChart{
				Employee{id: 1, name: "Lawrence", managerId: 0},
				Employee{id: 2, name: "Adrian", managerId: 1},
				Employee{id: 3, name: "Joshua", managerId: 2},
			},
		},
		"with example data (whitespace)": {
			input: `| Employee ID | Name | Manager ID |
			| 1 | Lawrence | |
			| 2 | Adrian | 1 |
			| 3 | Joshua | 2 |`,
			expectedResult: OrganisationChart{
				Employee{id: 1, name: "Lawrence", managerId: 0},
				Employee{id: 2, name: "Adrian", managerId: 1},
				Employee{id: 3, name: "Joshua", managerId: 2},
			},
		},
		"with missing rows": {
			input: `| Employee ID | Name | Manager ID |
			|1 | Lawrence | |
			|  |  |  |
			|3|Joshua|2|`,
			expectedResult: OrganisationChart{
				Employee{id: 1, name: "Lawrence", managerId: 0},
				Employee{id: 3, name: "Joshua", managerId: 2},
			},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {

			parser := setupParser(tc.input, t)
			result, err := parser.Parse()

			if err != nil {
				t.Fatalf("There was an error %v parsing the provided the input data.", err)
			}

			if slices.Equal(result, tc.expectedResult) == false {
				t.Errorf("The result %v was not the same as the expected result %v", result, tc.expectedResult)
			}
		})
	}
}

// func TestFailsToParseWhenOrgChartTextIsInvalid(t *testing.T) {
// 	type testCase struct {
// 		input         string
// 		expectedError error
// 	}

// 	testCases := map[string]testCase{
// 		"with empty input": {
// 			input:         "",
// 			expectedError: errParserEmptyInput,
// 		},
// with extra fields
// with too few fields
// 		"with missing header fields": {
// 			input: `
// 			| 1 | Lawrence | |
// 			| 2 | Adrian | 1 |
// 			`,
// 			expectedError: errParserInvalidHeader,
// 		},
// 		"with missing ID field": {
// 			input: `
// 			| Employee ID | Name | Manager ID |
// 			| 1 | Lawrence | |
// 			|  | Adrian | 1 |
// 			`,
// 			expectedError: errParserInvalidIdField,
// 		},
// 		"with non integer ids": {
// 			input: `
// 			| Employee ID | Name | Manager ID |
// 			| A | Lawrence | |
// 			`,
// 			expectedError: errParserInvalidIdField,
// 		},
// 		"with self referential data": {
// 			input: `
// 			| Employee ID | Name | Manager ID |
// 			| 1 | Lawrence | 1 |
// 			`,
// 			expectedError: errParserInvalidIdField,
// 		},
// 	}

// 	for desc, tc := range testCases {
// 		t.Run(desc, func(t *testing.T) {

// 			parser := setupParser(tc.input, t)
// 			result, err := parser.Parse()

// 			if err == nil {
// 				t.Fatalf("An error should have occurred while attempting to parse the data, but none did")
// 			}

// 			if err == tc.expectedError {
// 				t.Errorf("The expected error '%v' was not the same as the expected result '%v'", result, tc.expectedError)
// 			}
// 		})
// 	}
// }
