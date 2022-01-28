/*
Copyright 2022 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
 * Multiple tables with interleaved relationships and indexes
 */
CREATE TABLE Singers (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  BirthDate  DATE,
  ByteField BYTES(1025),
) PRIMARY KEY (SingerId);

CREATE INDEX SingersByFirstLastName ON Singers(FirstName, LastName);

CREATE TABLE Albums (
  SingerId     INT64 NOT NULL,
  AlbumId      INT64 NOT NULL,
  AlbumTitle   STRING(MAX),
  ReleaseDate  DATE,
  AlbumRating  NUMERIC,
) PRIMARY KEY (SingerId, AlbumId),
  INTERLEAVE IN PARENT Singers ON DELETE CASCADE;

CREATE INDEX AlbumsByAlbumTitle ON Albums(SingerId, AlbumTitle)
STORING (ReleaseDate),
  INTERLEAVE IN Singers;

CREATE TABLE Songs (
  SingerId   INT64 NOT NULL,
  AlbumId    INT64 NOT NULL,
  TrackId    INT64 NOT NULL,
  SongName   STRING(MAX),
  Duration   INT64,
) PRIMARY KEY (SingerId, AlbumId, TrackId),
  INTERLEAVE IN PARENT Albums ON DELETE CASCADE;

CREATE INDEX SongsBySingerAlbumSongNameDesc ON Songs(SingerId, AlbumId, SongName DESC),
  INTERLEAVE IN Albums;

CREATE TABLE Venues (
  VenueId       INT64 NOT NULL,
  VenueName     STRING(MAX),
  VenueCity     STRING(MAX),
  Capacities    ARRAY<INT64>, 
  TotalCapacity INT64
) PRIMARY KEY (VenueId);

CREATE TABLE Concerts (
  VenueId      INT64 NOT NULL,
  SingerId     INT64 NOT NULL,
  ConcertDate  DATE NOT NULL,
  BeginTime    TIMESTAMP,
  EndTime      TIMESTAMP,
  TicketPrices ARRAY<FLOAT64>,

  CONSTRAINT FKConcertsSingerId
    FOREIGN KEY (SingerId)
    REFERENCES Singers (SingerId)

) PRIMARY KEY (VenueId, SingerId, ConcertDate),
  INTERLEAVE IN PARENT Venues ON DELETE CASCADE;

CREATE INDEX ConcertsBySingerId ON Concerts(SingerId);
