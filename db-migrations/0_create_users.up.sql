CREATE TABLE users
(
    email    text  NOT NULL,
    password bytea NOT NULL,
    CONSTRAINT email_unique PRIMARY KEY (email)
);
