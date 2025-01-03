#!/usr/bin/env bash

# usage: file_env VAR [DEFAULT]
#    ie: file_env 'XYZ_DB_PASSWORD' 'example'
# (will allow for "$XYZ_DB_PASSWORD_FILE" to fill in the value of
#  "$XYZ_DB_PASSWORD" from a file, especially for Docker's secrets feature)
file_env() {
	local var="$1"
	local fileVar="${var}_FILE"
	local def="${2:-}"
	if [ "${!var:-}" ] && [ "${!fileVar:-}" ]; then
		printf >&2 'error: both %s and %s are set (but are exclusive)\n' "$var" "$fileVar"
		exit 1
	fi
	local val="$def"
	if [ "${!var:-}" ]; then
		val="${!var}"
	elif [ "${!fileVar:-}" ]; then
		val="$(< "${!fileVar}")"
	fi
	export "$var"="$val"
	unset "$fileVar"
}
echo "pulling out the password to local var."
file_env "JOBSEARCH_PASSWORD"
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER $JOBSEARCH_USER WITH PASSWORD '$JOBSEARCH_PASSWORD';
	CREATE DATABASE $JOBSEARCH_DB;
	GRANT ALL PRIVILEGES ON DATABASE $JOBSEARCH_DB TO $JOBSEARCH_USER;
	\c $JOBSEARCH_DB
        GRANT ALL ON SCHEMA public TO $JOBSEARCH_USER;
EOSQL

echo "creating users table and triggers on $JOBSEARCH_DB"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$JOBSEARCH_DB" < /sql/users.sql
echo "users table and trigger created"

echo "creating jobs table and triggers on $JOBSEARCH_DB"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$JOBSEARCH_DB" < /sql/jobs.sql
echo "jobs table and trigger created"

echo "creating password encryption and handling function on $JOBSEARCH_DB"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$JOBSEARCH_DB" < /sql/password_proc.sql
echo "password handlers created"

echo "Switching ownership on tables, procedures and triggert to $JOBSEARCH_USER"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	\c $JOBSEARCH_DB
	ALTER TABLE IF EXISTS public.users OWNER TO $JOBSEARCH_USER;
	ALTER TABLE IF EXISTS public.jobs OWNER TO $JOBSEARCH_USER;
	ALTER TABLE IF EXISTS public.password_store OWNER TO $JOBSEARCH_USER;
	ALTER PROCEDURE public.gen_user OWNER TO $JOBSEARCH_USER;
	ALTER FUNCTION public.verify_user OWNER TO $JOBSEARCH_USER;
EOSQL

echo "done creating database... sheeee!"
