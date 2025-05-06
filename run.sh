#!/bin/bash
trap "rm server;kill 0" EXIT

PROTOCOL=${PROTOCOL:-http}  # Default to http if not set

go build -o server
./server -port=8001 -protocol=$PROTOCOL &
./server -port=8002 -protocol=$PROTOCOL &
./server -port=8003 -api=1 -protocol=$PROTOCOL &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
