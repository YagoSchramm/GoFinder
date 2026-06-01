    INSERT INTO price_history (product_url, store, price)
        VALUES ($1, $2, $3)
        RETURNING id