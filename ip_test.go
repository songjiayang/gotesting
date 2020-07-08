package gotesting

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestIsIPV4WithoutTable(t *testing.T) {
	if IsIPV4("") {
		t.Errorf("IsIPV4(%s) should be false", "")
	}

	if IsIPV4("192.168.0") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.0")
	}

	if IsIPV4("192.168.x.1") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.x.1")
	}

	if IsIPV4("192.168.0.1.1") {
		t.Errorf("IsIPV4(%s) should be false", "192.168.0.1.1")
	}

	if !IsIPV4("127.0.0.1") {
		t.Errorf("IsIPV4(%s) should be true", "127.0.0.1")
	}

	if !IsIPV4("192.168.0.1") {
		t.Errorf("IsIPV4(%s) should be true", "192.168.0.1")
	}

	if !IsIPV4("255.255.255.255") {
		t.Errorf("IsIPV4(%s) should be true", "255.255.255.255")
	}

	if !IsIPV4("120.52.148.118") {
		t.Errorf("IsIPV4(%s) should be true", "120.52.148.118")
	}
}

func TestIsIPV4WithTable(t *testing.T) {
	testCases := []struct {
		IP    string
		valid bool
	}{
		{"", false},
		{"192.168.0", false},
		{"192.168.x.1", false},
		{"192.168.0.1.1", false},
		{"127.0.0.1", true},
		{"192.168.0.1", true},
		{"255.255.255.255", true},
		{"120.52.148.118", true},
	}

	for _, tc := range testCases {
		t.Run(tc.IP, func(t *testing.T) {
			if IsIPV4(tc.IP) != tc.valid {
				t.Errorf("IsIPV4(%s) should be %v", tc.IP, tc.valid)
			}
		})
	}
}

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

func TestIsIPV4WithGoconvey(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("gotesting.IsIPV4()", t, func() {
		Convey("should be invalid", func() {
			Convey("empty string", func() {
				So(IsIPV4(""), ShouldEqual, false)
			})

			Convey("with less length", func() {
				So(IsIPV4("192.0.1"), ShouldEqual, false)
			})

			Convey("with more length", func() {
				So(IsIPV4("192.168.1.0.1"), ShouldEqual, false)
			})

			Convey("with invalid character", func() {
				So(IsIPV4("192.168.x.1"), ShouldEqual, false)
			})
		})

		Convey("should be valid", func() {
			Convey("loopback address", func() {
				So(IsIPV4("127.0.0.1"), ShouldEqual, true)
			})

			Convey("extranet address", func() {
				So(IsIPV4("120.52.148.118"), ShouldEqual, true)
			})
		})
	})
}
