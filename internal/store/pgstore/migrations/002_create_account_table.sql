CREATE TABLE IF NOT EXISTS accounts(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(255) NOT NULL,
  balance FLOAT8 DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

---- create above / drop below ----

DROP TABLE IF EXISTS accounts;