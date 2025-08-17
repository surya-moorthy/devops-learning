package integer

import "testing"
	
func TestAdder(t *testing.T) {
	sum := Add(2,3)
	expected := 5
	
	if sum != int32(expected) {
	     t.Errorf("expected '%d' but got '%d'",expected,sum)	
	}
}