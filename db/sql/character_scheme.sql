DROP TABLE IF EXISTS "character_scheme";
CREATE TABLE "character_scheme" (
  "id" int4 NOT NULL DEFAULT nextval('character_buffs_save_list_id_seq'::regclass),
  "char_id" int4,
  "name" varchar(32) COLLATE "pg_catalog"."default"
)
;
