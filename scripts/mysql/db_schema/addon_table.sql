USE falcon_portal;
SET NAMES utf8;

/**
 */
DROP TABLE IF EXISTS user_defined_metrics;
CREATE TABLE user_defined_metrics
(
  id             INT UNSIGNED NOT NULL AUTO_INCREMENT,
  metric_name    VARCHAR(255) NOT NULL DEFAULT '',
  command     VARCHAR(500)  NOT NULL DEFAULT '',
  host_id  INT UNSIGNED  NOT NULL,
  status BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8
  COLLATE =utf8_unicode_ci;

