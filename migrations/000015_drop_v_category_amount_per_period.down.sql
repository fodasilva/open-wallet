CREATE OR REPLACE VIEW v_category_amount_per_period AS
WITH periods AS (
    SELECT DISTINCT 
        TO_CHAR(reference_date, 'YYYYMM') AS period,
        user_id
    FROM entries e
    JOIN transactions t ON e.transaction_id = t.id
),
category_period_combinations AS (
    SELECT 
        c.id,
        c.user_id,
        c.name,
        c.color,
        p.period
    FROM categories c
    CROSS JOIN LATERAL (
        SELECT DISTINCT period 
        FROM periods 
        WHERE periods.user_id = c.user_id
    ) p
),
actual_amounts AS (
    SELECT 
        t.category_id,
        t.user_id,
        TO_CHAR(e.reference_date, 'YYYYMM') AS period,
        SUM(e.amount) AS total_amount
    FROM entries e
    JOIN transactions t ON e.transaction_id = t.id
    GROUP BY 
        t.category_id,
        t.user_id,
        TO_CHAR(e.reference_date, 'YYYYMM')
)
SELECT 
    cpc.id,
    cpc.user_id,
    cpc.name,
    cpc.color,
    cpc.period,
    COALESCE(aa.total_amount, 0) AS total_amount
FROM category_period_combinations cpc
LEFT JOIN actual_amounts aa 
    ON cpc.id = aa.category_id 
    AND cpc.user_id = aa.user_id
    AND cpc.period = aa.period
ORDER BY cpc.user_id, cpc.period, cpc.name;