-- Seed SQL for orders-service
-- Order snapshots assume users/products seeded with ids 1 and 2
INSERT INTO orders (user_id, product_id, quantity, total_price, status, created_at, user_snapshot, product_snapshot) VALUES (
  1, 1, 2, 39.98, 'pending', '2026-01-21T00:00:00Z', '{"id":1,"username":"alice","email":"alice@example.com"}', '{"id":1,"name":"Widget","description":"A useful widget","price":19.99,"category":"Tools"}'
);

INSERT INTO orders (user_id, product_id, quantity, total_price, status, created_at, user_snapshot, product_snapshot) VALUES (
  2, 2, 1, 29.99, 'pending', '2026-01-21T00:00:00Z', '{"id":2,"username":"bob","email":"bob@example.com"}', '{"id":2,"name":"Gadget","description":"A fancy gadget","price":29.99,"category":"Electronics"}'
);
