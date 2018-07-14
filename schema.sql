-- Users
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "username" CHARACTER VARYING (100) NOT NULL,
  "score" DECIMAL NOT NULL DEFAULT 0,
  "created" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated" TIMESTAMP(0) WITH TIME ZONE
);
ALTER TABLE "users" ADD CONSTRAINT "users_id" PRIMARY KEY ("id");
ALTER TABLE "users" ADD CONSTRAINT "users_username" UNIQUE ("username");

-- User profile
CREATE TABLE "user_profile" (
  "user_id" INTEGER NOT NULL,
  "firstname" CHARACTER VARYING (100) NOT NULL,
  "lastname" CHARACTER VARYING (100) NOT NULL
);

ALTER TABLE "user_profile" ADD CONSTRAINT "user_id" UNIQUE ("user_id");
ALTER TABLE "user_profile"
  ADD CONSTRAINT fk_user_id
  FOREIGN KEY (user_id)
  REFERENCES users(id);

-- Actions
CREATE TABLE "actions" (
  "id" serial NOT NULL,
  "description" text NOT NULL,
  "created" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE "actions" ADD CONSTRAINT "actions_id" PRIMARY KEY ("id");
