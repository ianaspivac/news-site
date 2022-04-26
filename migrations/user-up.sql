CREATE TABLE IF NOT EXISTS user (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    uuid VARCHAR(255) UNIQUE NOT NULL,
    mail VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    type ENUM('BASIC','EDITOR','ADMIN')
) DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;