INSERT INTO apps(id, name, secret)
VALUES (1, 'test-app', 'test-secret')
ON CONFLICT (id) DO NOTHING;