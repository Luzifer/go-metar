package metar_test

import (
	"time"

	. "github.com/Luzifer/go-metar"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metar", func() {
	var (
		station = ""
		result  *Result
		err     error
	)

	JustBeforeEach(func() {
		result, err = FetchCurrentStationWeather(station)
	})

	Context("with station EDDH (HAM)", func() {
		BeforeEach(func() {
			station = "EDDH"
		})

		It("should not have errored", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should be at the expected position", func() {
			Expect(result.Latitude).To(Equal(53.63))
			Expect(result.Longitude).To(Equal(10.0))
		})

		It("should be a METAR station reporting", func() {
			Expect(result.MetarType).To(Equal("METAR"))
		})

		It("should be a fairly new result", func() {
			Expect(time.Since(result.ObservationTime) < 1*time.Hour).To(BeTrue())
		})

		It("should have information about SkyCover and FlightCategory", func() {
			Expect(result.SkyCondition.SkyCover).NotTo(Equal(SkyCover("")))
			Expect(result.FlightCategory).NotTo(Equal(FlightCategory("")))
		})
	})

})
