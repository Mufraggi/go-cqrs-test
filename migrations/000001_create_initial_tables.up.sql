CREATE TABLE "user" ("id" TEXT NOT NULL, "email_address" TEXT NOT NULL, "first_name" TEXT NOT NULL, "last_name" TEXT NOT NULL, PRIMARY KEY ("id"));

CREATE TABLE "article" ("id" TEXT NOT NULL, "user_id" TEXT NOT NULL, "title" TEXT NOT NULL, "content" TEXT NOT NULL, CONSTRAINT "article_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE, PRIMARY KEY ("id"));
CREATE INDEX "article_user_id_index" ON "article" ("user_id");

CREATE TABLE "clap" ("id" TEXT NOT NULL, "article_id" TEXT NOT NULL, "count" INTEGER NOT NULL, CONSTRAINT "clap_article_id_foreign" FOREIGN KEY("article_id") REFERENCES "article"("id") ON DELETE CASCADE ON UPDATE CASCADE, PRIMARY KEY ("id"));
CREATE INDEX "clap_article_id_index" ON "clap" ("article_id");