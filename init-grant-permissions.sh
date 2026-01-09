#!/bin/bash
# Grant schema permissions for studoto user
# This script runs after database creation

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "studoto" <<-EOSQL
    -- Grant usage and create privileges on public schema
    GRANT USAGE ON SCHEMA public TO studoto;
    GRANT CREATE ON SCHEMA public TO studoto;

    -- Grant all privileges on all tables in public schema (for existing and future tables)
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO studoto;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO studoto;

    -- Set default privileges for future objects created by postgres user
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO studoto;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO studoto;
EOSQL
