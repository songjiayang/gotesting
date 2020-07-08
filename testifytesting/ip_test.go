package testifytesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "gotesting/tabletesting"
)

func TestIsIPV4WithTestify(t *testing.T) {
	assertion := assert.New(t)

	assertion.False(IsIPV4(""))
	assertion.False(IsIPV4("192.168.0"))
	assertion.False(IsIPV4("192.168.x.1"))
	assertion.False(IsIPV4("192.168.0.1.1"))
	assertion.True(IsIPV4("127.0.0.1"))
	assertion.True(IsIPV4("192.168.0.1"))
	assertion.True(IsIPV4("255.255.255.255"))
	assertion.True(IsIPV4("120.52.148.118"))
}
