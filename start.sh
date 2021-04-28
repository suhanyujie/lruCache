#!/usr/bin/env sh

trap "rm server; kill 0" EXIT

go build -o server

./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2

echo ">>> start test"
curl "http://localhost:9001/api?key=Stu2" &
curl "http://localhost:9001/api?key=Stu3" &

wait
