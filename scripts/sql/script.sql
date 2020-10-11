DROP SCHEMA IF EXISTS auth;
CREATE SCHEMA auth;

USE auth;

CREATE TABLE ACCOUNTS
(
    ID            BIGINT AUTO_INCREMENT NOT NULL,
    EMAIL         VARCHAR(64)   UNIQUE        NOT NULL,
    PASSWORD      VARCHAR(256)          NOT NULL,
    FULLNAME      VARCHAR(128) DEFAULT NULL,
    ADDRESS       VARCHAR(256) DEFAULT NULL,
    PHONE         VARCHAR(32)  DEFAULT NULL,
    ACCOUNT_ID    BIGINT                NOT NULL,
    ACCOUNT_TYPE  VARCHAR(32)           NOT NULL,
    CREATION_DATE DATETIME              NOT NULL,
    PRIMARY KEY (ID)
);


CREATE TABLE SESSION_TOKENS
(
    TOKEN           VARCHAR(255) UNIQUE NOT NULL,
    ACCOUNT_ID      BIGINT              NOT NULL,
    FOREIGN KEY (ACCOUNT_ID) REFERENCES ACCOUNTS (ID)
);


CREATE TABLE FORGOT_PASSWORD_TOKENS
(
    TOKEN           VARCHAR(32) UNIQUE NOT NULL,
    ACCOUNT_ID      BIGINT             NOT NULL,
    EXPIRATION_DATE DATETIME           NOT NULL,
    FOREIGN KEY (ACCOUNT_ID) REFERENCES ACCOUNTS (ID)
);

