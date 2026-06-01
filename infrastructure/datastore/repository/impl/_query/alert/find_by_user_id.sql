   SELECT id, user_id, product_url, store, target_price, active, created_at
        FROM alerts
        WHERE user_id = $1
        ORDER BY created_at DESC