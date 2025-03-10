// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package generated

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminScale string

const (
	AdminScaleMinor AdminScale = "minor"
	AdminScaleMajor AdminScale = "major"
)

func (e *AdminScale) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AdminScale(s)
	case string:
		*e = AdminScale(s)
	default:
		return fmt.Errorf("unsupported scan type for AdminScale: %T", src)
	}
	return nil
}

type NullAdminScale struct {
	AdminScale AdminScale
	Valid      bool // Valid is true if AdminScale is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAdminScale) Scan(value interface{}) error {
	if value == nil {
		ns.AdminScale, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AdminScale.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAdminScale) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AdminScale), nil
}

type User struct {
	ID        uuid.UUID
	Username  string
	Pseudonym string
	FirstName string
	LastName  string
	IsDeleted bool
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type UsersAdmin struct {
	UserID    uuid.UUID
	Scale     AdminScale
	CreatedAt pgtype.Timestamp
}
