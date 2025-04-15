#!/bin/bash

# Available configurations
declare -A CONFIGS=(
    ["1"]="single-region:us-central1:STANDARD:1"
    ["2"]="single-region:us-central1:ENTERPRISE:1"
    ["3"]="multi-region:nam3:STANDARD:1"
    ["4"]="multi-region:nam3:ENTERPRISE:1"
    ["5"]="dual-region:australia-southeast1+australia-southeast2:STANDARD:1"
    ["6"]="dual-region:australia-southeast1+australia-southeast2:ENTERPRISE:1"
)

# Function to create a Spanner instance
create_instance() {
    local name=$1
    local config=$2
    local display_name=$3
    local processing_units=$4

    echo "Creating instance: $name"
    gcloud spanner instances create $name \
        --config=$config \
        --description="Test instance for $display_name" \
        --processing-units=$processing_units \
        --project=$PROJECT_ID
}

# Function to create a test database
create_database() {
    local instance=$1
    local database=$2

    echo "Creating database: $database in instance: $instance"
    gcloud spanner databases create $database \
        --instance=$instance \
        --database-dialect=GOOGLE_STANDARD_SQL \
        --project=$PROJECT_ID

    # Apply the test schema
    echo "Applying schema to database: $database"
    gcloud spanner databases ddl update $database \
        --instance=$instance \
        --ddl="$(cat <<EOF
CREATE TABLE TestTable (
    ID STRING(36) NOT NULL,
    Name STRING(MAX),
    Value INT64,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
    UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (ID);

CREATE INDEX TestTable_Value ON TestTable(Value);
EOF
)" \
        --project=$PROJECT_ID
}

# Function to display available configurations
show_configs() {
    echo "Available configurations:"
    echo "------------------------"
    for key in "${!CONFIGS[@]}"; do
        IFS=':' read -r type region version nodes <<< "${CONFIGS[$key]}"
        echo "$key) $type - $region - $version - $nodes node(s)"
    done
    echo "------------------------"
}

# Main script
echo "Spanner Instance Creation Script"
echo "==============================="

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "gcloud is not installed. Please install the Google Cloud SDK first."
    exit 1
fi

# Check if user is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q "@"; then
    echo "Please authenticate with gcloud first:"
    echo "gcloud auth login"
    exit 1
fi

# Get project ID if not set
if [ "$PROJECT_ID" == "your-project-id" ]; then
    echo "Current project: $(gcloud config get-value project)"
    read -p "Enter project ID (press Enter to use current): " input_project
    if [ ! -z "$input_project" ]; then
        PROJECT_ID=$input_project
    else
        PROJECT_ID=$(gcloud config get-value project)
    fi
fi

echo "Using project: $PROJECT_ID"

# Show available configurations
show_configs

# Get user selection
read -p "Enter configuration number(s) to create (comma-separated, e.g., 1,3,5): " selected_configs

# Process selected configurations
IFS=',' read -ra config_nums <<< "$selected_configs"
for num in "${config_nums[@]}"; do
    if [ -z "${CONFIGS[$num]}" ]; then
        echo "Invalid configuration number: $num"
        continue
    fi

    config="${CONFIGS[$num]}"
    IFS=':' read -r type region version nodes <<< "$config"
    instance_name="test-${type}-${version}-$(date +%Y%m%d)"
    database_name="testdb-${type}-${version}"
    
    echo "Creating $type instance with $version tier..."
    create_instance "$instance_name" "$region" "$type-$version" "$nodes"
    
    # Wait for instance to be ready
    echo "Waiting for instance to be ready..."
    sleep 30
    
    create_database "$instance_name" "$database_name"
    
    echo "Created instance: $instance_name with database: $database_name"
    echo "Configuration:"
    echo "  Type: $type"
    echo "  Region: $region"
    echo "  Version: $version"
    echo "  Nodes: $nodes"
    echo "----------------------------------------"
done

echo "All selected instances and databases created successfully!"
echo "You can now use these instances for load testing."
echo "Example gcsb command:"
echo "./gcsb --config test.yaml --project $PROJECT_ID --instance test-single-region-STANDARD-$(date +%Y%m%d) --database testdb-single-region-STANDARD" 