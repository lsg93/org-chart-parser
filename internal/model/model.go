package model

type Employee struct {
	Id        int
	Name      string
	ManagerId int
}

type OrganisationChart = []Employee
