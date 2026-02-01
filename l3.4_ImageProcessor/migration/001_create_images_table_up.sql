CREATE TABLE IF NOT EXISTS images (
id TEXT PRIMARY KEY,
original_path TEXT NOT NULL,
resized_path TEXT,
thumb_path TEXT,
watermark_path TEXT,
status TEXT NOT NULL,
created_at TIMESTAMP NOT NULL,
processed_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_images_status ON images(status);