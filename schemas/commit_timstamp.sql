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

CREATE TABLE Performances (
    SingerId        INT64 NOT NULL,
    VenueId         INT64 NOT NULL,
    EventDate       Date,
    Revenue         INT64,
    LastUpdateTime  TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
) PRIMARY KEY (SingerId, VenueId, EventDate),
  INTERLEAVE IN PARENT Singers ON DELETE CASCADE
