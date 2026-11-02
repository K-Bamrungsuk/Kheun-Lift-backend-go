package queries

const ListExercises = `
	SELECT "id", "name", "description"
	FROM "Exercise"
	ORDER BY "id" ASC
`

const GetExerciseByID = `
	SELECT "id", "name", "description"
	FROM "Exercise"
	WHERE "id" = $1
`
