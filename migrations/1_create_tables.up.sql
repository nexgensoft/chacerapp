CREATE TABLE account (
    id           UUID NOT NULL DEFAULT gen_random_uuid(),
    name         STRING NOT NULL,
    display_name STRING,
    quotas       JSONB,
    status       JSONB,
    created_time TIMESTAMP,
    updated_time TIMESTAMP,
    CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

CREATE TABLE location (
    id           UUID NOT NULL DEFAULT gen_random_uuid(),
    name         STRING NOT NULL,
    account      STRING NOT NULL,
    display_name STRING,
    description  STRING,
    created_time TIMESTAMP,
    updated_time TIMESTAMP,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX (account ASC)
);

CREATE TABLE room (
    id           UUID NOT NULL DEFAULT gen_random_uuid(),
    name         STRING NOT NULL,
    account      STRING NOT NULL,
    location      STRING NOT NULL,
    display_name STRING,
    description  STRING,
    created_time TIMESTAMP,
    updated_time TIMESTAMP,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX (account ASC)
);
