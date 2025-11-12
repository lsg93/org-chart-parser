package analysis

import (
	"io"

	"github.com/lsg93/org-chart-parser/internal/model"
)

type organisationChartAnalyser struct {
	chart  model.OrganisationChart
	output io.Writer
}

type OrganisationChartAnalysis struct{}

func NewOrganisationChartAnalyser(output io.Writer, chart model.OrganisationChart) *organisationChartAnalyser {
	return &organisationChartAnalyser{}
}

func (a *organisationChartAnalyser) analyse(name1 string, name2 string) error {
	return nil
}
