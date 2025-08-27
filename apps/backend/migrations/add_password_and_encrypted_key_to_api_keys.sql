-- Add password protection and encrypted key storage to API keys table
-- Migration: add_password_and_encrypted_key_to_api_keys

-- Add new columns to api_keys table
ALTER TABLE api_keys 
ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN encrypted_key TEXT NOT NULL DEFAULT '';

-- Create an index on password_hash for better query performance
CREATE INDEX idx_api_keys_password_hash ON api_keys(password_hash);

-- Update the comment to reflect the new structure
COMMENT ON TABLE api_keys IS 'API keys with password protection and encrypted storage for viewing';
COMMENT ON COLUMN api_keys.password_hash IS 'Bcrypt hash of the password required to view the API key';
COMMENT ON COLUMN api_keys.encrypted_key IS 'Encrypted version of the raw API key for secure viewing';
