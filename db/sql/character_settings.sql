-- ----------------------------
-- Table structure for character_settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."character_settings";
CREATE TABLE "public"."character_settings" (
  "char_id" int4 NOT NULL,
  "exp" bool,
  "sp" bool,
  "autoloot" bool,
  "trade_all" bool,
  "trade_party" bool,
  "trade_clan" bool,
  "trade_self_ip" bool,
  "party_all" bool,
  "party_clan" bool,
  "party_self_ip" bool,
  "soulshot_holo" bool
)
;
