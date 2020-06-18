PRAGMA foreign_keys = ON;

CREATE TABLE types ( id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name TEXT NOT NULL UNIQUE);

CREATE TABLE tea ( id INTEGER PRIMARY KEY AUTOINCREMENT,
                   name TEXT NOT NULL UNIQUE,
                   teaType INTEGER,
                   FOREIGN KEY (teaType) REFERENCES types (id)
                    ON UPDATE CASCADE
                    ON DELETE RESTRICT);

CREATE TABLE owner ( id INTEGER PRIMARY KEY AUTOINCREMENT,
                     name TEXT NOT NULL UNIQUE);

CREATE TABLE teaOwners ( teaID INTEGER,
                         ownerID INTEGER,
                         FOREIGN KEY (teaID) REFERENCES tea (id)
                            ON UPDATE CASCADE
                            ON DELETE RESTRICT,
                         FOREIGN KEY (ownerID) REFERENCES owner (id)
                            ON UPDATE CASCADE
                            ON DELETE RESTRICT);

INSERT INTO types (name) VALUES ('Black Tea'), ('Green Tea'), ('White Tea');

INSERT INTO owner (name) VALUES ('John'), ('Jane');

INSERT INTO tea (name, teaType) VALUES ('Snowball', 1), ('Nearly Nirvana', 3), ('Green Tea Orient', 2);

INSERT INTO teaOwners (teaID, ownerID) VALUES (1,1), (2,1), (2,2), (3, 2);