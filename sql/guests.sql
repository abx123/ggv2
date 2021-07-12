CREATE TABLE IF NOT EXISTS `guests` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `total_rsvp_guests` int(11) NOT NULL,
  `total_arrived_guests` int(11) DEFAULT NULL,
  `version` int(11) NOT NULL,
  `arrivaltime` varchar(45) DEFAULT NULL,
  `tableid` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
