CREATE TABLE location (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    country     CHAR(2),
    longitude   DECIMAL(9,5),
    latitude    DECIMAL(8,5),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE candidate_profile (
    candidate_id            SERIAL PRIMARY KEY,
    user_id                 INT NOT NULL UNIQUE REFERENCES user_auth(id) ON DELETE CASCADE,
    location_id             INT REFERENCES location(id) ON DELETE SET NULL,
    status                  SMALLINT NOT NULL DEFAULT 0,
    phone_country_code      VARCHAR(5),
    phone_number            VARCHAR(15),
    hunting_status          SMALLINT NOT NULL DEFAULT -1,
    np_days                 SMALLINT NOT NULL DEFAULT -1,
    mostly_worked_with      SMALLINT NOT NULL DEFAULT -1,
    social_profiles         JSONB NOT NULL DEFAULT '[]',
    is_phone_verified       BOOLEAN NOT NULL DEFAULT FALSE,
    open_to_relocate        BOOLEAN NOT NULL DEFAULT TRUE,
    current_salary          INT,
    min_expected_salary     INT,
    preferred_work_mode     SMALLINT NOT NULL DEFAULT -1,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE company_profile (
    company_id      SERIAL PRIMARY KEY,
    user_id         INT NOT NULL UNIQUE REFERENCES user_auth(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    funding         SMALLINT NOT NULL DEFAULT -1,
    revenue_range   SMALLINT NOT NULL DEFAULT -1,
    website_url     VARCHAR(500) NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    tagline         VARCHAR,
    team_size       SMALLINT NOT NULL DEFAULT -1,
    social_profiles JSONB NOT NULL DEFAULT '[]',
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);