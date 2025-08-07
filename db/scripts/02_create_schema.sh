#! /bin/bash
set -e

echo "Creating database '${POSTGRES_DB_NAME}'"

POSTGRES="psql -U ${POSTGRES_DB_NAME}"

function buildSchema {
cat <<EOF

    DROP SCHEMA IF EXISTS $1 CASCADE;
    CREATE SCHEMA $1;

    CREATE TABLE IF NOT EXISTS $1.users (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

        name TEXT NOT NULL, 
        email TEXT NOT NULL, 

        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
        updated_at TIMESTAMP WITH TIME ZONE
    );
    
    CREATE TABLE IF NOT EXISTS $1.posts (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

        user_id uuid NOT NULL,

        title TEXT NOT NULL, 
        content TEXT NOT NULL, 

        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
        updated_at TIMESTAMP WITH TIME ZONE,

        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES $1.users(id) ON DELETE CASCADE
    );
    
EOF
}

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    $(buildSchema 'public')
    $(buildSchema 'test')
EOSQL
