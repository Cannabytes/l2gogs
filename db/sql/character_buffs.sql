DROP TABLE IF EXISTS "character_buffs";
CREATE TABLE "character_buffs" (
  "id" int4 NOT NULL DEFAULT nextval('character_buffs_id_seq'::regclass),
  "char_id" int4,
  "skill_id" int4,
  "level" int4,
  "second" int4
);
