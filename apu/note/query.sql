CREATE TABLE IF NOT EXISTS public.authors
(
    author_id integer                           NOT NULL DEFAULT nextval('authors_author_id_seq'::regclass),
    name      text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    biography jsonb,
    CONSTRAINT authors_pkey PRIMARY KEY (author_id)
);

create type book_type as enum (
    'FICTION',
    'NONFICTION'
    );

create table books
(
    book_id   serial primary key,
    author_id integer                  not null references authors (author_id),
    isbn      text                     not null default '' unique,
    book_type book_type                not null default 'FICTION',
    title     text                     not null default '',
    year      integer                  not null default 2024,
    available timestamp with time zone not null default 'NOW()',
    tags      varchar[]                not null default '{}'
);

-- name: GetAuthor :one
SELECT *
FROM authors
WHERE author_id = $1;

-- name: CreateAuthor :one
INSERT INTO authors (name)
VALUES ($1)
RETURNING *;

-- name: CreateAuthors :exec
INSERT INTO authors
SELECT unnest(@ids::bigint[]) AS id,
       unnest(@names::text[]) as name;

-- name: CreateBook :batchone
INSERT INTO books (author_id,
                   isbn,
                   book_type,
                   title,
                   year,
                   available,
                   tags)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7)
RETURNING *;

-- name: UpdateBook :batchexec
update books
set title = $1, tags = $2
where book_id = $3;

-- name: BooksByYear :batchmany
SELECT * FROM books
WHERE year = $1;

-- name: DeleteBookExecResult :execresult
DELETE
FROM books
WHERE book_id = $1;

-- name: DeleteBook :batchexec
DELETE
FROM books
WHERE book_id = $1;

-- name: DeleteBookNamedFunc :batchexec
DELETE
FROM books
WHERE book_id = sqlc.arg(book_id);

-- name: DeleteBookNamedSign :batchexec
DELETE
FROM books
WHERE book_id = @book_id;