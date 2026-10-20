package queries

const ListExercises = `
	SELECT "id", "name", "description"
	FROM "Exercise"
	ORDER BY "id" ASC
`
