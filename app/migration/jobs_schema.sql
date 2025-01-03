CREATE TABLE IF NOT EXISTS "jobs"
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    company varchar(25) NOT NULL,
    position_desc varchar(100) NOT NULL, 
    remote boolean DEFAULT true,
    contract_type varchar(25) NOT NULL,
    contacted boolean DEFAULT false,
    general_status varchar(50) NOT NULL DEFAULT 'applied',
    created_at date DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    comments varchar(200),
    CONSTRAINT job_id_pk PRIMARY KEY (id),
    CONSTRAINT jobs_fk_user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

