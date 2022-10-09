package chirano

import "testing"

func TestHelpers_RandomString(t *testing.T) {

	var testHelps Helpers

	s := testHelps.RandomString(10)

	if len(s) != 10 {
		t.Error("Incorrect lenght received")
	}
}
