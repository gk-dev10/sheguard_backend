CREATE TABLE IF NOT EXISTS profiles (
  id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  name TEXT,
  phone_number TEXT,
  profile_image_url TEXT,
  blood_group TEXT DEFAULT 'Not Set' NOT NULL,
  allergies TEXT,
  medications TEXT,
  notification_preferences JSONB DEFAULT '{}'::jsonb,
  device_tokens JSONB DEFAULT '[]'::jsonb,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  last_login_at TIMESTAMPTZ
);