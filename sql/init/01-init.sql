CREATE TABLE IF NOT EXISTS audio_records (
    user_id BIGINT NOT NULL,
    phrase_id BIGINT NOT NULL,
    original_filename VARCHAR(255),
    original_format VARCHAR(10),
    original_file_uri VARCHAR(255),
    stored_file_uri VARCHAR(255),
    status INT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
    updated_at BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
    PRIMARY KEY (user_id, phrase_id),
    INDEX idx_audio_records_user_phrase (user_id, phrase_id)
);