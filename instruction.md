#!/usr/bin/env bash

DBHOST=host:5432
DBUSER=user
DBPASS=password
DBNAME=dbname
DBSCHEMA=schema

mkdir -p structs/${DBNAME}
COMMAND=$(cat <<EOF
-s ${DBSCHEMA} -p ${DBNAME}_${DBSCHEMA} -o structs/${DBNAME}/ postgres://${DBUSER}@${DBHOST}/${DBNAME}?sslmode=disable&password=${DBPASS}
EOF
)
echo $COMMAND

go run . $COMMAND
