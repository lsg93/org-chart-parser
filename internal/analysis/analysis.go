package analysis

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/lsg93/org-chart-parser/internal/model"
)

var (
	errAnalysisNoPathsFound          = errors.New("No suitable paths between the given employees could be found.")
	errAnalysisInvalidNameArgument   = errors.New("One, or both of the names provided as arguments do not exist in the organisation chart.")
	errAnalysisDuplicateNameArgument = errors.New("Tthe names provided as arguments are duplicate - there is no way to determine the path you would like to see.")
)

type organisationChartAnalyser struct {
	chart   model.OrganisationChart
	output  io.Writer
	adjList map[int][]int    // graph structure for BFS traversal
	nameMap map[string][]int // used to look up names when building string from path ID's
}

type OrganisationChartAnalysis struct{}

func NewOrganisationChartAnalyser(output io.Writer, chart model.OrganisationChart) *organisationChartAnalyser {
	analyser := &organisationChartAnalyser{
		chart:  chart,
		output: output,
	}

	analyser.adjList = analyser.mapEmployees()
	analyser.nameMap = analyser.mapEmployeeNames()

	return analyser
}

// Breadth-first search to traverse graph.
// If we wanted to make this code as optimal as possi
func (a *organisationChartAnalyser) Analyse(name1 string, name2 string) error {

	// Validate that the names actually exist
	err := a.validateNames(name1, name2)

	if err != nil {
		return err
	}

	/*
		There is no guarantee names are unique
		As such, the only thing we can really do to calculate the shortest path
		Is generate all possible paths for each instance of a duplicated name
		And use the shortest one in our final output
	*/

	startIds := a.nameMap[name1]
	targetIds := a.nameMap[name2]

	// Store all the paths so we can then sort them to find and return the shortest one
	allPaths := make([][]int, 0)

	for _, startId := range startIds {
		for _, targetId := range targetIds {
			pathIds := a.search(startId, targetId)
			path, err := a.constructPath(startId, targetId, pathIds)

			if err != nil {
				return errAnalysisNoPathsFound
			}

			allPaths = append(allPaths, path)
		}
	}

	// Sort all of our calculated paths, push the shortest one to the front.
	sort.Slice(allPaths, func(i int, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	// The shortest path gets turned into a string.
	shortestPath, err := a.pathToString(allPaths[0])

	if err != nil {
		return err
	}

	_, err = a.output.Write([]byte(shortestPath.String()))

	if err != nil {
		return err
	}

	return nil
}

// BFS algorithm.
func (a *organisationChartAnalyser) search(startId int, targetId int) map[int]int {
	queue := []int{startId}

	// Use a map for quicker lookup of seenIds.
	seenIds := map[int]bool{startId: true}
	// This map shows the actual 'hops' between nodes.
	pathIds := make(map[int]int)

	for len(queue) > 0 {
		currentId := queue[0]
		queue = queue[1:] // Shift current item off start of queue.

		for _, relationId := range a.adjList[currentId] {
			// If the next node hasn't been seen, then add it to the queue.
			// All nodes at a particular depth get added to the queue and get 'seen'.
			if !seenIds[relationId] {
				seenIds[relationId] = true
				pathIds[relationId] = currentId
				queue = append(queue, relationId)

				if relationId == targetId {
					return pathIds
				}
			}

		}

	}

	return pathIds
}

func (a *organisationChartAnalyser) constructPath(startId int, targetId int, pathMap map[int]int) ([]int, error) {

	// Iterate through map starting from targetId to build path.
	path := make([]int, 0)

	currentId := targetId

	for currentId != startId {

		prev, ok := pathMap[currentId]

		if !ok {
			// Can't find the id in the map - path is invalid
			return []int{}, errAnalysisNoPathsFound
		}

		path = append([]int{currentId}, path...)

		currentId = prev
	}

	path = append([]int{startId}, path...)

	return path, nil
}

// // In order to achieve the desired output, we have to analyse the direction of the data flow.
// // We iterate through the path, and determine whether the direction of travel is up/down based on whether the next item in the slice has a manager.
// // I think you could use a bi-directional BFS for this in future for better performance maybe if it was critical.
func (a *organisationChartAnalyser) pathToString(path []int) (strings.Builder, error) {
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

		_, err := stringsPath.WriteString(fmt.Sprintf("%s %s ", idMap[current], flowDirection))

		if err != nil {
			return stringsPath, err
		}
	}

	// Manually append last item in path, since we won't iterate over it.
	_, err := stringsPath.WriteString(idMap[path[len(path)-1]])

	if err != nil {
		return stringsPath, err
	}

	return stringsPath, nil
}

func (a *organisationChartAnalyser) validateNames(name1 string, name2 string) error {
	// Making an assumption here - I think working on duplicate name inputs is quite messy.
	if name1 == name2 {
		return errAnalysisDuplicateNameArgument
	}

	_, startExists := a.nameMap[name1]
	_, endExists := a.nameMap[name2]

	if !startExists || !endExists {
		return errAnalysisInvalidNameArgument
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
	idMap := make(map[int]string)
	for _, employee := range a.chart {
		idMap[employee.Id] = fmt.Sprintf("%s (%d)", employee.Name, employee.Id)
	}
	return idMap
}

// Create easy lookups to translate names to ID's - we need to have slices instead of a hashmap
// Because names might not be unique.
func (a *organisationChartAnalyser) mapEmployeeNames() map[string][]int {
	nameMap := make(map[string][]int)
	for _, employee := range a.chart {
		nameMap[employee.Name] = append(nameMap[employee.Name], employee.Id)
	}
	return nameMap
}
