CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    type VARCHAR NOT NULL,
    category VARCHAR,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);