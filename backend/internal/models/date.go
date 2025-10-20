package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Date is a custom type for handling date-only fields (YYYY-MM-DD)
type Date struct {
	time.Time
}

const dateFormat = "2006-01-02"

// UnmarshalJSON handles JSON unmarshaling for Date fields
func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	d.Time = t
	return nil
}

// MarshalJSON handles JSON marshaling for Date fields
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(d.Time.Format(dateFormat))
}

// Scan implements the sql.Scanner interface
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Date", value)
	}
}

// Value implements the driver.Valuer interface
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// NullDate is a custom type for nullable date fields
type NullDate struct {
	sql.NullTime
}

// UnmarshalJSON handles JSON unmarshaling for NullDate fields
func (nd *NullDate) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == nil || *s == "" {
		nd.Valid = false
		nd.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateFormat, *s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	nd.Time = t
	nd.Valid = true
	return nil
}

// MarshalJSON handles JSON marshaling for NullDate fields
func (nd NullDate) MarshalJSON() ([]byte, error) {
	if !nd.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nd.Time.Format(dateFormat))
}
