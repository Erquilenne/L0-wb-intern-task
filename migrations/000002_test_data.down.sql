DELETE FROM order_items WHERE order_id IN (SELECT order_id FROM orders WHERE order_uid IN ('order_uid_1', 'order_uid_2'));

-- Удаляем тестовые элементы в orders
DELETE FROM orders WHERE order_uid IN ('order_uid_1', 'order_uid_2', 'order_uid_3', 'order_uid_4', 'order_uid_5', 'order_uid_6', 'order_uid_7', 'order_uid_8');