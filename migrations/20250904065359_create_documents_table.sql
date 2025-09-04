-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,

  user_id UUID NOT NULL REFERENCES users(id) on DELETE CASCADE,

  -- File information
  original_filename VARCHAR(255) NOT NULL,
  file_path VARCHAR(255) NOT NULL,
  file_size BIGINT NOT NULL,
  mime_type VARCHAR(255) NOT NULL,
  file_hash VARCHAR(255) NOT NULL,

  -- Metadata
  title VARCHAR(255),
  description VARCHAR(1000),
  tags VARCHAR(500)
);

--
CREATE INDEX idx_documents_user_id ON documents(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_documents_user_id;
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
