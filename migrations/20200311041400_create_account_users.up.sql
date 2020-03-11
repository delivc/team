CREATE TABLE IF NOT EXISTS `{{ index .Options "Namespace" }}accounts_users` (
  `id` varchar(255) NOT NULL,
  `account_id` varchar(255) NOT NULL,
  `user_id` varchar(255) NOT NULL,
  `role_id` varchar(255) NOT NULL,
  `confirmed_at` timestamp NULL DEFAULT NULL,
  `invited_at` timestamp NULL DEFAULT NULL,
  `invited_by` varchar(255) NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`account_id`, `user_id`),
  FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
