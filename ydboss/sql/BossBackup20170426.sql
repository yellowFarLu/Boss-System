CREATE DATABASE  IF NOT EXISTS `backup_boss` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `backup_boss`;
-- MySQL dump 10.13  Distrib 5.7.17, for Win64 (x86_64)
--
-- Host: localhost    Database: backup_boss
-- ------------------------------------------------------
-- Server version	5.7.17-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `t_account`
--

DROP TABLE IF EXISTS `t_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_account` (
  `account_id` VARCHAR(64) NOT NULL,
  `agent_id` INT(11) NOT NULL,
  `pwd` VARCHAR(64) NOT NULL,
  `NAME` VARCHAR(64) NOT NULL,
  `gender` INT(11) NOT NULL,
  `mobile` VARCHAR(64),
  `enabled` INT(11) DEFAULT '0',
  `timex` DATETIME NOT NULL,
  PRIMARY KEY (`account_id`),
  KEY `Reft_agent34` (`agent_id`),
  CONSTRAINT `Reft_agent34` FOREIGN KEY (`agent_id`) REFERENCES `t_agent` (`agent_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_account`
--

LOCK TABLES `t_account` WRITE;
/*!40000 ALTER TABLE `t_account` DISABLE KEYS */;
INSERT INTO `t_account` VALUES ('peyton.li@xinda.im',2,'832d9f8e4f1524aa76c7963a8d274c844e6951f0','wu',0,'13527206719', 0,'2017-03-31 09:09:09');

INSERT INTO `t_account` VALUES ('bill.zeng@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','bill',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('jamine.chen@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','jamine',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('icey.lan@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','icey',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('may.li@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','may',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('rock.li@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','rock',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('jansen.sun@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','jansen',0,'13527206719', 0,'2017-03-31 09:09:09');
INSERT INTO `t_account` VALUES ('denise.wu@xinda.im',2,'9806fc19405e136195fab8da8e2f922faf019430','wu',0,'13527206719', 0,'2017-03-31 09:09:09');

/*!40000 ALTER TABLE `t_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_account_customers`
--

DROP TABLE IF EXISTS `t_account_customers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_account_customers` (
  `custom_id` INT(11) NOT NULL,
  `account_id` VARCHAR(64) NOT NULL,
  `assign_time` DATE NOT NULL,
  `STATUS` INT(11) DEFAULT NULL,
  PRIMARY KEY (`custom_id`),
  KEY `Reft_account14` (`account_id`),
  CONSTRAINT `Reft_account14` FOREIGN KEY (`account_id`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_customer13` FOREIGN KEY (`custom_id`) REFERENCES `t_customer` (`custom_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_account_customers`
--

LOCK TABLES `t_account_customers` WRITE;
/*!40000 ALTER TABLE `t_account_customers` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_account_customers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_account_roles`
--

DROP TABLE IF EXISTS `t_account_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_account_roles` (
  `role_id` INT(11) NOT NULL,
  `account_id` VARCHAR(64) NOT NULL,
  PRIMARY KEY (`role_id`,`account_id`),
  KEY `Reft_account23` (`account_id`),
  CONSTRAINT `Reft_account23` FOREIGN KEY (`account_id`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_role22` FOREIGN KEY (`role_id`) REFERENCES `t_role` (`role_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_account_roles`
--

LOCK TABLES `t_account_roles` WRITE;
/*!40000 ALTER TABLE `t_account_roles` DISABLE KEYS */;
INSERT INTO `t_account_roles` VALUES (0,'peyton.li@xinda.im');

INSERT INTO `t_account_roles` VALUES (2, 'bill.zeng@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'jamine.chen@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'icey.lan@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'may.li@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'rock.li@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'jansen.sun@xinda.im');
INSERT INTO `t_account_roles` VALUES (2, 'denise.wu@xinda.im');

/*!40000 ALTER TABLE `t_account_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_agent`
--

DROP TABLE IF EXISTS `t_agent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_agent` (
  `agent_id` INT(11) NOT NULL AUTO_INCREMENT,
  `NAME` VARCHAR(64) NOT NULL,
  `contacts` VARCHAR(64) NOT NULL,
  `mobile` VARCHAR(64),
  `mail` VARCHAR(64),
  `enabled` INT(11) DEFAULT '0',
  `note` TEXT,
  `timex` DATETIME NOT NULL,
  `account_id` VARCHAR(64) NOT NULL,
  PRIMARY KEY (`agent_id`)
) ENGINE=INNODB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_agent`
--

LOCK TABLES `t_agent` WRITE;
/*!40000 ALTER TABLE `t_agent` DISABLE KEYS */;
INSERT INTO `t_agent` VALUES (1,'未分配','超级管理员','13527206719','296947440@qq.com',0,'广东 珠海','2017-03-31', 'peyton.li@xinda.im'),(2,'信达九州','超级管理员','13527206719','296947440@qq.com',0,'广东 珠海','2017-03-31', 'peyton.li@xinda.im');
/*!40000 ALTER TABLE `t_agent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_biz`
--

DROP TABLE IF EXISTS `t_biz`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_biz` (
  `biz_id` INT(11) NOT NULL AUTO_INCREMENT,
  `account_id` VARCHAR(64) NOT NULL,
  `custom_id` INT(11) NOT NULL,
  `title` VARCHAR(128) NOT NULL,
  `content` TEXT,
  `amount` FLOAT(8,2) DEFAULT NULL,
  `status` INT(11) DEFAULT 0,
  `timex` DATETIME NOT NULL,
  `estimate_time` DATE,
  `real_time` DATE,
  `enabled` INT(11) DEFAULT '0',
  PRIMARY KEY (`biz_id`),
  KEY `Reft_account1` (`account_id`),
  KEY `Reft_customer2` (`custom_id`),
  CONSTRAINT `Reft_account1` FOREIGN KEY (`account_id`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_customer2` FOREIGN KEY (`custom_id`) REFERENCES `t_customer` (`custom_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_biz`
--

LOCK TABLES `t_biz` WRITE;
/*!40000 ALTER TABLE `t_biz` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_biz` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_biz_comments`
--

DROP TABLE IF EXISTS `t_biz_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_biz_comments` (
  `comment_id` INT(11) NOT NULL AUTO_INCREMENT,
  `biz_id` INT(11) NOT NULL,
  `committer` VARCHAR(64) NOT NULL,
  `comments` TEXT NOT NULL,
  `timex` DATETIME NOT NULL,
  `type` INT(11) DEFAULT 0,
  PRIMARY KEY (`comment_id`),
  KEY `Reft_biz31` (`biz_id`),
  KEY `Reft_account37` (`committer`),
  CONSTRAINT `Reft_account37` FOREIGN KEY (`committer`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_biz31` FOREIGN KEY (`biz_id`) REFERENCES `t_biz` (`biz_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_biz_comments`
--

LOCK TABLES `t_biz_comments` WRITE;
/*!40000 ALTER TABLE `t_biz_comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_biz_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_biz_tags`
--

DROP TABLE IF EXISTS `t_biz_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_biz_tags` (
  `biz_id` INT(11) NOT NULL,
  `tag_id` INT(11) NOT NULL,
  `timex` DATETIME NOT NULL,
  PRIMARY KEY (`biz_id`,`tag_id`),
  KEY `Reft_tag27` (`tag_id`),
  CONSTRAINT `Reft_biz26` FOREIGN KEY (`biz_id`) REFERENCES `t_biz` (`biz_id`),
  CONSTRAINT `Reft_tag27` FOREIGN KEY (`tag_id`) REFERENCES `t_tag` (`tag_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_biz_tags`
--

LOCK TABLES `t_biz_tags` WRITE;
/*!40000 ALTER TABLE `t_biz_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_biz_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_city`
--

DROP TABLE IF EXISTS `t_city`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_city` (
  `city_id` INT(11) NOT NULL AUTO_INCREMENT,
  `province_id` INT(11) NOT NULL,
  `city` VARCHAR(64) NOT NULL,
  `display_order` INT(11) NOT NULL,
  PRIMARY KEY (`city_id`),
  KEY `Reft_province4` (`province_id`),
  CONSTRAINT `Reft_province4` FOREIGN KEY (`province_id`) REFERENCES `t_province` (`province_id`)
) ENGINE=INNODB AUTO_INCREMENT=391 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_city`
--

LOCK TABLES `t_city` WRITE;
/*!40000 ALTER TABLE `t_city` DISABLE KEYS */;
INSERT INTO `t_city` VALUES (1,1,'北京',0),(2,2,'天津',0),(3,3,'上海',0),(4,4,'重庆',0),(5,5,'石家庄',0),(6,5,'唐山',0),(7,5,'秦皇岛',0),(8,5,'邯郸',0),(9,5,'邢台',0),(10,5,'保定',0),(11,5,'张家口',0),(12,5,'承德',0),(13,5,'沧州',0),(14,5,'廊坊',0),(15,5,'衡水',0),(16,6,'太原',0),(17,6,'大同',0),(18,6,'阳泉',0),(19,6,'长治',0),(20,6,'晋城',0),(21,6,'朔州',0),(22,6,'晋中',0),(23,6,'运城',0),(24,6,'忻州',0),(25,6,'临汾',0),(26,6,'吕梁',0),(27,7,'台北',0),(28,7,'高雄',0),(29,7,'基隆',0),(30,7,'台中',0),(31,7,'台南',0),(32,7,'新竹',0),(33,7,'嘉义',0),(34,7,'台北',0),(35,7,'宜兰',0),(36,7,'桃园',0),(37,7,'新竹',0),(38,7,'苗栗',0),(39,7,'台中',0),(40,7,'彰化',0),(41,7,'南投',0),(42,7,'云林',0),(43,7,'嘉义',0),(44,7,'台南',0),(45,7,'高雄',0),(46,7,'屏东',0),(47,7,'澎湖',0),(48,7,'台东',0),(49,7,'花莲',0),(50,8,'沈阳',0),(51,8,'大连',0),(52,8,'鞍山',0),(53,8,'抚顺',0),(54,8,'本溪',0),(55,8,'丹东',0),(56,8,'锦州',0),(57,8,'营口',0),(58,8,'阜新',0),(59,8,'辽阳',0),(60,8,'盘锦',0),(61,8,'铁岭',0),(62,8,'朝阳',0),(63,8,'葫芦岛',0),(64,9,'长春',0),(65,9,'吉林',0),(66,9,'四平',0),(67,9,'辽源',0),(68,9,'通化',0),(69,9,'白山',0),(70,9,'松原',0),(71,9,'白城',0),(72,9,'延边',0),(73,10,'哈尔滨',0),(74,10,'齐齐哈尔',0),(75,10,'鹤岗',0),(76,10,'双鸭山',0),(77,10,'鸡西',0),(78,10,'大庆',0),(79,10,'伊春',0),(80,10,'牡丹江',0),(81,10,'佳木斯',0),(82,10,'七台河',0),(83,10,'黑河',0),(84,10,'绥化',0),(85,10,'大兴安岭',0),(86,11,'南京',0),(87,11,'无锡',0),(88,11,'徐州',0),(89,11,'常州',0),(90,11,'苏州',0),(91,11,'南通',0),(92,11,'连云港',0),(93,11,'淮安',0),(94,11,'盐城',0),(95,11,'扬州',0),(96,11,'镇江',0),(97,11,'泰州',0),(98,11,'宿迁',0),(99,12,'杭州',0),(100,12,'宁波',0),(101,12,'温州',0),(102,12,'嘉兴',0),(103,12,'湖州',0),(104,12,'绍兴',0),(105,12,'金华',0),(106,12,'衢州',0),(107,12,'舟山',0),(108,12,'台州',0),(109,12,'丽水',0),(110,13,'合肥',0),(111,13,'芜湖',0),(112,13,'蚌埠',0),(113,13,'淮南',0),(114,13,'马鞍山',0),(115,13,'淮北',0),(116,13,'铜陵',0),(117,13,'安庆',0),(118,13,'黄山',0),(119,13,'滁州',0),(120,13,'阜阳',0),(121,13,'宿州',0),(122,13,'巢湖',0),(123,13,'六安',0),(124,13,'亳州',0),(125,13,'池州',0),(126,13,'宣城',0),(127,14,'福州',0),(128,14,'厦门',0),(129,14,'莆田',0),(130,14,'三明',0),(131,14,'泉州',0),(132,14,'漳州',0),(133,14,'南平',0),(134,14,'龙岩',0),(135,14,'宁德',0),(136,15,'南昌',0),(137,15,'景德镇',0),(138,15,'萍乡',0),(139,15,'九江',0),(140,15,'新余',0),(141,15,'鹰潭',0),(142,15,'赣州',0),(143,15,'吉安',0),(144,15,'宜春',0),(145,15,'抚州',0),(146,15,'上饶',0),(147,16,'济南',0),(148,16,'青岛',0),(149,16,'淄博',0),(150,16,'枣庄',0),(151,16,'东营',0),(152,16,'烟台',0),(153,16,'潍坊',0),(154,16,'济宁',0),(155,16,'泰安',0),(156,16,'威海',0),(157,16,'日照',0),(158,16,'莱芜',0),(159,16,'临沂',0),(160,16,'德州',0),(161,16,'聊城',0),(162,16,'滨州',0),(163,16,'菏泽',0),(164,17,'郑州',0),(165,17,'开封',0),(166,17,'洛阳',0),(167,17,'平顶山',0),(168,17,'安阳',0),(169,17,'鹤壁',0),(170,17,'新乡',0),(171,17,'焦作',0),(172,17,'濮阳',0),(173,17,'许昌',0),(174,17,'漯河',0),(175,17,'三门峡',0),(176,17,'南阳',0),(177,17,'商丘',0),(178,17,'信阳',0),(179,17,'周口',0),(180,17,'驻马店',0),(181,17,'济源',0),(182,18,'武汉',0),(183,18,'黄石',0),(184,18,'十堰',0),(185,18,'荆州',0),(186,18,'宜昌',0),(187,18,'襄樊',0),(188,18,'鄂州',0),(189,18,'荆门',0),(190,18,'孝感',0),(191,18,'黄冈',0),(192,18,'咸宁',0),(193,18,'随州',0),(194,18,'仙桃',0),(195,18,'天门',0),(196,18,'潜江',0),(197,18,'神农架',0),(198,18,'恩施',0),(199,19,'长沙',0),(200,19,'株洲',0),(201,19,'湘潭',0),(202,19,'衡阳',0),(203,19,'邵阳',0),(204,19,'岳阳',0),(205,19,'常德',0),(206,19,'张家界',0),(207,19,'益阳',0),(208,19,'郴州',0),(209,19,'永州',0),(210,19,'怀化',0),(211,19,'娄底',0),(212,19,'湘西',0),(213,20,'广州',0),(214,20,'深圳',0),(215,20,'珠海',0),(216,20,'汕头',0),(217,20,'韶关',0),(218,20,'佛山',0),(219,20,'江门',0),(220,20,'湛江',0),(221,20,'茂名',0),(222,20,'肇庆',0),(223,20,'惠州',0),(224,20,'梅州',0),(225,20,'汕尾',0),(226,20,'河源',0),(227,20,'阳江',0),(228,20,'清远',0),(229,20,'东莞',0),(230,20,'中山',0),(231,20,'潮州',0),(232,20,'揭阳',0),(233,20,'云浮',0),(234,21,'兰州',0),(235,21,'金昌',0),(236,21,'白银',0),(237,21,'天水',0),(238,21,'嘉峪关',0),(239,21,'武威',0),(240,21,'张掖',0),(241,21,'平凉',0),(242,21,'酒泉',0),(243,21,'庆阳',0),(244,21,'定西',0),(245,21,'陇南',0),(246,21,'临夏',0),(247,21,'甘南',0),(248,22,'成都',0),(249,22,'自贡',0),(250,22,'攀枝花',0),(251,22,'泸州',0),(252,22,'德阳',0),(253,22,'绵阳',0),(254,22,'广元',0),(255,22,'遂宁',0),(256,22,'内江',0),(257,22,'乐山',0),(258,22,'南充',0),(259,22,'眉山',0),(260,22,'宜宾',0),(261,22,'广安',0),(262,22,'达州',0),(263,22,'雅安',0),(264,22,'巴中',0),(265,22,'资阳',0),(266,22,'阿坝',0),(267,22,'甘孜',0),(268,22,'凉山',0),(269,24,'贵阳',0),(270,24,'六盘水',0),(271,24,'遵义',0),(272,24,'安顺',0),(273,24,'铜仁',0),(274,24,'毕节',0),(275,24,'黔西南',0),(276,24,'黔东南',0),(277,24,'黔南',0),(278,25,'海口',0),(279,25,'三亚',0),(280,25,'五指山',0),(281,25,'琼海',0),(282,25,'儋州',0),(283,25,'文昌',0),(284,25,'万宁',0),(285,25,'东方',0),(286,25,'澄迈',0),(287,25,'定安',0),(288,25,'屯昌',0),(289,25,'临高',0),(290,25,'白沙',0),(291,25,'昌江',0),(292,25,'乐东',0),(293,25,'陵水',0),(294,25,'保亭',0),(295,25,'琼中',0),(296,26,'昆明',0),(297,26,'曲靖',0),(298,26,'玉溪',0),(299,26,'保山',0),(300,26,'昭通',0),(301,26,'丽江',0),(302,26,'思茅',0),(303,26,'临沧',0),(304,26,'文山',0),(305,26,'红河',0),(306,26,'西双版纳',0),(307,26,'楚雄',0),(308,26,'大理',0),(309,26,'德宏',0),(310,26,'怒江',0),(311,26,'迪庆',0),(312,27,'西宁',0),(313,27,'海东',0),(314,27,'海北',0),(315,27,'黄南',0),(316,27,'海南',0),(317,27,'果洛',0),(318,27,'玉树',0),(319,27,'海西',0),(320,28,'西安',0),(321,28,'铜川',0),(322,28,'宝鸡',0),(323,28,'咸阳',0),(324,28,'渭南',0),(325,28,'延安',0),(326,28,'汉中',0),(327,28,'榆林',0),(328,28,'安康',0),(329,28,'商洛',0),(330,29,'南宁',0),(331,29,'柳州',0),(332,29,'桂林',0),(333,29,'北海',0),(334,29,'防城港',0),(335,29,'钦州',0),(336,29,'贵港',0),(337,29,'玉林',0),(338,29,'百色',0),(339,29,'贺州',0),(340,29,'河池',0),(341,29,'来宾',0),(342,29,'崇左',0),(343,30,'拉萨',0),(344,30,'那曲',0),(345,30,'昌都',0),(346,30,'山南',0),(347,30,'日喀则',0),(348,30,'阿里',0),(349,30,'林芝',0),(350,31,'银川',0),(351,31,'石嘴山',0),(352,31,'吴忠',0),(353,31,'固原',0),(354,31,'中卫',0),(355,32,'乌鲁木齐',0),(356,32,'克拉玛依',0),(357,32,'石河子　',0),(358,32,'阿拉尔',0),(359,32,'图木舒克',0),(360,32,'五家渠',0),(361,32,'吐鲁番',0),(362,32,'阿克苏',0),(363,32,'喀什',0),(364,32,'哈密',0),(365,32,'和田',0),(366,32,'阿图什',0),(367,32,'库尔勒',0),(368,32,'昌吉　',0),(369,32,'阜康',0),(370,32,'米泉',0),(371,32,'博乐',0),(372,32,'伊宁',0),(373,32,'奎屯',0),(374,32,'塔城',0),(375,32,'乌苏',0),(376,32,'阿勒泰',0),(377,33,'呼和浩特',0),(378,33,'包头',0),(379,33,'乌海',0),(380,33,'赤峰',0),(381,33,'通辽',0),(382,33,'鄂尔多斯',0),(383,33,'呼伦贝尔',0),(384,33,'巴彦淖尔',0),(385,33,'乌兰察布',0),(386,33,'锡林郭勒盟',0),(387,33,'兴安盟',0),(388,33,'阿拉善盟',0),(389,34,'澳门',0),(390,35,'香港',0);
/*!40000 ALTER TABLE `t_city` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_customer`
--

DROP TABLE IF EXISTS `t_customer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_customer` (
  `custom_id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `rtx_number` VARCHAR(8) DEFAULT NULL,
  `contacts` VARCHAR(64) DEFAULT NULL,
  `phone` VARCHAR(64) DEFAULT NULL,
  `mobile` VARCHAR(64) DEFAULT NULL,
  `qq` VARCHAR(64) DEFAULT NULL,
  `mail` CHAR(64) DEFAULT NULL,
  `agent_id` INT(11) NOT NULL,
  `city_id` INT(11) DEFAULT NULL,
  `timex` DATETIME DEFAULT NULL,
  `assign_status` INT(11) NOT NULL,
  `last_follow_time` DATETIME DEFAULT NULL,
  `note` TEXT,
  `enabled` INT(11) DEFAULT '0',
  `emp_account` INT(11) DEFAULT '0',
  PRIMARY KEY (`custom_id`),
  KEY `Reft_agent35` (`agent_id`),
  KEY `Reft_city36` (`city_id`),
  CONSTRAINT `Reft_agent35` FOREIGN KEY (`agent_id`) REFERENCES `t_agent` (`agent_id`),
  CONSTRAINT `Reft_city36` FOREIGN KEY (`city_id`) REFERENCES `t_city` (`city_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_customer`
--

LOCK TABLES `t_customer` WRITE;
/*!40000 ALTER TABLE `t_customer` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_customer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_customer_assign_history`
--

DROP TABLE IF EXISTS `t_customer_assign_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_customer_assign_history` (
  `assign_id` INT(11) NOT NULL AUTO_INCREMENT,
  `custom_id` INT(11) NOT NULL,
  `assigner` VARCHAR(64) NOT NULL,
  `agent_id` INT(11) NOT NULL,
  `assignee` VARCHAR(64) DEFAULT NULL,
  `timex` DATETIME NOT NULL,
  PRIMARY KEY (`assign_id`),
  KEY `Reft_customer20` (`custom_id`),
  KEY `Reft_agent40` (`agent_id`),
  KEY `Reft_account41` (`assigner`),
  KEY `Reft_account42` (`assignee`),
  CONSTRAINT `Reft_account41` FOREIGN KEY (`assigner`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_account42` FOREIGN KEY (`assignee`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_agent40` FOREIGN KEY (`agent_id`) REFERENCES `t_agent` (`agent_id`),
  CONSTRAINT `Reft_customer20` FOREIGN KEY (`custom_id`) REFERENCES `t_customer` (`custom_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_customer_assign_history`
--

LOCK TABLES `t_customer_assign_history` WRITE;
/*!40000 ALTER TABLE `t_customer_assign_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_customer_assign_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_customer_comments`
--

DROP TABLE IF EXISTS `t_customer_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_customer_comments` (
  `comment_id` INT(11) NOT NULL AUTO_INCREMENT,
  `custom_id` INT(11) NOT NULL,
  `committer` VARCHAR(64) NOT NULL,
  `comments` TEXT NOT NULL,
  `timex` DATETIME NOT NULL,
  `type` INT(11) DEFAULT 0,
  PRIMARY KEY (`comment_id`),
  KEY `Reft_customer21` (`custom_id`),
  KEY `Reft_account38` (`committer`),
  CONSTRAINT `Reft_account38` FOREIGN KEY (`committer`) REFERENCES `t_account` (`account_id`),
  CONSTRAINT `Reft_customer21` FOREIGN KEY (`custom_id`) REFERENCES `t_customer` (`custom_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_customer_comments`
--

LOCK TABLES `t_customer_comments` WRITE;
/*!40000 ALTER TABLE `t_customer_comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_customer_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_customer_tags`
--

DROP TABLE IF EXISTS `t_customer_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_customer_tags` (
  `custom_id` INT(11) NOT NULL,
  `tag_id` INT(11) NOT NULL,
  `timex` DATETIME NOT NULL,
  PRIMARY KEY (`custom_id`,`tag_id`),
  KEY `Reft_tag25` (`tag_id`),
  CONSTRAINT `Reft_customer24` FOREIGN KEY (`custom_id`) REFERENCES `t_customer` (`custom_id`),
  CONSTRAINT `Reft_tag25` FOREIGN KEY (`tag_id`) REFERENCES `t_tag` (`tag_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_customer_tags`
--

LOCK TABLES `t_customer_tags` WRITE;
/*!40000 ALTER TABLE `t_customer_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_customer_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_province`
--

DROP TABLE IF EXISTS `t_province`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_province` (
  `province_id` INT(11) NOT NULL AUTO_INCREMENT,
  `province` VARCHAR(64) NOT NULL,
  `display_order` INT(11) NOT NULL,
  PRIMARY KEY (`province_id`)
) ENGINE=INNODB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_province`
--

LOCK TABLES `t_province` WRITE;
/*!40000 ALTER TABLE `t_province` DISABLE KEYS */;
INSERT INTO `t_province` VALUES (1,'北京',0),(2,'天津',0),(3,'上海',0),(4,'重庆',0),(5,'河北',0),(6,'山西',0),(7,'台湾',0),(8,'辽宁',0),(9,'吉林',0),(10,'黑龙江',0),(11,'江苏',0),(12,'浙江',0),(13,'安徽',0),(14,'福建',0),(15,'江西',0),(16,'山东',0),(17,'河南',0),(18,'湖北',0),(19,'湖南',0),(20,'广东',0),(21,'甘肃',0),(22,'四川',0),(24,'贵州',0),(25,'海南',0),(26,'云南',0),(27,'青海',0),(28,'陕西',0),(29,'广西',0),(30,'西藏',0),(31,'宁夏',0),(32,'新疆',0),(33,'内蒙古',0),(34,'澳门',0),(35,'香港',0);
/*!40000 ALTER TABLE `t_province` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_role`
--

DROP TABLE IF EXISTS `t_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_role` (
  `role_id` INT(11) NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_role`
--

LOCK TABLES `t_role` WRITE;
/*!40000 ALTER TABLE `t_role` DISABLE KEYS */;
INSERT INTO `t_role` VALUES (0,'超级管理员'),(1,'信达管理员'),(2,'渠道'),(3,'员工');
/*!40000 ALTER TABLE `t_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_tag`
--

DROP TABLE IF EXISTS `t_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_tag` (
  `tag_id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `type` VARCHAR(20) DEFAULT NULL,
  `note` TEXT,
  PRIMARY KEY (`tag_id`)
) ENGINE=INNODB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_tag`
--

LOCK TABLES `t_tag` WRITE;
/*!40000 ALTER TABLE `t_tag` DISABLE KEYS */;
INSERT INTO `t_tag` VALUES (1,'国企','unit_nature',''),(2,'私企','unit_nature',''),(3,'股份制','unit_nature',''),(4,'其他','unit_nature',''),(5,'教育','bs_domain',''),
(6,'金融','bs_domain',''),(7,'石油石化','bs_domain',''),(8,'能源','bs_domain',''),(9,'钢铁','bs_domain',''),(10,'餐饮','bs_domain',''),(11,'电讯业','bs_domain',''),
(12,'房地产','bs_domain',''),(13,'媒体','bs_domain',''),(14,'出版社','bs_domain',''),(15,'医疗','bs_domain',''),(16,'互联网','bs_domain',''),(17,'体育运动','bs_domain',''),
(18,'政府机关','bs_domain',''),(19,'其他','bs_domain',''),(20,'A类','cu_level',''),(21,'B类','cu_level',''),(22,'C类','cu_level',''),(23,'有度销售','biz_category',''),
(24,'RTX销售','biz_category',''),(25,'RTX短信','biz_category',''),(26,'企业微信项目','biz_category',''),(27,'有度项目','biz_category','');
/*!40000 ALTER TABLE `t_tag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'backup_boss'
--

--
-- Dumping routines for database 'backup_boss'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-04-06  9:28:10