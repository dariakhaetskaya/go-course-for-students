package fizzbuzz

import "strconv"

func FizzBuzz(i int) string {
	fizz := i%3 == 0
	buzz := i%5 == 0
	if fizz && buzz {
		return "FizzBuzz"
	}
	if fizz {
		return "Fizz"
	}
	if buzz {
		return "Buzz"
	}

	return strconv.Itoa(i)
}
