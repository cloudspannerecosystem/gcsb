/*
 * MaxBounds is a test table for ensuring MAX is properly converted to 1024
 *
 * Data generation for max is extremely expensive. Thusly we clip to a lower size. 
 * As an example, MAX (~2.6MB) @ 1000 rows means gcsb needs to generate 2.6GB of string data.
 * 
 * Strings: https://cloud.google.com/spanner/docs/data-definition-language#string
 * Bytes: https://cloud.google.com/spanner/docs/data-definition-language#bytes
 */
CREATE TABLE MaxBounds (
  SingerId   INT64 NOT NULL,
  StringField  STRING(MAX),
  ByteFieldMax BYTES(MAX),
) PRIMARY KEY (SingerId);