CREATE UNLOGGED TABLE IF NOT EXISTS cache_entries (
  domain TEXT,
  key TEXT,
  value JSONB,
  count INT,
  expires_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (domain, key)
);
