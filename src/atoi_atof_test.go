package main

import "testing"

func TestAtoiAndAtof_ParseCommonFormsAndOverflow(t *testing.T) {
	if got := Atoi(" 42 "); got != 42 {
		t.Fatalf("Atoi() = %v, want 42", got)
	}
	if got := Atoi("-17"); got != -17 {
		t.Fatalf("Atoi() = %v, want -17", got)
	}
	if got := Atoi("2147483648"); got != IMax {
		t.Fatalf("Atoi(overflow+) = %v, want %v", got, IMax)
	}
	if got := Atoi("-2147483649"); got != IErr {
		t.Fatalf("Atoi(overflow-) = %v, want %v", got, IErr)
	}

	if got := Atof(" 3.5 "); got != 3.5 {
		t.Fatalf("Atof() = %v, want 3.5", got)
	}
	if got := Atof("-2.5e2"); got != -250 {
		t.Fatalf("Atof(exp) = %v, want -250", got)
	}
	if got := Atof("1.2.3"); got != 1.2 {
		t.Fatalf("Atof(extra dot) = %v, want 1.2", got)
	}
}
