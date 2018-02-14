-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tsconfigs (
  module text,
  target text,
  source_map boolean,
  excluding integer
);

INSERT INTO tsconfigs (module, target, source_map, excluding) VALUES
('admin', 'admins', true, 1),
('user', 'users', true, 1),
('customer', 'customers', true, 100),
('vendor', 'vendors', true, 33),
('admin', 'vendors', true, 0);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tsconfigs;