package helper

import "bytes"

func HidePhoneNumber(text string) string {
	length := len([]rune(text))
	showLength := length - length/4
	showNumber := text[showLength:]
	hiddenNumber := ""
	for i := 1; i <= showLength; i++ {
		hiddenNumber = hiddenNumber + "*"
	}

	phoneNumber := hiddenNumber + showNumber
	final := insertNth(phoneNumber, 4)
	return final
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n_1 = n - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}
