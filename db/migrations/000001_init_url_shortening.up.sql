CREATE TABLE "urls" (
  "id" BIGSERIAL PRIMARY KEY,
  "original_url" TEXT NOT NULL,
  "short_url" VARCHAR(7) UNIQUE NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "expired_at" TIMESTAMP NOT NULL
);

CREATE INDEX idx_short_url ON urls (short_url);
CREATE INDEX idx_expired_at ON urls (expired_at);

