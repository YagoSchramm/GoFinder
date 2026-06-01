     SELECT id, user_id, product_url, store, target_price, active, created_at
        FROM alerts
        WHERE id = $1