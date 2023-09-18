package cmd

const groupAnnotation = "group"

func group(value string) map[string]string {
	return annotate(groupAnnotation, value)
}

func annotate(key, value string) map[string]string {
	return map[string]string{key: value}
}
