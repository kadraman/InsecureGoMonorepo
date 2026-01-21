-- Seed data for demo database
-- Users
INSERT INTO users (username, email, password) VALUES ('seeduser', 'seeduser@example.com', 'md5:482c811da5d5b4bc6d497ffa98491e38');

-- Products
INSERT INTO products (name, description, price, category) VALUES ('Seed Product', 'A seeded product', 19.99, 'Seed');

-- (Optional) Orders can be inserted by services that snapshot data.
