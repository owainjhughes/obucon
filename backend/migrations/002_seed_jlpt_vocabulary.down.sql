-- Rollback: Remove seeded JLPT vocabulary
-- Only deletes pre-populated dictionary entries (jlpt_level IS NOT NULL)

DELETE FROM japanese_dictionary WHERE jlpt_level IS NOT NULL;
