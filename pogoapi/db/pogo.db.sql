BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"comment"	TEXT NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"user_id"	INTEGER NOT NULL,
	FOREIGN KEY("post_id") REFERENCES "posts"("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"body"	TEXT NOT NULL,
	"created_at"	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at"	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"user_id"	INTEGER NOT NULL,
	"likes"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id"),
	FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"username"	TEXT NOT NULL,
	"email"	TEXT NOT NULL,
	"password"	TEXT NOT NULL
);
INSERT INTO "comments" VALUES (1,'Comment One',1,2);
INSERT INTO "comments" VALUES (2,'Comment Two',1,2);
INSERT INTO "comments" VALUES (3,'Comment Three',2,1);
INSERT INTO "posts" VALUES (1,'Post One','Post Body One','2020-03-20 20:25:54','2020-03-20 20:25:54',1,34);
INSERT INTO "posts" VALUES (2,'Post Two','Post Body Two','2020-03-20 20:26:23','2020-03-20 20:26:23',1,168);
INSERT INTO "posts" VALUES (3,'Post Three','Post Body Three','2020-03-20 20:27:09','2020-03-20 20:27:09',1,46);
INSERT INTO "posts" VALUES (4,'Post Four','Post Body Four','2020-03-20 20:27:35','2020-03-20 20:27:35',2,98);
INSERT INTO "posts" VALUES (5,'Post Five','Post Body Five','2020-03-20 20:27:55','2020-03-20 20:27:55',2,450);
INSERT INTO "users" VALUES (1,'cawis','cawis@test.com','1234');
INSERT INTO "users" VALUES (2,'zidak','zidak@test.com','1111');
COMMIT;
