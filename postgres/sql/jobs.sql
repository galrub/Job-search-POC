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

CREATE TABLE IF NOT EXISTS public.jobs
(
    id uuid NOT NULL DEFAULT (uuid_generate_v4()),
    user_id uuid NOT NULL,
    company VARCHAR(25) COLLATE pg_catalog."default" NOT NULL,
    position_desc VARCHAR(100) COLLATE pg_catalog."default" NOT NULL,
    remote boolean DEFAULT true,
    contract_type VARCHAR(25) COLLATE pg_catalog."default" DEFAULT 'Contract'::character varying,
    contacted boolean DEFAULT false,
    general_status VARCHAR(50) NOT NULL DEFAULT 'applied',
    created_at date DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    comments VARCHAR(200) COLLATE pg_catalog."default",
    CONSTRAINT job_id_pk PRIMARY KEY (id),
    CONSTRAINT jobs_fk_user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION refresh_update_at_function() RETURNS trigger AS 
$$
BEGIN
  -- Check if the 'price' column was updated
  NEW.updated_at := now();
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER refresh_jobs_update_at_trigger
BEFORE UPDATE ON jobs
FOR EACH ROW 
EXECUTE PROCEDURE refresh_update_at_function();
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
