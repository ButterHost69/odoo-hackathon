package db

import (
	"github.com/lib/pq"
)

type ManagerInfo struct {
	ManagerEmail string
	ManagerName  string
}

type ManagerInfoSlice []ManagerInfo

// Scan implements the sql.Scanner interface for ApproverInfoSlice.
// This tells the Go SQL driver how to convert the array from the database
// into our custom Go slice type.
func (m ManagerInfo) Scan(src interface{}) error {
	return pq.Array(m).Scan(src)
}

type User struct {
	Email        string
	Name         string
	Role         string
	ManagerEmail string
	ManagerName  string
	CompanyID    int
}

type ApproverInfo struct {
	ApproverEmail    string
	ApprovalRequired bool
}

type ApproverInfoSlice []ApproverInfo

// Scan implements the sql.Scanner interface for ApproverInfoSlice.
// This tells the Go SQL driver how to convert the array from the database
// into our custom Go slice type.
func (a *ApproverInfoSlice) Scan(src interface{}) error {
	return pq.Array(a).Scan(src)
}

type Rules struct {
	EmployeeEmail        string
	IsManagerApprover    bool
	MinApprovalPercent   int
	IsApprovalSequential bool
	Approvers            ApproverInfoSlice
}
