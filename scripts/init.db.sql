-- Create a database
CREATE DATABASE app_db;

-- Switch to the new database
\c app_db;

-- Create a user with a password
CREATE USER app_user WITH PASSWORD 'app_password';

-- Grant all privileges to the new user on the database
GRANT ALL PRIVILEGES ON DATABASE app_db TO app_user;

-- Create a sample table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    due_date DATETIME,
    completed BOOLEAN
);

-- Insert sample data
INSERT INTO tasks (id, title, description, due_date, completed)
VALUES 
(1, 'Buy groceries', 'Purchase milk, eggs, bread, and fruits.', '2024-11-19 10:00:00', false),
(2, 'Complete project report', 'Finalize and submit the quarterly project report.', '2024-11-20 15:00:00', false),
(3, 'Plan birthday party', 'Organize the venue, decorations, and guest list for the birthday party.', '2024-11-25 18:00:00', false),
(4, 'Attend team meeting', 'Participate in the weekly team sync-up meeting.', '2024-11-21 09:00:00', true),
(5, 'Car maintenance', 'Take the car for oil change and tire rotation.', '2024-11-22 11:00:00', false);
