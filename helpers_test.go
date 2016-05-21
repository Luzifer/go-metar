package metar_test

import (
	. "github.com/Luzifer/go-metar"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {

	It("should convert values into expected results", func() {
		Expect(KtsToMs(1)).To(Equal(0.514444))
		Expect(InHgTohPa(1)).To(Equal(33.8638866667))
		Expect(StatMileToKm(1)).To(Equal(1.60934))
		Expect(MbTohPa(1)).To(Equal(0.1))
		Expect(KtsToBft(5)).To(Equal(2))
	})

})
