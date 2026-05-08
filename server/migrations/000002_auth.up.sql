CREATE TABLE user_auth (
    id               SERIAL PRIMARY KEY,
    first_name       VARCHAR(100) NOT NULL,
    last_name        VARCHAR(100) NOT NULL,
    email            VARCHAR(255) NOT NULL,
    hashed_password  VARCHAR(255) NOT NULL,
    profile_image    VARCHAR(500),
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_user_auth_email UNIQUE (email)
);