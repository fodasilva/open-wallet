CREATE TABLE recurrences (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL REFERENCES users(id),
    name VARCHAR NOT NULL,
    note VARCHAR(400),
    amount NUMERIC NOT NULL,
    day_of_month INTEGER NOT NULL CHECK (day_of_month BETWEEN 1 AND 31),
    start_period VARCHAR(6) NOT NULL,
    end_period VARCHAR(6),
    category_id VARCHAR REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
