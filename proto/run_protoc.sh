#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(realpath "$0")")

# Run the protoc command with the proper absolute paths
protoc -I="$SCRIPT_DIR" \
  -I="$SCRIPT_DIR/google/api" \
  -I="$SCRIPT_DIR/validate" \
  --go_out="$SCRIPT_DIR" \
  --go-grpc_out="$SCRIPT_DIR" \
  "$SCRIPT_DIR/tdtm.proto"

# Check if the command was successful
if [ $? -eq 0 ]; then
    echo "Protobuf generation successful!"
else
    echo "Protobuf generation failed!"
    exit 1
fi
