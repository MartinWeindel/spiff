package yaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("YAML Parser", func() {

	Context("value is a map", func() {
		It("parses maps as strings mapping to Nodes", func() {
			parsed, err := Parse("test", []byte(`foo: "fizz \"buzz\""`))
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed).To(Equal(node(map[string]Node{"foo": node(`fizz "buzz"`)})))
		})

		It("parses maps with block string values", func() {
			parsesAs("foo: |\n  sup\n  :3", map[string]Node{"foo": node("sup\n:3")})
			parsesAs("foo: >\n  sup\n  :3", map[string]Node{"foo": node("sup :3")})
		})

		Context("keys are not strings", func() {
			It("fails", func() {
				_, err := Parse("test", []byte("1: foo"))
				Expect(err).To(BeAssignableToTypeOf(NonStringKeyError{}))
			})
		})
	})

	Context("value is a list", func() {
		It("parses with Node contents", func() {
			parsesAs("- 1\n- two", []Node{node(1), node("two")})
		})
	})

	Context("value is an integer", func() {
		It("parses as ints", func() {
			parsesAs("1", 1)
			parsesAs("-1", -1)
		})
	})

	Context("value is a float", func() {
		It("parses as float64s", func() {
			parsesAs("1.0", 1.0)
			parsesAs("-1.0", -1.0)
		})
	})

	Context("value is a boolean", func() {
		It("parses as bools", func() {
			parsesAs("true", true)
			parsesAs("false", false)
		})
	})

	Context("parsing multi documents", func() {
		It("returns all documents", func() {
			sourceName := "test"
			source := []byte(`
doc1:
---
doc2:
`)

			docs, err := ParseMulti(sourceName, source)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(docs)).To(Equal(2))
		})
	})
})

func parsesAs(source string, expr interface{}) {
	parsed, err := Parse("test", []byte(source))
	Expect(err).NotTo(HaveOccurred())
	Expect(parsed).To(Equal(node(expr)))
}
