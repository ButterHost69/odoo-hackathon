package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/lib/pq" // registers driver
)

var db *sql.DB

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

	_, err = stmt.Exec(session_token, email)
	if err != nil {
		fmt.Println("[db.UpdateSessionTokenInDB] Error Occured: \n", err.Error())
		return err
	}

	fmt.Println("[LOG] [db.UpdateSessionTokenInDB]  Update Session Token For User: ", email)
	return nil
}

func InsertNewRecordInAuthDB(email string, password string) error {
	// query := "INSERT INTO auth (email, password, session_token) VALUE (?,?,?)"

	query := `INSERT INTO auth
              (email, password, session_token) 
              VALUES ($1, $2, $3)`

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

func InsertNewCompany(name, country, currency, adminEmail string, managers []ManagerInfo) (int, error) {
	query := `INSERT INTO company 
              (company_name, country, currency, admin_email, managers) 
              VALUES ($1, $2, $3, $4, $5)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("[db.InsertNewCompany] Prepare Error: ", err.Error())
		return -1, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, country, currency, adminEmail, pq.Array(managers))
	if err != nil {
		fmt.Println("[db.InsertNewCompany] Exec Error: ", err.Error())
		return -1, err
	}

	fmt.Println("[LOG] [db.InsertNewCompany] Added New Company Record For: ", name)

	query = "SELECT company_id FROM company WHERE admin_email = $1"

	var company_id int
	err = db.QueryRow(query, adminEmail).Scan(&company_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errs.ErrAdminEmailNotFound
		}
		// For all other potential errors (connection issues, etc.), return the original error.
		fmt.Printf("[db.InsertNewCompany] Database error: %v\n", err)
		return -1, err
	}

	fmt.Printf("[LOG] [db.InsertNewCompany] Successfully found company id for the given admin email.\n")

	return company_id, nil
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

func GetAllUsersDetailsUsingCompanyID(company_id int) ([]User, error) {
	query := "SELECT email, name, role, manager_email, manager_name, company_id FROM user_account WHERE company_id = $1"

	rows, err := db.Query(query, company_id)
	if err != nil {
		fmt.Printf("[db.GetAllUsersDetailsUsingCompanyID] Query Error: %v\n", err)
		return nil, err // Return nil slice and the error.
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
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

func GetRulesUsingUserEmail(email string) (Rules, error) {
	query := `SELECT
                employee_email,
                is_manager_approver,
                min_approval_percent,
                is_approval_sequential,
                approvers
              FROM rules WHERE employee_email = $1`

	var rules Rules

	err := db.QueryRow(query, email).Scan(
		&rules.EmployeeEmail,
		&rules.IsManagerApprover,
		&rules.MinApprovalPercent,
		&rules.IsApprovalSequential,
		&rules.Approvers,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[db.GetRulesUsingUserEmail] No rules found for email: %s\n", email)
		} else {
			fmt.Printf("[db.GetRulesUsingUserEmail] Error scanning rules: %v\n", err)
		}
		// Return a zero-value struct and the error.
		return Rules{}, err
	}

	fmt.Printf("[LOG] [db.GetRulesUsingUserEmail] Successfully fetched rules for email: %s\n", email)
	return rules, nil
}

func GetAllManagerListUsingCompanyID(company_id int) ([]ManagerInfo, error) {
	query := "SELECT managers FROM company WHERE company_id = $1"

	var managers ManagerInfoSlice

	// db.QueryRow is used because we expect only one row (the company record).
	// The Scan method will use our custom ManagerInfoSlice.Scan method automatically.
	err := db.QueryRow(query, company_id).Scan(&managers)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[db.GetAllManagerListUsingCompanyID] No company found for ID: %d\n", company_id)
		} else {
			fmt.Printf("[db.GetAllManagerListUsingCompanyID] Error scanning managers list: %v\n", err)
		}
		// On any error, return a nil slice and the error itself.
		return nil, err
	}

	fmt.Printf("[LOG] [db.GetAllManagerListUsingCompanyID] Fetched %d managers for company ID %d\n", len(managers), company_id)

	// On success, return the populated slice and a nil error.
	return managers, nil
}

func UpdateManagerListInCompanyUsingCompanyID(company_id int, managers []ManagerInfo) error {

	query := "UPDATE company SET managers = $1 WHERE company_id = $2"

	// db.Exec is used for statements that don't return rows (UPDATE, INSERT, DELETE).
	// We wrap the 'managers' slice with pq.Array() to handle the conversion.
	result, err := db.Exec(query, pq.Array(managers), company_id)
	if err != nil {
		fmt.Printf("[db.UpdateManagerList] Error executing update: %v\n", err)
		return err
	}

	// It's crucial to check if any rows were actually updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("[db.UpdateManagerList] Error checking rows affected: %v\n", err)
		return err
	}

	// If no rows are affected, it means no company with the given ID was found.
	if rowsAffected == 0 {
		return fmt.Errorf("no company found with ID %d to update", company_id)
	}

	fmt.Printf("[LOG] [db.UpdateManagerList] Successfully updated managers for company ID: %d\n", company_id)
	return nil

}

func UpdateRulesUsingEmailID(email string, rules Rules) error {
	query := `UPDATE rules SET
                is_manager_approver = $1,
                min_approval_percent = $2,
                is_approval_sequential = $3,
                approvers = $4
              WHERE employee_email = $5`

	result, err := db.Exec(query,
		rules.IsManagerApprover,
		rules.MinApprovalPercent,
		rules.IsApprovalSequential,
		pq.Array(rules.Approvers), // This works perfectly with ApproverInfoSlice
		email,
	)
	if err != nil {
		fmt.Printf("[db.UpdateRulesUsingEmailID] Error executing update: %v\n", err)
		return err
	}

	// Check if any rows were actually modified.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("[db.UpdateRulesUsingEmailID] Error checking rows affected: %v\n", err)
		return err
	}

	// If no rows are affected, the email was not found in the table.
	if rowsAffected == 0 {
		return fmt.Errorf("no rules found for email %s to update", email)
	}

	fmt.Printf("[LOG] [db.UpdateRulesUsingEmailID] Successfully updated rules for email: %s\n", email)
	return nil
}

func GetUserDetailsUsingEmail(email string) (User, error) {
	query := `SELECT
                email, name, role, manager_email, manager_name, company_id
              FROM user_account WHERE email = $1`

	var user User

	err := db.QueryRow(query, email).Scan(
		&user.Email,
		&user.Name,
		&user.Role,
		&user.ManagerEmail,
		&user.ManagerName,
		&user.CompanyID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[db.GetUserDetailsUsingEmail] No user found for email: %s\n", email)
			return User{}, errs.ErrUserEmailDoesNotExist
		}
		fmt.Printf("[db.GetUserDetailsUsingEmail] Error scanning user details: %v\n", err)
		return User{}, err
	}

	fmt.Printf("[LOG] [db.GetUserDetailsUsingEmail] Successfully fetched details for user: %s\n", email)

	return user, nil
}

func GetEmailUsingSessionToken(session_token string) (string, error) {
	query := "SELECT email FROM auth WHERE session_token = $1"

	var email string

	err := db.QueryRow(query, session_token).Scan(&email)
	if err != nil {
		// Check if the error is specifically because no record was found.
		if err == sql.ErrNoRows {
			return "", errs.ErrSessionTokenDoesNotExist
		}
		fmt.Printf("[db.GetEmailUsingSessionToken] Database error: %v\n", err)
		return "", err
	}

	fmt.Printf("[LOG] [db.GetEmailUsingSessionToken] Successfully found email for the given session token.\n")

	return email, nil
}

func GetExpenseIDsByManagerEmail(managerEmail string) ([]int, error) {
    query := "SELECT expense_id FROM approval_status WHERE manager_email = $1"

    rows, err := db.Query(query, managerEmail)
    if err != nil {
        fmt.Printf("[db.GetExpenseIDsByManagerEmail] Query error: %v\n", err)
        return nil, err
    }
    defer rows.Close()

    var expenseIDs []int

    for rows.Next() {
        var id int
        if err := rows.Scan(&id); err != nil {
            fmt.Printf("[db.GetExpenseIDsByManagerEmail] Scan error: %v\n", err)
            return nil, err
        }
        expenseIDs = append(expenseIDs, id)
    }

    if err = rows.Err(); err != nil {
        fmt.Printf("[db.GetExpenseIDsByManagerEmail] Rows iteration error: %v\n", err)
        return nil, err
    }

    fmt.Printf("[LOG] [db.GetExpenseIDsByManagerEmail] Found %d expense(s) for manager %s\n", len(expenseIDs), managerEmail)
    
    return expenseIDs, nil
}