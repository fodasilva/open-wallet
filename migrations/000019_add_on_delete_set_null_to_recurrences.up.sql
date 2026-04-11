ALTER TABLE recurrences
DROP CONSTRAINT IF EXISTS recurrences_category_id_fkey,
ADD CONSTRAINT recurrences_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;
