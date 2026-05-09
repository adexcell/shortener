-- for daily/monthly fast search
CREATE INDEX IF NOT EXISTS idx_clicks_code_time ON clicks(short_code, clicked_at);
CREATE INDEX IF NOT EXISTS idx_clicks_code_ua ON clicks(short_code, user_agent);
CREATE INDEX IF NOT EXISTS idx_clicks_time ON clicks(clicked_at DESC);
