USE falcon_portal;
SET NAMES utf8;

DROP TABLE IF EXISTS metric;
CREATE TABLE metric
(
  id            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name          VARCHAR(255) NOT NULL,
  alias         VARCHAR(255) NOT NULL DEFAULT '',
  command       VARCHAR(500)  NOT NULL DEFAULT '',
  step          INT UNSIGNED NOT NULL DEFAULT 60 COMMENT 'in second',
  metric_type   VARCHAR(10) NOT NULL DEFAULT 'GAUGE' COMMENT 'GAUGE|COUNTER|DERIVE',
  value_type    VARCHAR(10) NOT NULL DEFAULT 'int' COMMENT 'int|float',
  unit          VARCHAR(50) NOT NULL DEFAULT '',
  built_in      BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY       KEY (id),
  UNIQUE (name)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS namespace;
CREATE TABLE namespace
(
  id            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name          VARCHAR(255) NOT NULL,
  PRIMARY       KEY (id),
  UNIQUE (name)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

DROP TABLE IF EXISTS namespace_metric;
CREATE TABLE namespace_metric
(
  namespace_id  INT UNSIGNED NOT NULL,
  metric_id INT UNSIGNED NOT NULL,
  KEY idx_namespace_metric_namespace_id (namespace_id),
  KEY idx_namespace_metric_metric_id (metric_id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

DROP TABLE IF EXISTS host_metric;
CREATE TABLE host_metric
(
  host_id  INT UNSIGNED NOT NULL,
  metric_id INT UNSIGNED NOT NULL,
  KEY idx_host_metric_host_id (host_id),
  KEY idx_host_metric_metric_id (metric_id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;


DROP TABLE IF EXISTS user_defined_metric;
CREATE TABLE user_defined_metric
(
  id            INT UNSIGNED NOT NULL AUTO_INCREMENT,
  name          VARCHAR(255) NOT NULL,
  command       VARCHAR(500)  NOT NULL,
  step          INT UNSIGNED NOT NULL,
  metric_type   VARCHAR(10) NOT NULL,
  value_type    VARCHAR(10) NOT NULL,
  host_id       INT UNSIGNED  NOT NULL,
  status        BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY       KEY (id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

