CREATE TABLE IF NOT EXISTS sync_state (
  sync_key text PRIMARY KEY,
  synced_at timestamptz NOT NULL
);
