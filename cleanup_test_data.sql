-- 清理测试数据
USE rpa_middleware;

-- 删除测试采购请求
DELETE FROM purchase_requests WHERE request_id LIKE 'TEST-PR-%';

-- 删除测试采购订单
DELETE FROM purchase_orders WHERE po_number LIKE 'TEST-PO-%';

-- 删除测试供应商
DELETE FROM suppliers WHERE code LIKE 'SUP%' AND name LIKE '%Test%';

-- 删除测试流程监控数据
DELETE FROM process_monitor WHERE request_id LIKE 'TEST-PR-%' OR request_id LIKE 'TEST-PO-%';

-- 删除测试AI对话数据
DELETE FROM ai_conversations WHERE id LIKE 'conv_%';

-- 删除测试AI消息数据
DELETE FROM ai_messages WHERE conversation_id LIKE 'conv_%';

-- 显示清理后的数据统计
SELECT 'Purchase Requests' as table_name, COUNT(*) as count FROM purchase_requests
UNION ALL
SELECT 'Purchase Orders', COUNT(*) FROM purchase_orders
UNION ALL
SELECT 'Suppliers', COUNT(*) FROM suppliers
UNION ALL
SELECT 'Process Monitor', COUNT(*) FROM process_monitor
UNION ALL
SELECT 'AI Conversations', COUNT(*) FROM ai_conversations
UNION ALL
SELECT 'AI Messages', COUNT(*) FROM ai_messages;
