        INSERT INTO searches (user_id, query)
        VALUES ($1, $2)
        RETURNING id