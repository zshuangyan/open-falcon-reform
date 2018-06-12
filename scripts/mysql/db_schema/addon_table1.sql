USE falcon_portal;
SET NAMES utf8;

DROP TABLE IF EXISTS metric;
CREATE TABLE metric
(
  id            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name          VARCHAR(255) NOT NULL,
  command       VARCHAR(500)  NOT NULL DEFAULT '',
  step          INT UNSIGNED NOT NULL DEFAULT 60 COMMENT 'in second',
  metric_type   VARCHAR(10) NOT NULL DEFAULT 'GAUGE' COMMENT 'GAUGE|COUNTER|DERIVE',
  value_type    VARCHAR(10) NOT NULL DEFAULT 'int' COMMENT 'int|float',
  built_in      BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY       KEY (id),
  UNIQUE (name)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

