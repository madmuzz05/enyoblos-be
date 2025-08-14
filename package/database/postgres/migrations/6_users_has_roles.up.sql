CREATE TABLE IF NOT EXISTS users_has_roles (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    role_id INT NOT NULL REFERENCES roles(id),
    organization_id INT NOT NULL REFERENCES organizations(id)
);