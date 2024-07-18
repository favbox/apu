-- name: CreateOriginalURL :execlastid
INSERT INTO original_url (url) VALUES
(?)
ON DUPLICATE KEY UPDATE updated_at=now();