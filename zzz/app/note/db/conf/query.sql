-- name: CreateNotes :batchone
-- 批量创建笔记，重复仅更新时间以标记
insert into notes (title, description, tag_ids, post_time, type, is_privacy,
                   source_type, source_url)
values (@title, @description, @tag_ids, @post_time, @type, @is_privacy,
        @source_type, @source_url)
on conflict (source_type, source_url) do update
    set update_time = now()
returning id;


-- name: GetNote :one
SELECT *
FROM notes
WHERE id = $1
LIMIT 1;

-- name: ListImagesByNoteID :many
SELECT *
FROM images
WHERE note_id = $1
order by sort;