-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mongodb (
  domain  TEXT,
  mongodb BOOLEAN,
  host    TEXT,
  port    TEXT
);

INSERT INTO mongodb (domain, mongodb, host, port) VALUES
('mydom', TRUE, 'localhost', '8080'),
('testdom', TRUE, '127.0.0.1', '9090'),
('remote', TRUE, '227.255.255.1', '8090'),
('asia', TRUE, '217.155.155.1', '8081');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE mongodb;