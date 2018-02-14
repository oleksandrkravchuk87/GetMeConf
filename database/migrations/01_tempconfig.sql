-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tempconfig (
  rest_api_root text,
  host text,
  port text,
  remoting text,
  legasy_explorer boolean
);

INSERT INTO tempconfig (rest_api_root, host, port, remoting, legasy_explorer) VALUES
('', '', '', '', false),
('/', 'localhost', '8080', 'rem', true),
('/home', 'europa', '9080', 'local', true),
('/home', 'asia', '8080', 'local_uk', true);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tempconfig;