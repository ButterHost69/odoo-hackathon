package db

import (
	"database/sql/driver"
	"fmt"
	"time"

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
// Used to Read Values
func (m ManagerInfo) Scan(src interface{}) error {
	return pq.Array(m).Scan(src)
}

// Value implements the driver.Valuer interface.
// This tells the pq driver how to format the struct when writing to the database.
// This is used when inserting values
func (m ManagerInfo) Value() (driver.Value, error) {
	// Formats the struct as a string literal "(email,name)".
	return fmt.Sprintf("(%s,%s)", m.ManagerEmail, m.ManagerName), nil
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

// Value implements the driver.Valuer interface, telling the pq driver
// how to format the struct when writing to the database.
// Value implements the driver.Valuer interface.
func (a ApproverInfo) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%t)", a.ApproverEmail, a.ApprovalRequired), nil
}

type Rules struct {
	EmployeeEmail        string
	IsManagerApprover    bool
	MinApprovalPercent   int
	IsApprovalSequential bool
	Approvers            ApproverInfoSlice
}

type Expense struct {
    ExpenseID     int
    EmployeeEmail string
    Description   string
    ExpenseDate   time.Time
    Category      string
    Amount        int
    Remarks       string
    Status        string
}

// ApprovalStatus represents a record in the 'approval_status' table.
type ApprovalStatus struct {
    ExpenseID         int
    ManagerEmail      string
    ApprovalTimestamp time.Time
    Status            string
}