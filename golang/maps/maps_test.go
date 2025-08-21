package maps

import "testing"

func TestSearch(t *testing.T) {

	dictionary := Dictionary{"test" : "this is test 1"}

	t.Run("known word",func(t *testing.T) {

		got,_ := dictionary.Search("test")
		want := "this is test 1"

		assertStrings(t,got,want)
	})

	t.Run("unknown word",func(t *testing.T){
		_,err := dictionary.Search("tasty")
		
		want := "could not find the word here"


		if err == nil {
          t.Fatal("expected an error here.")
		}

		assertStrings(t,err.Error(),want)
	})
}

func assertStrings(t testing.TB,got string,want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q given %q",got , want , "test")
	}
}

