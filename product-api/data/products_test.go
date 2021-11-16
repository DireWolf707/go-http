package data

import "testing"

func TestValidator(t *testing.T) {
	p := &Product{Name: "test coffee", Price: 20}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
