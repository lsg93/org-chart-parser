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

func setupTestAnalyser(chart model.OrganisationChart, t *testing.T) (*organisationChartAnalyser, *testWriter) {
	// Output needs to go to a writer to make assertions against.
	tw := &testWriter{}
	return NewOrganisationChartAnalyser(tw, exampleOrgChart), tw
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
		"handles duplicates": {
			input: model.OrganisationChart{
				model.Employee{Id: 1, Name: "Minion", ManagerId: 1},
				model.Employee{Id: 2, Name: "Minion", ManagerId: 1},
			},
			employee1:      "Minion",
			employee2:      "Minion",
			expectedOutput: "Minion (1) -> Minion (2)",
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			analyser, writer := setupTestAnalyser(tc.input, t)
			err := analyser.analyse(tc.employee1, tc.employee2)
			output := writer.contents

			if err != nil {
				t.Fatalf("There was an error analysing the given input.")
			}

			if output != tc.expectedOutput {
				t.Errorf("The received output %s was not equal to the expected output %s", output, tc.expectedOutput)
			}
		})
	}

}

// func TestAnalysisHandlesDuplicateNames() {
// 	// With contiguous whitespace
// 	// With non-contiguous whitespace
// 	// With different casing
// }

// func testAnalysisHandlesEmptyOrganisationChartGracefully(t *testing.T) {

// }

// func testAnalysisErrorsIfInputIsNonExistentNameisSupplied(t *testing.T) {

// }

// func testAnalysisHandlesNonExistentIds(t *testing.T) {

// }

// func TestAnalysisErrorsWhenCyclicalPathFound() (t *testing.T) {

// }
