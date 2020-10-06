/*
SQLyog 企业版 - MySQL GUI v8.14 
MySQL - 5.7.30 : Database - im
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`im` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `im`;

/*Table structure for table `group` */

DROP TABLE IF EXISTS `group`;

CREATE TABLE `group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_owner` int(10) unsigned NOT NULL COMMENT '群主',
  `members` varchar(300) NOT NULL DEFAULT '' COMMENT '成员',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=DYNAMIC COMMENT='用户群';

/*Table structure for table `group_message` */

DROP TABLE IF EXISTS `group_message`;

CREATE TABLE `group_message` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `content_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '消息类型',
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '消息内容',
  `from` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '发送人ID',
  `to` int(11) NOT NULL DEFAULT '0' COMMENT '接收群ID',
  `len` int(11) NOT NULL DEFAULT '0' COMMENT '消息长度',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=DYNAMIC COMMENT='群消息';

/*Table structure for table `login_recode` */

DROP TABLE IF EXISTS `login_recode`;

CREATE TABLE `login_recode` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '登陆时间',
  `login_device` varchar(120) NOT NULL DEFAULT '' COMMENT '登陆设备',
  `login_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '登陆IP',
  `login_chn` tinyint(4) NOT NULL DEFAULT '0' COMMENT '登陆渠道',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=DYNAMIC COMMENT='用户登录记录';

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(60) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(31) NOT NULL DEFAULT '' COMMENT '密码',
  `nike_name` varchar(60) NOT NULL DEFAULT '' COMMENT '昵称',
  `icon` varchar(60) NOT NULL DEFAULT '' COMMENT '图标',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=DYNAMIC COMMENT='用户表';

/*Table structure for table `user_message` */

DROP TABLE IF EXISTS `user_message`;

CREATE TABLE `user_message` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `content_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '消息类型',
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '消息内容',
  `from` int(11) NOT NULL DEFAULT '0' COMMENT '发送人',
  `to` int(11) NOT NULL DEFAULT '0' COMMENT '接收人',
  `message_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '消息类型',
  `is_read` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已读',
  `is_send` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否已接收',
  `len` int(11) NOT NULL DEFAULT '0' COMMENT '消息长度',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `deleted_at` int(11) NOT NULL DEFAULT '0' COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=DYNAMIC COMMENT='用户消息表';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
