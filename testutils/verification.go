package testutils

import (
	"github.com/Woody1193/goutils/utils"
	. "github.com/onsi/gomega"
)

// ItemIsNil is a function that verifies that the data sent to it is nil
func ItemIsNil[T any](item T) {
	Expect(item).Should(BeNil())
}

// NoInnerError verifies that an error did not occur
func NoInnerError() func(error) {
	return func(err error) {
		Expect(err).ShouldNot(HaveOccurred())
	}
}

// InnerErrorVerifier verifies an error message
func InnerErrorVerifier(message string) func(error) {
	return func(err error) {
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal(message))
	}
}

// ErrorVerifier verifies the fields on a backend Error
func ErrorVerifier(env string, pkg string, file string, class string,
	function string, line int, innerVerifier func(error), message string,
	fullMsg string) func(*utils.Error) {
	return func(err *utils.Error) {
		Expect(err.Class).Should(Equal(class))
		Expect(err.Environment).Should(Equal(env))
		Expect(err.File).Should(Equal(file))
		Expect(err.Function).Should(Equal(function))
		Expect(err.GeneratedAt).ShouldNot(BeNil())
		Expect(err.LineNumber).Should(Equal(line))
		Expect(err.Message).Should(Equal(message))
		Expect(err.Package).Should(Equal(pkg))
		Expect(err.Error()).Should(HaveSuffix(fullMsg))
		innerVerifier(err.Inner)
	}
}
