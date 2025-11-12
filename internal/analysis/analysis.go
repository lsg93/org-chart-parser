package analysis

import (
	"fmt"
	"io"
	"strings"

	"github.com/lsg93/org-chart-parser/internal/model"
)

type organisationChartAnalyser struct {
	chart  model.OrganisationChart
	output io.Writer
}

type OrganisationChartAnalysis struct{}

func NewOrganisationChartAnalyser(output io.Writer, chart model.OrganisationChart) *organisationChartAnalyser {
	return &organisationChartAnalyser{chart: chart, output: output}
}

// Breadth-first search to traverse graph.
// If we wanted to make this code as optimal as possi
func (a *organisationChartAnalyser) Analyse(name1 string, name2 string) error {

	// Need to validate here before any operations.

	adjList := a.mapEmployees()
	nameMap := a.mapEmployeeNames()

	startId := nameMap[name1]
	targetId := nameMap[name2]

	queue := []int{startId}

	// Use a map for quicker lookup of seenIds.
	seenIds := map[int]bool{startId: true}
	// This map shows the actual 'hops' between nodes.
	pathIds := make(map[int]int)

	for len(queue) > 0 {
		currentId := queue[0]
		queue = queue[1:] // Shift current item off start of queue.

		for _, relationId := range adjList[currentId] {
			if !seenIds[relationId] {
				seenIds[relationId] = true
				pathIds[relationId] = currentId
				queue = append(queue, relationId)

				if relationId == targetId {
					return a.constructPath(startId, targetId, pathIds)
				}
			}

		}

	}

	return a.constructPath(startId, targetId, pathIds)
}

func (a *organisationChartAnalyser) constructPath(startId int, targetId int, pathMap map[int]int) error {

	// Iterate through map starting from targetId to build path.
	path := make([]int, 0)

	currentId := targetId

	for currentId != startId {
		path = append([]int{currentId}, path...)
		currentId = pathMap[currentId]
	}

	path = append([]int{startId}, path...)

	// Now that we have the shortest path, we have to analyse the direction of the data flow.
	// We iterate through the path, and determine whether the direction of travel is up/down based on whether the next item in the slice has a manager.
	// I think you could use a bi-directional BFS for this in future for better performance maybe if it was critical.

	idMap := a.mapEmployeeIds()
	managerMap := a.mapManagement()

	var stringsPath strings.Builder
	flowDirection := ""

	for i := 0; i < len(path)-1; i++ {
		current := path[i]
		next := path[i+1]

		if managerMap[current] == next {
			flowDirection = "->"
		} else if managerMap[next] == current {
			flowDirection = "<-"
		}

		stringsPath.WriteString(fmt.Sprintf("%s %s ", idMap[current], flowDirection))
	}

	// Manually append last item in path, since we won't iterate over it.
	stringsPath.WriteString(idMap[path[len(path)-1]])

	fmt.Println(stringsPath.String())

	_, err := a.output.Write([]byte(stringsPath.String()))

	if err != nil {
		return nil
	}

	return nil
}

// Build adjacency list structure - this is what we'll iterate over with the breadth-first-search.
// We want Employee > [List of employees related to them]
func (a *organisationChartAnalyser) mapEmployees() map[int][]int {
	adjList := make(map[int][]int)

	for _, employee := range a.chart {
		adjList[employee.Id] = make([]int, 0)

		if employee.ManagerId != 0 {
			adjList[employee.Id] = append(adjList[employee.Id], employee.ManagerId)
			adjList[employee.ManagerId] = append(adjList[employee.ManagerId], employee.Id)
		}

	}

	return adjList
}

// Having a map of each employee and their direct report helps us determine the direction of the data flow.
func (a *organisationChartAnalyser) mapManagement() map[int]int {
	managerMap := make(map[int]int)

	for _, employee := range a.chart {
		managerMap[employee.Id] = employee.ManagerId
	}

	return managerMap
}

// Create easy lookups to translate ID's to names.
func (a *organisationChartAnalyser) mapEmployeeIds() map[int]string {
	nameMap := make(map[int]string)
	for _, employee := range a.chart {
		nameMap[employee.Id] = fmt.Sprintf("%s (%d)", employee.Name, employee.Id)
	}
	return nameMap
}

// Create easy lookups to translate names to ID's.
func (a *organisationChartAnalyser) mapEmployeeNames() map[string]int {
	nameMap := make(map[string]int)
	for _, employee := range a.chart {
		nameMap[employee.Name] = employee.Id
	}
	return nameMap
}
