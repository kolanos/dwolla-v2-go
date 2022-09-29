package dwolla

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMicroDeposit_CreatedTime(t *testing.T) {
	microDeposit := &MicroDeposit{Created: "2015-09-03T23:56:10.023Z"}
	assert.Equal(t, int64(1441324570023), microDeposit.CreatedTime().UnixMilli())
}
