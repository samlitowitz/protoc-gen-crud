
DROP TABLE IF EXISTS "primary_key_enum";
CREATE TABLE IF NOT EXISTS "primary_key_enum" (
    "id" INTEGER PRIMARY KEY,
    "value" TEXT
);

INSERT INTO "primary_key_enum" ("id", "value") VALUES
    (0, 'ZERO'),
    (1, 'ONE'),
    (2, 'TWO'),
    (3, 'THREE')
;

DROP TABLE IF EXISTS "sa_enum";
CREATE TABLE IF NOT EXISTS "sa_enum" (
    "id" INTEGER /* references "primary_key_enum"."id" */,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "sa_int_32";
CREATE TABLE IF NOT EXISTS "sa_int_32" (
    "id" INTEGER,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "sa_int_64";
CREATE TABLE IF NOT EXISTS "sa_int_64" (
    "id" INTEGER,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "sa_uint_32";
CREATE TABLE IF NOT EXISTS "sa_uint_32" (
    "id" INTEGER,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "sa_uint_64";
CREATE TABLE IF NOT EXISTS "sa_uint_64" (
    "id" INTEGER,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "sa_string";
CREATE TABLE IF NOT EXISTS "sa_string" (
    "id" TEXT,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);

DROP TABLE IF EXISTS "primary_key_enum";
CREATE TABLE IF NOT EXISTS "primary_key_enum" (
    "id" INTEGER PRIMARY KEY,
    "value" TEXT
);

INSERT INTO "primary_key_enum" ("id", "value") VALUES
    (0, 'ZERO'),
    (1, 'ONE'),
    (2, 'TWO'),
    (3, 'THREE')
;

DROP TABLE IF EXISTS "ma_all";
CREATE TABLE IF NOT EXISTS "ma_all" (
    "id_uint_64" INTEGER,
    "id_string" TEXT,
    "id_enum" INTEGER /* references "primary_key_enum"."id" */,
    "id_int_32" INTEGER,
    "id_int_64" INTEGER,
    "id_uint_32" INTEGER,
    "data" TEXT,

    PRIMARY KEY (
        "id_uint_64",
        "id_string",
        "id_enum",
        "id_int_32",
        "id_int_64",
        "id_uint_32"
    )
);
