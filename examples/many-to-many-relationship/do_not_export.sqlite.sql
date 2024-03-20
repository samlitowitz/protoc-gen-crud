DROP TABLE IF EXISTS "user_role";
CREATE TABLE IF NOT EXISTS "user_role" (
    "id" TEXT,
    "user_id" TEXT,
    "role_id" TEXT,

    PRIMARY KEY ("id")
);
