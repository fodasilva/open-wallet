CREATE VIEW v_summaries AS
SELECT
    user_id,
    period,
    SUM(CASE WHEN category != 'income' THEN amount ELSE 0 END) AS total_expense,
    SUM(CASE WHEN category = 'income' THEN amount ELSE 0 END) AS total_income,
    SUM(CASE WHEN category = 'income' THEN amount ELSE -amount END) AS total_balance
FROM
    v_entries
GROUP BY
    user_id,
    period;
