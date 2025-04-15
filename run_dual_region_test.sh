#!/bin/bash

# Check if gcloud is installed and authenticated
if ! command -v gcloud &> /dev/null; then
    echo "Google Cloud SDK is not installed. Please install it first."
    exit 1
fi

if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
    echo "Please authenticate with gcloud first:"
    echo "gcloud auth login"
    exit 1
fi

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo "yq is not installed. Installing yq..."
    # For macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if ! command -v brew &> /dev/null; then
            echo "Homebrew is not installed. Please install Homebrew first:"
            echo "/bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            exit 1
        fi
        brew install yq
    # For Linux
    else
        sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
        sudo chmod a+x /usr/local/bin/yq
    fi
fi

# Get project ID
PROJECT_ID=$(gcloud config get-value project)
if [ -z "$PROJECT_ID" ]; then
    echo "No project ID set. Please set your project ID:"
    echo "gcloud config set project YOUR_PROJECT_ID"
    exit 1
fi

echo "Using project: $PROJECT_ID"

# Create the dual-region instance with a valid name
echo "Creating instance..."
INSTANCE_NAME="test-dual-region-$(date +%Y%m%d)"
DB_NAME="testdb-dual-region"

# Check if instance exists
if gcloud spanner instances describe $INSTANCE_NAME --project=$PROJECT_ID &>/dev/null; then
    echo "Instance $INSTANCE_NAME already exists. Using existing instance."
else
    echo "Creating new instance $INSTANCE_NAME..."
    gcloud spanner instances create $INSTANCE_NAME \
        --config=dual-region-australia1 \
        --description="Dual-region test" \
        --processing-units=100 \
        --project=$PROJECT_ID \
        --edition=ENTERPRISE_PLUS

    # Wait for instance to be ready
    echo "Waiting for instance to be ready..."
    sleep 30
fi

# Check if database exists
if gcloud spanner databases describe $DB_NAME --instance=$INSTANCE_NAME --project=$PROJECT_ID &>/dev/null; then
    echo "Database $DB_NAME already exists. Using existing database."
else
    echo "Creating database $DB_NAME..."
    gcloud spanner databases create $DB_NAME \
        --instance=$INSTANCE_NAME \
        --database-dialect=GOOGLE_STANDARD_SQL \
        --project=$PROJECT_ID

    # Extract and validate schema
    echo "Extracting schema from YAML..."
    SCHEMA_DDL=$(yq '.schema' dual_region_test.yaml)
    if [ -z "$SCHEMA_DDL" ]; then
        echo "Error: Could not extract schema from YAML file"
        exit 1
    fi

    # Validate schema syntax
    echo "Validating schema syntax..."
    if ! echo "$SCHEMA_DDL" | grep -q "CREATE TABLE"; then
        echo "Error: Invalid schema - no CREATE TABLE statements found"
        exit 1
    fi

    # Apply the schema
    echo "Applying schema..."
    if ! gcloud spanner databases ddl update $DB_NAME \
        --instance=$INSTANCE_NAME \
        --ddl="$SCHEMA_DDL" \
        --project=$PROJECT_ID; then
        echo "Error: Failed to apply schema"
        exit 1
    fi
fi

# Update the YAML file with the correct instance and database names
echo "Updating YAML configuration..."
yq -i ".instance = \"$INSTANCE_NAME\"" dual_region_test.yaml
yq -i ".database = \"$DB_NAME\"" dual_region_test.yaml
yq -i ".project = \"$PROJECT_ID\"" dual_region_test.yaml

# Verify schema was applied
echo "Verifying schema..."
if ! gcloud spanner databases execute-sql $DB_NAME \
    --instance=$INSTANCE_NAME \
    --project=$PROJECT_ID \
    --sql="SELECT * FROM information_schema.tables WHERE table_name = 'TestTable'" | grep -q "TestTable"; then
    echo "Error: Schema verification failed - TestTable not found"
    exit 1
fi

# Run initial data load
echo "Loading initial data..."
if ! ./gcsb load --config dual_region_test.yaml -t TestTable; then
    echo "Error: Failed to load initial data"
    exit 1
fi

# Run the benchmark
echo "Starting benchmark..."
RESULTS_FILE="benchmark_results_$(date +%Y%m%d_%H%M%S).txt"
if ! ./gcsb run --config dual_region_test.yaml -t TestTable | tee $RESULTS_FILE; then
    echo "Error: Benchmark failed"
    exit 1
fi

echo "Test completed successfully!"
echo "Instance: $INSTANCE_NAME"
echo "Database: $DB_NAME"
echo "Results saved to: $RESULTS_FILE"
echo "You can view the results in the Cloud Console:"
echo "https://console.cloud.google.com/spanner/instances/$INSTANCE_NAME/databases/$DB_NAME/overview?project=$PROJECT_ID" 