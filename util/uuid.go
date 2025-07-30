package util

import (
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

func ConvertToPgUUID(id uuid.UUID) pgtype.UUID {
    var pgUUID pgtype.UUID
	copy(pgUUID.Bytes[:], id[:]) // copy bytes into fixed array
    pgUUID.Valid = true
    return pgUUID
}
