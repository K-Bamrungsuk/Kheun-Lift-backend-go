package queries

const GetUserByID = `
	SELECT "id", "username", "profileImage", "gender", "bodyWeight", "height"
	FROM "User"
	WHERE "id" = $1
`

const DeleteUser = `DELETE FROM "User" WHERE "id" = $1`

const UpdateUser = `
	UPDATE "User"
	SET
		"username" = COALESCE($2, "username"),
		"bodyWeight" = COALESCE($3, "bodyWeight"),
		"height" = COALESCE($4, "height"),
		"gender" = COALESCE($5, "gender"),
		"dateOfBirth" = COALESCE($6, "dateOfBirth"),
		"profileImage" = COALESCE($7, "profileImage"),
		"updatedAt" = CURRENT_TIMESTAMP
	WHERE "id" = $1
	RETURNING
		"id", "username", "email", "profileImage",
		"bodyWeight", "height", "gender", "dateOfBirth",
		"createdAt", "updatedAt"
`
