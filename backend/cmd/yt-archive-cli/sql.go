package main

func execRowsAffected(query string, args ...any) (int64, error) {
	r, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}
