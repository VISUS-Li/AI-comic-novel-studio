/*
 * ScriptStudio base tables to support ScriptProject, WorkflowInstance, WorkflowRun, NodeOutput and Asset.
 */

CREATE TABLE IF NOT EXISTS `script_projects` (
  `id` VARCHAR(64) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `type` VARCHAR(64) NOT NULL DEFAULT 'novel',
  `owner_id` VARCHAR(64) NOT NULL,
  `description` TEXT,
  `worldbuilding` JSON,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_script_projects_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflow_instances` (
  `id` VARCHAR(64) NOT NULL,
  `project_id` VARCHAR(64) NOT NULL,
  `workflow_id` VARCHAR(64) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `mode` VARCHAR(64) NOT NULL,
  `last_context` JSON,
  `status` VARCHAR(32) NOT NULL DEFAULT 'active',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_workflow_instances_project_id` (`project_id`),
  KEY `idx_workflow_instances_workflow_id` (`workflow_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `workflow_runs` (
  `id` VARCHAR(64) NOT NULL,
  `instance_id` VARCHAR(64) NOT NULL,
  `started_at` DATETIME(3) NOT NULL,
  `finished_at` DATETIME(3),
  `status` VARCHAR(32) NOT NULL,
  `input_context` JSON,
  `output_context` JSON,
  `error_message` TEXT,
  PRIMARY KEY (`id`),
  KEY `idx_workflow_runs_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `node_outputs` (
  `id` VARCHAR(64) NOT NULL,
  `run_id` VARCHAR(64) NOT NULL,
  `node_id` VARCHAR(128) NOT NULL,
  `output_data` JSON,
  `asset_ids` JSON,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_node_outputs_run_id` (`run_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `assets` (
  `id` VARCHAR(64) NOT NULL,
  `project_id` VARCHAR(64) NOT NULL,
  `workflow_instance_id` VARCHAR(64),
  `node_id` VARCHAR(128),
  `type` VARCHAR(32) NOT NULL,
  `url` TEXT NOT NULL,
  `text_content` TEXT,
  `meta` JSON,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_assets_project_id` (`project_id`),
  KEY `idx_assets_workflow_instance_id` (`workflow_instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
