CREATE TABLE IF NOT EXISTS roles
(
    id          SERIAL PRIMARY KEY NOT NULL,
    role_name   citext UNIQUE      NOT NULL,
    level       INT DEFAULT 0,
    description TEXT
);

CREATE UNIQUE INDEX roles_roles_name_idx ON roles (role_name);

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY NOT NULL,
    role_id    INT                NOT NULL REFERENCES roles (id),
    username   VARCHAR(50) UNIQUE NOT NULL,
    email      citext UNIQUE      NOT NULL,
    password   bytea              NOT NULL,
    created_at timestamptz default now(),
    last_login timestamptz
);

CREATE UNIQUE INDEX users_email_idx ON users (email);
CREATE UNIQUE INDEX users_user_name_idx ON users (username);

CREATE TABLE IF NOT EXISTS currencies
(
    id      SERIAL PRIMARY KEY NOT NULL,
    code    VARCHAR(3) UNIQUE  NOT NULL,
    name    VARCHAR(50),
    iconURL text
);

CREATE UNIQUE INDEX currencies_code_idx ON currencies (code);

CREATE TABLE IF NOT EXISTS exchange_rates
(
    id                 SERIAL PRIMARY KEY NOT NULL,
    base_currency_id   INT                NOT NULL REFERENCES currencies (id),
    target_currency_id INT                NOT NULL REFERENCES currencies (id),
    rate               DECIMAL(18, 8)     NOT NULL,
    last_update        TIMESTAMP          NOT NULL,
    next_update        TIMESTAMP          NOT NULL,

    -- Unique key pair
    UNIQUE (base_currency_id, target_currency_id)
);

CREATE TABLE IF NOT EXISTS transactions
(
    id                 SERIAL PRIMARY KEY NOT NULL,
    user_id            INT                NOT NULL REFERENCES users (id),
    base_currency_id   INT                NOT NULL REFERENCES currencies (id),
    target_currency_id INT                NOT NULL REFERENCES currencies (id),
    amount             DECIMAL(18, 8)     NOT NULL,
    converted_amount   DECIMAL(18, 8)     NOT NULL,
    converted_rate     DECIMAL(18, 8)     NOT NULL,
    created_at         TIMESTAMP default now()
);




