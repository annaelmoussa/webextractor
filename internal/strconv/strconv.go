package strconv

import (
	"errors"
)

// Atoi convertit une chaîne en entier
func Atoi(s string) (int, error) {
	if s == "" {
		return 0, errors.New("invalid syntax")
	}

	neg := false
	i := 0

	// On gère les nombres négatifs
	if s[0] == '-' {
		if len(s) == 1 {
			return 0, errors.New("invalid syntax")
		}
		neg = true
		i = 1
	} else if s[0] == '+' {
		if len(s) == 1 {
			return 0, errors.New("invalid syntax")
		}
		i = 1
	}

	result := 0
	for ; i < len(s); i++ {
		char := s[i]
		if char < '0' || char > '9' {
			return 0, errors.New("invalid syntax")
		}

		digit := int(char - '0')

		// Vérification de débordement (simplifiée)
		if result > (1<<31-1-digit)/10 {
			if neg {
				return -1 << 31, errors.New("value out of range")
			}
			return 1<<31 - 1, errors.New("value out of range")
		}

		result = result*10 + digit
	}

	if neg {
		return -result, nil
	}
	return result, nil
}
