package strconv

import (
	"errors"
)

// Atoi converts a string to an integer
func Atoi(s string) (int, error) {
	if s == "" {
		return 0, errors.New("invalid syntax")
	}
	
	neg := false
	i := 0
	
	// Handle negative numbers
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
		
		// Check for overflow (simplified)
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

// Itoa converts an integer to a string
func Itoa(i int) string {
	if i == 0 {
		return "0"
	}
	
	neg := i < 0
	if neg {
		i = -i
	}
	
	// Count digits
	temp := i
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}
	
	// Build string
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

// ParseInt parses a string in the given base and returns the corresponding value
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