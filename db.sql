--database bank

CREATE DATABASE IF NOT EXISTS bank;

CREATE TABLE IF NOT EXISTS accounts 
(
    id           INT, 
    currency     char(3)       DEFAULT 'tjs',
    external_ref varchar(64)   NOT NULL,
    balance      NUMERIC(18,2) DEFAULT 0,
    name         VARCHAR(100),
    regdate      TIMESTAMPTZ   DEFAULT now(),
    stsdate      TIMESTAMPTZ   DEFAULT now(),
    status       VARCHAR(20)   DEFAULT 'pending',
    PRIMARY KEY(id, currency)
);

CREATE TABLE IF NOT EXISTS txns
(
    id BIGSERIAL PRIMARY KEY    NOT NULL,
    service_name varchar(64)    NOT NULL,
    external_ref varchar(64)    NOT NULL,
    txn_type     varchar(20)    NOT NULL,
    account      int            NOT NULL,
    currency     char(3)        DEFAULT 'tjs',
    amount       numeric(18, 2) NOT NULL,
    fee          numeric(18, 2) default 0,
    description  varchar        DEFAULT '',
    balance      NUMERIC(18,2)  DEFAULT 0,
    regdate      timestamptz    DEFAULT now(),
    stsdate      timestamptz    DEFAULT now(),
    status       varchar(20)    DEFAULT 'pending',
    err_code     int            DEFAULT 0,
    err          varchar        DEFAULT '',
    UNIQUE (service_name, external_ref, txn_type)
);

CREATE INDEX ON txns (account, currency, amount);