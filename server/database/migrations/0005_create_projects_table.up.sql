CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    band_id TEXT NOT NULL,
    release_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- CD, tapes, merch, vinyl
    status TEXT NOT NULL, -- in-progress, done, pending
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (band_id) REFERENCES bands (id) ON DELETE CASCADE,
    FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE
);
