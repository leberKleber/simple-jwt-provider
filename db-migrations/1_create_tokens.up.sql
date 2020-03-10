CREATE TABLE tokens
(
    id         serial      NOT NULL,
    email      text        NOT NULL,
    token      text        NOT NULL,
    type       text        NOT NULL,
    created_at timestamptz NOT NULL,
    CONSTRAINT id_unique PRIMARY KEY (id),
    CONSTRAINT tokens_email_fkey FOREIGN KEY (email) REFERENCES  users(email)
);
