CREATE DATABASE datahara;

\c datahara;

CREATE SCHEMA IF NOT EXISTS geodata;

CREATE TABLE geodata.kondisi (
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	map JSONB
	
);

-- GRANT ROLES FOR TOOLS
CREATE ROLE anon NOLOGIN;
CREATE ROLE authenticator NOINHERIT LOGIN PASSWORD 'postgres';

GRANT USAGE ON SCHEMA geodata TO anon;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA geodata TO anon;

GRANT anon TO authenticator;

