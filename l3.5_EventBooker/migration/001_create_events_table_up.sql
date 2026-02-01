CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    capacity INT NOT NULL,
    free_seats INT NOT NULL,
    payment_ttl INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);