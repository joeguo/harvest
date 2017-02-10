CREATE TABLE IF NOT EXISTS domains(
    domain VARCHAR(255) PRIMARY KEY NOT NULL,
    da INT DEFAULT 0,
    crawled TINYINT DEFAULT 0,
    available TINYINT DEFAULT 0,
    KEY `da_crawled_available_index` ( `da`,`crawled`,`available` )
)ENGINE=InnoDB  DEFAULT CHARSET=utf8  ;

CREATE TABLE IF NOT EXISTS `urls` (
  `url` varchar(255) NOT NULL,
  `category` varchar(255) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`url`),
  KEY `language` (`status`,`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `backlinks` (
  `id` int(10) unsigned NOT NULL,
  `url` varchar(255) NOT NULL,
  `target` varchar(255) NOT NULL,
  `active` varchar(255) NOT NULL DEFAULT 'Yes',
  `pr` int(11) NOT NULL DEFAULT '-1',
  `title` varchar(255) NOT NULL DEFAULT '',
  `anchor` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '0',
  `first` varchar(255) NOT NULL DEFAULT '',
  `citation` int(11) NOT NULL DEFAULT '0',
  `trust` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `fail` int(11) NOT NULL DEFAULT '0',
  KEY `id` (`id`),
  KEY `pr` (`pr`),
  KEY `status` (`status`,`fail`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `proxy` (
  `ip` varchar(255) NOT NULL,
  `success` int(11) NOT NULL DEFAULT '0',
  `fail` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ip`),
  KEY `success` (`success`,`fail`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `qualified` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) NOT NULL,
  `tld` varchar(16) NOT NULL,
  `time` datetime NOT NULL,
  `category` varchar(255) NOT NULL DEFAULT '',
  `da` int(11) NOT NULL DEFAULT '0',
  `pa` int(11) NOT NULL DEFAULT '0',
  `www` int(11) NOT NULL DEFAULT '0' COMMENT 'www:1,domain:0',
  `safe` varchar(255) NOT NULL DEFAULT '',
  `indexed` int(11) NOT NULL DEFAULT '0',
  `pr` int(11) NOT NULL DEFAULT '-1',
  `ppr` int(11) NOT NULL DEFAULT '-1',
  `semrank` int(11) NOT NULL DEFAULT '-1',
  `semkeyword` int(11) NOT NULL DEFAULT '-1',
  `semtraffic` int(11) NOT NULL DEFAULT '-1',
  `semcost` int(11) NOT NULL DEFAULT '-1',
  `semdomains` int(11) NOT NULL DEFAULT '-1',
  `semips` int(11) NOT NULL DEFAULT '-1',
  `semfollows` int(11) NOT NULL DEFAULT '-1',
  `semnofollows` int(11) NOT NULL DEFAULT '-1',
  `anchors`  text default '',
  `citation` int(11) NOT NULL DEFAULT '0',
  `trust` int(11) NOT NULL DEFAULT '0',
  `domains` int(11) NOT NULL DEFAULT '0',
  `links` int(11) NOT NULL DEFAULT '0',
  `edu` int(11) NOT NULL DEFAULT '0',
  `gov` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `fail` int(11) NOT NULL DEFAULT '0',
  `pr3` int(11) NOT NULL DEFAULT '0',
  `apr3` int(11) NOT NULL DEFAULT '0',
  `pr4` int(11) NOT NULL DEFAULT '0',
  `pr5` int(11) NOT NULL DEFAULT '0',
  `pr6` int(11) NOT NULL DEFAULT '0',
  `pr7` int(11) NOT NULL DEFAULT '0',
  `pr8` int(11) NOT NULL DEFAULT '0',
  `pr9` int(11) NOT NULL DEFAULT '0',
  `apr4` int(11) NOT NULL DEFAULT '0',
  `apr5` int(11) NOT NULL DEFAULT '0',
  `apr6` int(11) NOT NULL DEFAULT '0',
  `apr7` int(11) NOT NULL DEFAULT '0',
  `apr8` int(11) NOT NULL DEFAULT '0',
  `apr9` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain` (`domain`),
  KEY `fail` (`fail`),
  KEY `time` (`time`),
  KEY `category` (`category`),
  KEY `sortindex` (`category`,`citation`,`trust`,`domains`,`links`,`edu`,`gov`,`da`,`indexed`,`ppr`,`pr`),
  KEY `pr` (`pr4`,`pr5`,`pr6`,`pr7`,`pr8`,`pr9`,`apr4`,`apr5`,`apr6`,`apr7`,`apr8`,`apr9`),
  KEY `pa` (`pa`,`safe`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8  ;

CREATE TABLE IF NOT EXISTS `statics` (
  `uncrawled` int(11) NOT NULL DEFAULT '0',
  `crawled` int(11) NOT NULL DEFAULT '0',
  `time` datetime NOT NULL,  
  PRIMARY KEY (`time`)
 
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

 insert into domains(domain,da) values ('wordpress.org',100),('github.com',96),('joomla.org',100),('drupal.org',99),('name.com',87),('godaddy.com',96);
