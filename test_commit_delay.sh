#!/bin/bash

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo "yq is not installed. Please install it first:"
    echo "brew install yq"
    exit 1
fi

# Extract values from YAML file
PROJECT_ID=$(yq '.project' test.yaml)
INSTANCE_ID=$(yq '.instance' test.yaml)
DATABASE_ID=$(yq '.database' test.yaml)

# Validate extracted values
if [ -z "$PROJECT_ID" ] || [ -z "$INSTANCE_ID" ] || [ -z "$DATABASE_ID" ]; then
    echo "Error: Could not extract required values from test.yaml"
    echo "Please ensure project, instance, and database are set in the YAML file"
    exit 1
fi

echo "Using configuration:"
echo "Project: $PROJECT_ID"
echo "Instance: $INSTANCE_ID"
echo "Database: $DATABASE_ID"

# Create a temporary config file without commit delay
sed "s/commitDelay: 100ms/#commitDelay: 100ms/" test.yaml > test_no_delay.yaml

# Create the database
echo "Creating database..."
gcloud spanner databases create $DATABASE_ID \
    --instance=$INSTANCE_ID \
    --database-dialect=GOOGLE_STANDARD_SQL

# Extract and apply the schema
echo "Applying schema..."
SCHEMA=$(yq '.schema' test.yaml | sed 's/^|//' | sed 's/^  //')
gcloud spanner databases ddl update $DATABASE_ID \
    --instance=$INSTANCE_ID \
    --ddl="$SCHEMA"

# Run test without commit delay
echo "Running test without commit delay..."
./gcsb --config test_no_delay.yaml

# Run test with commit delay
echo "Running test with commit delay..."
./gcsb --config test.yaml

# Clean up
echo "Cleaning up..."
rm test_no_delay.yaml
gcloud spanner databases delete $DATABASE_ID \
    --instance=$INSTANCE_ID \
    --quiet 