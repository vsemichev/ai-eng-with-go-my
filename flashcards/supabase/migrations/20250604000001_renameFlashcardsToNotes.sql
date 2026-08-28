ALTER TABLE gocourse.flashcards RENAME TO notes;

ALTER INDEX gocourse.idx_flashcards_created_at RENAME TO idx_notes_created_at;
