package paths

import (
	"reflect"
	"testing"
)

func Test_MatchAndExtract(t *testing.T) {

	// Test with valid paths
	{
		path := "/wallets/{wallet_id}/transaction_history"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected_result := &MatchResult{
			MatchFound:       true,
			WildcardSegments: make(map[string]string, maximum_number_of_wildcard_segments),
			KeyValuePairs:    make(map[string]string, maximum_number_of_keyvalue_pairs),
		}
		expected_result.WildcardSegments["wallet_id"] = "wallet_id"
		if !reflect.DeepEqual(*result, *expected_result) {
			t.Error("Expected: ", *expected_result, ", Got: ", result)
		}
	}
	{
		path := "/wallets/{wallet_id}/transaction_history?key1=value1&key2=value2"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected_result := &MatchResult{
			MatchFound:       true,
			WildcardSegments: make(map[string]string, maximum_number_of_wildcard_segments),
			KeyValuePairs:    make(map[string]string, maximum_number_of_keyvalue_pairs),
		}
		expected_result.WildcardSegments["wallet_id"] = "wallet_id"
		expected_result.KeyValuePairs["key1"] = "value1"
		expected_result.KeyValuePairs["key2"] = "value2"
		if !reflect.DeepEqual(*result, *expected_result) {
			t.Error("Expected: ", *expected_result, ", Got: ", result)
		}
	}
	{
		path := "/wallets/{wallet_id}/transaction_history?key1=&key2="
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected_result := &MatchResult{
			MatchFound:       true,
			WildcardSegments: make(map[string]string, maximum_number_of_wildcard_segments),
			KeyValuePairs:    make(map[string]string, maximum_number_of_keyvalue_pairs),
		}
		expected_result.WildcardSegments["wallet_id"] = "wallet_id"
		expected_result.KeyValuePairs["key1"] = ""
		expected_result.KeyValuePairs["key2"] = ""
		if !reflect.DeepEqual(*result, *expected_result) {
			t.Error("Expected: ", *expected_result, ", Got: ", result)
		}
	}

	// Test with invalid paths
	{
		path := "/wall2ts/{wallet_id}/transaction_history?key1=value1&key2=value2"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
	{
		path := "/wall2ts/{wallet_id}/"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
	{
		path := "/wall2ts/{wallet_id}/transaction_history?key1=value1&key2=value2"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
	{
		path := "/wallets/{wallet_id}/transaction_history?keyyyy"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
	{
		path := "/wallets/{wallet_id}/transaction_history?key1=value1&keyy2"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
	{
		path := "/wallets/{wallet_id}"
		result := MatchAndExtract(path, Wallets_transaction_history)
		expected := false
		if result.MatchFound != expected {
			t.Error("Expected: ", expected, ", Got: ", result.MatchFound)
		}
	}
}
