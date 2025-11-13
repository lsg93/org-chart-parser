package analysis

import (
	"testing"

	"github.com/lsg93/org-chart-parser/internal/model"
)

var exampleOrgChart = model.OrganisationChart{
	model.Employee{Id: 1, Name: "Dangermouse"},
	model.Employee{Id: 2, Name: "Gonzo the Great", ManagerId: 1},
	model.Employee{Id: 3, Name: "Invisible Woman", ManagerId: 1},
	model.Employee{Id: 6, Name: "Black Widow", ManagerId: 2},
	model.Employee{Id: 12, Name: "Hit Girl", ManagerId: 3},
	model.Employee{Id: 15, Name: "Super Ted", ManagerId: 3},
	model.Employee{Id: 16, Name: "Batman", ManagerId: 6},
	model.Employee{Id: 17, Name: "Catwoman", ManagerId: 6},
}

func setupTestAnalyser(chart model.OrganisationChart) (*organisationChartAnalyser, *testWriter) {
	// Output needs to go to a writer to make assertions against.
	tw := &testWriter{}
	return NewOrganisationChartAnalyser(tw, chart), tw
}

type testWriter struct {
	contents string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.contents = tw.contents + string(p)
	return len(p), nil
}

func TestAnalysisReturnsShortestPathAsString(t *testing.T) {

	type testCase struct {
		input          model.OrganisationChart
		employee1      string
		employee2      string
		expectedOutput string
	}

	testCases := map[string]testCase{
		"first documented example": {
			input:          exampleOrgChart,
			employee1:      "Batman",
			employee2:      "Super Ted",
			expectedOutput: "Batman (16) -> Black Widow (6) -> Gonzo the Great (2) -> Dangermouse (1) <- Invisible Woman (3) <- Super Ted (15)",
		},
		"second documented example": {
			input:          exampleOrgChart,
			employee1:      "Batman",
			employee2:      "Catwoman",
			expectedOutput: "Batman (16) -> Black Widow (6) <- Catwoman (17)",
		},
		"handles duplicates gracefully": {
			input: model.OrganisationChart{
				model.Employee{Id: 1, Name: "CEO", ManagerId: 0},
				model.Employee{Id: 2, Name: "Boss", ManagerId: 1},
				model.Employee{Id: 10, Name: "Minion", ManagerId: 1}, // This should be the one referenced.
				model.Employee{Id: 3, Name: "Boss", ManagerId: 0},
				model.Employee{Id: 20, Name: "Minion", ManagerId: 2}, // This path would have 2 jumps.
			},
			employee1:      "Minion",
			employee2:      "CEO",
			expectedOutput: "Minion (10) -> CEO (1)",
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			analyser, writer := setupTestAnalyser(tc.input)
			err := analyser.Analyse(tc.employee1, tc.employee2)
			output := writer.contents

			if err != nil {
				t.Fatalf("There was an error '%s' analysing the given input.", err.Error())
			}

			if output != tc.expectedOutput {
				t.Errorf("The received output '%s' was not equal to the expected output '%s'", output, tc.expectedOutput)
			}
		})
	}

}

func TestAnalysisErrorsWithInvalidArguments(t *testing.T) {
	type testCase struct {
		input         model.OrganisationChart
		employee1     string
		employee2     string
		expectedError error
	}

	testCases := map[string]testCase{
		"Non-existent name as argument": {
			input:         exampleOrgChart,
			employee1:     "Superman",
			employee2:     "Batman",
			expectedError: errAnalysisInvalidNameArgument,
		},
		"Two duplicate names as arguments": {
			input:         exampleOrgChart,
			employee1:     "Batman",
			employee2:     "Batman",
			expectedError: errAnalysisDuplicateNameArgument,
		},
		"When an invalid path is attempted": {
			input: model.OrganisationChart{
				model.Employee{Id: 1, Name: "CEO", ManagerId: 0},
				model.Employee{Id: 2, Name: "VP", ManagerId: 1},
				model.Employee{Id: 3, Name: "CTO", ManagerId: 0},
				model.Employee{Id: 4, Name: "SWE", ManagerId: 3},
			},
			employee1:     "SWE",
			employee2:     "VP",
			expectedError: errAnalysisNoPathsFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			analyser, writer := setupTestAnalyser(tc.input)
			err := analyser.Analyse(tc.employee1, tc.employee2)
			output := writer.contents

			if err == nil {
				t.Fatalf("There was no error analysing the given input when one should have occurred.")
			}

			if err != tc.expectedError {
				t.Errorf("The received error '%s' was not equal to the expected output '%s'", output, tc.expectedError)
			}
		})
	}

}
