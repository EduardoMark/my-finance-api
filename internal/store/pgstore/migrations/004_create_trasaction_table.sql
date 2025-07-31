-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  description VARCHAR NOT NULL,
  amount FLOAT8 NOT NULL,
  date DATE NOT NULL,
  type transaction_type NOT NULL,
  account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

---- create above / drop below ----

DROP TABLE IF NOT EXISTS transaction;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
