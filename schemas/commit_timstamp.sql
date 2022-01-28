CREATE TABLE Performances (
    SingerId        INT64 NOT NULL,
    VenueId         INT64 NOT NULL,
    EventDate       Date,
    Revenue         INT64,
    LastUpdateTime  TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (SingerId, VenueId, EventDate),
  INTERLEAVE IN PARENT Singers ON DELETE CASCADE