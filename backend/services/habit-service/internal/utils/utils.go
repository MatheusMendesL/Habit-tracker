package utils

import "database/sql"

func ToNullString(param string) sql.NullString {
	return sql.NullString{
		String: param,
		Valid:  true,
	}
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
