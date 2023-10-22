package protocol

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalID(t *testing.T) {
	var id ID
	copy(id[:], "012345678901234567890")
	assert.Equal(t, "20:01234567890123456789", string(bencode.MustMarshal(id)))
}

func TestID_Distance(t *testing.T) {
	tests := []struct {
		a, b, distance string
	}{
		{
			"ffbfba8945192d408d3dcc52ba24903a00000000",
			"ffbfba8945192d408d3dcc52ba24903a00000001",
			"0000000000000000000000000000000000000001",
		},
		{
			"ffbfba8945192d408d3dcc52ba24903a00000000",
			"ffbfba8945192d408d3dcc52ba24903a00000002",
			"0000000000000000000000000000000000000002",
		},
		{
			"ffcfba8945192d408d3dcc52ba24903a00000000",
			"ffbfba8945192d408d3dcc52ba24903a00000001",
			"0070000000000000000000000000000000000001",
		},
		{
			"0fffffffffffffffffffffffffffff0000000000",
			"0fffffffffffffffffffffffffffff0000000001",
			"0000000000000000000000000000000000000001",
		},
		{
			"0fffffffffffffffffffffffffffff0000000000",
			"0fffffffffffffffffffffffffffff0000000002",
			"0000000000000000000000000000000000000002",
		},
		{
			"0fffffffffffffffffffffffffffff0000000000",
			"0fffffffffffffffffffffffffffff0010000000",
			"0000000000000000000000000000000010000000",
		},
	}
	for _, test := range tests {
		a := MustParseID(test.a)
		b := MustParseID(test.b)
		distance := MustParseID(test.distance)
		if a.Distance(b) != distance {
			t.Errorf("Distance(%v, %v) != %v (%v)", a, b, distance, a.Distance(b))
		}
		t.Log(distance)
		t.Log(distance.LeadingZeros())
	}
}
