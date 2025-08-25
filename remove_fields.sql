-- 删除字段脚本：移除 metadata, callback, attachments 字段
-- 执行前请备份数据库

USE middleware;

-- 1. 删除指定字段
ALTER TABLE purchase_requests 
DROP COLUMN IF EXISTS metadata,
DROP COLUMN IF EXISTS callback,
DROP COLUMN IF EXISTS attachments;

-- 2. 验证表结构
DESCRIBE purchase_requests;

-- 3. 显示当前数据
SELECT * FROM purchase_requests LIMIT 5; 