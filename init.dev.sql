-- Create studoto database
-- This script runs automatically on first database initialization
-- Note: This will only run when the PostgreSQL data directory is empty

CREATE DATABASE studoto WITH OWNER = postgres ENCODING = 'UTF8' CONNECTION LIMIT = -1;

-- Create studoto role/user
CREATE ROLE studoto WITH 
    LOGIN 
    NOSUPERUSER 
    NOCREATEDB 
    NOCREATEROLE 
    NOINHERIT 
    NOREPLICATION 
    CONNECTION LIMIT -1 
    PASSWORD 'studoto';

-- Grant all privileges on the studoto database to the studoto role
GRANT ALL PRIVILEGES ON DATABASE studoto TO studoto;