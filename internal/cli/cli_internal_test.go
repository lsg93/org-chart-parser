package cli

import (
	"flag"
	"os"
	"testing"
)

func TestParsingValidCliArgumentsReturnsInputType(t *testing.T) {
	filepathArg := "path/to/file.txt"
	firstEmployeeNameArg := "Joshua"
	secondEmployeeNameArg := "Lawrence"

	mockArgs := []string{"test", filepathArg, firstEmployeeNameArg, secondEmployeeNameArg}

	expectedResult := OrgChartParserInput{
		filepath:           filepathArg,
		firstEmployeeName:  firstEmployeeNameArg,
		secondEmployeeName: secondEmployeeNameArg,
	}

	flag.CommandLine = flag.NewFlagSet(mockArgs[0], flag.ExitOnError)
	t.Cleanup(func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	})

	originalArgs := os.Args
	os.Args = mockArgs
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	result, err := parseArguments()

	if err != nil {
		t.Fatalf("An error '%s' was returned when none was expected", err)
	}

	if result != expectedResult {
		t.Errorf("The struct %v returned was not equal to the expected value %v", result, expectedResult)
	}
}

func TestParsingInvalidCliArgumentsErrors(t *testing.T) {
	type testCase struct {
		input           []string
		expectedMessage string
	}

	testCases := map[string]testCase{
		"When any/all arguments are empty": {
			input:           []string{"", " ", " "},
			expectedMessage: errArgValidationBlankArgumentProvided.Error(),
		},
		"When there are too few arguments": {
			input:           []string{"/path/to/file.txt", "Joshua"},
			expectedMessage: errArgValidationIncorrectArgumentAmount.Error(),
		},
		"When there are too many arguments": {
			input:           []string{"path/to/file.txt", "Joshua", "Adrian", "Lawrence"},
			expectedMessage: errArgValidationIncorrectArgumentAmount.Error(),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			err := validateArguments(tc.input)

			if err == nil {
				t.Fatalf("Validation for arguments passed, when it should have failed.")
			}

			if err.Error() != tc.expectedMessage {
				t.Errorf("An error '%s' was returned from validation, but it was not the expected error '%s'", err, tc.expectedMessage)
			}

		})
	}
}
