package utilities

import "errors"

/*
Money is represented using a 64 bit integer in the database to avoid precision loss.

Customers can deposit, withdraw or transfer money up to a hundredth of a currency.

	SGD 0.01

The amount of money recorded with all deposits, withdrawals and transfers is multiplied
by a 100 before updating the database.

The amount of money shown with balance or transaction history requests needs to be
divided by 100 before display to the user.
*/

func Convert_database_to_display_format(internal_amount int64) string {

	// Database format: 150
	// Display format: 1.50

	if internal_amount == 0 {
		return "0.00"
	}

	// Conversion is not done using floating point arithmetic to avoid precision loss
	// Conversion is done by manually adding a decimal point after 2 digits from the right

	// Time complexity O(dx) where d is the number of digits, x is the time complexity of Go's
	// string concatenation. This could be improved to O(d) but is left as future work.

	// Get number of significant digits and place value of most significant digit
	number_of_digits := 1
	largest_place_value := int64(1)
	internal_amount_copy := internal_amount
	internal_amount_copy /= 10
	for internal_amount_copy > 0 {
		number_of_digits++
		largest_place_value *= 10
		internal_amount_copy /= 10
	}

	// Convert integer representation to string
	var result string = ""
	internal_amount_copy = internal_amount
	for number_of_digits > 0 {
		if number_of_digits == 2 {
			if len(result) == 0 {
				result += "0."
			} else {
				result += "."
			}
		}
		if number_of_digits == 1 {
			if len(result) == 0 {
				result += "0.0"
			}
		}
		largest_digit := internal_amount_copy / largest_place_value
		result += string(byte(largest_digit) + '0')
		internal_amount_copy -= (largest_digit * largest_place_value)
		largest_place_value /= 10
		number_of_digits--
	}
	return result
}

func Convert_display_to_database_format(display_amount string) (int64, error) {

	// Display format: 01, 1, 1.0, 1.00,
	// Invalid display format: 1.000
	// The 3rd decimal place will be truncated. In practice, client devices forbid
	// entering more than 2 decimal places.
	// Database format: 100

	// Conversion is not done using floating point arithmetic to avoid precision loss
	// Conversion is done by removing the decimal place and adding more zeros to
	// the right if necessary
	// XX.XX -> XXXX
	// XX.X -> XXX0
	// XX -> XX00
	// Leading zeros are ignored
	// 0XX.XX -> XXXX

	// Time complexity: O(d), where d = Number of digits

	after_decimal_place := false
	number_of_digits_after_decimal_place := 0
	var sum int64 = 0
	for i := 0; i < len(display_amount); i++ {
		character := display_amount[i]
		if character == '.' {
			after_decimal_place = true
			continue
		} else {
			digit := int64(character - '0')
			sum = sum*10 + digit
			if after_decimal_place {
				number_of_digits_after_decimal_place++
			}
		}
	}

	if number_of_digits_after_decimal_place > 2 {
		return 0, errors.New("Using more than 2 decimal places is forbidden!")
	}

	for number_of_digits_after_decimal_place < 2 {
		sum *= 10
		number_of_digits_after_decimal_place++
	}

	return sum, nil
}
