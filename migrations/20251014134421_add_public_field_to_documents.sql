-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD is_public BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents
DROP COLUMN is_public;
-- +goose StatementEnd
