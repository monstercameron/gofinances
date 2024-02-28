--drop table debts

-- Add extra dummy data here

select * from users

delete from users

select * from menus

delete from menus

DELETE FROM users WHERE id=1

select * from debts
VALUES (1, 'Car Note', 'Cam', '2022-01-01', '2028-01-01', 48995.99, 37689.95, 'Sccu Ford Escape');


--drop table debts
DROP TABLE IF EXISTS debts;


CREATE TABLE IF NOT EXISTS debts (
    id INTEGER PRIMARY KEY,
    debt_name TEXT,
    debt_owner TEXT,
    start_date TEXT,
    end_date TEXT,
    initial_amount REAL, -- was 'initial'
    current_amount REAL, -- was 'current'
    interest_rate REAL,  -- was 'interest'
    credit_limit REAL,   -- was 'limit'
    notes TEXT
);



INSERT INTO debts (id, debt_name, debt_owner, start_date, end_date, initial_amount, current_amount, interest_rate, credit_limit, notes) VALUES
(1, 'Car Loan', 'John Doe', '2023-01-01', '2028-01-01', 20000.00, 15000.00, 3.50, 25000.00, 'Monthly payments'),
(2, 'Credit Card', 'Alice Johnson', '2023-02-15', '2024-02-15', 5000.00, 4500.00, 19.99, 7000.00, 'High interest rate'),
(3, 'Mortgage', 'Emily Smith', '2022-06-30', '2042-06-30', 300000.00, 290000.00, 2.75, 500000.00, 'Fixed rate for 20 years'),
(4, 'Student Loan', 'Michael Brown', '2021-09-01', '2031-09-01', 25000.00, 24300.00, 4.25, 0.00, 'Deferred interest for 2 years'),
(5, 'Personal Loan', 'Lisa Davis', '2023-03-01', '2026-03-01', 10000.00, 9800.00, 5.00, 15000.00, 'No early repayment penalty');


select * from debts