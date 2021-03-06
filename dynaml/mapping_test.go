package dynaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mapping expressions", func() {
	It("prints mapping expression", func() {
		desc := MappingExpr{
			ReferenceExpr{[]string{"list"}},
			ConcatenationExpr{
				ReferenceExpr{[]string{"x"}},
				StringExpr{".*"},
			},
			MapToListContext,
		}.String()
		Expect(desc).To(Equal("map[list|x \".*\"]"))
	})

	It("simplifies lambda mapping expression", func() {
		desc := MappingExpr{
			ReferenceExpr{[]string{"list"}},
			LambdaExpr{
				[]string{"x"},
				ConcatenationExpr{
					ReferenceExpr{[]string{"x"}},
					StringExpr{".*"},
				},
			},
			MapToListContext,
		}.String()
		Expect(desc).To(Equal("map[list|x|->x \".*\"]"))
	})
})
