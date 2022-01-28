# GCSB

- [GCSB](#gcsb)
  - [Quickstart](#quickstart)
    - [Create a test table](#create-a-test-table)
    - [Load data into table](#load-data-into-table)
    - [Run a load test](#run-a-load-test)
  - [Operations](#operations)
    - [Load](#load)
      - [Single table load](#single-table-load)
      - [Mulitiple table load](#mulitiple-table-load)
      - [Loading into interleaved tables](#loading-into-interleaved-tables)
    - [Run](#run)
      - [Single table run](#single-table-run)
      - [Mulitiple table run](#mulitiple-table-run)
      - [Running against interleaved tables](#running-against-interleaved-tables)
  - [Distributed testing](#distributed-testing)
  - [Configuration](#configuration)
  - [Roadmap](#roadmap)
    - [Not Supported (yet)](#not-supported-yet)
  - [Development](#development)
    - [Build](#build)
    - [Test](#test)

It's like YCSB but with more Google. A simple tool meant to generate load against Google Spanner databases. The primary goals of the project are

- Write randomized data to tables in order to
  - Facilitate load testing
  - Initiate [database splits](https://cloud.google.com/spanner/docs/schema-and-data-model#database-splits)
- Generate read/write load against user provided schemas

## Quickstart

To initiate a simple load test against your spanner instance using one of our test schemas

### Create a test table

You can use your own schema if you'd prefer, but we provide a few test schemas to help you get started. To get started, create a table named `SingleSingers`

```sh
gcloud spanner databases ddl update YOUR_DATABASE_ID --instance=YOUR_INSTANCE_ID --ddl-file=schemas/single_table.sql
```

### Load data into table

Load some data into the table to seed the upcoming laod test. In the example below, we are loading 10,000 rows of random data into the table `SingleSingers`

```sh
gcsb load -p YOUR_GCP_PROJECT_ID -i YOUR_INSTANCE_ID -d YOUR_DATABASE_ID -t SingleSingers -o 10000
```

### Run a load test

Now you can perform a load test operation by using the `run` sub command. The below command will generate 10,000 operations. Of these operations, 75% will be READ operations and 25% will be writes. These operations will be performed over 50 threads.

```sh
gcsb run -p YOUR_GCP_PROJECT_ID -i YOUR_INSTANCE_ID -d YOUR_DATABASE_ID -t SingleSingers -o 10000 --reads 75 --writes 25 --threads 50
```

## Operations

The tool usage is generally broken down into two categories, `load` and `run` operations.

### Load

GCSB proides batched loading functionality to facilitate load testing as well as assist with performing [database splits](https://cloud.google.com/spanner/docs/schema-and-data-model#database-splits) on your tables.

At runtime, GCSB will detect the schema of the table you're loading data for an create data generators that are appropriate for the column types you have in your database. Each type of generator has some configurable funcationality that allows you to refine the type, length, or range of data the tool generates. For in depth information on the various configuration values, please read the comments in [example_gcsb.yaml](example_gcsb.yaml).

#### Single table load

By default, GCSB will detect the table schema and create default random data generators based on the columns it finds. In order to tune the values the generator creates, you must create override configurations in the gcsb.yaml file. Please see that file's documentation for more information.

```sh
gcsb load -t TABLE_NAME -o NUM_ROWS
```

Additionally, please see `gcsb load --help` for additional configuration options.

#### Mulitiple table load

Similar to the above [Single table load](#single-table-load), you may specify multiple tables by repeating the `-t TABLE_NAME` argument. By default, the number of operations is applied to each table. For example, specifying 2 tables with 1000 operations, will yield 2000 total operations. 1000 per table.

```sh
gcsb load -t TABLE1 -t TABLE2 -o NUM_ROWS
```

Operations per table can be configured in the yaml configuration. For example

```yaml
tables:
  - name: TABLE1
    operations:
      total: 500
  - name: TABLE2
    operations:
      total: 500
```

#### Loading into interleaved tables

Loading data into interleaved tables is supported but has some behavioral side effects that should be known. When loading data into an interleaved table, GCSB will detect all tables in the hierarchy and begin loading data at the familial apex. The configured number of operations applies to this apex table. By default, the number of tables for each child table is multiplied by 5. For example:

Using our [test INTERLEAVE schema](schemas/multi_table.sql), we see an INTERLEAVE relationship between the `Singers`, `Albums`, and `Songs` tables.

If we execute a load operation against these tables with total operations set to `10` we will see the following occur

```sh
gcsb load -t Songs -o 10

+---------+------------+------+-------+---------+
|  TABLE  | OPERATIONS | READ | WRITE | CONTEXT |
+---------+------------+------+-------+---------+
| Singers |         10 | N/A  | N/A   | LOAD    |
| Albums  |         50 | N/A  | N/A   | LOAD    |
| Songs   |        250 | N/A  | N/A   | LOAD    |
+---------+------------+------+-------+---------+
```

In this case, for each child table we take the number of operations for the parent and multiply it by the default value of `5`.

To change this multiplier, we use the yaml configuration file for the table we want. The `operations.total` value becomes a multiplier. 

```yaml
tables:
  - name: Albums
    operations:
      total: 10
  - name: Songs
    operations:
      total: 20
```

At present, GCSB will sort it's operations from the apex down. Meaning it will populate the `Singers` table first and then it's child, and then the next child. Multiple table operations are not mixed within the same transaction. 

### Run

#### Single table run

By default, GCSB will detect the table schema and create default random data generators based on the columns it finds. In order to tune the values the generator creates, you must create override configurations in the gcsb.yaml file. Please see that file's documentation for more information.

```sh
gcsb run -p YOUR_GCP_PROJECT_ID -i YOUR_INSTANCE_ID -d YOUR_DATABASE_ID -t SingleSingers -o 10000 --reads 75 --writes 25 --threads 50
```

Additionally, please see `gcsb run --help` for additional configuration options.

#### Mulitiple table run

Similar to the above [Single table run](#single-table-run), you may specify multiple tables by repeating the `-t TABLE_NAME` argument. By default, the number of operations is applied to each table. For example, specifying 2 tables with 1000 operations, will yield 2000 total operations. 1000 per table.

```sh
gcsb run -t TABLE1 -t TABLE2 -o NUM_ROWS
```

Operations per table can be configured in the yaml configuration. For example

```yaml
tables:
  - name: TABLE1
    operations:
      total: 500
  - name: TABLE2
    operations:
      total: 500
```

#### Running against interleaved tables

Run operations against interleaved tables are only supported on the APEX table.

Using our [test INTERLEAVE schema](schemas/multi_table.sql), we see an INTERLEAVE relationship between the `Singers`, `Albums`, and `Songs` tables.

You will notice that if we try to run against any child table an error will occur.

```sh
gcsb run -t Songs -o 10

unable to execute run operation: can only execute run against apex table (try 'Singers')
```

## Distributed testing

GCSB is intended to run in a stateless mannger. This design choice was to allow massive horizontal scaling of gcsb to stress your database to it's absolute limits. During development we've identified kubernetes as the prefered tool for the job. We've provided two separate tutorials for running gcsb inside of kubernetes

- [GKE](docs/GKE.md) - For running GCSB inside GKE using a service account key. This can be used for non-GKE clusters as well as it contains instructions for mounting a service account key into the container.
- [GKE with Workload Identity](docs/workload_identity.md) - For running GCSB inside a GKE cluster that has workload identity turned on. This is most useful in organizations that have security policies preventing you from generating or downloading a service account key.

## Configuration

The tool can receive configuration input in several different ways. The tool will load the file `gcsb.yaml` if it detects it in the current working directory. Alternatively you can use the global flag `-c` to specify a path to the configuration file. Each sub-command has a number of configuration flags that are relevant to that operation. These values are bound to their counterparts in the yaml configuration file and take precedent over the config file. Think of them as overrides. The same is true for environment variables.

Please note, at present, the yaml conifguration file is the only way to specify generator overrides for data loading and write operations. Without this file, the tool will use a random data generator that is appropriate for the table schema it detects at runtime.

For in depth information on the various configuration values, please read the comments in [example_gcsb.yaml](example_gcsb.yaml)

## Roadmap

### Not Supported (yet)

- [ ] Generating read operations utilizing [ReadByIndex](https://cloud.google.com/spanner/docs/samples/spanner-read-data-with-index#spanner_read_data_with_index-go)
- [ ] Generating NULL values for load operations. If a column is NULLable, gcsb will still generate a value for it.
- [x] ~~JSON column types~~
- [ ] STRUCT Objects.
- [ ] VIEWS issue #26 is WIP to support 100% READ tests against views but is not currently functioning
- [x] ~~INTERLEAVE tables. Support for INTERLEAVE is currently WIP, and it will only support apex test targets~~
  - [ ] INTERLEAVE talbes are now supported for both LOAD and RUN. RUN operations only support the table at the APEX of the hierarchy
- [ ] Inserting data across multiple tables in the same transaction
- [ ] No SCAN or DELETE operations are supported at this time
- [ ] Tables with foreign key relationships
- [ ] Testing multiple tables at once
- [x] ~~Tables containing columns with `allow_commit_timestamp=TRUE`~~

## Development

### Build

```sh
make build
```

### Test

```sh
make test
```
