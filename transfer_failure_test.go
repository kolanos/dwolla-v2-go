package dwolla

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferFailureCreatedTime(t *testing.T) {
	transferFailure := &TransferFailure{Created: "2015-09-03T23:56:10.023Z"}
	assert.Equal(t, int64(1441324570023), transferFailure.CreatedTime().UnixMilli())
}
