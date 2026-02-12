SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `goods`;
CREATE TABLE `goods`  (
  `gid` int(11) NOT NULL AUTO_INCREMENT,
  `gName` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gShelf` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gCnt` int(11) NULL DEFAULT 0,
  `gPrice` decimal(10, 2) NULL DEFAULT NULL,
  PRIMARY KEY (`gid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

INSERT INTO `goods` VALUES (6, '粉笔', '3号', 4, 5.00);
INSERT INTO `goods` VALUES (8, '双面奶', '2号', 3, 0.30);
INSERT INTO `goods` VALUES (12, '水杯', '2号', 13, 5.00);
INSERT INTO `goods` VALUES (20, '大宝', '2号', 12, 10.90);
INSERT INTO `goods` VALUES (26, '曲奇饼干', '4号', 3, 0.90);
INSERT INTO `goods` VALUES (29, '小刚比', '1号', 12, 12.10);
INSERT INTO `goods` VALUES (31, '比比', '23号', 10, 1.10);
INSERT INTO `goods` VALUES (32, '白萝卜', '12号', 12, 1.10);

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `userid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `userpwd` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`userid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

INSERT INTO `user` VALUES ('小强', '111');
INSERT INTO `user` VALUES ('小猪', '222');
INSERT INTO `user` VALUES ('小白', '111');
INSERT INTO `user` VALUES ('小美', '111');
INSERT INTO `user` VALUES ('张三', '1234');
INSERT INTO `user` VALUES ('李四', '111');

SET FOREIGN_KEY_CHECKS = 1;
