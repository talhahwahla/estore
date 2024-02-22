#!/bin/bash
set -o errexit \
    -o pipefail \
    -o nounset \
    -o noglob

MYSQL_USER="root"
MYSQL_PASSWORD=""
MYSQL_HOST="localhost"
MYSQL_PORT="3306"
MYSQL_DATABASE="estore"

DB_URL="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}"

function main () {
    echo 'create_user'
    create_user

    echo 'create_database'
    create_database

    echo 'run_migrations'
    run_migrations
} && readonly -f main

function create_user () {
    mysql --host="$MYSQL_HOST" --port="$MYSQL_PORT" --user="$MYSQL_USER" --password="$MYSQL_PASSWORD" -e "
        CREATE USER IF NOT EXISTS '$MYSQL_USER'@'$MYSQL_HOST' IDENTIFIED BY '$MYSQL_PASSWORD';
        GRANT ALL PRIVILEGES ON *.* TO '$MYSQL_USER'@'$MYSQL_HOST';
        FLUSH PRIVILEGES;
    "
} && readonly -f create_user

function create_database () {
    mysql --host="$MYSQL_HOST" --port="$MYSQL_PORT" --user="$MYSQL_USER" --password="$MYSQL_PASSWORD" -e "
        CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE;
    "
} && readonly -f create_database

function run_migrations () {
    local -r files=$(find ./ -type f -name '??-*.sql' | sort)

    for migration in $files; do
        mysql_file "$migration"
    done
} && readonly -f run_migrations

function mysql_file () {
    mysql --host="$MYSQL_HOST" --port="$MYSQL_PORT" --user="$MYSQL_USER" --password="$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$1"
} && readonly -f mysql_file

main
