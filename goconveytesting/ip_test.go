package goconveytesting

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "gotesting/tabletesting"
)

func TestIsIPV4WithGoconvey(t *testing.T) {
	Convey("ip.IsIPV4()", t, func() {
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
