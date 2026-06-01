    SELECT id, search_id, title, price, store, url, thumbnail, found_at
        FROM products
        WHERE search_id = $1
        ORDER BY price ASC