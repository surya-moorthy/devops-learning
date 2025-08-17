package hello

import "fmt"

const HelloString = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return HelloString + name
}

func main() {
	fmt.Println(Hello("Chris"))
}