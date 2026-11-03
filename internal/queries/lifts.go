package queries

const liftColumns = `
	SELECT l."id", l."userId", l."exerciseId", l."weightClassId",
	       l."weight", l."reps", l."caption", l."videoUrl", l."status",
	       l."createdAt", l."updatedAt", u."username", e."name", w."name", w."gender"
`

const liftJoins = `
	JOIN "User" u ON u."id" = l."userId"
	JOIN "Exercise" e ON e."id" = l."exerciseId"
	JOIN "WeightClass" w ON w."id" = l."weightClassId"
`

const ListLiftsByUser = liftColumns + `FROM "LiftRecord" l` + liftJoins + `
	WHERE l."userId" = $1
	ORDER BY l."createdAt" DESC, l."id" DESC
`

const GetLiftByID = liftColumns + `FROM "LiftRecord" l` + liftJoins + `
	WHERE l."id" = $1
`

const LockLiftUser = `
	SELECT "bodyWeight", "gender"
	FROM "User"
	WHERE "id" = $1
	FOR UPDATE
`

const LiftExists = `
	SELECT EXISTS (
		SELECT 1 FROM "LiftRecord"
		WHERE "userId" = $1 AND "exerciseId" = $2
		  AND "weight" = $3 AND "reps" = $4
	)
`

const LiftExerciseExists = `SELECT EXISTS (SELECT 1 FROM "Exercise" WHERE "id" = $1)`

const FindLiftWeightClass = `
	SELECT "id" FROM "WeightClass"
	WHERE "gender" = $1 AND "minWeight" < $2
	  AND ("maxWeight" >= $2 OR "maxWeight" IS NULL)
	ORDER BY "minWeight" DESC
	LIMIT 1
`

const InsertLift = `
	WITH inserted AS (
		INSERT INTO "LiftRecord" (
			"userId", "exerciseId", "weightClassId", "weight", "reps",
			"caption", "videoUrl", "status", "updatedAt"
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', CURRENT_TIMESTAMP)
		RETURNING *
	)
` + liftColumns + `FROM inserted l` + liftJoins

const UpdateLiftCaption = `
	WITH updated AS (
		UPDATE "LiftRecord"
		SET "caption" = $3, "updatedAt" = CURRENT_TIMESTAMP
		WHERE "id" = $1 AND "userId" = $2
		RETURNING *
	)
` + liftColumns + `FROM updated l` + liftJoins

const DeleteLift = `DELETE FROM "LiftRecord" WHERE "id" = $1 AND "userId" = $2`
