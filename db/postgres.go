package db

import (
	"crypto/rand"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/ButterHost69/odoo-hackathon/models"
	"github.com/lib/pq" // registers driver
)

var db *sql.DB

// ManagerInfo corresponds to the 'manager_info' SQL type
type ManagerInfo struct {
	ManagerEmail string
	ManagerName  string
}

// ApproverInfo corresponds to the 'approver_info' SQL type
type ApproverInfo struct {
	ApproverEmail    string
	ApprovalRequired bool
}

// Value implements the driver.Valuer interface for ManagerInfo.
// This tells the pq driver how to format the struct for the database.
func (m ManagerInfo) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", m.ManagerEmail, m.ManagerName), nil
}

// Value implements the driver.Valuer interface for ApproverInfo.
func (a ApproverInfo) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%t)", a.ApproverEmail, a.ApprovalRequired), nil
}

func InitDB() error {
	fmt.Println("[Log] Connecting to Postgress")
	sqlPassword := os.Getenv("POSTGRES_PASSWORD")
	if sqlPassword == "" {
		fmt.Print("Enter Your My POSTGRES Database Password: ")
		fmt.Scan(&sqlPassword)
	}

	sqlUsername := os.Getenv("POSTGRES_USERNAME")
	sqlDBName := os.Getenv("POSTGRES_DBNAME")

	sqlDBIP := os.Getenv("POSTGRES_IP")
	sqlDBPort := os.Getenv("POSTGRES_PORT")

	dblink := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		sqlUsername, sqlPassword, sqlDBIP, sqlDBPort, sqlDBName,
	)

	var err error
	db, err = sql.Open("postgres", dblink)
	if err != nil {
		fmt.Println("[db.InitDB] error in connecting to postgress db: ", err)
		return err
	}

	fmt.Println("[Log] Connected to postgress db is succesful !!")

	// Pinging The Database To Verify The Connection
	if err = db.Ping(); err != nil {
		fmt.Println("[db.InitDB] error in pinging to postgress db: ", err)
		return err
	}

	fmt.Println("[Log] Pinging to postgress db is succesful !!")

	return nil
}

func GenerateToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("[db.GenerateToken] Error Occured : \n", err.Error())
		return "", err
	}

	customEncoding := base64.RawURLEncoding
	token := customEncoding.EncodeToString(randomBytes)

	return token, nil
}

func UpdateSessionTokenInAuthDB(email string, session_token string) error {
	query := "UPDATE auth SET session_token = $1 WHERE email = $2"
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.UpdateSessionTokenInDB] Error Occured : \n", err.Error())
		return err
	}

	_, err = stmt.Exec(email, session_token)
	if err != nil {
		fmt.Println("[db.UpdateSessionTokenInDB] Error Occured: \n", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.UpdateSessionTokenInDB]  Update Session Token For User: ", email)
	return nil
}

func InsertNewRecordInAuthDB(email string, password string) error {
	query := "INSERT INTO auth (email, password, session_token) VALUE (?,?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewRecordAuthInDB] Error Occured : \n", err.Error())
		return err
	}

	_, err = stmt.Exec(email, password, "")
	if err != nil {
		fmt.Println("[db.InsertNewRecordAuthInDB] Error Occured: \n", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewRecordAuthInDB] Added New User Record To AUTH For ID: ", email)
	return nil
}

func InsertNewCompany(name, country, currency, adminEmail string, managers []ManagerInfo) error {
	query := `INSERT INTO company 
              (company_name, country, currency, admin_email, managers) 
              VALUES ($1, $2, $3, $4, $5)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewCompany] Prepare Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, country, currency, adminEmail, pq.Array(managers))
	if err != nil {
		fmt.Println("[db.InsertNewCompany] Exec Error: ", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewCompany] Added New Company Record For: ", name)
	return nil
}

func InsertNewUserAccount(email, name, role, managerEmail, managerName string, companyID int) error {
	query := `INSERT INTO user_account 
              (email, name, role, manager_email, manager_name, company_id) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewUserAccount] Prepare Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, name, role, managerEmail, managerName, companyID)
	if err != nil {
		fmt.Println("[db.InsertNewUserAccount] Exec Error: ", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewUserAccount] Added New User Account Record For: ", email)
	return nil
}

func InsertNewRule(employeeEmail string, isManagerApprover, isApprovalSequential bool, minApprovalPercent int, approvers []ApproverInfo) error {
	query := `INSERT INTO rules 
              (employee_email, is_manager_approver, min_approval_percent, is_approval_sequential, approvers) 
              VALUES ($1, $2, $3, $4, $5)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewRule] Prepare Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(employeeEmail, isManagerApprover, minApprovalPercent, isApprovalSequential, pq.Array(approvers))
	if err != nil {
		fmt.Println("[db.InsertNewRule] Exec Error: ", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewRule] Added New Rule Record For: ", employeeEmail)
	return nil
}

func InsertNewExpense(employeeEmail, description, category, remarks, status string, amount int, expenseDate time.Time) error {
	query := `INSERT INTO expenses 
              (employee_email, description, expense_date, category, amount, remarks, status) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewExpense] Prepare Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(employeeEmail, description, expenseDate, category, amount, remarks, status)
	if err != nil {
		fmt.Println("[db.InsertNewExpense] Exec Error: ", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewExpense] Added New Expense Record For: ", employeeEmail)
	return nil
}

func InsertNewApprovalStatus(expenseID int, managerEmail, status string, approvalTimestamp time.Time) error {
	query := `INSERT INTO approval_status 
              (expense_id, manager_email, approval_timestamp, status) 
              VALUES ($1, $2, $3, $4)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewApprovalStatus] Prepare Error: ", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(expenseID, managerEmail, approvalTimestamp, status)
	if err != nil {
		fmt.Println("[db.InsertNewApprovalStatus] Exec Error: ", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.InsertNewApprovalStatus] Added New Approval Status For Expense ID: ", expenseID)
	return nil
}

func GetPasswordByEmailFromAuth(email string) string {
	query := "SELECT password FROM auth WHERE email = $1"

	var password string

	err := db.QueryRow(query, email).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[db.GetPasswordByEmailFromAuth] No record found for email: %s\n", email)
		} else {
			fmt.Printf("[db.GetPasswordByEmailFromAuth] Error Occurred: %s\n", err.Error())
		}
		return ""
	}

	fmt.Printf("[LOG] [db.GetPasswordByEmailFromAuth] Fetched password for email: %s\n", email)
	return password
}

func GetSessionTokenByCredentials(email string, password string) string {
	query := "SELECT session_token FROM auth WHERE email = $1 AND password = $2"

	var sessionToken sql.NullString

	err := db.QueryRow(query, email, password).Scan(&sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[db.GetSessionTokenByCredentials] Invalid credentials for email: %s\n", email)
		} else {
			fmt.Printf("[db.GetSessionTokenByCredentials] Error Occurred: %s\n", err.Error())
		}
		return ""
	}

	if sessionToken.Valid {
		fmt.Printf("[LOG] [db.GetSessionTokenByCredentials] Fetched session token for email: %s\n", email)
		return sessionToken.String
	}

	return ""
}

func GetCompanyIDByAdminEmail(adminEmail string) (int, error) {
	query := "SELECT company_id FROM company WHERE admin_email = $1"

	var companyID int

	err := db.QueryRow(query, adminEmail).Scan(&companyID)
	if err != nil {
		fmt.Printf("[db.GetCompanyIDByAdminEmail] Error fetching ID for admin '%s': %v\n", adminEmail, err)
		return 0, err
	}

	fmt.Printf("[LOG] [db.GetCompanyIDByAdminEmail] Found company_id %d for admin '%s'\n", companyID, adminEmail)
	return companyID, nil
}

func GetAllUsersDetailsUsingCompanyID(company_id int) ([]models.User, error) {
	query := "SELECT email, name, role, manager_email, manager_name, company_id FROM user_account WHERE company_id = $1"

	rows, err := db.Query(query, company_id)
	if err != nil {
		fmt.Printf("[db.GetAllUsersDetailsUsingCompanyID] Query Error: %v\n", err)
		return nil, err // Return nil slice and the error.
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Email, &user.Name, &user.Role, &user.ManagerEmail, &user.ManagerName, &user.CompanyID); err != nil {
			fmt.Printf("[db.GetAllUsersDetailsUsingCompanyID] Scan Error: %v\n", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("[db.GetAllUsersDetailsUsingCompanyID] Rows Iteration Error: %v\n", err)
		return nil, err
	}

	fmt.Printf("[LOG] [db.GetAllUsersDetailsUsingCompanyID] Fetched %d users for company_id: %d\n", len(users), company_id)

	return users, nil
}
