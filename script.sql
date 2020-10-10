DROP SCHEMA IF EXISTS  AUTH;
CREATE SCHEMA AUTH;                                                             

USE AUTH;

CREATE TABLE  ACCOUNTS (
ID BIGINT AUTO_INCREMENT NOT NULL,
USERNAME VARCHAR(64) UNIQUE NOT NULL,
PASSWORD VARCHAR(256) NOT NULL,
ACCOUNT_TYPE VARCHAR(32) NOT NULL,
CREATION_DATE DATE NOT NULL,
PRIMARY KEY (ID)
);


CREATE  TABLE USER_INFO (
ID BIGINT AUTO_INCREMENT NOT NULL,
FULLNAME VARCHAR(128) DEFAULT NULL,
ADDRESS VARCHAR(256) DEFAULT NULL,
PHONE VARCHAR(32) DEFAULT NULL,
ACCOUNT_ID BIGINT NOT NULL,
PRIMARY KEY (ID),
FOREIGN KEY (ACCOUNT_ID) REFERENCES ACCOUNTS(ID)
);



CREATE TABLE SESSION_TOKEN(
TOKEN VARCHAR(255) NOT NULL,
EXPIRATION_DATE DATE NOT NULL,
ACCOUNT_ID BIGINT NOT NULL,
FOREIGN KEY (ACCOUNT_ID) REFERENCES ACCOUNTS(ID)
)

