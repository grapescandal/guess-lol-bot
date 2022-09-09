package helper

func FilterInput(message string, stringToTrim string) string {
	if len(message) == len(stringToTrim) {
		return ""
	}
	inputMessage := message[len(stringToTrim)+1:]
	return inputMessage
}
