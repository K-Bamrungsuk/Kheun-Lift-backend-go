package queries

const CreateUser = `
	INSERT INTO "User" ("username", "email", "password", "updatedAt")
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
	RETURNING
		"id", "username", "email", "profileImage",
		"bodyWeight", "height", "gender", "dateOfBirth",
		"createdAt", "updatedAt"
`

const FindLoginUser = `
	SELECT "id", "username", "email", "password"
	FROM "User"
	WHERE "email" = $1
`

const CreateSession = `
	INSERT INTO "Session" ("tokenHash", "userId", "expiresAt")
	VALUES ($1, $2, $3)
`

const FindSessionUser = `
	SELECT
		u."id", u."username", u."email", u."profileImage",
		u."bodyWeight", u."height", u."gender", u."dateOfBirth",
		u."createdAt", u."updatedAt"
	FROM "Session" s
	JOIN "User" u ON u."id" = s."userId"
	WHERE s."tokenHash" = $1
	  AND s."expiresAt" > CURRENT_TIMESTAMP
`

const DeleteSession = `
	DELETE FROM "Session"
	WHERE "tokenHash" = $1
`
