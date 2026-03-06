DROP VIEW IF EXISTS v_entries;

ALTER TABLE transactions
ADD COLUMN recurrence_id VARCHAR REFERENCES recurrences(id);

ALTER TABLE transactions
ADD CONSTRAINT chk_transactions_recurrence_id CHECK (
    (category = 'recurrence' AND recurrence_id IS NOT NULL) OR
    (category != 'recurrence' AND recurrence_id IS NULL)
);

CREATE OR REPLACE VIEW v_entries AS
SELECT
    e.id,
    e.transaction_id,
    t.name,
    t.description,
    e.amount,
    left(regexp_replace(e.reference_date::text, '[^0-9]', '', 'g'), 6) AS period,
    t.user_id,
    t.category,
    sum(e.amount) OVER (PARTITION BY e.transaction_id) AS total_amount,
    row_number() OVER (PARTITION BY e.transaction_id ORDER BY e.reference_date) AS installment,
    count(*) OVER (PARTITION BY e.transaction_id) AS total_installments,
    e.created_at,
    e.reference_date,
    c.id AS category_id,
    c.name AS category_name,
    c.color AS category_color,
    t.recurrence_id
FROM
    entries e
JOIN transactions t ON e.transaction_id = t.id
LEFT JOIN categories c ON t.category_id = c.id AND t.user_id = c.user_id;
