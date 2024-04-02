
DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user" (
    "id" TEXT,
    "username" TEXT,
    "password" TEXT,

    PRIMARY KEY (
        "id"
    )
);
