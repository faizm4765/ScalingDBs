-- ========================================
-- 1️⃣ Create users table and populate 50 users
-- ========================================

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL
);

INSERT INTO users (user_name) VALUES
('Alice Johnson'),
('Bob Smith'),
('Charlie Brown'),
('David Williams'),
('Eva Davis'),
('Frank Miller'),
('Grace Wilson'),
('Hannah Moore'),
('Ian Taylor'),
('Julia Anderson'),
('Kevin Thomas'),
('Laura Jackson'),
('Michael White'),
('Nina Harris'),
('Oscar Martin'),
('Paula Thompson'),
('Quinn Garcia'),
('Rachel Martinez'),
('Samuel Robinson'),
('Tina Clark'),
('Umar Lewis'),
('Victoria Lee'),
('William Walker'),
('Xander Hall'),
('Yara Allen'),
('Zachary Young'),
('Amber King'),
('Brian Scott'),
('Catherine Green'),
('Derek Adams'),
('Elena Baker'),
('Felix Gonzalez'),
('Gabriella Nelson'),
('Henry Carter'),
('Isla Mitchell'),
('Jack Perez'),
('Katherine Roberts'),
('Liam Turner'),
('Mia Phillips'),
('Noah Campbell'),
('Olivia Parker'),
('Peter Evans'),
('Queenie Edwards'),
('Ryan Collins'),
('Sophia Stewart'),
('Thomas Sanchez'),
('Uma Flores'),
('Violet Rivera'),
('Walter Morris'),
('Ximena Rogers');

-- ========================================
-- 2️⃣ Create seats table
-- ========================================

CREATE TABLE seats (
    seat_id SERIAL PRIMARY KEY,
    seat_number INT UNIQUE,
    user_id INT REFERENCES users(user_id)  -- nullable, initially empty
);

-- ========================================
-- 3️⃣ Populate seats table with 20 empty seats
-- ========================================

INSERT INTO seats (seat_number, user_id)
SELECT generate_series(1, 20), NULL;

-- ========================================
-- 4️⃣ Verify tables
-- ========================================

-- Check users
SELECT * FROM users;

-- Check seats
SELECT * FROM seats;
