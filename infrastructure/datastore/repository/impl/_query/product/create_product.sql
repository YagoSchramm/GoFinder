       INSERT INTO products (search_id, title, price, store, url, thumbnail)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id