CREATE TABLE "characters" (
  "login" varchar(25) COLLATE "pg_catalog"."default" NOT NULL,
  "object_id" int4 NOT NULL,
  "char_name" varchar(16) COLLATE "pg_catalog"."default" NOT NULL,
  "level" int2 NOT NULL DEFAULT 1,
  "max_hp" int4 NOT NULL DEFAULT 100,
  "cur_hp" int4 NOT NULL DEFAULT 100,
  "max_mp" int4 NOT NULL DEFAULT 100,
  "cur_mp" int4 NOT NULL DEFAULT 100,
  "face" int2 NOT NULL,
  "hair_style" int2 NOT NULL,
  "hair_color" int2 NOT NULL,
  "sex" int2 NOT NULL,
  "x" int4 NOT NULL,
  "y" int4 NOT NULL,
  "z" int4 NOT NULL,
  "exp" int8 NOT NULL DEFAULT 0,
  "sp" int8 NOT NULL DEFAULT 0,
  "karma" int4 NOT NULL DEFAULT 0,
  "pvp_kills" int4 NOT NULL DEFAULT 0,
  "pk_kills" int4 NOT NULL DEFAULT 0,
  "clan_id" int4 NOT NULL DEFAULT 0,
  "race" int2 NOT NULL,
  "class_id" int4 NOT NULL,
  "base_class" int4 NOT NULL DEFAULT 0,
  "title" varchar(16) COLLATE "pg_catalog"."default",
  "online_time" int4 NOT NULL DEFAULT 0,
  "nobless" int4 NOT NULL DEFAULT 0,
  "vitality" int4 NOT NULL DEFAULT 20000
);

-- ----------------------------
-- Indexes structure for table characters
-- ----------------------------
CREATE UNIQUE INDEX "table_name_char_id_uindex_copy1" ON "public"."characters" USING btree (
  "object_id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table characters
-- ----------------------------
ALTER TABLE "characters" ADD CONSTRAINT "characters_copy1_pkey" PRIMARY KEY ("object_id");
