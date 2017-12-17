// Package word ...
package word

import "unicode"

// // IsPalindromeP ...
// func IsPalindromeP(s string) bool {
// 	for i := range s {
// 		if s[i] != s[len(s)-1-i] {
// 			return false
// 		}
// 	}
// 	return true
// }

// IsPalindrome ...
func IsPalindrome(s string) bool {
	// var letters []rune // 优化前
	letters := make([]rune, 0, len(s)) // 优化后
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	// 优化前
	// for i := range letters {
	// 	if letters[i] != letters[len(letters)-1-i] {
	// 		return false
	// 	}
	// }
	// return true

	// 优化后
	n := len(letters) / 2
	for i := 0; i < n; i++ {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
