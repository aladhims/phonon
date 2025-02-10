CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status INT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_users_username (username),
    INDEX idx_users_email (email)
);

CREATE TABLE IF NOT EXISTS audio_records (
    user_id BIGINT NOT NULL,
    phrase_id BIGINT NOT NULL,
    storage_uri TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (user_id, phrase_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (phrase_id) REFERENCES phrases(id),
    INDEX idx_audio_records_user_phrase (user_id, phrase_id)
);