package luhn

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

// Calculate returns luhn check digit and the provided string number with its luhn check digit appended.
func Calculate(number string) (string, string, error) {
	p := (len(number) + 1) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return "", "", nil
	}

	luhn := sum % 10
	if luhn != 0 {
		luhn = 10 - luhn
	}

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	return strconv.FormatInt(luhn, 10), fmt.Sprintf("%s%d", number, luhn), nil
}

// Generate will generate a valid luhn number of the provided length
func Generate(length int, count int) []string {

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		log.Fatal(err.Error())
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	resArr := make([]string, 0, count)
	var s strings.Builder
	for j := 0; j < count; j++ {
		s.Reset()
		for i := 0; i < length-1; i++ {
			s.WriteString(strconv.Itoa(rand.Intn(9)))
		}
		_, res, _ := Calculate(s.String()) //ignore error because this will always be valid
		resArr = append(resArr, res)
	}

	return resArr
}
func calculateLuhnSum(number string, parity int) (int64, error) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, errors.New("invalid digit")
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum, nil
}
