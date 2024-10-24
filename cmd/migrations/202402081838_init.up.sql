-- migrate Up
CREATE TABLE IF NOT EXISTS users (
       uuid UUID PRIMARY KEY,
       username VARCHAR(255) NOT NULL UNIQUE,
       password VARCHAR(255) NOT NULL
);