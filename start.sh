#!/bin/bash

GREEN='\033[0;32m'
NO_COLOR='\033[0m'

server="server"
counter="counter"
gauge="gauge"
dir="bin"

if [[ $# -ne 1 ]]; then
    echo "Error: Path to the prometheus binary must be supplied to the script"
    exit 1
fi

prometheus=$1

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    # Running on windows; attaching .exe extension.
    server="server.exe"
    counter="counter.exe"
    gauge="gauge.exe"
fi

if [[ ! -d "$dir" ]]; then
    # Create the directory if it doesn't exist
    mkdir "$dir"
fi

# Build the Golang apps
if [[ ! -f $dir/$server ]]; then
  go build -o $dir/$server  main.go
else
  echo "Skipping server building..."
fi
if [[ ! -f $dir/$counter ]]; then
  go build -o $dir/$counter clients/counter.go
else
  echo "Skipping counter client building..."
fi
if [[ ! -f $dir/$gauge ]]; then
  go build -o $dir/$gauge   clients/gauge.go
else
  echo "Skipping gauge client building..."
fi

# Start the server in the background
./"$dir"/"$server" &
echo -e "${GREEN}Server started with PID $!${NO_COLOR}"


# Start the Prometheus server in the background
./"$prometheus" --config.file=prometheus.yml &
echo -e "${GREEN}Prometheus started with PID $!${NO_COLOR}"

sleep 10

# Start counter client in the background
$dir/$counter &
echo -e "${GREEN}Counter client started with PID $!${NO_COLOR}"