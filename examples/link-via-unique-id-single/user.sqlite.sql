
DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user" (
    "id" TEXT,
    "username" TEXT,
    "password" TEXT,
    "profile_id" TEXT,

    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "profile";
CREATE TABLE IF NOT EXISTS "profile" (
    "id" TEXT,
    "name" TEXT,

    PRIMARY KEY ("id")
);
