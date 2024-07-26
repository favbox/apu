CREATE TYPE public.note_type AS ENUM
    ('normal', 'video');

CREATE TYPE public.source_type AS ENUM
    ('weixin', 'zhimo', 'xhs', 'behance', 'archdaily', 'unset');

CREATE TABLE IF NOT EXISTS public.notes
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    title character varying(64) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    description character varying(1000) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    tag_ids bigint[] NOT NULL DEFAULT '{}'::bigint[],
    post_time timestamp without time zone,
    type note_type NOT NULL DEFAULT 'normal'::note_type,
    user_id bigint NOT NULL DEFAULT 0,
    video_id bigint NOT NULL DEFAULT 0,
    is_privacy boolean NOT NULL DEFAULT false,
    is_enabled boolean NOT NULL DEFAULT false,
    create_time timestamp without time zone NOT NULL DEFAULT now(),
    update_time timestamp without time zone,
    source_type source_type NOT NULL DEFAULT 'unset'::source_type,
    source_url character varying(256) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    CONSTRAINT notes_pkey PRIMARY KEY (id),
    CONSTRAINT notes_uk_source UNIQUE (source_type, source_url)
);


CREATE TABLE IF NOT EXISTS public.images
(
    id bigint NOT NULL,
    file_id character varying(128) COLLATE pg_catalog."default" NOT NULL,
    width integer NOT NULL,
    height integer NOT NULL,
    extra_info jsonb NOT NULL,
    note_id bigint NOT NULL,
    sort smallint NOT NULL DEFAULT 0,
    is_cover boolean NOT NULL DEFAULT false,
    CONSTRAINT images_pkey PRIMARY KEY (id)
)