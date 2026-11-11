-- Mock accounts for local/development use only. All passwords: password123.
-- Plain INSERT intentionally fails on username/email conflicts without changing existing users.
-- Accounts (username | email | password | gender | bodyWeight kg | weight class):
-- Bob | bob@example.com | password123 | male | 58 | 59 kg
-- Mark | mark@example.com | password123 | male | 57 | 59 kg
-- Jacob | jacob@example.com | password123 | male | 56 | 59 kg
-- John | john@example.com | password123 | male | 65 | 66 kg
-- Mike | mike@example.com | password123 | male | 64 | 66 kg
-- Alex | alex@example.com | password123 | male | 63 | 66 kg
-- Ben | ben@example.com | password123 | male | 73 | 74 kg
-- Chris | chris@example.com | password123 | male | 72 | 74 kg
-- Daniel | daniel@example.com | password123 | male | 71 | 74 kg
-- Ethan | ethan@example.com | password123 | male | 82 | 83 kg
-- Finn | finn@example.com | password123 | male | 81 | 83 kg
-- George | george@example.com | password123 | male | 80 | 83 kg
-- Harry | harry@example.com | password123 | male | 92 | 93 kg
-- Isaac | isaac@example.com | password123 | male | 91 | 93 kg
-- Jack | jack@example.com | password123 | male | 90 | 93 kg
-- Kevin | kevin@example.com | password123 | male | 104 | 105 kg
-- Leo | leo@example.com | password123 | male | 103 | 105 kg
-- Max | max@example.com | password123 | male | 102 | 105 kg
-- Noah | noah@example.com | password123 | male | 119 | 120 kg
-- Owen | owen@example.com | password123 | male | 118 | 120 kg
-- Paul | paul@example.com | password123 | male | 117 | 120 kg
-- Ryan | ryan@example.com | password123 | male | 123 | 120+ kg
-- Sam | sam@example.com | password123 | male | 124 | 120+ kg
-- Tom | tom@example.com | password123 | male | 125 | 120+ kg
-- Anna | anna@example.com | password123 | female | 46 | 47 kg
-- Bella | bella@example.com | password123 | female | 45 | 47 kg
-- Chloe | chloe@example.com | password123 | female | 44 | 47 kg
-- Daisy | daisy@example.com | password123 | female | 51 | 52 kg
-- Emma | emma@example.com | password123 | female | 50 | 52 kg
-- Faith | faith@example.com | password123 | female | 49 | 52 kg
-- Grace | grace@example.com | password123 | female | 56 | 57 kg
-- Hannah | hannah@example.com | password123 | female | 55 | 57 kg
-- Ivy | ivy@example.com | password123 | female | 54 | 57 kg
-- Jade | jade@example.com | password123 | female | 62 | 63 kg
-- Kate | kate@example.com | password123 | female | 61 | 63 kg
-- Lily | lily@example.com | password123 | female | 60 | 63 kg
-- Mia | mia@example.com | password123 | female | 68 | 69 kg
-- Nora | nora@example.com | password123 | female | 67 | 69 kg
-- Olivia | olivia@example.com | password123 | female | 66 | 69 kg
-- Paige | paige@example.com | password123 | female | 75 | 76 kg
-- Quinn | quinn@example.com | password123 | female | 74 | 76 kg
-- Rose | rose@example.com | password123 | female | 73 | 76 kg
-- Sarah | sarah@example.com | password123 | female | 83 | 84 kg
-- Tina | tina@example.com | password123 | female | 82 | 84 kg
-- Uma | uma@example.com | password123 | female | 81 | 84 kg
-- Vera | vera@example.com | password123 | female | 87 | 84+ kg
-- Wendy | wendy@example.com | password123 | female | 88 | 84+ kg
-- Zoe | zoe@example.com | password123 | female | 89 | 84+ kg

-- +goose Up
INSERT INTO "User" ("username", "email", "password", "gender", "bodyWeight", "updatedAt")
VALUES
    ('Bob', 'bob@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 58, CURRENT_TIMESTAMP),
    ('Mark', 'mark@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 57, CURRENT_TIMESTAMP),
    ('Jacob', 'jacob@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 56, CURRENT_TIMESTAMP),
    ('John', 'john@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 65, CURRENT_TIMESTAMP),
    ('Mike', 'mike@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 64, CURRENT_TIMESTAMP),
    ('Alex', 'alex@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 63, CURRENT_TIMESTAMP),
    ('Ben', 'ben@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 73, CURRENT_TIMESTAMP),
    ('Chris', 'chris@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 72, CURRENT_TIMESTAMP),
    ('Daniel', 'daniel@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 71, CURRENT_TIMESTAMP),
    ('Ethan', 'ethan@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 82, CURRENT_TIMESTAMP),
    ('Finn', 'finn@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 81, CURRENT_TIMESTAMP),
    ('George', 'george@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 80, CURRENT_TIMESTAMP),
    ('Harry', 'harry@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 92, CURRENT_TIMESTAMP),
    ('Isaac', 'isaac@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 91, CURRENT_TIMESTAMP),
    ('Jack', 'jack@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 90, CURRENT_TIMESTAMP),
    ('Kevin', 'kevin@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 104, CURRENT_TIMESTAMP),
    ('Leo', 'leo@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 103, CURRENT_TIMESTAMP),
    ('Max', 'max@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 102, CURRENT_TIMESTAMP),
    ('Noah', 'noah@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 119, CURRENT_TIMESTAMP),
    ('Owen', 'owen@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 118, CURRENT_TIMESTAMP),
    ('Paul', 'paul@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 117, CURRENT_TIMESTAMP),
    ('Ryan', 'ryan@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 123, CURRENT_TIMESTAMP),
    ('Sam', 'sam@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 124, CURRENT_TIMESTAMP),
    ('Tom', 'tom@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'male', 125, CURRENT_TIMESTAMP),
    ('Anna', 'anna@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 46, CURRENT_TIMESTAMP),
    ('Bella', 'bella@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 45, CURRENT_TIMESTAMP),
    ('Chloe', 'chloe@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 44, CURRENT_TIMESTAMP),
    ('Daisy', 'daisy@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 51, CURRENT_TIMESTAMP),
    ('Emma', 'emma@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 50, CURRENT_TIMESTAMP),
    ('Faith', 'faith@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 49, CURRENT_TIMESTAMP),
    ('Grace', 'grace@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 56, CURRENT_TIMESTAMP),
    ('Hannah', 'hannah@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 55, CURRENT_TIMESTAMP),
    ('Ivy', 'ivy@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 54, CURRENT_TIMESTAMP),
    ('Jade', 'jade@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 62, CURRENT_TIMESTAMP),
    ('Kate', 'kate@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 61, CURRENT_TIMESTAMP),
    ('Lily', 'lily@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 60, CURRENT_TIMESTAMP),
    ('Mia', 'mia@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 68, CURRENT_TIMESTAMP),
    ('Nora', 'nora@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 67, CURRENT_TIMESTAMP),
    ('Olivia', 'olivia@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 66, CURRENT_TIMESTAMP),
    ('Paige', 'paige@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 75, CURRENT_TIMESTAMP),
    ('Quinn', 'quinn@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 74, CURRENT_TIMESTAMP),
    ('Rose', 'rose@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 73, CURRENT_TIMESTAMP),
    ('Sarah', 'sarah@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 83, CURRENT_TIMESTAMP),
    ('Tina', 'tina@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 82, CURRENT_TIMESTAMP),
    ('Uma', 'uma@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 81, CURRENT_TIMESTAMP),
    ('Vera', 'vera@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 87, CURRENT_TIMESTAMP),
    ('Wendy', 'wendy@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 88, CURRENT_TIMESTAMP),
    ('Zoe', 'zoe@example.com', '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA', 'female', 89, CURRENT_TIMESTAMP);

-- +goose Down
-- Match the seeded hash as well as email; preserve accounts whose password has changed.
-- Existing lift records prevent deletion through the foreign key; do not delete user activity.
DELETE FROM "User"
WHERE "password" = '$argon2id$v=19$m=65536,t=3,p=4$A8L3/OotcnWEvshotN9fAQ$VRksD+6u1x0HixGyFO0sU6vIl6+uVq4cnCtX8SkyuoA'
    AND "email" IN (
        'bob@example.com',
        'mark@example.com',
        'jacob@example.com',
        'john@example.com',
        'mike@example.com',
        'alex@example.com',
        'ben@example.com',
        'chris@example.com',
        'daniel@example.com',
        'ethan@example.com',
        'finn@example.com',
        'george@example.com',
        'harry@example.com',
        'isaac@example.com',
        'jack@example.com',
        'kevin@example.com',
        'leo@example.com',
        'max@example.com',
        'noah@example.com',
        'owen@example.com',
        'paul@example.com',
        'ryan@example.com',
        'sam@example.com',
        'tom@example.com',
        'anna@example.com',
        'bella@example.com',
        'chloe@example.com',
        'daisy@example.com',
        'emma@example.com',
        'faith@example.com',
        'grace@example.com',
        'hannah@example.com',
        'ivy@example.com',
        'jade@example.com',
        'kate@example.com',
        'lily@example.com',
        'mia@example.com',
        'nora@example.com',
        'olivia@example.com',
        'paige@example.com',
        'quinn@example.com',
        'rose@example.com',
        'sarah@example.com',
        'tina@example.com',
        'uma@example.com',
        'vera@example.com',
        'wendy@example.com',
        'zoe@example.com'
    );
