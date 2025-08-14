INSERT INTO organizations (name, short_name, address)
VALUES ('Superadmin Organization', 'SuperadminOrg', '789 Superadmin Blvd, Superadmin City, SA 11223')
ON CONFLICT (name) DO NOTHING;