BEGIN;

ALTER TABLE portfolio_initiative_membership ADD PRIMARY KEY (portfolio_id, initiative_id);

COMMIT;