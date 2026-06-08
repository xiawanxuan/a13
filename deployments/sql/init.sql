CREATE DATABASE IF NOT EXISTS `task_scheduler` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `task_scheduler`;

DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL,
  `type` TINYINT NOT NULL DEFAULT 0 COMMENT '0: cron, 1: onetime',
  `cron_expr` VARCHAR(64) DEFAULT '',
  `payload` TEXT,
  `max_retry_times` INT NOT NULL DEFAULT 3,
  `retry_interval` INT NOT NULL DEFAULT 60 COMMENT 'retry interval in seconds',
  `timeout` INT NOT NULL DEFAULT 300 COMMENT 'timeout in seconds',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0: pending, 1: running, 2: success, 3: failed, 4: paused',
  `retry_times` INT NOT NULL DEFAULT 0,
  `next_execute_time` DATETIME DEFAULT NULL,
  `last_execute_time` DATETIME DEFAULT NULL,
  `worker_id` VARCHAR(64) DEFAULT '',
  `version` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_next_execute_time` (`next_execute_time`),
  KEY `idx_worker_id` (`worker_id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `task_logs`;
CREATE TABLE `task_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `task_id` BIGINT UNSIGNED NOT NULL,
  `task_name` VARCHAR(128) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0: running, 1: success, 2: failed',
  `retry_times` INT NOT NULL DEFAULT 0,
  `worker_id` VARCHAR(64) DEFAULT '',
  `start_time` DATETIME DEFAULT NULL,
  `end_time` DATETIME DEFAULT NULL,
  `duration` BIGINT NOT NULL DEFAULT 0 COMMENT 'duration in seconds',
  `result` TEXT,
  `error_msg` TEXT,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_status` (`status`),
  KEY `idx_worker_id` (`worker_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `tasks` (`name`, `type`, `cron_expr`, `payload`, `max_retry_times`, `retry_interval`, `timeout`, `status`, `next_execute_time`) VALUES
('test-cron-task', 0, '*/5 * * * *', '{"action":"test_cron","data":"hello world"}', 3, 60, 300, 0, DATE_ADD(NOW(), INTERVAL 1 MINUTE)),
('test-onetime-task', 1, '', '{"action":"test_once","data":"one time task"}', 3, 60, 300, 0, NOW());
