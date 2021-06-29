CREATE TABLE IF NOT EXISTS `table` (
  `tableid` int(11) NOT NULL AUTO_INCREMENT,
  `capacity` int(11) NOT NULL DEFAULT '0',
  `pcapacity` int(11) NOT NULL DEFAULT '0',
  `acapacity` int(11) NOT NULL DEFAULT '0',
  `version` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`tableid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
