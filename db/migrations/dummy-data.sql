INSERT INTO company (
    company_name,
    country,
    currency,
    admin_email,
    managers
) VALUES
(
    'Nemu Corp',
    'India',
    '₹',
    'admin@nemu.com',
    ARRAY[
        ROW('palash@nemu.com', 'Palash')::manager_info,
        ROW('prachin@nemu.com', 'Prachin')::manager_info
    ]
),
(
    'OpenAI',
    'United States',
    '$',
    'admin@openai.com',
    ARRAY[
        ROW('palash@openai.com', 'Palash')::manager_info,
        ROW('konark@openai.com', 'Konark')::manager_info
    ]
);

INSERT INTO user_account (email, name, role, manager_email, manager_name, company_id) VALUES
('palash@nemu.com', 'Palash', 'manager', 'admin@nemu.com', 'Admin', 1),
('prachin@nemu.com', 'Prachin', 'manager', 'admin@nemu.com', 'Admin', 1),
('parth@nemu.com', 'Parth', 'employee', 'palash@nemu.com', 'Palash', 1);

INSERT INTO user_account (email, name, role, manager_email, manager_name, company_id) VALUES
('palash@openai.com', 'Palash', 'manager', 'admin@openai.com', 'Admin', 2),
('konark@openai.com', 'Konark', 'manager', 'admin@openai.com', 'Admin', 2);

INSERT INTO auth (email, password, session_token) VALUES
('palash@nemu.com', 'pass123', 'token123'),
('prachin@nemu.com', 'pass456', 'token456'),
('parth@nemu.com', 'pass789', 'token789'),
('palash@openai.com', 'pass321', 'token321'),
('konark@openai.com', 'pass654', 'token654');

INSERT INTO rules (
    empployee_email,
    is_manager_approver,
    min_approval_percent,
    is_approval_sequential,
    approvers
) VALUES
(
    'parth@nemu.com',
    TRUE,
    75,
    FALSE,
    ARRAY[
        ROW('palash@nemu.com', TRUE)::approver_info,
        ROW('prachin@nemu.com', FALSE)::approver_info
    ]
),
(
    'konark@openai.com',
    FALSE,
    100,
    TRUE,
    ARRAY[
        ROW('palash@openai.com', TRUE)::approver_info
    ]
);

INSERT INTO expenses (
    employee_email,
    description,
    expense_date,
    category,
    amount,
    remarks,
    status
) VALUES
('parth@nemu.com', 'Flight to Mumbai', '2025-09-20', 'Travel', 15000, 'Client meeting', 'submitted'),
('parth@nemu.com', 'Hotel Stay', '2025-09-21', 'Accommodation', 8000, '2 nights stay', 'approved'),
('konark@openai.com', 'Conference Ticket', '2025-09-15', 'Training', 12000, 'AI Summit 2025', 'draft');

INSERT INTO approval_status (
    expense_id,
    manager_email,
    approval_timestamp,
    status
) VALUES
(1, 'palash@nemu.com', '2025-09-21 10:15:00', 'rejected'),
(2, 'prachin@nemu.com', '2025-09-22 09:00:00', 'approved'),
(3, 'palash@openai.com', '2025-09-16 14:30:00', 'pending');

