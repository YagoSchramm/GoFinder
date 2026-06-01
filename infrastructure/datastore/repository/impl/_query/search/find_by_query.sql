        SELECT id, user_id, query, created_at
        FROM searches
        WHERE user_id = $1
          AND query = $2
          AND created_at >= $3
        ORDER BY created_at DESC
        LIMIT 1