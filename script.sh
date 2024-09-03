#!/usr/bin/env bash

set -e

let COUNT=10000
PERCENT="0.021"

go run cmd/simulate/main.go --csv-file=simulation.csv --simulation-time=36h --spread-percent=$PERCENT --object-count=$COUNT --overwrite-csv-file

for n in {00,03,06,09,12,15,18,21,23,26,29,32}; do go run cmd/graph/main.go --csv-file=simulation.csv --graph-start-time=${n}h --graph-length=4h --image-file=time-plus-$n-hours.png --overwrite-image-file; done 
