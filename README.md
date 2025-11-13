This is my attempt at using a breadth first search to find the shortest path between data in an adjacency list in Go via the command line.

# Example usage:

The application requires exactly three arguments: [filepath] [start name] [target name]

You can clone this repo, and in your terminal run the command with your desired arguments, for example:
- `go run main.go example.txt Catwoman "Invisible Woman"`

Alternatively, you can clone the repo, build the binary, and run it in a similar fashion:
- `go build -o org-chart-parser main.go`
- `./org-chart-parser [filepath] "Employee A" "Employee B"`

Tests can be run in the root of the repo with the command `go test ./... -v`

# Output

The result with arrows indicating the direction of management flow:
- `Employee (ID) -> Manager (ID) <- Employee (ID)`
