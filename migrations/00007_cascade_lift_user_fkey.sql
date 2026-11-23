-- +goose Up
ALTER TABLE "LiftRecord" DROP CONSTRAINT "LiftRecord_userId_fkey",
    ADD CONSTRAINT "LiftRecord_userId_fkey" FOREIGN KEY ("userId")
    REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose Down
ALTER TABLE "LiftRecord" DROP CONSTRAINT "LiftRecord_userId_fkey",
    ADD CONSTRAINT "LiftRecord_userId_fkey" FOREIGN KEY ("userId")
    REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
