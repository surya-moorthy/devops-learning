package hello

import "testing"

// testing.T is to get the state of the test case 

func TestHello(T *testing.T  ) {   // name must be in Testxxx
	// subsets 
   T.Run("saying hello with name", func(t *testing.T) {
	    got := Hello("Chris")
		want := "Hello, Chris"
        asseetErrorMessage(T,got,want)
		
   })

   //  subsets
   T.Run("say Hello, World when there is no string", func(t *testing.T) {
	   got := Hello("")
	   want := "Hello, World"
       asseetErrorMessage(T,got,want)
   })
}

func asseetErrorMessage(t testing.TB, got , want string)  {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q",got , want)
	}
}