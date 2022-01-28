CREATE VIEW SingerNames
SQL SECURITY INVOKER
AS SELECT
  Singers.SingerId AS SingerId,
  Singers.FirstName || " " || Singers.LastName AS Name
FROM Singers;