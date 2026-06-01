        SELECT id, user_id, query, created_at
        FROM searches
        WHERE id = $1