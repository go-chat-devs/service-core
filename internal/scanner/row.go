package scanner

import "github.com/jackc/pgx/v5"

func Row[T Scannable](row pgx.Row) (T, error) {
	var t T
	err := t.FromRow(row)
	return t, err
}
