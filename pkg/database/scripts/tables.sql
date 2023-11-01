CREATE TABLE IF NOT EXISTS Game (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  rows INTEGER NOT NULL,
  symbols TEXT NOT NULL,
  word TEXT NOT NULL,
  won INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

