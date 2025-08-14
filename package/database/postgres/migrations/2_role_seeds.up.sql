INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator role'), ('superadmin', 'Super administrator role')
ON CONFLICT (name) DO NOTHING;