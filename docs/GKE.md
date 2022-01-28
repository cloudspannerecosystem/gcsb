# Distributed load testing with GKE

- [Distributed load testing with GKE](#distributed-load-testing-with-gke)
  - [Setup Environment](#setup-environment)
    - [Create a GCP Project](#create-a-gcp-project)
    - [Enable services](#enable-services)
    - [Create a service account](#create-a-service-account)
    - [Grant access to spanner](#grant-access-to-spanner)
    - [Create key for service account](#create-key-for-service-account)
  - [Setup Spanner Database](#setup-spanner-database)
    - [Create a Spanner Instance](#create-a-spanner-instance)
    - [Create a database](#create-a-database)
  - [Setup GKE](#setup-gke)
    - [Create GKE Cluster](#create-gke-cluster)
    - [Import Service Account Key](#import-service-account-key)
    - [Build Docker Container](#build-docker-container)
  - [Run the tool](#run-the-tool)
    - [CLI Only](#cli-only)
      - [Single instance load operation](#single-instance-load-operation)
      - [Multi instance load operation](#multi-instance-load-operation)
      - [Single instance run operation](#single-instance-run-operation)
      - [Multi instance run operation](#multi-instance-run-operation)
    - [Custom Configuration](#custom-configuration)
      - [Create Configmap](#create-configmap)
      - [Multi instance load operation](#multi-instance-load-operation-1)
      - [Multi instance run operation](#multi-instance-run-operation-1)
  - [Troubleshooting](#troubleshooting)
    - [Kubectl errors](#kubectl-errors)

> **NOTE** The below instructions to  for multi instance load operations deploy as a [kubernetes deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).  This means that kubernetes will continuously restart your load or run instances until you manually stop them. It is important that you do not leave them running indefinitely.

## Setup Environment

```sh
export PROJECT_ID=spanner-test
export SPANNER_INSTANCE=test-instance
export SPANNER_DATABASE=test-database
export GKE_CLUSTER_NAME=test-cluster
export GCP_REGION=us-west2
```

### Create a GCP Project

```sh
gcloud projects create $PROJECT_ID
```

### Enable services

```sh
gcloud services enable spanner.googleapis.com --project $PROJECT_ID
gcloud services enable cloudbuild.googleapis.com --project $PROJECT_ID
gcloud services enable container.googleapis.com --project $PROJECT_ID
gcloud services enable artifactregistry.googleapis.com --project $PROJECT_ID
```

### Create a service account

```sh
gcloud iam service-accounts create gcsb-test-sa \
    --description="GCSB Test Account" \
    --display-name="gcsb" \
    --project $PROJECT_ID
```

### Grant access to spanner

```sh
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:gcsb-test-sa@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/spanner.databaseUser"
```

### Create key for service account

```sh
gcloud iam service-accounts keys create key.json --iam-account=gcsb-test-sa@${PROJECT_ID}.iam.gserviceaccount.com
```

## Setup Spanner Database

### Create a Spanner Instance

```sh
gcloud alpha spanner instances create $SPANNER_INSTANCE --config=regional-us-east1 --processing-units=1000 --project $PROJECT_ID
```

### Create a database

```sh
gcloud spanner databases create $SPANNER_DATABASE --instance=$SPANNER_INSTANCE --project $PROJECT_ID
```

## Setup GKE

### Create GKE Cluster

Please see `gcloud compute machine-types list` for a list of machine types

```sh
gcloud container clusters create $GKE_CLUSTER_NAME \
  --project $PROJECT_ID \
  --region $GCP_REGION \
  --machine-type=n1-standard-8 \
  --num-nodes 3
```

### Import Service Account Key

```sh
kubectl create secret generic gcsb-sa-key --from-file=key.json=./key.json
```

### Build Docker Container

```sh
gcloud builds submit --tag gcr.io/$PROJECT_ID/gcsb .
```

## Run the tool

### CLI Only

To run inside Kubernetes, you will need to mount the ServiceAccount key secret into the container. Additionally, you will need to set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` pointing to this key file. Below, there are instructions for single container launches versus parallel launches.

#### Single instance load operation

If you want to simply run a single instance of gcsb to perform a load operation, you can launch a POD via the `kubectl run` command.

```sh
kubectl run gcsb-load \
  --image=gcr.io/$PROJECT_ID/gcsb \
  --replicas=1 \
  --restart=Never \
  --overrides='{
     "apiVersion": "v1",
     "spec": {
        "containers": [{
          "name": "gcsb-load",
          "image": "gcr.io/'$PROJECT_ID'/gcsb",
          "command": [ "/gcsb", "load" ],
          "args": [ "--project='$PROJECT_ID'", "--instance='$SPANNER_INSTANCE'", "--database='$SPANNER_DATABASE'", "--table=SingleSingers", "--operations=1000000"],
          "volumeMounts": [{"mountPath": "/var/secrets/google", "name": "google-cloud-key"}],
          "env": [{ "name": "GOOGLE_APPLICATION_CREDENTIALS", "value": "/var/secrets/google/key.json" }]
        }],
        "volumes": [ { "name": "google-cloud-key", "secret": { "secretName": "gcsb-sa-key" } } ]
     }
  }'
```

#### Multi instance load operation

To create a load operation named 'gcsb-load', you must edit the [gke_load.yaml](gke_load.yaml) file, supplying your spanner information.

For example, everywhere you se the comment `EDIT:` you must specify your information.

```yaml
          - --project=YOUR_PROJECT_ID   # EDIT: Your GCP Project ID
          - --instance=YOUR_INSTANCE_ID # EDIT: Your Spanner Instance ID
          - --database=YOUR_DATABASE    # EDIT: Your Spanner Database Name
          - --table=YOUR_TABLE          # EDIT: Your Table Name
          - --operations=1000000        # EDIT: Number of Operations
          - --threads=10                # EDIT: Number of Threads
```

Once you have completed the necessary file edits, run the following.

```sh
kubectl apply -f docs/gke_load.yaml
```

to stop the test

```sh
kubectl delete deploy gcsb-load
```

#### Single instance run operation

If you want to simply run a single instance of gcsb to perform a run operation, you can launch a POD via the `kubectl run` command.

```sh
kubectl run gcsb-load \
  --image=gcr.io/$PROJECT_ID/gcsb \
  --replicas=1 \
  --restart=Never \
  --overrides='{
     "apiVersion": "v1",
     "spec": {
        "containers": [{
          "name": "gcsb-load",
          "image": "gcr.io/'$PROJECT_ID'/gcsb",
          "command": [ "/gcsb", "run" ],
          "args": [ "--project='$PROJECT_ID'", "--instance='$SPANNER_INSTANCE'", "--database='$SPANNER_DATABASE'", "--table=SingleSingers", "--operations=1000000"],
          "volumeMounts": [{"mountPath": "/var/secrets/google", "name": "google-cloud-key"}],
          "env": [{ "name": "GOOGLE_APPLICATION_CREDENTIALS", "value": "/var/secrets/google/key.json" }]
        }],
        "volumes": [ { "name": "google-cloud-key", "secret": { "secretName": "gcsb-sa-key" } } ]
     }
  }'
```

#### Multi instance run operation

To create a load operation named 'gcsb-run', you must edit the [gke_run.yaml](gke_run.yaml) file, supplying your spanner information.

For example, everywhere you se the comment `EDIT:` you must specify your information.

```yaml
          - --project=YOUR_PROJECT_ID   # EDIT: Your GCP Project ID
          - --instance=YOUR_INSTANCE_ID # EDIT: Your Spanner Instance ID
          - --database=YOUR_DATABASE    # EDIT: Your Spanner Database Name
          - --table=YOUR_TABLE          # EDIT: Your Table Name
          - --operations=1000000        # EDIT: Number of Operations
          - --threads=10                # EDIT: Number of Threads
          - --reads=50                  # EDIT: Read Weight (Example: 50 = 50% reads)
          - --writes=50                 # EDIT: Write Weight (Example: 50 = 50% writes)
          - --sample-size=5             # EDIT: Percentage of table to sample for generating reads (Example: 5 = 5% of the rows in the table)
```

Once you have completed the necessary file edits, run the following.

```sh
kubectl apply -f docs/gke_run.yaml
```

to stop the test

```sh
kubectl delete deploy gcsb-run
```

### Custom Configuration

If you prefer to mount your custom gcsb file into the container, you should follow these instructions

#### Create Configmap

The below example assumes your conig file is named `gcsb.yaml`

```sh
kubectl create configmap gcsb-config --from-file=gcsb.yaml
```

#### Multi instance load operation

To create a load operation named 'gcsb-load', you must edit the [gke_load_custom.yaml](gke_load_custom.yaml) file, supplying your spanner information.

For example, everywhere you se the comment `EDIT:` you must specify your information.

```yaml
          - --table=YOUR_TABLE              # EDIT: Your Table Name
```

Once you have completed the necessary file edits, run the following.

```sh
kubectl apply -f docs/gke_load_custom.yaml
```

to stop the test

```sh
kubectl delete deploy gcsb-load
```

#### Multi instance run operation

To create a load operation named 'gcsb-run', you must edit the [gke_run_custom.yaml](gke_run_custom.yaml) file, supplying your spanner information.

For example, everywhere you se the comment `EDIT:` you must specify your information.

```yaml
          - --table=YOUR_TABLE              # EDIT: Your Table Name
```

Once you have completed the necessary file edits, run the following.

```sh
kubectl apply -f docs/gke_run_custom.yaml
```

to stop the test

```sh
kubectl delete deploy gcsb-run
```

## Troubleshooting

### Kubectl errors

`SchemaError(io.k8s.api.autoscaling.v2beta2.MetricTarget): invalid object doesn't have additional properties`

There are issues with some kubectl installations from homebrew. Please relink your kubectl installation by following [these instructions](https://stackoverflow.com/a/55564032/145479).
