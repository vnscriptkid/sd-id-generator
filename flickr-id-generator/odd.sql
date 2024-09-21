-- On TicketServer2 (odd IDs)
SET GLOBAL auto_increment_increment = 2;
SET GLOBAL auto_increment_offset = 2;

CREATE TABLE IF NOT EXISTS `Tickets64` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `stub` char(1) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `stub` (`stub`)
) ENGINE=InnoDB;
