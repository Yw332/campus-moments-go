-- phpMyAdmin SQL Dump
-- version 5.0.2
-- https://www.phpmyadmin.net/
--
-- 主机： 127.0.0.1:3306
-- 生成日期： 2025-11-25 08:58:27
-- 服务器版本： 5.7.31
-- PHP 版本： 7.3.21

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `campus_moments`
--
CREATE DATABASE IF NOT EXISTS `campus_moments` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `campus_moments`;

-- --------------------------------------------------------

--
-- 表的结构 `admins`
--

DROP TABLE IF EXISTS `admins`;
CREATE TABLE IF NOT EXISTS `admins` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` varchar(20) DEFAULT 'admin' COMMENT '角色',
  `last_login_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

--
-- 转存表中的数据 `admins`
--

INSERT INTO `admins` VALUES
(1, 'admin', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', NULL, '2025-11-24 22:14:03', '2025-11-24 22:14:03');

-- --------------------------------------------------------

--
-- 表的结构 `comments`
--

DROP TABLE IF EXISTS `comments`;
CREATE TABLE IF NOT EXISTS `comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `post_id` int(11) NOT NULL COMMENT '所属内容ID',
  `user_id` int(11) NOT NULL COMMENT '评论者ID',
  `parent_id` int(11) DEFAULT NULL COMMENT '父评论ID（回复功能）',
  `content` text NOT NULL COMMENT '评论内容',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态：1-正常，0-违规删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `parent_id` (`parent_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

--
-- 转存表中的数据 `comments`
--

INSERT INTO `comments` VALUES
(1, 1, 2, NULL, '笔记做得好详细！可以借我参考一下吗？', 1, '2025-11-24 09:15:00'),
(2, 1, 3, NULL, '学习氛围真不错，明天我也要去图书馆', 1, '2025-11-24 09:20:00'),
(3, 1, 4, NULL, '坚持就是胜利！加油💪', 1, '2025-11-24 09:30:00'),
(4, 2, 1, NULL, '看起来好香！中午就去尝尝', 1, '2025-11-24 12:35:00'),
(5, 2, 3, NULL, '多少钱一份呀？求推荐口味', 1, '2025-11-24 12:40:00'),
(6, 3, 1, NULL, '太帅了！那个三分球绝杀！', 1, '2025-11-24 16:50:00'),
(7, 3, 2, NULL, '恭喜夺冠！🎉', 1, '2025-11-24 16:55:00'),
(8, 4, 5, NULL, '日出好美！这是哪座山？', 1, '2025-11-23 07:00:00'),
(9, 5, 1, NULL, '我也在学Vue！可以一起交流', 1, '2025-11-24 20:30:00');

-- --------------------------------------------------------

--
-- 表的结构 `friends`
--

DROP TABLE IF EXISTS `friends`;
CREATE TABLE IF NOT EXISTS `friends` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `friend_id` int(11) NOT NULL COMMENT '好友ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0-待确认，1-已好友，2-已拒绝',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`,`friend_id`),
  KEY `idx_friend_id` (`friend_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';

--
-- 转存表中的数据 `friends`
--

INSERT INTO `friends` VALUES
(1, 1, 2, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(2, 1, 3, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(3, 1, 4, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(4, 2, 3, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(5, 2, 5, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(6, 3, 4, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03'),
(7, 3, 5, 1, '2025-11-24 22:14:03', '2025-11-24 22:14:03');

-- --------------------------------------------------------

--
-- 表的结构 `likes`
--

DROP TABLE IF EXISTS `likes`;
CREATE TABLE IF NOT EXISTS `likes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `post_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_user` (`post_id`,`user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

--
-- 转存表中的数据 `likes`
--

INSERT INTO `likes` VALUES
(1, 1, 2, '2025-11-24 22:14:03'),
(2, 1, 3, '2025-11-24 22:14:03'),
(3, 1, 4, '2025-11-24 22:14:03'),
(4, 1, 5, '2025-11-24 22:14:03'),
(5, 2, 1, '2025-11-24 22:14:03'),
(6, 2, 3, '2025-11-24 22:14:03'),
(7, 2, 4, '2025-11-24 22:14:03'),
(8, 2, 5, '2025-11-24 22:14:03'),
(9, 3, 1, '2025-11-24 22:14:03'),
(10, 3, 2, '2025-11-24 22:14:03'),
(11, 3, 4, '2025-11-24 22:14:03'),
(12, 3, 5, '2025-11-24 22:14:03'),
(13, 4, 1, '2025-11-24 22:14:03'),
(14, 4, 2, '2025-11-24 22:14:03'),
(15, 4, 3, '2025-11-24 22:14:03'),
(16, 4, 5, '2025-11-24 22:14:03'),
(17, 5, 1, '2025-11-24 22:14:03'),
(18, 5, 2, '2025-11-24 22:14:03'),
(19, 6, 2, '2025-11-24 22:14:03'),
(20, 6, 3, '2025-11-24 22:14:03'),
(21, 6, 4, '2025-11-24 22:14:03'),
(22, 6, 5, '2025-11-24 22:14:03'),
(23, 7, 1, '2025-11-24 22:14:03'),
(24, 7, 3, '2025-11-24 22:14:03'),
(25, 7, 5, '2025-11-24 22:14:03');

-- --------------------------------------------------------

--
-- 表的结构 `posts`
--

DROP TABLE IF EXISTS `posts`;
CREATE TABLE IF NOT EXISTS `posts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '发布者ID',
  `title` varchar(200) NOT NULL,
  `content` text COMMENT '文本内容',
  `images` json DEFAULT NULL COMMENT '图片URL数组，如["url1","url2"]',
  `video` varchar(500) DEFAULT NULL COMMENT '视频URL',
  `visibility` tinyint(4) DEFAULT '1' COMMENT '可见性：1-公开，2-好友，3-私密',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态：1-正常，0-违规删除',
  `like_count` int(11) DEFAULT '0' COMMENT '点赞数',
  `comment_count` int(11) DEFAULT '0' COMMENT '评论数',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COMMENT='内容表';

--
-- 转存表中的数据 `posts`
--

INSERT INTO `posts` VALUES
(1, 1, '图书馆学习日记', '今天在图书馆学习了一整天，收获满满！分享我的学习笔记～', '[\"https://example.com/study1.jpg\", \"https://example.com/study2.jpg\"]', NULL, 1, 1, 4, 3, '2025-11-24 09:00:00', '2025-11-25 10:15:00'),
(2, 2, '食堂美食推荐', '学校食堂新开的麻辣烫超级好吃！强烈推荐给大家！', '[\"https://example.com/food1.jpg\"]', NULL, 1, 1, 4, 2, '2025-11-24 12:30:00', '2025-11-25 10:15:00'),
(3, 3, '篮球赛夺冠时刻', '下午的篮球赛太精彩了，我们班赢得了冠军！🎉', '[]', 'https://example.com/basketball.mp4', 1, 1, 4, 2, '2025-11-24 16:45:00', '2025-11-25 10:15:00'),
(4, 4, '登山看日出', '周末去爬山看到的日出，真的太美了！分享给大家', '[\"https://example.com/sunrise1.jpg\", \"https://example.com/sunrise2.jpg\"]', NULL, 1, 1, 4, 1, '2025-11-23 06:20:00', '2025-11-25 10:15:00'),
(5, 5, 'Vue.js学习交流', '最近在学习Vue.js，有没有一起学习的小伙伴？', '[\"https://example.com/code1.jpg\"]', NULL, 1, 1, 2, 1, '2025-11-24 20:15:00', '2025-11-25 10:15:00'),
(6, 1, '夜跑打卡第7天', '晚上在操场跑步5公里，坚持运动第7天！', '[\"https://example.com/running1.jpg\"]', NULL, 1, 1, 4, 0, '2025-11-23 21:00:00', '2025-11-25 10:15:00'),
(7, 2, '学习效率APP分享', '分享一个超好用的学习APP，提高效率神器！', '[\"https://example.com/app1.jpg\", \"https://example.com/app2.jpg\"]', NULL, 1, 1, 3, 0, '2025-11-22 15:20:00', '2025-11-25 10:15:00');

-- --------------------------------------------------------

--
-- 表的结构 `post_tags`
--

DROP TABLE IF EXISTS `post_tags`;
CREATE TABLE IF NOT EXISTS `post_tags` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `post_id` int(11) NOT NULL,
  `tag_id` int(11) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_tag` (`post_id`,`tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COMMENT='内容标签关联表';

--
-- 转存表中的数据 `post_tags`
--

INSERT INTO `post_tags` VALUES
(1, 1, 2, '2025-11-24 22:14:03'),
(2, 1, 6, '2025-11-24 22:14:03'),
(3, 2, 3, '2025-11-24 22:14:03'),
(4, 2, 6, '2025-11-24 22:14:03'),
(5, 3, 4, '2025-11-24 22:14:03'),
(6, 3, 8, '2025-11-24 22:14:03'),
(7, 4, 5, '2025-11-24 22:14:03'),
(8, 4, 6, '2025-11-24 22:14:03'),
(9, 5, 2, '2025-11-24 22:14:03'),
(10, 5, 7, '2025-11-24 22:14:03'),
(11, 6, 4, '2025-11-24 22:14:03'),
(12, 6, 6, '2025-11-24 22:14:03'),
(13, 7, 2, '2025-11-24 22:14:03'),
(14, 7, 7, '2025-11-24 22:14:03');

-- --------------------------------------------------------

--
-- 表的结构 `tags`
--

DROP TABLE IF EXISTS `tags`;
CREATE TABLE IF NOT EXISTS `tags` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '标签名',
  `usage_count` int(11) DEFAULT '0' COMMENT '使用次数',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

--
-- 转存表中的数据 `tags`
--

INSERT INTO `tags` VALUES
(1, '校园生活', 0, '2025-11-24 22:14:03'),
(2, '学习分享', 3, '2025-11-24 22:14:03'),
(3, '美食探店', 1, '2025-11-24 22:14:03'),
(4, '运动健身', 2, '2025-11-24 22:14:03'),
(5, '旅行游记', 1, '2025-11-24 22:14:03'),
(6, '日常碎片', 4, '2025-11-24 22:14:03'),
(7, '技术交流', 2, '2025-11-24 22:14:03'),
(8, '娱乐休闲', 1, '2025-11-24 22:14:03');

-- --------------------------------------------------------

--
-- 表的结构 `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '加密密码',
  `nickname` varchar(50) NOT NULL COMMENT '昵称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态：1-正常，0-禁用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

--
-- 转存表中的数据 `users`
--

INSERT INTO `users` VALUES
(1, 'zhangsan', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '张三', 'https://example.com/avatar1.jpg', '13800138001', 1, NULL, '2025-11-20 10:00:00', '2025-11-24 22:14:03'),
(2, 'lisi', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '李四', 'https://example.com/avatar2.jpg', '13800138002', 1, NULL, '2025-11-21 11:00:00', '2025-11-24 22:14:03'),
(3, 'wangwu', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '王五', 'https://example.com/avatar3.jpg', '13800138003', 1, NULL, '2025-11-22 12:00:00', '2025-11-24 22:14:03'),
(4, 'xiaoming', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '小明', 'https://example.com/avatar4.jpg', '13800138004', 1, NULL, '2025-11-23 13:00:00', '2025-11-24 22:14:03'),
(5, 'xiaohong', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '小红', 'https://example.com/avatar5.jpg', '13800138005', 1, NULL, '2025-11-24 14:00:00', '2025-11-24 22:14:03');

--
-- 限制导出的表
--

--
-- 限制表 `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `comments_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `comments_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `comments_ibfk_3` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE;

--
-- 限制表 `friends`
--
ALTER TABLE `friends`
  ADD CONSTRAINT `friends_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `friends_ibfk_2` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- 限制表 `likes`
--
ALTER TABLE `likes`
  ADD CONSTRAINT `likes_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `likes_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- 限制表 `posts`
--
ALTER TABLE `posts`
  ADD CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

--
-- 限制表 `post_tags`
--
ALTER TABLE `post_tags`
  ADD CONSTRAINT `post_tags_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `post_tags_ibfk_2` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
