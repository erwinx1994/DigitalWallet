package paths

const (
	Wallets_deposits            string = "/wallets/{wallet_id}/deposits"
	Wallets_withdrawals         string = "/wallets/{wallet_id}/withdrawals"
	Wallets_balance             string = "/wallets/{wallet_id}/balance"
	Wallets_transaction_history string = "/wallets/{wallet_id}/transaction_history"
	Transfer                    string = "/transfer"
	Test                        string = "/test"
)

// These can be determined by the patterns specified above
const (
	maximum_number_of_wildcard_segments int = 1
	maximum_number_of_keyvalue_pairs    int = 2
)

// A wildcard segment is denoted by {} in the pattern
// Key value pairs are specified in the query string section of a URL
type MatchResult struct {
	MatchFound       bool
	WildcardSegments map[string]string
	KeyValuePairs    map[string]string
}

/*
The Parser.Match function tries to match the path with the specified pattern.
It returns true if a match was found. Otherwise, it returns false.
It also extracts the wildcard segments and the key values pairs in the path.

It assumes that both the path and pattern contains only these characters:
1-9, a-z, A-Z

Using unicode characters in the path and pattern will cause erroneous results.

This function completes in O((n + m)x) time complexity where n is the length of the path,
m is the length of the pattern and x is the overhead of string concatenation
for extracting wildcard segments, keys and values. Go's standard string
concatenation operator += was used for its simplicity.

It is possible to further reduce the time complexity to O(n + m) by using a custom
string builder (block of bytes with fixed maximum size and a pointer) instead of
Go's simple string concatenation operator +=. That is left as future work.
*/
func MatchAndExtract(path, pattern string) *MatchResult {

	result := MatchResult{
		MatchFound:       true,
		WildcardSegments: make(map[string]string, maximum_number_of_wildcard_segments),
		KeyValuePairs:    make(map[string]string, maximum_number_of_keyvalue_pairs),
	}

	path_index := 0
	pattern_index := 0
	wildcard_segment := ""
	wildcard_segment_name := ""

	// Try matching each character of path with pattern
	for path_index < len(path) && pattern_index < len(pattern) {

		// Mismatching characters found
		if path[path_index] != pattern[pattern_index] {
			result.MatchFound = false
			result.WildcardSegments = nil
			result.KeyValuePairs = nil
			return &result
		}

		// Found start of wildcard pattern as indicated by a left curly brace {
		// Move the path_index and pattern_index to the closing right curly brace }
		if path[path_index] == '{' && pattern[pattern_index] == '{' {

			path_index++
			pattern_index++

			for path_index < len(path) && path[path_index] != '}' {
				wildcard_segment += string(path[path_index])
				path_index++
			}

			for pattern_index < len(pattern) && pattern[pattern_index] != '}' {
				wildcard_segment_name += string(pattern[pattern_index])
				pattern_index++
			}

			if path_index < len(path) && path[path_index] == '}' {
				result.WildcardSegments[wildcard_segment_name] = wildcard_segment
				wildcard_segment = ""
				wildcard_segment_name = ""
			} else {
				// It is invalid for path_index to reach the end of path without finding
				// a closing curling brace }
				result.MatchFound = false
				result.WildcardSegments = nil
				result.KeyValuePairs = nil
				return &result
			}

			// At this point, it is guaranteed that both path_index and pattern_index
			// point to the right curly brace }
		}

		// Both path[path_index] and pattern[pattern_index] are currently.
		// Increment both indices.
		path_index++
		pattern_index++
	}

	// Pattern has not been fully matched. Invalid!
	// Usually because path is shorter than pattern.
	if pattern_index != len(pattern) {
		result.MatchFound = false
		result.WildcardSegments = nil
		result.KeyValuePairs = nil
		return &result
	}

	// Pattern is fully matched.
	// Other than the remaining query string, it is invalid for the path
	// to be longer than the pattern
	if path_index < len(path) && path[path_index] != '?' {
		result.MatchFound = false
		result.WildcardSegments = nil
		result.KeyValuePairs = nil
		return &result
	}

	// Extract key value pairs from query string
	path_index++
	for path_index < len(path) {

		// Extract key
		key := ""
		for path_index < len(path) && path[path_index] != '=' {
			key += string(path[path_index])
			path_index++
		}

		if !(path_index < len(path) && path[path_index] == '=') {
			// Invalid query string found
			result.MatchFound = false
			result.WildcardSegments = nil
			result.KeyValuePairs = nil
			return &result
		}

		path_index++

		// Extract value
		value := ""
		for path_index < len(path) && path[path_index] != '&' {
			value += string(path[path_index])
			path_index++
		}

		if path_index < len(path) && path[path_index] == '&' {
			path_index++
		}

		// Store result
		result.KeyValuePairs[key] = value
	}
	return &result
}
