package ginkgotesting_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gotesting/tabletesting"
)

var _ = Describe("Ip", func() {
	Describe("IsIPV4()", func() {
		// fore content level prepare
		BeforeEach(func() {
			// prepare data before every case
		})

		AfterEach(func() {
			// clear data after every case
		})

		Context("should be invalid", func() {
			It("empty string", func() {
				Expect(IsIPV4("")).To(Equal(false))
			})

			It("with less length", func() {
				Expect(IsIPV4("192.0.1")).To(Equal(false))
			})

			It("with more length", func() {
				Expect(IsIPV4("192.168.1.0.1")).To(Equal(false))
			})

			It("with invalid character", func() {
				Expect(IsIPV4("192.168.x.1")).To(Equal(false))
			})
		})

		Context("should be valid", func() {
			It("loopback address", func() {
				Expect(IsIPV4("127.0.0.1")).To(Equal(true))
			})

			It("extranet address", func() {
				Expect(IsIPV4("120.52.148.118")).To(Equal(true))
			})
		})
	})
})
