package converter

import "github.com/jackc/pgx/v5/pgtype"

func ToFloat8(v float64) pgtype.Float8 {
	return pgtype.Float8{
		Float64: v,
		Valid:   true,
	}
}
