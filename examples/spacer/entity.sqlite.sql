DROP TABLE IF EXISTS "entity_tags";
CREATE TABLE IF NOT EXISTS "entity_tags" (
    "entity_id" TEXT
    "tag_id" TEXT
);

DROP TABLE IF EXISTS "entity";
CREATE TABLE IF NOT EXISTS "entity" (
    "id" TEXT,
    "description" TEXT,

    PRIMARY KEY ("id")
);
