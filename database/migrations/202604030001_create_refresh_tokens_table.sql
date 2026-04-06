-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_refresh_tokens (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_id VARCHAR(255) NOT NULL, -- JTI
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, token_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_refresh_tokens;
-- +goose StatementEnd
