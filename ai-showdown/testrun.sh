#!/bin/bash

set -e
set -x

i="0"
while [ $i -lt 300 ]
do
    ./ai-showdown
    sqlite3 smart-results.db < update-weights.sql
    i=$[$i+1]
done
