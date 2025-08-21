package survey

import "fmt"

// General bird interface
type Bird interface {
	Walk() string
}

// Interface for flying birds
type FlyingBird interface {
	Bird
	Fly() string
}

type Sparrow struct{}

func (s Sparrow) Walk() string { return "Sparrow is walking" }
func (s Sparrow) Fly() string  { return "Sparrow is flying" }

type Penguin struct{}

func (p Penguin) Walk() string { return "Penguin is walking" }
func (p Penguin) Swim() string { return "Penguin is swimming" }

func makeBirdWalk(b Bird) {
	fmt.Println(b.Walk())
}

func makeBirdFly(f FlyingBird) {
	fmt.Println(f.Fly())
}

func Birdies() {
	s := Sparrow{}
	p := Penguin{}

	makeBirdWalk(s) // ✅ works
	makeBirdFly(s)  // ✅ works

	makeBirdWalk(p) // ✅ works
	// makeBirdFly(p) // ❌ won't compile (and that's GOOD)
}

