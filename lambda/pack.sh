#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

# Set environment variables
export GOOS=linux
export GOARCH=amd64

# require 1 cmd arg, specifying the filename of the Go program
if [ $# -ne 1 ]; then
  echo "Usage: $0 <filename>"
  exit 1
fi

# strip the .go extension from the filename, if present
if [[ $1 == *".go" ]]; then
  filename=$(echo $1 | cut -f 1 -d '.')
else
  filename=$1
fi

# Build the Go program
go build -o bootstrap $filename.go

# Check if the build was successful
if [ $? -eq 0 ]; then
  # Create a zip file containing the bootstrap binary
  zip $filename\_packed.zip bootstrap
  echo "Build and zip successful."
else
  echo "Build failed."
fi

# Clean up
rm bootstrap