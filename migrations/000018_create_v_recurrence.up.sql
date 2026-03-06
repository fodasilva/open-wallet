create or replace view v_recurrences as
select
    r.id,
    r.user_id,
    r.name,
    r.note,
    r.amount,
    r.day_of_month,
    r.start_period,
    r.end_period,
    r.category_id,
    r.created_at,
    c.name as category_name,
    c.color as category_color
from
    recurrences r
left join categories c on
    r.category_id = c.id
    and r.user_id = c.user_id;
