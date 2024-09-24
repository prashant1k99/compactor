#!/bin/bash

# Ensure GoReleaser is installed
if ! command -v goreleaser &> /dev/null
then
    echo "GoReleaser is not installed. Please install it first."
    exit 1
fi

# Check the configuration
echo "Checking GoReleaser configuration..."
goreleaser check

if [ $? -eq 0 ]; then
    echo "Configuration check passed."
else
    echo "Configuration check failed. Please fix the issues and try again."
    exit 1
fi

# Run a snapshot release
echo "Running a snapshot release..."
# Remove the dist folder if it exists
rm -rf dist

# Run goreleaser
goreleaser release --snapshot

if [ $? -eq 0 ]; then
    echo "Snapshot release completed successfully."
    echo "Check the ./dist directory for the built artifacts."
else
    echo "Snapshot release failed. Please check the output for errors."
    exit 1
fi
