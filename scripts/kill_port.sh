#!/bin/bash

PORT=$1

if [ -z "$PORT" ]; then
  echo "Usage: $0 <port_number>"
  exit 1
fi

PID=$(lsof -t -i:$PORT)

if [ -z "$PID" ]; then
  echo "No process found running on port $PORT."
else
  echo "Killing process $PID running on port $PORT..."
  kill -9 $PID
  echo "Process $PID killed."
fi