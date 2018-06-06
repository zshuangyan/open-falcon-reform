USE falcon_portal;
SET NAMES utf8;

/**
 */
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

