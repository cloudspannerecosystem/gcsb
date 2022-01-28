/*
 * Single table with a composite primary key
 */
CREATE TABLE SingleSingersMultiKey (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  BirthDate  DATE,
  ByteField BYTES(1025),
  FloatField FLOAT64,
  TSField    TIMESTAMP,
  NumericField NUMERIC,
) PRIMARY KEY (SingerId, FirstName);
