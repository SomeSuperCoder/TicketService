-- Mock Data Script for Testing
-- This script populates the database with realistic test data
-- ensuring all API parameters can return non-default values

-- Clean existing data (in reverse order of dependencies)
TRUNCATE TABLE ticket_history CASCADE;
TRUNCATE TABLE ticket_tags CASCADE;
TRUNCATE TABLE comments CASCADE;
TRUNCATE TABLE complaint_details CASCADE;
TRUNCATE TABLE tickets CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE subcategories CASCADE;
TRUNCATE TABLE categories CASCADE;
TRUNCATE TABLE departments CASCADE;
TRUNCATE TABLE tags CASCADE;
TRUNCATE TABLE sources CASCADE;
TRUNCATE TABLE statuses CASCADE;

-- Reset sequences
ALTER SEQUENCE categories_id_seq RESTART WITH 1;
ALTER SEQUENCE subcategories_id_seq RESTART WITH 1;
ALTER SEQUENCE departments_id_seq RESTART WITH 1;
ALTER SEQUENCE tags_id_seq RESTART WITH 1;
ALTER SEQUENCE sources_id_seq RESTART WITH 1;
ALTER SEQUENCE statuses_id_seq RESTART WITH 1;

-- ============================================
-- 1. DICTIONARIES
-- ============================================

-- Categories
INSERT INTO categories (name) VALUES 
  ('ЖКХ'),
  ('Дороги'),
  ('Благоустройство'),
  ('Безопасность и правопорядок'),
  ('Связь и телевидение');

-- Subcategories
INSERT INTO subcategories (category_id, name) VALUES 
  -- ЖКХ
  (1, 'Плата за услуги ЖКУ'),
  (1, 'Качество услуг'),
  (1, 'Отопление'),
  (1, 'Водоснабжение'),
  -- Дороги
  (2, 'Ремонт дорог'),
  (2, 'Уборка дорог зимой'),
  (2, 'Освещение дорог'),
  -- Благоустройство
  (3, 'Уборка дворов'),
  (3, 'Детские площадки'),
  (3, 'Озеленение'),
  -- Безопасность
  (4, 'Приюты для животных'),
  (4, 'Освещение улиц'),
  -- Связь
  (5, 'Качество интернета'),
  (5, 'Мобильная связь');

-- Departments
INSERT INTO departments (name) VALUES 
  ('Администрация города Чебоксары'),
  ('Отдел ЖКХ'),
  ('Дорожный отдел'),
  ('Отдел благоустройства'),
  ('Отдел безопасности');

-- Tags
INSERT INTO tags (name) VALUES 
  ('Срочно'),
  ('Требует проверки'),
  ('Дубликат'),
  ('Важное'),
  ('Массовая проблема');

-- Sources
INSERT INTO sources (name) VALUES 
  ('Веб-форма'),
  ('Госуслуги'),
  ('Горячая линия'),
  ('Мобильное приложение'),
  ('Социальные сети');

-- Statuses (for future use)
INSERT INTO statuses (name) VALUES 
  ('Новое'),
  ('В работе'),
  ('Решено'),
  ('Отклонено');

-- ============================================
-- 2. USERS
-- ============================================

INSERT INTO users (id, email, role, status, department_id, first_name, last_name, middle_name) VALUES 
  -- Admin
  ('11111111-1111-1111-1111-111111111111', 'admin@cheboksary.ru', 'admin', 'active', 1, 'Иван', 'Иванов', 'Иванович'),
  
  -- Organizations (РОИ)
  ('22222222-2222-2222-2222-222222222222', 'roi.gkh@cheboksary.ru', 'org', 'active', 2, 'Петр', 'Петров', 'Петрович'),
  ('33333333-3333-3333-3333-333333333333', 'roi.roads@cheboksary.ru', 'org', 'active', 3, 'Сергей', 'Сергеев', 'Сергеевич'),
  ('44444444-4444-4444-4444-444444444444', 'roi.improvement@cheboksary.ru', 'org', 'blocked', 4, 'Мария', 'Смирнова', 'Александровна'),
  
  -- Executors
  ('55555555-5555-5555-5555-555555555555', 'executor1@cheboksary.ru', 'executor', 'active', 2, 'Алексей', 'Алексеев', 'Алексеевич'),
  ('66666666-6666-6666-6666-666666666666', 'executor2@cheboksary.ru', 'executor', 'active', 2, 'Ольга', 'Орлова', 'Олеговна'),
  ('77777777-7777-7777-7777-777777777777', 'executor3@cheboksary.ru', 'executor', 'active', 3, 'Дмитрий', 'Дмитриев', 'Дмитриевич'),
  ('88888888-8888-8888-8888-888888888888', 'executor4@cheboksary.ru', 'executor', 'blocked', 4, 'Елена', 'Егорова', 'Евгеньевна');

-- ============================================
-- 3. TICKETS
-- ============================================

-- Helper function to generate random embeddings
CREATE OR REPLACE FUNCTION random_embedding() RETURNS vector AS $$
DECLARE
    result float[];
    i int;
BEGIN
    result := ARRAY[]::float[];
    FOR i IN 1..768 LOOP
        result := array_append(result, (random() * 2 - 1)::float);
    END LOOP;
    RETURN result::vector;
END;
$$ LANGUAGE plpgsql;

-- CLOSED tickets (old, resolved)
INSERT INTO tickets (id, status, description, subcategory_id, department_id, embedding, created_at) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'closed', 'Нет отопления в квартире на улице Ленина 45', 3, 2, random_embedding(), NOW() - INTERVAL '30 days'),
  ('a0000000-0000-0000-0000-000000000002', 'closed', 'Яма на дороге по улице Гагарина', 5, 3, random_embedding(), NOW() - INTERVAL '25 days'),
  ('a0000000-0000-0000-0000-000000000003', 'closed', 'Не работает уличное освещение', 7, 3, random_embedding(), NOW() - INTERVAL '20 days'),
  ('a0000000-0000-0000-0000-000000000004', 'closed', 'Мусор во дворе не убирается', 8, 4, random_embedding(), NOW() - INTERVAL '15 days'),
  ('a0000000-0000-0000-0000-000000000005', 'closed', 'Сломана детская площадка', 9, 4, random_embedding(), NOW() - INTERVAL '10 days');

-- OPEN tickets (recent, in progress)
INSERT INTO tickets (id, status, description, subcategory_id, department_id, embedding, created_at) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'open', 'Холодные батареи в доме на Ленина 67', 3, 2, random_embedding(), NOW() - INTERVAL '5 days'),
  ('b0000000-0000-0000-0000-000000000002', 'open', 'Плохое качество воды из крана', 4, 2, random_embedding(), NOW() - INTERVAL '4 days'),
  ('b0000000-0000-0000-0000-000000000003', 'open', 'Дорога разбита на улице Мира', 5, 3, random_embedding(), NOW() - INTERVAL '3 days'),
  ('b0000000-0000-0000-0000-000000000004', 'open', 'Не убирают снег во дворе', 6, 3, random_embedding(), NOW() - INTERVAL '2 days'),
  ('b0000000-0000-0000-0000-000000000005', 'open', 'Нужна новая детская площадка', 9, 4, random_embedding(), NOW() - INTERVAL '1 day');

-- OVERDUE tickets (open for more than 7 days)
INSERT INTO tickets (id, status, description, subcategory_id, department_id, embedding, created_at) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'open', 'Протечка крыши в подъезде', 2, 2, random_embedding(), NOW() - INTERVAL '15 days'),
  ('c0000000-0000-0000-0000-000000000002', 'open', 'Большая яма на дороге, опасно', 5, 3, random_embedding(), NOW() - INTERVAL '12 days'),
  ('c0000000-0000-0000-0000-000000000003', 'open', 'Бродячие собаки во дворе', 11, 5, random_embedding(), NOW() - INTERVAL '10 days'),
  ('c0000000-0000-0000-0000-000000000004', 'init', 'Нет освещения в парке', 12, NULL, random_embedding(), NOW() - INTERVAL '9 days'),
  ('c0000000-0000-0000-0000-000000000005', 'init', 'Плохой интернет в районе', 13, NULL, random_embedding(), NOW() - INTERVAL '8 days');

-- INIT tickets (new, not assigned)
INSERT INTO tickets (id, status, description, subcategory_id, department_id, embedding, created_at) VALUES
  ('d0000000-0000-0000-0000-000000000001', 'init', 'Высокие счета за ЖКУ', 1, NULL, random_embedding(), NOW() - INTERVAL '2 hours'),
  ('d0000000-0000-0000-0000-000000000002', 'init', 'Нужно озеленение двора', 10, NULL, random_embedding(), NOW() - INTERVAL '1 hour'),
  ('d0000000-0000-0000-0000-000000000003', 'init', 'Плохая мобильная связь', 14, NULL, random_embedding(), NOW() - INTERVAL '30 minutes');

-- HIDDEN ticket (for testing is_hidden filter)
INSERT INTO tickets (id, status, description, subcategory_id, department_id, embedding, is_hidden, created_at) VALUES
  ('e0000000-0000-0000-0000-000000000001', 'init', 'Спам сообщение', 1, NULL, random_embedding(), TRUE, NOW() - INTERVAL '1 day');

-- ============================================
-- 4. COMPLAINT DETAILS
-- ============================================

-- For closed tickets
INSERT INTO complaint_details (ticket, description, sender_name, sender_phone, sender_email, geo_location) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'Нет отопления в квартире на улице Ленина 45', 'Иванов Иван Иванович', '+7 900 123 45 67', 'ivanov@example.com', ST_SetSRID(ST_MakePoint(47.2501, 56.1324), 4326)),
  ('a0000000-0000-0000-0000-000000000002', 'Яма на дороге по улице Гагарина', 'Петрова Мария', '+7 900 234 56 78', NULL, ST_SetSRID(ST_MakePoint(47.2520, 56.1340), 4326)),
  ('a0000000-0000-0000-0000-000000000003', 'Не работает уличное освещение', 'Сидоров Петр', NULL, 'sidorov@example.com', ST_SetSRID(ST_MakePoint(47.2480, 56.1310), 4326)),
  ('a0000000-0000-0000-0000-000000000004', 'Мусор во дворе не убирается', 'Анонимный', '+7 900 345 67 89', NULL, ST_SetSRID(ST_MakePoint(47.2510, 56.1330), 4326)),
  ('a0000000-0000-0000-0000-000000000005', 'Сломана детская площадка', 'Козлова Ольга', '+7 900 456 78 90', 'kozlova@example.com', ST_SetSRID(ST_MakePoint(47.2490, 56.1320), 4326));

-- For open tickets
INSERT INTO complaint_details (ticket, description, sender_name, sender_phone, sender_email, geo_location) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'Холодные батареи в доме на Ленина 67', 'Смирнов Алексей', '+7 900 567 89 01', 'smirnov@example.com', ST_SetSRID(ST_MakePoint(47.2505, 56.1325), 4326)),
  ('b0000000-0000-0000-0000-000000000002', 'Плохое качество воды из крана', 'Новикова Елена', NULL, 'novikova@example.com', ST_SetSRID(ST_MakePoint(47.2515, 56.1335), 4326)),
  ('b0000000-0000-0000-0000-000000000003', 'Дорога разбита на улице Мира', 'Волков Дмитрий', '+7 900 678 90 12', NULL, ST_SetSRID(ST_MakePoint(47.2525, 56.1345), 4326)),
  ('b0000000-0000-0000-0000-000000000004', 'Не убирают снег во дворе', 'Морозова Анна', '+7 900 789 01 23', 'morozova@example.com', ST_SetSRID(ST_MakePoint(47.2485, 56.1315), 4326)),
  ('b0000000-0000-0000-0000-000000000005', 'Нужна новая детская площадка', 'Лебедев Сергей', '+7 900 890 12 34', NULL, ST_SetSRID(ST_MakePoint(47.2495, 56.1305), 4326));

-- For overdue tickets
INSERT INTO complaint_details (ticket, description, sender_name, sender_phone, sender_email, geo_location) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'Протечка крыши в подъезде', 'Соколова Татьяна', '+7 900 901 23 45', 'sokolova@example.com', ST_SetSRID(ST_MakePoint(47.2500, 56.1300), 4326)),
  ('c0000000-0000-0000-0000-000000000002', 'Большая яма на дороге, опасно', 'Григорьев Игорь', NULL, 'grigoriev@example.com', ST_SetSRID(ST_MakePoint(47.2530, 56.1350), 4326)),
  ('c0000000-0000-0000-0000-000000000003', 'Бродячие собаки во дворе', 'Федорова Наталья', '+7 900 012 34 56', NULL, ST_SetSRID(ST_MakePoint(47.2475, 56.1295), 4326)),
  ('c0000000-0000-0000-0000-000000000004', 'Нет освещения в парке', 'Михайлов Андрей', '+7 900 123 45 67', 'mikhailov@example.com', ST_SetSRID(ST_MakePoint(47.2540, 56.1360), 4326)),
  ('c0000000-0000-0000-0000-000000000005', 'Плохой интернет в районе', 'Павлова Ирина', NULL, 'pavlova@example.com', ST_SetSRID(ST_MakePoint(47.2470, 56.1290), 4326));

-- For init tickets
INSERT INTO complaint_details (ticket, description, sender_name, sender_phone, sender_email, geo_location) VALUES
  ('d0000000-0000-0000-0000-000000000001', 'Высокие счета за ЖКУ', 'Романов Виктор', '+7 900 234 56 78', 'romanov@example.com', ST_SetSRID(ST_MakePoint(47.2512, 56.1328), 4326)),
  ('d0000000-0000-0000-0000-000000000002', 'Нужно озеленение двора', 'Кузнецова Людмила', '+7 900 345 67 89', NULL, ST_SetSRID(ST_MakePoint(47.2488, 56.1318), 4326)),
  ('d0000000-0000-0000-0000-000000000003', 'Плохая мобильная связь', 'Николаев Олег', NULL, 'nikolaev@example.com', ST_SetSRID(ST_MakePoint(47.2522, 56.1342), 4326));

-- For hidden ticket
INSERT INTO complaint_details (ticket, description, sender_name, sender_phone, sender_email, geo_location) VALUES
  ('e0000000-0000-0000-0000-000000000001', 'Спам сообщение', 'Спамер', '+7 900 000 00 00', NULL, ST_SetSRID(ST_MakePoint(47.2500, 56.1320), 4326));

-- ============================================
-- 5. TICKET TAGS
-- ============================================

-- Add tags to various tickets
INSERT INTO ticket_tags (ticket, tag) VALUES
  -- Urgent overdue tickets
  ('c0000000-0000-0000-0000-000000000001', 1), -- Срочно
  ('c0000000-0000-0000-0000-000000000002', 1), -- Срочно
  ('c0000000-0000-0000-0000-000000000002', 4), -- Важное
  
  -- Tickets requiring verification
  ('b0000000-0000-0000-0000-000000000002', 2), -- Требует проверки
  ('d0000000-0000-0000-0000-000000000001', 2), -- Требует проверки
  
  -- Duplicate
  ('b0000000-0000-0000-0000-000000000001', 3), -- Дубликат (similar to a0000000-0000-0000-0000-000000000001)
  
  -- Mass problem
  ('c0000000-0000-0000-0000-000000000003', 5), -- Массовая проблема
  ('b0000000-0000-0000-0000-000000000004', 5); -- Массовая проблема

-- ============================================
-- 6. COMMENTS
-- ============================================

INSERT INTO comments (ticket, message) VALUES
  -- Closed tickets have resolution comments
  ('a0000000-0000-0000-0000-000000000001', 'Проблема принята в работу'),
  ('a0000000-0000-0000-0000-000000000001', 'Отопление восстановлено, проблема решена'),
  
  ('a0000000-0000-0000-0000-000000000002', 'Яма заасфальтирована'),
  
  ('a0000000-0000-0000-0000-000000000003', 'Освещение восстановлено'),
  
  -- Open tickets have work in progress comments
  ('b0000000-0000-0000-0000-000000000001', 'Заявка передана в отдел ЖКХ'),
  ('b0000000-0000-0000-0000-000000000001', 'Специалист выехал на место'),
  
  ('b0000000-0000-0000-0000-000000000003', 'Дорожные работы запланированы на следующую неделю'),
  
  -- Overdue tickets have multiple follow-ups
  ('c0000000-0000-0000-0000-000000000001', 'Заявка зарегистрирована'),
  ('c0000000-0000-0000-0000-000000000001', 'Ожидаем материалы для ремонта'),
  ('c0000000-0000-0000-0000-000000000001', 'Материалы задерживаются'),
  
  ('c0000000-0000-0000-0000-000000000002', 'Требуется срочный ремонт'),
  ('c0000000-0000-0000-0000-000000000002', 'Ожидаем бригаду');

-- ============================================
-- 7. TICKET HISTORY
-- ============================================

-- History for closed tickets
INSERT INTO ticket_history (ticket_id, action, new_value, user_name, user_email, created_at) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'created', '{"status": "init"}', NULL, NULL, NOW() - INTERVAL '30 days'),
  ('a0000000-0000-0000-0000-000000000001', 'status_changed', '{"status": "open"}', 'Петр Петров', 'roi.gkh@cheboksary.ru', NOW() - INTERVAL '29 days'),
  ('a0000000-0000-0000-0000-000000000001', 'department_changed', '{"department_id": 2}', 'Петр Петров', 'roi.gkh@cheboksary.ru', NOW() - INTERVAL '29 days'),
  ('a0000000-0000-0000-0000-000000000001', 'comment_added', '{"comment_id": "uuid", "message": "Проблема принята в работу"}', 'Алексей Алексеев', 'executor1@cheboksary.ru', NOW() - INTERVAL '28 days'),
  ('a0000000-0000-0000-0000-000000000001', 'status_changed', '{"status": "closed"}', 'Алексей Алексеев', 'executor1@cheboksary.ru', NOW() - INTERVAL '25 days');

-- History for open tickets
INSERT INTO ticket_history (ticket_id, action, new_value, user_name, created_at) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'created', '{"status": "init"}', NULL, NOW() - INTERVAL '5 days'),
  ('b0000000-0000-0000-0000-000000000001', 'status_changed', '{"status": "open"}', 'Петр Петров', NOW() - INTERVAL '4 days'),
  ('b0000000-0000-0000-0000-000000000001', 'tags_added', '{"tags": [3]}', 'Петр Петров', NOW() - INTERVAL '4 days');

-- History for overdue tickets
INSERT INTO ticket_history (ticket_id, action, new_value, created_at) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'created', '{"status": "init"}', NOW() - INTERVAL '15 days'),
  ('c0000000-0000-0000-0000-000000000001', 'status_changed', '{"status": "open"}', NOW() - INTERVAL '14 days'),
  ('c0000000-0000-0000-0000-000000000001', 'tags_added', '{"tags": [1]}', NOW() - INTERVAL '10 days');

-- History for init tickets
INSERT INTO ticket_history (ticket_id, action, new_value, created_at) VALUES
  ('d0000000-0000-0000-0000-000000000001', 'created', '{"status": "init"}', NOW() - INTERVAL '2 hours'),
  ('d0000000-0000-0000-0000-000000000002', 'created', '{"status": "init"}', NOW() - INTERVAL '1 hour'),
  ('d0000000-0000-0000-0000-000000000003', 'created', '{"status": "init"}', NOW() - INTERVAL '30 minutes');

-- ============================================
-- CLEANUP
-- ============================================

-- Drop the helper function
DROP FUNCTION random_embedding();

-- ============================================
-- SUMMARY
-- ============================================

-- Display summary
DO $$
BEGIN
    RAISE NOTICE '===========================================';
    RAISE NOTICE 'Mock Data Loaded Successfully!';
    RAISE NOTICE '===========================================';
    RAISE NOTICE 'Categories: %', (SELECT COUNT(*) FROM categories);
    RAISE NOTICE 'Subcategories: %', (SELECT COUNT(*) FROM subcategories);
    RAISE NOTICE 'Departments: %', (SELECT COUNT(*) FROM departments);
    RAISE NOTICE 'Tags: %', (SELECT COUNT(*) FROM tags);
    RAISE NOTICE 'Users: %', (SELECT COUNT(*) FROM users);
    RAISE NOTICE 'Tickets: %', (SELECT COUNT(*) FROM tickets);
    RAISE NOTICE '  - Closed: %', (SELECT COUNT(*) FROM tickets WHERE status = 'closed');
    RAISE NOTICE '  - Open: %', (SELECT COUNT(*) FROM tickets WHERE status = 'open');
    RAISE NOTICE '  - Init: %', (SELECT COUNT(*) FROM tickets WHERE status = 'init');
    RAISE NOTICE '  - Overdue: %', (SELECT COUNT(*) FROM v_ticket_overdue_status WHERE is_overdue = TRUE);
    RAISE NOTICE 'Comments: %', (SELECT COUNT(*) FROM comments);
    RAISE NOTICE 'History Entries: %', (SELECT COUNT(*) FROM ticket_history);
    RAISE NOTICE '===========================================';
END $$;
