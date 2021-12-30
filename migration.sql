CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
);
CREATE TABLE organisations (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL
);
CREATE TABLE organisation_members (
  user_id INTEGER NOT NULL,
  organisation_id INTEGER NOT NULL,
  role TEXT NOT NULL DEFAULT('viewer'),

  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (organisation_id) REFERENCES organisations(id)
);
CREATE TABLE devices (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  owner_id INTEGER NOT NULL,

  FOREIGN KEY (owner_id) REFERENCES organisations(id)
);

-- Session data
CREATE TABLE sessions (
  key VARCHAR(32) NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL
);

-- Mock data

INSERT INTO users (id,username,password) VALUES (1, 'john', 'password'),(2, 'bob', 'password'),(3, 'alice', 'password');
INSERT INTO organisations (id, name) VALUES (1, 'pzld'), (2, 'pollex');
INSERT INTO organisation_members (user_id, organisation_id, role)
  VALUES (1, 1, 'viewer'), (1, 2, 'administrator'), (2, 1, 'editor'), (3, 1, 'administrator');
INSERT INTO devices (id, name, owner_id)
  VALUES (1, 'temp_livingroom', 1), (2, 'temp_bedroom', 1), (3, 'light_kitchen', 2), (4, 'light_toilet', 3);