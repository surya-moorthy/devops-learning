package structs

import "testing"

// func TestPerimeter(t *testing.T) {
// 	t.Run("get valid perimeter of the given dimension",func(t *testing.T) {
// 		dimension := Rectangle{1.0,2.0}
// 		got := Perimeter(dimension)
// 		want := 6.0

// 		if got != want {
// 			t.Errorf("got %.2f want %.2f",got,want)
// 		}
// 	})	
// }

func  TestArea(t *testing.T) {
   t.Run("Test Area of the Rectangle",func(t *testing.T) {
	 dimension := Rectangle{1.0,2.0}
	 got := dimension.Area()
	 want := 2.0

	 if got != want {
		t.Errorf("got %g want %g",got,want)
	 }
   })

   t.Run("Test Area of the Triangle",func(t *testing.T) {
	   dimension := Circle{10}
	   got := dimension.Area()

	   want := 314.1592653589793

	   if got != want {
		 t.Errorf("got %g want %g",got,want)
	   }
   })
}