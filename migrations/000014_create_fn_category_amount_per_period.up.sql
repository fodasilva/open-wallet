CREATE OR REPLACE FUNCTION fn_category_amount_per_period(
  p_period text
)
RETURNS TABLE (
  id  text,
  user_id      text,
  name         text,
  color        text,
  period       text,
  total_amount numeric
)
LANGUAGE sql
STABLE
AS $$
  SELECT
    c.id,
    c.user_id::text   AS user_id,
    c.name,
    c.color,
    p_period          AS period,
    COALESCE(SUM(e.amount), 0::numeric) AS total_amount
  FROM categories c
  LEFT JOIN transactions t
    ON t.category_id = c.id
   AND t.user_id::text = c.user_id::text
  LEFT JOIN entries e
    ON e.transaction_id = t.id
   AND to_char(e.reference_date, 'YYYYMM') = p_period
  GROUP BY c.id, c.user_id, c.name, c.color
  ORDER BY c.user_id::text, c.name;
$$;