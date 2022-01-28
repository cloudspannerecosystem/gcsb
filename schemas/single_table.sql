/*
 * Single table with a single column primary key
 */
CREATE TABLE SingleSingers (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  BirthDate  DATE,
  ByteField BYTES(1025),
  FloatField FLOAT64,
  ArrayField ARRAY<INT64>, 
  TSField    TIMESTAMP,
  NumericField NUMERIC,
) PRIMARY KEY (SingerId);
