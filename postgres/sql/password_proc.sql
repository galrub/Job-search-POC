-- Add up migration script here
-- --------------------------------------------------------
-- Host:                         192.168.1.13
-- Server version:               PostgreSQL 15.0 on x86_64-pc-linux-musl, compiled by gcc (Alpine 11.2.1_git20220219) 11.2.1 20220219, 64-bit
-- Server OS:                    
-- HeidiSQL Version:             11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES  */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE TABLE IF NOT EXISTS "password_store" (
	"id" SERIAL NOT NULL,
	"user_id" UUID NOT NULL,
	"pword" VARCHAR(255) NOT NULL,
	"salt" VARCHAR(100) NOT NULL,
	CONSTRAINT "FK__user_entity" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS pgcrypto;

create or replace procedure gen_user(vEmail VARCHAR(255), vPword VARCHAR(255) )
LANGUAGE plpgsql 
AS $$
DECLARE
	userId UUID;
	salt VARCHAR(100) = gen_salt('bf');
BEGIN
	INSERT INTO users(email) VALUES (vEmail);
	COMMIT;
	SELECT id INTO userId FROM users WHERE email = vEmail;
	INSERT INTO password_store(user_id, salt, pword) 
	VALUES (userId, salt, crypt(vPword, salt));
	COMMIT;
END;$$;

create or replace function verify_user(vUserId UUID, vPword VARCHAR(255))
RETURNS BOOL
LANGUAGE plpgsql 
AS $$
DECLARE
	vVerified  BOOL;
	vSalt      VARCHAR(255);
	vStoredPW  VARCHAR(255);
BEGIN
	SELECT salt  INTO vSalt     FROM password_store WHERE user_id = vUserId;
	SELECT pword INTO vStoredPW FROM password_store WHERE user_id = vUserId;
	SELECT crypt(vPword, vSalt) = vStoredPW INTO vVerified;
	RETURN vVerified;
END;$$;


/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
