Create TABLE  if not exists Users(
    Id serial primary key ,
    guid uuid default uuid_generate_v4(),
    Telegram_id varchar not null unique ,
    NickName varchar,
    FirstName varchar,
    LastName varchar,
    idea varchar,
    created_at timestamp default Current_timestamp
    );
create table if not  exists Questions(
    Id serial primary key,
    guid uuid default uuid_generate_v4(),
    User_id int ,
    Url varchar ,
    created_at timestamp default Current_timestamp
);
Create table if not exists Label(
    id serial primary key,
    guid uuid default uuid_generate_v4(),
    entry varchar not null,
    name varchar,
    question_id int,
    is_multi bool default false
);
Create table if not exists Choices(
    id serial primary key,
    guid uuid default uuid_generate_v4(),
    choice varchar,
    probability int default 0,
    label_id int,
    updated_at timestamp
    );
Create table if not exists Payment(
    id serial primary key,
    guid uuid default uuid_generate_v4(),
    user_id int,
    amount numeric(15,2),
    quantity int,
    created_at timestamp default Current_timestamp,
    status varchar ,
    updated_at timestamp default null
    );
-- Добавление индексов в таблице Users
CREATE INDEX IF NOT EXISTS idx_users_guid ON Users (guid);

-- Добавление индексов в таблице Questions
CREATE INDEX IF NOT EXISTS idx_questions_guid ON Questions (guid);

-- Добавление индексов в таблице Label
CREATE INDEX IF NOT EXISTS idx_label_guid ON Label (guid);

-- Добавление индексов в таблице Choices
CREATE INDEX IF NOT EXISTS idx_choices_guid ON Choices (guid);

-- Добавление индексов в таблице Payment
CREATE INDEX IF NOT EXISTS idx_payment_guid ON Payment (guid);

-- Добавление внешних ключей

-- Внешний ключ User_id в таблице Questions, ссылается на столбец Id таблицы Users
ALTER TABLE IF EXISTS Questions
    ADD CONSTRAINT fk_questions_user_id
    FOREIGN KEY (User_id) REFERENCES Users (Id);

-- Внешний ключ question_id в таблице Label, ссылается на столбец Id таблицы Questions
ALTER TABLE IF EXISTS Label
    ADD CONSTRAINT fk_label_question_id
    FOREIGN KEY (question_id) REFERENCES Questions (Id);

-- Внешний ключ label_id в таблице Choices, ссылается на столбец Id таблицы Label
ALTER TABLE IF EXISTS Choices
    ADD CONSTRAINT fk_choices_label_id
    FOREIGN KEY (label_id) REFERENCES Label (Id);

-- Внешний ключ user_id в таблице Payment, ссылается на столбец Id таблицы Users
ALTER TABLE IF EXISTS Payment
    ADD CONSTRAINT fk_payment_user_id
    FOREIGN KEY (user_id) REFERENCES Users (Id);
