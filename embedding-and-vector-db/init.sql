-- INIT EXTENSION

CREATE EXTENSION IF NOT EXISTS vector;

-- INIT TABLE

CREATE TABLE IF NOT EXISTS points (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  color TEXT NOT NULL,
  point2d vector(2) NOT NULL
);

-- INSERT DATA

INSERT INTO points (name, color, point2d) VALUES ('A', 'brown', '[1,1]');
INSERT INTO points (name, color, point2d) VALUES ('B', 'orange', '[2,1]');
INSERT INTO points (name, color, point2d) VALUES ('C', 'red', '[-3,3]');
INSERT INTO points (name, color, point2d) VALUES ('D', 'green', '[3,1]');
INSERT INTO points (name, color, point2d) VALUES ('E', 'purple', '[-2,2]');

-- EUCLIDEAN DISTANCE

SELECT *, (point2d <-> '[2,2]') AS distance
FROM points
ORDER BY distance;

-- COSINE SIMILARITY

SELECT *, (point2d <=> '[2,2]') AS distance
FROM points
ORDER BY distance;