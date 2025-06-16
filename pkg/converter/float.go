package converter

import "github.com/jackc/pgx/v5/pgtype"

func Float64(v pgtype.Float8) float64 {
	if v.Valid {
		return v.Float64
	}
	return 0
}

func ToFloat8(v float64) pgtype.Float8 {
	return pgtype.Float8{
		Float64: v,
		Valid:   true,
	}
}
