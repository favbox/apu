CREATE TABLE `original_url` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` enum('','doc','img') NOT NULL COMMENT '1:doc,2:img',
  `url` varchar(512) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk` (`url`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='原文或原图的网址';