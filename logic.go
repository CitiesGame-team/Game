package main

var Dictionary map[string]([2]rune)

//TODO
func fillDictionary() {

}

func isTown(word string) bool {
	_, ok := Dictionary[word]
	if ok {
		return true
	}
	return checkInBase(word)
}

//TODO
func checkInBase(word string) bool {
	return false
}

func correctAncessor(word string, r rune) bool {
	return isTown(word) && Dictionary[word][1] == r
}
