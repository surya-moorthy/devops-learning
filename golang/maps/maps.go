package maps

import "errors"

type Dictionary map[string]string

func (D Dictionary) Search(word string) (string,error) {

	dictionary, ok := D[word]

	if !ok {
		return "",errors.New("could not find the word here")
	}

	return dictionary, nil
}

