-- DEFAULT CHARACTER SET utf8
-- DEFAULT COLLATE utf8_general_ci;
CREATE USER falcon@localhost IDENTIFIED BY '1234';
GRANT ALL PRIVILEGES ON falcon.* TO  falcon@localhost identified by '1234';
GRANT ALL PRIVILEGES ON alarm.* TO  falcon@localhost identified by '1234';
GRANT ALL PRIVILEGES ON idx.* TO  falcon@localhost identified by '1234';
FLUSH PRIVILEGES;
