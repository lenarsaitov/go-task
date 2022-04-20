DROP TABLE IF EXISTS account_records;
DROP INDEX IF EXISTS operation_id_idx;
DROP INDEX IF EXISTS delta_on_date_idx;

CREATE TABLE account_records (
             id           		    BIGSERIAL PRIMARY KEY,
             account_id 		    BIGINT,
             operation_id 		    BIGINT DEFAULT NULL,
             balance_delta 		    MONEY,
             balance_after 		    MONEY,
             balance_updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX operation_id_idx ON account_records (operation_id);
CREATE INDEX delta_on_date_idx ON account_records (account_id, balance_updated_at);

INSERT INTO account_records (account_id, operation_id, balance_delta, balance_after, balance_updated_at)
SELECT round(random() * 100), round(random() * 100), i+100, i+10000, (now() - interval '30 day' * round(random() * 100))
FROM generate_series(1, 1000000) AS i;

EXPLAIN (ANALYZE) SELECT balance_after FROM account_records WHERE operation_id = 340000;

EXPLAIN (ANALYZE) SELECT balance_delta FROM account_records WHERE account_id = 895001 AND balance_updated_at = '2022-04-20';

EXPLAIN (ANALYZE) INSERT INTO account_records (account_id, operation_id, balance_delta, balance_after)
VALUES (4, 1, 4000.22223, 5000.22223);