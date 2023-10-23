CREATE TABLE IF NOT EXISTS people (
    people_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255),
    age INT,
    gender VARCHAR(6),
    country_id VARCHAR(5)
);