CREATE TABLE `sessions` (
    `id` varchar(255) PRIMARY KEY NOT NULL,
    `user_email` varchar(255) NOT NULL,
    `refresh_token` varchar(512) NOT NULL,
    `is_revoked` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` datetime DEFAULT (now()),
    `expires_at` datetime
);