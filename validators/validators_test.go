package validators

import "testing"

// TODO refactor tests to parameterized form, add more cases
func TestIsValidNumberId_Number(t *testing.T) {
	expected := true
	res := IsValidNumberId("12")
	if expected != res {
		t.Errorf("Wrong result:\nexpected: %v\ngot: %v", expected, res)
	}
}

func TestIsValidNumberId_NotNumber(t *testing.T) {
	expected := false
	res := IsValidNumberId("abc")
	if expected != res {
		t.Errorf("Wrong result:\nexpected: %v\ngot: %v", expected, res)
	}
}

func TestIsValidNumberId_EmpryString(t *testing.T) {
	expected := false
	res := IsValidNumberId("")
	if expected != res {
		t.Errorf("Wrong result:\nexpected: %v\ngot: %v", expected, res)
	}
}
