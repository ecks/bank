package transactions

import "testing"

func TestProcessPAIN(t *testing.T) {
	data := []string{"", ""}
	_, err := ProcessPAIN(data)
	if err == nil {
		t.Errorf("ProcessPAIN does not pass. Looking for %v, got %v", "Not all data is present. Run pain~help to check for needed PAIN data", nil)
	}

	data = []string{"", "", "not integer"}
	_, err = ProcessPAIN(data)
	if err == nil {
		t.Errorf("ProcessPAIN CheckTransactionType does not pass. Looking for %v, got %v", "Could not get type of PAIN transaction", nil)
	}

	data = []string{"", "", "1"}
	_, err = ProcessPAIN(data)
	if err == nil {
		t.Errorf("ProcessPAIN PainType1 does not pass. Looking for %v, got %v", "Not all data is present. Run pain~help to check for needed PAIN data", nil)
	}

	data = []string{"", "", "1000"}
	_, err = ProcessPAIN(data)
	if err == nil {
		t.Errorf("ProcessPAIN PainType1000 does not pass. Looking for %v, got %v", "Not all data is present. Run pain~help to check for needed PAIN data", nil)
	}
}

func BenchmarkProcessPAIN(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// None of these pass/do inserts into transaction table
		data := []string{"", ""}
		_, _ = ProcessPAIN(data)

		data = []string{"", "", "not integer"}
		_, _ = ProcessPAIN(data)

		data = []string{"", "", "1"}
		_, _ = ProcessPAIN(data)

		data = []string{"", "", "1000"}
		_, _ = ProcessPAIN(data)
	}
}
