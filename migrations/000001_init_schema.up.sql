CREATE TABLE IF NOT EXISTS exercises (
  id bigserial PRIMARY KEY,
  external_uuid text NOT NULL,
  name text NOT NULL,
  description text NOT NULL DEFAULT '',
  muscle_groups text [] NOT NULL DEFAULT '{}',
  difficulty text NOT NULL,
  category text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS exercises_external_uuid_idx ON exercises (external_uuid);
CREATE TABLE IF NOT EXISTS sync_state (
  sync_key text PRIMARY KEY,
  synced_at timestamptz NOT NULL
);
