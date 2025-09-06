-- Default Categories Seed Data for Taiwan Market
-- This file contains default expense and income categories for new users

-- Default Expense Categories (支出類別)
INSERT INTO expense_categories (id, user_id, name, created_at, updated_at) VALUES 
('default-expense-1', 'system-default', '餐飲', NOW(), NOW()),
('default-expense-2', 'system-default', '交通', NOW(), NOW()),
('default-expense-3', 'system-default', '購物', NOW(), NOW()),
('default-expense-4', 'system-default', '娛樂', NOW(), NOW()),
('default-expense-5', 'system-default', '醫療', NOW(), NOW()),
('default-expense-6', 'system-default', '教育', NOW(), NOW()),
('default-expense-7', 'system-default', '居住', NOW(), NOW()),
('default-expense-8', 'system-default', '其他', NOW(), NOW());

-- Default Expense Subcategories
INSERT INTO expense_subcategories (id, parent_id, name) VALUES
-- 餐飲子類別
('default-expense-sub-1-1', 'default-expense-1', '早餐'),
('default-expense-sub-1-2', 'default-expense-1', '午餐'),
('default-expense-sub-1-3', 'default-expense-1', '晚餐'),
('default-expense-sub-1-4', 'default-expense-1', '飲料'),
('default-expense-sub-1-5', 'default-expense-1', '外食'),

-- 交通子類別
('default-expense-sub-2-1', 'default-expense-2', '捷運/公車'),
('default-expense-sub-2-2', 'default-expense-2', '計程車'),
('default-expense-sub-2-3', 'default-expense-2', '停車費'),
('default-expense-sub-2-4', 'default-expense-2', '油費'),
('default-expense-sub-2-5', 'default-expense-2', '汽機車維修'),

-- 購物子類別
('default-expense-sub-3-1', 'default-expense-3', '生活用品'),
('default-expense-sub-3-2', 'default-expense-3', '服飾'),
('default-expense-sub-3-3', 'default-expense-3', '3C產品'),
('default-expense-sub-3-4', 'default-expense-3', '書籍'),

-- 娛樂子類別
('default-expense-sub-4-1', 'default-expense-4', '電影'),
('default-expense-sub-4-2', 'default-expense-4', '遊戲'),
('default-expense-sub-4-3', 'default-expense-4', '運動'),
('default-expense-sub-4-4', 'default-expense-4', '旅遊'),

-- 醫療子類別
('default-expense-sub-5-1', 'default-expense-5', '看診費'),
('default-expense-sub-5-2', 'default-expense-5', '藥費'),
('default-expense-sub-5-3', 'default-expense-5', '健康檢查'),

-- 教育子類別
('default-expense-sub-6-1', 'default-expense-6', '學費'),
('default-expense-sub-6-2', 'default-expense-6', '補習費'),
('default-expense-sub-6-3', 'default-expense-6', '教材'),

-- 居住子類別
('default-expense-sub-7-1', 'default-expense-7', '房租'),
('default-expense-sub-7-2', 'default-expense-7', '水電費'),
('default-expense-sub-7-3', 'default-expense-7', '網路費'),
('default-expense-sub-7-4', 'default-expense-7', '家具'),

-- 其他子類別
('default-expense-sub-8-1', 'default-expense-8', '雜項支出');

-- Default Income Categories (收入類別)
INSERT INTO income_categories (id, user_id, name, created_at, updated_at) VALUES 
('default-income-1', 'system-default', '薪資', NOW(), NOW()),
('default-income-2', 'system-default', '投資', NOW(), NOW()),
('default-income-3', 'system-default', '副業', NOW(), NOW()),
('default-income-4', 'system-default', '其他收入', NOW(), NOW());

-- Default Income Subcategories
INSERT INTO income_subcategories (id, parent_id, name) VALUES
-- 薪資子類別
('default-income-sub-1-1', 'default-income-1', '本薪'),
('default-income-sub-1-2', 'default-income-1', '獎金'),
('default-income-sub-1-3', 'default-income-1', '加班費'),
('default-income-sub-1-4', 'default-income-1', '年終獎金'),

-- 投資子類別
('default-income-sub-2-1', 'default-income-2', '股票股利'),
('default-income-sub-2-2', 'default-income-2', '基金收益'),
('default-income-sub-2-3', 'default-income-2', '租金收入'),
('default-income-sub-2-4', 'default-income-2', '利息收入'),

-- 副業子類別
('default-income-sub-3-1', 'default-income-3', '兼職'),
('default-income-sub-3-2', 'default-income-3', '接案'),
('default-income-sub-3-3', 'default-income-3', '網拍'),
('default-income-sub-3-4', 'default-income-3', '教學'),

-- 其他收入子類別
('default-income-sub-4-1', 'default-income-4', '發票中獎'),
('default-income-sub-4-2', 'default-income-4', '禮金'),
('default-income-sub-4-3', 'default-income-4', '退稅'),
('default-income-sub-4-4', 'default-income-4', '其他');