package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Duration time.Duration

func (d Duration) Value() (driver.Value, error) {
	return time.Duration(d).String(), nil
}

func (d *Duration) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		// Convert format of hh:mm:ss into format parseable by time.ParseDuration()
		v = strings.Replace(v, ":", "h", 1)
		v = strings.Replace(v, ":", "m", 1)
		v += "s"
		dur, err := time.ParseDuration(v)
		if err != nil {
			return err
		}
		*d = Duration(dur)
		return nil
	default:
		return fmt.Errorf("cannot sql.Scan() PgDuration from: %#v", v)
	}
}
