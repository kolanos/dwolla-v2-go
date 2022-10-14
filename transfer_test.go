package dwolla

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferCreatedTime(t *testing.T) {
	transfer := &Transfer{Created: "2015-09-03T23:56:10.023Z"}
	assert.Equal(t, int64(1441324570023), transfer.CreatedTime().UnixMilli())
}
