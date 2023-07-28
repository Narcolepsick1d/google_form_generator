-- Удаление индексов

-- Индекс в таблице Users
DROP INDEX IF EXISTS idx_users_guid;

-- Индекс в таблице Questions
DROP INDEX IF EXISTS idx_questions_guid;

-- Индекс в таблице Label
DROP INDEX IF EXISTS idx_label_guid;

-- Индекс в таблице Choices
DROP INDEX IF EXISTS idx_choices_guid;

-- Индекс в таблице Payment
DROP INDEX IF EXISTS idx_payment_guid;

-- Удаление внешних ключей

-- Внешний ключ User_id в таблице Questions
ALTER TABLE IF EXISTS Questions
DROP CONSTRAINT IF EXISTS fk_questions_user_id;

-- Внешний ключ question_id в таблице Label
ALTER TABLE IF EXISTS Label
DROP CONSTRAINT IF EXISTS fk_label_question_id;

-- Внешний ключ label_id в таблице Choices
ALTER TABLE IF EXISTS Choices
DROP CONSTRAINT IF EXISTS fk_choices_label_id;

-- Внешний ключ user_id в таблице Payment
ALTER TABLE IF EXISTS Payment
DROP CONSTRAINT IF EXISTS fk_payment_user_id;

    -- Удаление таблиц

-- Удаление таблицы Payment
DROP TABLE IF EXISTS Payment;

-- Удаление таблицы Choices
DROP TABLE IF EXISTS Choices;

-- Удаление таблицы Label
DROP TABLE IF EXISTS Label;

-- Удаление таблицы Questions
DROP TABLE IF EXISTS Questions;

-- Удаление таблицы Users
DROP TABLE IF EXISTS Users;
