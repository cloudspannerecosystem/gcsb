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
