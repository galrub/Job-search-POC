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


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
        id UUID NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4()),
        name VARCHAR(100),
        email VARCHAR(255) NOT NULL UNIQUE,
        photo VARCHAR(255) NOT NULL DEFAULT 'default.png',
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

CREATE OR REPLACE FUNCTION refresh_update_at_function()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
  -- Check if the 'price' column was updated
  NEW.updated_at := now();
  RETURN NEW;
END;
$BODY$;

CREATE OR REPLACE TRIGGER refresh_users_update_at_trigger
BEFORE UPDATE ON "users"
FOR EACH ROW 
EXECUTE PROCEDURE refresh_update_at_function();
-- Dumping data for table public.user_entity: -1 rows
/*!40000 ALTER TABLE "user_entity" DISABLE KEYS */;
/*!40000 ALTER TABLE "user_entity" ENABLE KEYS */;

