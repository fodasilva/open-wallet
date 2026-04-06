drop view if exists v_entries;

alter table entries add column period varchar(6);

-- Restore data from reference_date to keep it consistent
update entries set period = left(regexp_replace(reference_date::text, '[^0-9]', '', 'g'), 6);

-- Make period not null as it was originally
alter table entries alter column period set not null;

create or replace view v_entries as
select
	e.id,
    e.transaction_id,
	t.name,
	t.description,
	e.amount,
	e.period,
	t.user_id,
	t.category,
	sum(e.amount) over (partition by e.transaction_id) as total_amount,
	row_number() over (partition by e.transaction_id order by e.period) as installment,
	count(*) over (partition by e.transaction_id) as total_installments,
    e.created_at,
    e.reference_date
from
	entries e
join transactions t on
	e.transaction_id = t.id;
