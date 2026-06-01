       SELECT id, user_id, query, created_at
        FROM searches
        WHERE user_id = $1
        ORDER BY created_at DESC