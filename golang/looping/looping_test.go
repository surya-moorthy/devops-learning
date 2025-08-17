package looping

import "testing"

func TestLooping(t *testing.T) {
	got := Repeat("a")
	expected := "aaaaaa"
	asseetErrorMessage(t,got,expected)
}

func asseetErrorMessage(t testing.TB, got , want string)  {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' want '%s'",got , want)
	}
}

func BenchmarkLooping(b *testing.B) {
   for b.Loop() {
	  Repeat("a")
   }	
}