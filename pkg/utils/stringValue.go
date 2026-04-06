package utils

// StringValue возвращает строку по указателю или пустую строку, если указатель nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
