# # include custom config file

# echo "include = '/etc/postgresql.conf'" >> $PGDATA/postgresql.conf


#!/bin/bash


echo "Creating db"
POSTGRES_DB= docker_process_sql <<-EOSQL
  CREATE DATABASE pgcdc;
EOSQL

POSTGRES_DB= docker_process_sql --dbname pgcdc <<-EOSQL
  create table protobuf (id int, proto bytea);
EOSQL

cat >> ${PGDATA}/pg_hba.conf <<EOF
host replication all 0.0.0.0/0 md5
EOF


cat >> ${PGDATA}/postgresql.conf <<EOF
wal_level = logical
max_wal_senders=5
max_replication_slots=5
EOF


# echo "Creating publication"
# POSTGRES_DB= docker_process_sql <<-EOSQL
#   CREATE PUBLICATION percpub FOR ALL TABLES;
# EOSQL