package queries

const LeaderboardWeightClassGender = `
	SELECT "gender" FROM "WeightClass" WHERE "id" = $1
`

// shortcut: timestamps are stored as UTC; normalize legacy rows before supporting other database time zones.
const ListLeaderboard = `
	WITH best_lifts AS (
		SELECT DISTINCT ON (l."userId")
		       l."id", l."userId", l."exerciseId", l."weightClassId",
		       l."weight", l."reps", l."caption", l."videoUrl", l."status",
		       l."createdAt", l."updatedAt",
		       jsonb_build_object(
		           'id', u."id", 'username', u."username",
		           'profileImage', u."profileImage", 'gender', u."gender",
		           'bodyWeight', u."bodyWeight"
		       ) AS public_user
		FROM "LiftRecord" l
		JOIN "User" u ON u."id" = l."userId"
		WHERE l."status" = 'verified'
		  AND l."exerciseId" = $1
		  AND ($2::integer IS NULL OR l."weightClassId" = $2)
		  AND ($3::"Gender" IS NULL OR u."gender" = $3)
		ORDER BY l."userId", l."weight" DESC, l."reps" DESC,
		         l."createdAt", l."id"
	)
	SELECT jsonb_build_object(
		'id', l."id", 'userId', l."userId",
		'exerciseId', l."exerciseId", 'weightClassId', l."weightClassId",
		'weight', l."weight", 'reps', l."reps",
		'caption', l."caption", 'videoUrl', l."videoUrl", 'status', l."status",
		'createdAt', to_char(l."createdAt", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),
		'updatedAt', to_char(l."updatedAt", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),
		'rank', row_number() OVER ranking,
		'user', l.public_user,
		'exercise', jsonb_build_object(
			'id', e."id", 'name', e."name", 'description', e."description"
		),
		'weightClass', jsonb_build_object(
			'id', w."id", 'name', w."name", 'gender', w."gender",
			'minWeight', w."minWeight", 'maxWeight', w."maxWeight",
			'createdAt', to_char(w."createdAt", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),
			'updatedAt', to_char(w."updatedAt", 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"')
		)
	)
	FROM best_lifts l
	JOIN "Exercise" e ON e."id" = l."exerciseId"
	JOIN "WeightClass" w ON w."id" = l."weightClassId"
	WINDOW ranking AS (ORDER BY l."weight" DESC, l."reps" DESC, l."createdAt", l."id")
	ORDER BY row_number() OVER ranking
`

const RandomLeaderboardPair = `
	SELECT pairs."exerciseId", pairs."weightClassId"
	FROM (
		SELECT DISTINCT "exerciseId", "weightClassId"
		FROM "LiftRecord"
		WHERE "status" = 'verified'
	) AS pairs
	ORDER BY random()
	LIMIT 1
`
