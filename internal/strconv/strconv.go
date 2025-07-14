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

// Itoa convertit un entier en chaîne
func Itoa(i int) string {
	if i == 0 {
		return "0"
	}
	
	neg := i < 0
	if neg {
		i = -i
	}
	
	// Compte les chiffres
	temp := i
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}
	
	// On construit la chaîne
	var result []byte
	if neg {
		result = make([]byte, digits+1)
		result[0] = '-'
		for j := digits; j > 0; j-- {
			result[j] = byte('0' + i%10)
			i /= 10
		}
	} else {
		result = make([]byte, digits)
		for j := digits - 1; j >= 0; j-- {
			result[j] = byte('0' + i%10)
			i /= 10
		}
	}
	
	return string(result)
}

// ParseInt analyse une chaîne dans la base donnée et retourne la valeur correspondante
func ParseInt(s string, base, bitSize int) (int64, error) {
	if base != 10 {
		return 0, errors.New("unsupported base")
	}
	
	if bitSize != 0 && bitSize != 64 {
		return 0, errors.New("unsupported bit size")
	}
	
	i, err := Atoi(s)
	return int64(i), err
}