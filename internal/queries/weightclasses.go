package queries

const ListWeightClasses = `
	SELECT "id", "gender", "name", "minWeight", "maxWeight",
	       "createdAt", "updatedAt"
	FROM "WeightClass"
	WHERE "gender" = $1
	ORDER BY "minWeight" ASC
`
