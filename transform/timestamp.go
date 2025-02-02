package transform

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qntfy/jsonparser"
)

// This is necessary for testing purposes
var now = time.Now

const unixFormat = "$unix"

// Timestamp parses and formats timestamp strings using the golang syntax
func Timestamp(spec *Config, data []byte) ([]byte, error) {
	for k, v := range *spec.Spec {
		assertedV, vErr := v.(map[string]interface{})
		if !vErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get value for key: %s", k))
		}
		inputFormat, inputErr := assertedV["inputFormat"].(string)
		if !inputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"inputFormat\" for key: %s", k))
		}
		outputFormat, outputErr := assertedV["outputFormat"].(string)
		if !outputErr {
			return nil, SpecError(fmt.Sprintf("Warn: Invalid spec. Unable to get \"outputFormat\" for key: %s", k))
		}
		// check if an array wildcard is present and if it is, treat it the
		// same as a key with an array
		//if k[len(k)-2] == '*' {
		//	k = k[:len(k)-3]
		//}
		var dataForV []byte
		var err error

		if inputFormat == "$now" {
			t, errParse := now().MarshalText()
			if errParse != nil {
				return nil, errParse
			}
			dataForV = bookend(t, '"', '"')
			// this is the standard format that `time.Now().String()` uses
			inputFormat = time.RFC3339
		} else {
			// grab the data
			dataForV, err = getJSONRaw(data, k, spec.Require, spec.KeySeparator)
			if err != nil {
				return nil, err
			}
		}
		// if the key is missing bail and keep iterating
		if bytes.Equal(dataForV, []byte("null")) {
			continue
		}
		// can only parse and format strings and arrays of strings, check the
		// value type and handle accordingly
		switch dataForV[0] {
		case '"':
			var errFormat error
			formattedItem, errFormat := parseAndFormatValue(inputFormat, outputFormat, string(dataForV[1:len(dataForV)-1]))
			if errFormat != nil {
				return nil, errFormat
			}
			data, errFormat = setJSONRaw(data, []byte(formattedItem), k, spec.KeySeparator)
			if errFormat != nil {
				return nil, errFormat
			}
		case '[':
			var unformattedItems []string
			_, err = jsonparser.ArrayEach(dataForV, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				unformattedItems = append(unformattedItems, string(value))
			})
			if err != nil {
				return nil, err
			}
			for idx, unformattedItem := range unformattedItems {
				formattedItem, err := parseAndFormatValue(inputFormat, outputFormat, unformattedItem)
				if err != nil {
					return nil, err
				}
				// replacing the wildcard here feels hacky, but seems to be the
				// quickest way to achieve the outcome we want
				data, err = setJSONRaw(data, []byte(formattedItem), strings.Replace(k, "*", strconv.Itoa(idx), -1), spec.KeySeparator)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, ParseError(fmt.Sprintf("Warn: Unknown type in message for key: %s", v))
		}
	}
	return data, nil
}

// parseAndFormatValue generates a properly formatted timestamp
func parseAndFormatValue(inputFormat, outputFormat, unformattedItem string) (string, error) {
	var (
		parsedItem    time.Time
		formattedItem string
		err           error
		i             int64
	)

	if inputFormat == unixFormat {
		i, err = strconv.ParseInt(unformattedItem, 10, 64)
		if err != nil {
			return "", err
		}
		parsedItem = time.Unix(i, 0)
	} else {
		parsedItem, err = time.Parse(inputFormat, unformattedItem)
		if err != nil {
			return "", err
		}
	}

	if outputFormat == unixFormat {
		formattedItem = strconv.FormatInt(parsedItem.Unix(), 10)
	} else {
		formattedItem = parsedItem.Format(outputFormat)
	}
	return strings.Join([]string{"\"", formattedItem, "\""}, ""), nil
}
