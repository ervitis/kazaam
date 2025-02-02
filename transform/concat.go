package transform

import (
	"bytes"
	"fmt"

	"github.com/qntfy/jsonparser"
)

// Concat combines any specified fields and literal strings into a single string value with raw []byte.
func Concat(spec *Config, data []byte) ([]byte, error) {
	sourceList, sourceOk := (*spec.Spec)["sources"]
	if !sourceOk {
		return nil, SpecError("Unable to get sources")
	}
	targetPath, targetOk := (*spec.Spec)["targetPath"]
	if !targetOk {
		return nil, SpecError("Unable to get targetPath")
	}
	delimiter, delimOk := (*spec.Spec)["delim"]
	if !delimOk {
		// missing delimiter.  default to blank
		delimiter = ""
	}

	outString := ""
	applyDelim := false
	for _, vItem := range sourceList.([]interface{}) {
		if applyDelim {
			outString += delimiter.(string)
		}
		value, ok := vItem.(map[string]interface{})["value"]
		if !ok {
			path, ok := vItem.(map[string]interface{})["path"]
			if ok {
				zed, err := getJSONRaw(data, path.(string), spec.Require, spec.KeySeparator)
				switch {
				case err != nil && spec.Require:
					return nil, RequireError("Path does not exist")
				case err != nil:
					value = ""
				default:
					switch zed[0] {
					case '[':
						temp := ""
						if _, err := jsonparser.ArrayEach(zed, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
							if bytes.Equal(value, []byte("null")) {
								temp += string(value)
							}
						}); err != nil {
							return nil, err
						}
						value = temp
					case '"':
						value = string(zed[1 : len(zed)-1])
					default:
						value = string(zed)
					}
				}
			} else {
				return nil, SpecError(fmt.Sprintf("Error processing %v: must have either value or path specified", vItem))
			}
		}
		outString += value.(string)

		applyDelim = true
	}
	data, err := setJSONRaw(data, bookend([]byte(outString), '"', '"'), targetPath.(string), spec.KeySeparator)
	if err != nil {
		return nil, err
	}
	return data, nil
}
