CREATE TABLE `users` (
    `telegram_id` INTEGER NOT NULL UNIQUE,
    `username` TEXT,
    `name` TEXT,
    `age` INTEGER,
    `sex` TEXT,
    `about` TEXT,
    `hobbies` TEXT,
    `work` TEXT,
    `education` TEXT,
    `cover_letter` TEXT,
    `contacts` TEXT,
    PRIMARY KEY (`telegram_id`)
);