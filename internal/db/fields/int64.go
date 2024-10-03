package fields

import (
	"encoding/json"
	"strconv"

	"github.com/gocql/gocql"
)

type (
	// Int64 is a type alias for int64. Although gocql supports int64, during our application we needed conversions to
	// and from string and int64. This type alias allows us to define custom methods on the int64 type. Marshaling and
	// unmarshalling to and from JSON and CQL is also supported to make it easy to work with the type.
	Int64 int64
)

// String returns the string representation of the Int64 value.
func (i Int64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Int64 returns the int64 value of the Int64.
func (i Int64) Int64() int64 {
	return int64(i)
}

// MarshalJSON implements the json.Marshaler interface.
func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *Int64) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

// MarshalCQL implements the gocql.Marshaler interface.
func (i Int64) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, i.Int64())
}

// UnmarshalCQL implements the gocql.Unmarshaler interface.
func (i *Int64) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	var v int64
	if err := gocql.Unmarshal(info, data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}