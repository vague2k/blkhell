CREATE TABLE IF NOT EXISTS releases (
    id TEXT PRIMARY KEY,
    band_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- album, ep, single, compilation
    number TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (band_id) REFERENCES bands (id) ON DELETE CASCADE
);
