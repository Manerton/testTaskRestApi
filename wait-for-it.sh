#!/usr/bin/env bash

# wait-for-it.sh
# Use this script to test if a given TCP host/port are available

set -e

TIMEOUT=15
QUIET=0

echoerr() { if [ "$QUIET" -ne 1 ]; then echo "$@" 1>&2; fi }

usage() {
    cat << USAGE >&2
Usage:
  $0 host:port [-t timeout] [-q]

Options:
  host:port   Host and port to test
  -t timeout  Timeout in seconds (default: $TIMEOUT)
  -q          Quiet mode (suppress output)

USAGE
    exit 1
}

while [[ $# -gt 0 ]]
do
    case "$1" in
        *:* )
        HOST=$(echo "$1" | cut -d ":" -f 1)
        PORT=$(echo "$1" | cut -d ":" -f 2)
        shift 1
        ;;
        -t)
        TIMEOUT="$2"
        shift 2
        ;;
        -q)
        QUIET=1
        shift 1
        ;;
        *)
        usage
        ;;
    esac
done

if [ -z "$HOST" ] || [ -z "$PORT" ]; then
    echo "Error: you need to provide a host and port."
    usage
fi

for i in $(seq 1 "$TIMEOUT"); do
    nc -z "$HOST" "$PORT" >/dev/null 2>&1 && break
    if [ "$i" -eq "$TIMEOUT" ]; then
        echoerr "Error: timeout reached while waiting for $HOST:$PORT"
        exit 1
    fi
    sleep 1
done

echo "Success: $HOST:$PORT is available."
