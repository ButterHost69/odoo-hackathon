CREATE DATABASE podoo;

CREATE TYPE manager_info AS (
	manager_email TEXT,
	manager_name TEXT
);

CREATE TABLE company (
	company_id SERIAL PRIMARY KEY,
	company_name TEXT,
	country TEXT,
	currency TEXT,
	admin_email TEXT,
	managers manager_info[]
);

CREATE TABLE user_account (
	email TEXT PRIMARY KEY,
	name TEXT,
	role TEXT,
	manager_email TEXT,
	manager_name TEXT,
	company_id INT,
	FOREIGN KEY (company_id) REFERENCES company(company_id)
);

CREATE TABLE auth (
	email TEXT,
	password TEXT,
	session_token TEXT,
	FOREIGN KEY (email) REFERENCES user_account(email)
);

CREATE TYPE approver_info AS (
    approver_email TEXT,
    approval_required BOOL
);

CREATE TABLE rules (
	empployee_email TEXT,
	is_manager_approver BOOL,
	min_approval_percent INT,
	is_approval_sequential BOOL,
	approvers approver_info[]
);

CREATE TABLE expenses (
	expense_id SERIAL PRIMARY KEY,
	employee_email TEXT,
	description TEXT,
	expense_date DATE,
	category TEXT,
	amount INT,
	remarks TEXT,
	status TEXT,
	FOREIGN KEY (employee_email) REFERENCES user_account(email)	
);

CREATE TABLE approval_status (
	expense_id INT,
	manager_email TEXT,
	approval_timestamp TIMESTAMP,
	status TEXT,
	FOREIGN KEY (expense_id) REFERENCES expenses(expense_id),
	FOREIGN KEY (manager_email) REFERENCES user_account(email)
);
