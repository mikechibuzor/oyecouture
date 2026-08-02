WITH ranked_attempts AS (
    SELECT id,
        ROW_NUMBER() OVER (
            PARTITION BY order_id
            ORDER BY updated_at DESC, created_at DESC
        ) AS attempt_rank
    FROM payment_attempts
    WHERE status = 'initialized'
)
UPDATE payment_attempts
SET status = 'failed',
    updated_at = NOW()
WHERE id IN (
    SELECT id
    FROM ranked_attempts
    WHERE attempt_rank > 1
);

CREATE UNIQUE INDEX payment_attempts_one_initialized_attempt_per_order_idx
    ON payment_attempts (order_id)
    WHERE status = 'initialized';
