FROM postgres:15.0-alpine

# Copy the database initialize script: 
# Contents of /docker-entrypoint-initdb.d are run on startup

COPY pg_hba.conf /usr/share/postgresql/15/
COPY postgresql.conf /usr/share/postgresql/15/
ADD  docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
ADD  sql/ /sql/

# Default values for passwords and database name. Can be overridden on docker run
