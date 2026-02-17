CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  phone_number VARCHAR(15) NOT NULL UNIQUE,
  pin_hash TEXT NOT NULL,

  name VARCHAR(255) NOT NULL,
  email VARCHAR(320),
  profile_image_url TEXT,

  blood_group VARCHAR(5) DEFAULT 'Not Set' NOT NULL,
  allergies TEXT,
  medications TEXT,

  notification_preferences JSONB DEFAULT '{}'::jsonb,
  device_tokens JSONB DEFAULT '[]'::jsonb,

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  last_login_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone_number);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
