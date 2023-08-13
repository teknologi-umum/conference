#!/bin/bash

DBSTRING="host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME"

goose postgres "$DBSTRING" up