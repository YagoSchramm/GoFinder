 INSERT INTO alerts (user_id, product_url, store, target_price)
        VALUES ($1, $2, $3, $4)
        RETURNING id