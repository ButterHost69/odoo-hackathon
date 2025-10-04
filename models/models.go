package models

type User struct {
	Email        string
	Name         string
	Role         string
	ManagerEmail string
	ManagerName  string
	CompanyID    int
}
