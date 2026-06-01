        SELECT id, search_id, title, price, store, url, thumbnail, found_at
        FROM products
        WHERE url = $1
        LIMIT 1