package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
GetJSON sends an HTTP GET request to the specified URL
and decodes the response JSON data into the provided target value.

The target argument must be a pointer to a value of the same type as the JSON data
to be decoded. If the request is successful and the JSON data is valid,
the target value will be populated with the decoded data.

If the HTTP request fails, or the response status code is not 200 OK, an error will
be returned. The error message will provide information about the reason for the failure.

Example usage:

	var data MyData
	err := GetJSON("https://example.com/data.json", &data)
	if err != nil {
	  log.Fatal(err)
	}

Parameters:
- url: The URL to fetch JSON data from.
- target: A pointer to the target value to decode JSON data into.

Returns:
  - error: An error if the HTTP request fails or the JSON decoding fails, or nil if
    the request is successful and the JSON data is decoded.
*/
func GetJSON[T any](url string, target *T) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get data from %s: %w", url, err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get data from %s: %s", url, response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&target); err != nil {
		return fmt.Errorf("failed to decode JSON data from %s: %w", url, err)
	}

	return nil
}
