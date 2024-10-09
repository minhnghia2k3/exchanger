CREATE DATABASE exchanger;
\c exchanger;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE USER admin WITH PASSWORD 'secret';
GRANT ALL privileges ON DATABASE exchanger TO admin;
