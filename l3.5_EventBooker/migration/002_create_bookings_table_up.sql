CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id),
    seats INT NOT NULL,
    paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    expires_at TIMESTAMP NOT NULL
);  


