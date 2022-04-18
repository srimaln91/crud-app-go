#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASS';
  CREATE DATABASE $APP_DB_NAME;
  GRANT ALL PRIVILEGES ON DATABASE $APP_DB_NAME TO $APP_DB_USER;
  \connect $APP_DB_NAME $APP_DB_USER
  BEGIN;
    CREATE SCHEMA events;
    CREATE TABLE events.events (
		event_id uuid NOT NULL,
		addr_nbr varchar(255) NULL,
		client_id varchar(255) NULL,
		event_cnt int4 NOT NULL,
		location_cd varchar(255) NULL,
		location_id1 varchar(255) NULL,
		location_id2 varchar(255) NULL,
		rc_num varchar(255) NULL,
		trans_id varchar(255) NULL,
		trans_tms varchar(255) NULL,
		CONSTRAINT events_pkey PRIMARY KEY (event_id)
	);
  COMMIT;
EOSQL