CREATE TABLE IF NOT EXISTS `{{ index .Options "Namespace" }}accounts` (
  `instance_id` varchar(255) DEFAULT NULL,
  `id` varchar(255) NOT NULL,
  `aud` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `billing_name` varchar(255) DEFAULT NULL,
  `billing_email` varchar(255) DEFAULT NULL,
  `billing_details` varchar(255) DEFAULT NULL,
  `billing_period` varchar(255) DEFAULT NULL,
  `payment_method_id` varchar(255) NOT NULL,
  `raw_owner_ids` JSON NULL DEFAULT NULL,
  `raw_account_meta_data` JSON NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;