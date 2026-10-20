-- +goose Up
CREATE TABLE "Session" (
    "tokenHash" BYTEA PRIMARY KEY,
    "userId" INTEGER NOT NULL,
    "expiresAt" TIMESTAMPTZ NOT NULL,

    CONSTRAINT "Session_tokenHash_length"
        CHECK (octet_length("tokenHash") = 32),

    CONSTRAINT "Session_userId_fkey"
        FOREIGN KEY ("userId") REFERENCES "User"("id")
        ON DELETE CASCADE
);

CREATE INDEX "Session_userId_idx" ON "Session"("userId");
CREATE INDEX "Session_expiresAt_idx" ON "Session"("expiresAt");

-- +goose Down
DROP TABLE "Session";
