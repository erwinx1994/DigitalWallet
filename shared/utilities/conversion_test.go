package utilities

import "testing"

func Test_DatabaseToDisplayConversion(t *testing.T) {

	// Test converting database to display format
	{
		// 23 -> 0.23
		database_value := int64(23)
		display_value := Convert_database_to_display_format(database_value)
		expected := "0.23"
		if display_value != expected {
			t.Error("Expected: ", expected, ", Got: ", display_value)
		}
	}
	{
		// 123 -> 1.23
		database_value := int64(123)
		display_value := Convert_database_to_display_format(database_value)
		expected := "1.23"
		if display_value != expected {
			t.Error("Expected: ", expected, ", Got: ", display_value)
		}
	}
	{
		// 3 -> 0.03
		database_value := int64(3)
		display_value := Convert_database_to_display_format(database_value)
		expected := "0.03"
		if display_value != expected {
			t.Error("Expected: ", expected, ", Got: ", display_value)
		}
	}
	{
		// 100 -> 1.00
		database_value := int64(100)
		display_value := Convert_database_to_display_format(database_value)
		expected := "1.00"
		if display_value != expected {
			t.Error("Expected: ", expected, ", Got: ", display_value)
		}
	}
	{
		// 10000 -> 100.00
		database_value := int64(10000)
		display_value := Convert_database_to_display_format(database_value)
		expected := "100.00"
		if display_value != expected {
			t.Error("Expected: ", expected, ", Got: ", display_value)
		}
	}
}

func Test_DisplayToDatabaseConversion(t *testing.T) {

	// Test converting display to database format
	{
		// 1.11 -> 111
		display_value := "1.11"
		database_value, err := Convert_display_to_database_format(display_value)
		if err != nil {
			t.Fatal(err)
		}
		expected_int := int64(111)
		if database_value != expected_int {
			t.Error("Expected: ", expected_int, ", Got: ", database_value)
		}
	}
	{
		// 1.1 -> 110
		display_value := "1.1"
		database_value, err := Convert_display_to_database_format(display_value)
		if err != nil {
			t.Fatal(err)
		}
		expected_int := int64(110)
		if database_value != expected_int {
			t.Error("Expected: ", expected_int, ", Got: ", database_value)
		}
	}
	{
		// 1 -> 100
		display_value := "1"
		database_value, err := Convert_display_to_database_format(display_value)
		if err != nil {
			t.Fatal(err)
		}
		expected_int := int64(100)
		if database_value != expected_int {
			t.Error("Expected: ", expected_int, ", Got: ", database_value)
		}
	}
	{
		// 00101.1 -> 10110
		display_value := "00101.1"
		database_value, err := Convert_display_to_database_format(display_value)
		if err != nil {
			t.Fatal(err)
		}
		expected_int := int64(10110)
		if database_value != expected_int {
			t.Error("Expected: ", expected_int, ", Got: ", database_value)
		}
	}
	{
		// 00101.12 -> 10110
		display_value := "00101.12"
		database_value, err := Convert_display_to_database_format(display_value)
		if err != nil {
			t.Fatal(err)
		}
		expected_int := int64(10112)
		if database_value != expected_int {
			t.Error("Expected: ", expected_int, ", Got: ", database_value)
		}
	}
	{
		// 00101.123 -> error
		display_value := "00101.123"
		_, err := Convert_display_to_database_format(display_value)
		if err == nil {
			t.Fatal("Expected an error. Got: nil.")
		}
	}
}
