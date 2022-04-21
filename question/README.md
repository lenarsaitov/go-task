### Таблица account_records (изменение баланса)
Изменение баланса может быть у операции, тогда в поле ___operation_id___ будет вставлен __id__ операции

Предполагается наличие следующих полей:
1. __id__ - Primary Key
2. __account_id__ - ид счета
3. __operation_id__ default null - ид операции
4. __balance_delta__ - изменение баланса на сумму с учетом знака
5. __balance_after__ - значение баланса после обновления
6. __balance_updated_at__ - время изменения баланса timestampz

Используемые запросы к данной таблице:
 - найти баланс у определенной операции (по __operation_id__)
 - изменение баланса у счета (__account_id__) за определенную дату (__balance_updated_at__)

### Создание таблицы account_records 
Написанные ниже запросы к PostgreSQL имеются в файле **_schema.sql_**

Предполагаю наличие следующих типов у каждого столбца:

___BIGSERIAL, BIGINT, MONEY___

В соответствии с этим создание таблицы осуществляется следующим образом:
```sql
DROP TABLE IF EXISTS account_records;

CREATE TABLE account_records (
    id           		BIGSERIAL PRIMARY KEY,
    account_id                  BIGINT NOT NULL,
    operation_id 		BIGINT DEFAULT 0 NOT NULL,
    balance_delta 		MONEY,
    balance_after 		MONEY,
    balance_updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
```

### Запросы 
Описанные выше запросы можно представить следующим образом:
```sql
SELECT balance_after FROM account_records WHERE operation_id = 111;
```

```sql
SELECT balance_delta FROM account_records WHERE account_id = 22 AND balance_updated_at >= to_timestamp('2022-04-20','YYYY-MM-DD') AND balance_updated_at < to_timestamp('2022-04-21','YYYY-MM-DD');
```

### Оптимизация 

__Предполагается высокая частотность использования вышеописанных запросов__.

В нашем случае оптимально использовать индексы типа __B-дерево__.

```sql
CREATE INDEX operation_id_idx ON account_records (operation_id);
```

По нескольким столбцам (___составные индексы___):
```sql
CREATE INDEX delta_on_date_idx ON account_records (account_id, balance_updated_at);
```
### Генерация данных

Сгенерируем тестовые данные на 1 миллион строк:
```sql
INSERT INTO account_records (account_id, operation_id, balance_delta, balance_after, balance_updated_at)
SELECT round(random() * 100), round(random() * 100), i+100, i+10000, (now() - interval '30 day' * round(random() * 100))
FROM generate_series(1, 1000000) AS i;
```

### Explain-анализ без использования индексов
```sql
EXPLAIN (ANALYZE) SELECT balance_after FROM account_records WHERE operation_id = 340000;
```
```bash
                                                           QUERY PLAN                                                            
---------------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..16054.45 rows=5000 width=8) (actual time=131.658..135.994 rows=0 loops=1)
   Workers Planned: 2
   Workers Launched: 2
   ->  Parallel Seq Scan on account_records  (cost=0.00..14554.45 rows=2083 width=8) (actual time=124.721..124.722 rows=0 loops=3)
         Filter: (operation_id = 340000)
         Rows Removed by Filter: 333333
 Planning Time: 0.065 ms
 Execution Time: 136.008 ms
```
```sql
EXPLAIN (ANALYZE) SELECT balance_delta FROM account_records WHERE account_id = 895001 AND balance_updated_at >= to_timestamp('2022-04-20','YYYY-MM-DD') AND balance_updated_at < to_timestamp('2022-04-21','YYYY-MM-DD');
```
```bash
                                                           QUERY PLAN                                                            
---------------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..19723.71 rows=25 width=8) (actual time=125.073..128.630 rows=0 loops=1)
   Workers Planned: 2
   Workers Launched: 2
   ->  Parallel Seq Scan on account_records  (cost=0.00..18721.21 rows=10 width=8) (actual time=116.687..116.687 rows=0 loops=3)
         Filter: ((account_id = 895001) AND (balance_updated_at >= to_timestamp('2022-Dec-20'::text, 'YYYY-Mon-DD'::text)) AND (balance_updated_at < to_timestamp('2022-Dec-21'::text, 'YYYY-Mon-DD'::text)))
         Rows Removed by Filter: 333333
 Planning Time: 0.124 ms
 Execution Time: 128.649 ms
```
### Explain-анализ с использованием индексов
```sql
EXPLAIN (ANALYZE) SELECT balance_after FROM account_records WHERE operation_id = 500001;
```
```bash
                                                        QUERY PLAN                                                         
---------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on account_records  (cost=59.17..8212.52 rows=5000 width=8) (actual time=0.009..0.009 rows=0 loops=1)
   Recheck Cond: (operation_id = 340000)
   ->  Bitmap Index Scan on operation_id_idx  (cost=0.00..57.92 rows=5000 width=0) (actual time=0.008..0.008 rows=0 loops=1)
         Index Cond: (operation_id = 340000)
 Planning Time: 0.078 ms
 Execution Time: 0.023 ms
```
```sql
EXPLAIN (ANALYZE) SELECT balance_delta FROM account_records WHERE account_id = 895001 AND balance_updated_at >= to_timestamp('2022-04-20','YYYY-MM-DD') AND balance_updated_at < to_timestamp('2022-04-21','YYYY-MM-DD');
```
```bash
                                                        QUERY PLAN                                                         
---------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on account_records  (cost=4.75..101.43 rows=25 width=8) (actual time=0.015..0.015 rows=0 loops=1)
   Recheck Cond: ((account_id = 895001) AND (balance_updated_at >= to_timestamp('2022-Dec-20'::text, 'YYYY-Mon-DD'::text)) AND (balance_updated_at < to_timestamp('2022-Dec-21'::text, 'YYYY-Mon-DD'::text)))
   ->  Bitmap Index Scan on delta_on_date_idx  (cost=0.00..4.74 rows=25 width=0) (actual time=0.014..0.014 rows=0 loops=1)
         Index Cond: ((account_id = 895001) AND (balance_updated_at >= to_timestamp('2022-Dec-20'::text, 'YYYY-Mon-DD'::text)) AND (balance_updated_at < to_timestamp('2022-Dec-21'::text, 'YYYY-Mon-DD'::text)))
 Planning Time: 0.077 ms
 Execution Time: 0.031 ms
```

### Выполнение INSERT запросов
```sql
EXPLAIN (ANALYZE) INSERT INTO account_records (account_id, operation_id, balance_delta, balance_after)
VALUES (4, 1, 4000.22223, 5000.22223);
```
Без использования индексов:
```bash
                                              QUERY PLAN                                                           
-------------------------------------------------------------------------------------------------------------------------------
 Insert on account_records  (cost=0.00..0.02 rows=1 width=48) (actual time=0.026..0.027 rows=0 loops=1)
   ->  Result  (cost=0.00..0.02 rows=1 width=48) (actual time=0.010..0.010 rows=1 loops=1)
 Planning Time: 0.024 ms
 Execution Time: 0.040 ms
```
С использованием вышеописанных индексов:
```bash
                                               QUERY PLAN                                               
--------------------------------------------------------------------------------------------------------
 Insert on account_records  (cost=0.00..0.02 rows=1 width=48) (actual time=0.047..0.047 rows=0 loops=1)
   ->  Result  (cost=0.00..0.02 rows=1 width=48) (actual time=0.010..0.011 rows=1 loops=1)
 Planning Time: 0.022 ms
 Execution Time: 0.066 ms
```

### Анализ результатов
Как видно применение индексов существенно увеличивает производительность выполнения вышеописанных **SELECT** запросов (с **136.008** и **128.649** _**ms**_ до **0.023** и **0.031** **_ms_**, соответственно)

Однако, в то же время, имеется негативное влияние на запросы вида **INSERT**, **UPDATE**, **DELETE** (прежде всего на первое). В случае **INSERT** наблюдается ухудшение производительности с разницей в **полтора раза**. 

Поэтому подробно учитывая логику сервиса, а также ожидаемые свойства, стоит дополнительно исследовать вопрос компромисса (к примеру, использовать только один из двух индексов).
Либо, в случае отсутствия необходимости ожидать завершения INSERT запроса, полностью отдавать приоритет вышеописанному индексированию.
