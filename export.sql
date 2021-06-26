/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80020
 Source Host           : localhost:3306
 Source Schema         : export

 Target Server Type    : MySQL
 Target Server Version : 80020
 File Encoding         : 65001

 Date: 26/06/2021 22:26:51
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for export_file
-- ----------------------------
DROP TABLE IF EXISTS `export_file`;
CREATE TABLE `export_file`  (
  `id` int(0) NOT NULL,
  `hash_key` varchar(23) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `create_at` datetime(0) NULL DEFAULT NULL,
  `update_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of export_file
-- ----------------------------

-- ----------------------------
-- Table structure for export_log
-- ----------------------------
DROP TABLE IF EXISTS `export_log`;
CREATE TABLE `export_log`  (
  `id` int(0) NOT NULL,
  `hash_key` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '参数哈希',
  `title` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '导出标题',
  `ext_type` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '导出类型(文件后缀)',
  `source_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '数据源类型',
  `param` json NULL COMMENT '请求参数（json）',
  `user_id` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户id',
  `callback` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '回调地址',
  `create_at` datetime(0) NULL DEFAULT NULL,
  `update_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of export_log
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
