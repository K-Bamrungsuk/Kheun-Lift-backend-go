package queries

const FindUserWeightClass = `
	SELECT "id", "gender", "name", "minWeight", "maxWeight",
	       "createdAt", "updatedAt"
	FROM "WeightClass"
	WHERE "gender" = $1 AND "minWeight" < $2
	  AND ("maxWeight" >= $2 OR "maxWeight" IS NULL)
	ORDER BY "minWeight" DESC, "id" ASC
	LIMIT 1
`

const ListWeightClasses = `
	SELECT "id", "gender", "name", "minWeight", "maxWeight",
	       "createdAt", "updatedAt"
	FROM "WeightClass"
	WHERE "gender" = $1
	ORDER BY "minWeight" ASC
`
