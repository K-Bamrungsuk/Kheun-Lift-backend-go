-- +goose Up
CREATE TYPE "Gender" AS ENUM ('male', 'female');
CREATE TYPE "LiftRecordStatus" AS ENUM ('pending', 'verified', 'rejected');

CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "profileImage" TEXT,
    "gender" "Gender",
    "dateOfBirth" TIMESTAMP(3),
    "height" DOUBLE PRECISION,
    "bodyWeight" DOUBLE PRECISION,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "Exercise" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,

    CONSTRAINT "Exercise_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "WeightClass" (
    "id" SERIAL NOT NULL,
    "gender" "Gender" NOT NULL,
    "name" TEXT NOT NULL,
    "minWeight" DOUBLE PRECISION NOT NULL,
    "maxWeight" DOUBLE PRECISION,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "WeightClass_pkey" PRIMARY KEY ("id")
);

CREATE TABLE "LiftRecord" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "exerciseId" INTEGER NOT NULL,
    "weightClassId" INTEGER NOT NULL,
    "weight" DOUBLE PRECISION NOT NULL,
    "reps" INTEGER NOT NULL,
    "caption" TEXT,
    "videoUrl" TEXT,
    "status" "LiftRecordStatus" NOT NULL DEFAULT 'pending',
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "LiftRecord_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "User_username_key" ON "User"("username");
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");
CREATE UNIQUE INDEX "Exercise_name_key" ON "Exercise"("name");
CREATE UNIQUE INDEX "WeightClass_gender_name_key"
    ON "WeightClass"("gender", "name");

CREATE INDEX "WeightClass_gender_idx" ON "WeightClass"("gender");
CREATE INDEX "LiftRecord_userId_idx" ON "LiftRecord"("userId");
CREATE INDEX "LiftRecord_exerciseId_weightClassId_idx"
    ON "LiftRecord"("exerciseId", "weightClassId");
CREATE INDEX "LiftRecord_status_idx" ON "LiftRecord"("status");

ALTER TABLE "LiftRecord"
    ADD CONSTRAINT "LiftRecord_userId_fkey"
    FOREIGN KEY ("userId") REFERENCES "User"("id")
    ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "LiftRecord"
    ADD CONSTRAINT "LiftRecord_exerciseId_fkey"
    FOREIGN KEY ("exerciseId") REFERENCES "Exercise"("id")
    ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "LiftRecord"
    ADD CONSTRAINT "LiftRecord_weightClassId_fkey"
    FOREIGN KEY ("weightClassId") REFERENCES "WeightClass"("id")
    ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose Down
DROP TABLE "LiftRecord";
DROP TABLE "WeightClass";
DROP TABLE "Exercise";
DROP TABLE "User";

DROP TYPE "LiftRecordStatus";
DROP TYPE "Gender";
