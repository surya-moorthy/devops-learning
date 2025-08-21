package survey

type Quesiton interface {
    SetTitle(title string)
	AddOption(option string)
}

type TextInputQuestion struct {
	title string
}

func (q *TextInputQuestion) SetTitle(title string) {

}

