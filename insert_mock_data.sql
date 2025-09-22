-- 使用数据库
USE middleware;

-- 插入采购请求模拟数据
INSERT INTO purchase_requests (
    request_id, doc_type, plant, quantity, unit_price, material, delivery_date,
    vendor_code, short_text, material_group, unit_type, requester,
    purchase_organization, currency, total_amount, priority, status, urgency,
    comments, retry_count, created_at, updated_at
) VALUES
('PR-20231201-001', 'NB', '1000', 100, 25.50, 'MAT001', '2023-12-15 00:00:00',
 'VENDOR001', 'Office Supplies Procurement', 'MG001', 'EA', 'John Smith',
 'PO001', 'CNY', 2550.00, 2, 'pending', 'normal',
 'Urgent office supplies needed', 0, NOW(), NOW()),

('PR-20231201-002', 'NB', '1000', 50, 100.00, 'MAT002', '2023-12-20 00:00:00',
 'VENDOR002', 'Computer Equipment', 'MG002', 'EA', 'Alice Johnson',
 'PO001', 'CNY', 5000.00, 3, 'approved', 'urgent',
 'New computer equipment for office', 0, NOW(), NOW()),

('PR-20231201-003', 'UB', '2000', 200, 15.75, 'MAT003', '2023-12-25 00:00:00',
 'VENDOR003', 'Raw Materials', 'MG003', 'KG', 'Bob Wilson',
 'PO002', 'CNY', 3150.00, 1, 'processing', 'normal',
 'Raw materials for production', 0, NOW(), NOW()),

('PR-20231201-004', 'NB', '1000', 10, 500.00, 'MAT004', '2023-12-30 00:00:00',
 'VENDOR004', 'Industrial Equipment', 'MG004', 'EA', 'Carol Davis',
 'PO001', 'CNY', 5000.00, 4, 'completed', 'critical',
 'Critical equipment for production line', 0, NOW(), NOW()),

('PR-20231201-005', 'NB', '3000', 75, 30.00, 'MAT005', '2024-01-05 00:00:00',
 'VENDOR005', 'Safety Equipment', 'MG005', 'EA', 'David Brown',
 'PO003', 'CNY', 2250.00, 2, 'failed', 'normal',
 'Safety equipment for warehouse', 2, NOW(), NOW());

-- 插入采购订单模拟数据
INSERT INTO purchase_orders (
    po_number, status, supplier_name, supplier_code, total_amount, currency,
    created_by, created_at, sent_at, confirmed_at, notes
) VALUES
('PO-20231201-001', 'sent', 'ABC Manufacturing Co., Ltd.', 'SUP001', 50000.00, 'CNY',
 'John Smith', NOW(), NOW(), NULL, 'Urgent procurement'),

('PO-20231201-002', 'confirmed', 'XYZ Electronics Inc.', 'SUP002', 35000.00, 'CNY',
 'Alice Johnson', NOW(), NOW(), NOW(), 'Electronics equipment order'),

('PO-20231201-003', 'delivered', 'DEF Materials Corp.', 'SUP003', 20000.00, 'CNY',
 'Bob Wilson', NOW(), NOW(), NOW(), 'Raw materials delivered'),

('PO-20231201-004', 'draft', 'GHI Industrial Ltd.', 'SUP004', 15000.00, 'CNY',
 'Carol Davis', NOW(), NULL, NULL, 'Draft order for review'),

('PO-20231201-005', 'cancelled', 'JKL Services Inc.', 'SUP005', 8000.00, 'CNY',
 'David Brown', NOW(), NOW(), NULL, 'Order cancelled due to budget constraints');

-- 插入供应商模拟数据
INSERT INTO suppliers (
    id, code, name, type, status, contact_person, phone, email, address,
    categories, notes, created_at
) VALUES
('SUP_001', 'SUP001', 'ABC Manufacturing Co., Ltd.', 'manufacturer', 'active',
 'John Smith', '13800138001', 'john.smith@abc.com', 'No.123, Industrial Road, Chaoyang District, Beijing',
 '["Raw Materials", "Equipment"]', 'Long-term partner', NOW()),

('SUP_002', 'SUP002', 'XYZ Electronics Inc.', 'manufacturer', 'active',
 'Alice Johnson', '13800138002', 'alice.johnson@xyz.com', 'No.456, Tech Street, Haidian District, Beijing',
 '["Electronics", "Computer Equipment"]', 'Reliable electronics supplier', NOW()),

('SUP_003', 'SUP003', 'DEF Materials Corp.', 'distributor', 'active',
 'Bob Wilson', '13800138003', 'bob.wilson@def.com', 'No.789, Material Avenue, Fengtai District, Beijing',
 '["Raw Materials", "Chemicals"]', 'Material distribution specialist', NOW()),

('SUP_004', 'SUP004', 'GHI Industrial Ltd.', 'manufacturer', 'active',
 'Carol Davis', '13800138004', 'carol.davis@ghi.com', 'No.321, Industrial Park, Shunyi District, Beijing',
 '["Industrial Equipment", "Machinery"]', 'Industrial equipment manufacturer', NOW()),

('SUP_005', 'SUP005', 'JKL Services Inc.', 'service', 'inactive',
 'David Brown', '13800138005', 'david.brown@jkl.com', 'No.654, Service Center, Dongcheng District, Beijing',
 '["Maintenance", "Consulting"]', 'Service provider', NOW()),

('SUP_006', 'SUP006', 'MNO Logistics Co.', 'service', 'active',
 'Eva Green', '13800138006', 'eva.green@mno.com', 'No.987, Logistics Hub, Tongzhou District, Beijing',
 '["Logistics", "Transportation"]', 'Logistics and transportation services', NOW()),

('SUP_007', 'SUP007', 'PQR Technology Ltd.', 'manufacturer', 'pending',
 'Frank Miller', '13800138007', 'frank.miller@pqr.com', 'No.147, Innovation Street, Changping District, Beijing',
 '["Technology", "Software"]', 'Technology solutions provider', NOW());

-- 插入流程监控模拟数据
INSERT INTO process_monitor (
    request_id, process_type, current_step, responsible, department,
    blocked_duration, priority, reason, status, created_at, updated_at
) VALUES
('PR-20231201-001', 'pr', 'Department Approval', 'Mike Johnson', 'Procurement Department',
 '2 days', 4, 'Waiting for department manager approval', 'blocked', NOW(), NOW()),

('PR-20231201-002', 'pr', 'Budget Review', 'Sarah Wilson', 'Finance Department',
 '1 day', 3, 'Budget allocation pending', 'blocked', NOW(), NOW()),

('PR-20231201-003', 'pr', 'Technical Review', 'David Chen', 'IT Department',
 '3 days', 2, 'Technical specifications under review', 'blocked', NOW(), NOW()),

('PO-20231201-001', 'po', 'Supplier Confirmation', 'Supplier', 'External',
 '1 day', 3, 'Waiting for supplier response', 'blocked', NOW(), NOW()),

('PO-20231201-002', 'po', 'Delivery', 'Supplier', 'External',
 '2 days', 2, 'Delivery in progress', 'processing', NOW(), NOW());

-- 插入AI对话模拟数据
INSERT INTO ai_conversations (
    id, user_id, department, context, created_at, updated_at
) VALUES
('conv_001', 'user_001', 'Procurement', '{"role": "procurement_manager", "level": "senior"}', NOW(), NOW()),
('conv_002', 'user_002', 'Finance', '{"role": "finance_manager", "level": "manager"}', NOW(), NOW()),
('conv_003', 'user_003', 'IT', '{"role": "it_specialist", "level": "senior"}', NOW(), NOW());

-- 插入AI消息模拟数据
INSERT INTO ai_messages (
    conversation_id, message, response, response_type, suggestions, created_at
) VALUES
('conv_001', 'Create a new PR for office supplies', 'I\'ll help you create a purchase request for office supplies. Please provide the following information...', 'text', '["Create PR", "Query PO Status", "View Process Bottlenecks"]', NOW()),

('conv_002', 'Check the status of PR-20231201-001', 'The PR PR-20231201-001 is currently pending department approval. It has been in this status for 2 days.', 'text', '["View Details", "Update Status", "Export Data"]', NOW()),

('conv_003', 'Show me the process workflow for PR', 'Here is the PR process workflow with current status...', 'process', '["View Process Details", "Check Blocked Items", "Process Statistics"]', NOW());