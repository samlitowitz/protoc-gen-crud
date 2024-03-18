DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user" (
    "id" TEXT,
    "username" TEXT,
    "password" TEXT,

    PRIMARY KEY ("id")
);
DROP TABLE IF EXISTS "role";
CREATE TABLE IF NOT EXISTS "role" (
    "id" TEXT,
    "name" TEXT,

    PRIMARY KEY ("id")
);
