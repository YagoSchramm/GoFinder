  SELECT id, product_url, store, price, captured_at
        FROM price_history
        WHERE product_url = $1
        ORDER BY captured_at DESC
        LIMIT 1