-- +goose Up
INSERT INTO "Exercise" ("name", "description")
VALUES
    ('Bench Press', 'Barbell bench press'),
    ('Squat', 'Barbell back squat'),
    ('Deadlift', 'Conventional deadlift');

INSERT INTO "WeightClass"
    ("gender", "name", "minWeight", "maxWeight", "updatedAt")
VALUES
    ('male', '59 kg', 0, 59, CURRENT_TIMESTAMP),
    ('male', '66 kg', 59, 66, CURRENT_TIMESTAMP),
    ('male', '74 kg', 66, 74, CURRENT_TIMESTAMP),
    ('male', '83 kg', 74, 83, CURRENT_TIMESTAMP),
    ('male', '93 kg', 83, 93, CURRENT_TIMESTAMP),
    ('male', '105 kg', 93, 105, CURRENT_TIMESTAMP),
    ('male', '120 kg', 105, 120, CURRENT_TIMESTAMP),
    ('male', '120+ kg', 120, NULL, CURRENT_TIMESTAMP),
    ('female', '47 kg', 0, 47, CURRENT_TIMESTAMP),
    ('female', '52 kg', 47, 52, CURRENT_TIMESTAMP),
    ('female', '57 kg', 52, 57, CURRENT_TIMESTAMP),
    ('female', '63 kg', 57, 63, CURRENT_TIMESTAMP),
    ('female', '69 kg', 63, 69, CURRENT_TIMESTAMP),
    ('female', '76 kg', 69, 76, CURRENT_TIMESTAMP),
    ('female', '84 kg', 76, 84, CURRENT_TIMESTAMP),
    ('female', '84+ kg', 84, NULL, CURRENT_TIMESTAMP);

-- +goose Down
DELETE FROM "WeightClass"
WHERE
    (
        "gender" = 'male'
        AND "name" IN (
            '59 kg', '66 kg', '74 kg', '83 kg',
            '93 kg', '105 kg', '120 kg', '120+ kg'
        )
    )
    OR
    (
        "gender" = 'female'
        AND "name" IN (
            '47 kg', '52 kg', '57 kg', '63 kg',
            '69 kg', '76 kg', '84 kg', '84+ kg'
        )
    );

DELETE FROM "Exercise"
WHERE "name" IN ('Bench Press', 'Squat', 'Deadlift');
