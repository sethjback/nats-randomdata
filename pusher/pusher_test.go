package pusher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomData(t *testing.T) {
	data := newRandomMessage()
	assert.NotEmpty(t, data)
}
