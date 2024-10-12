INSERT INTO roles(id, role_name, level, description)
VALUES (1, 'user', 1, 'an user can only perform convert currency');

INSERT INTO roles(id, role_name, level, description)
VALUES (2, 'moderator', 2, 'a moderator can perform add, update');

INSERT INTO roles(id, role_name, level, description)
VALUES (3, 'admin', 3, 'an admin can perform add, update and delete');
