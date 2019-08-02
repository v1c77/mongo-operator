package failover

import (
	"fmt"
	"reflect"
)

// getMapStringKeys return map keys or raise an error.
func getMapStringKeys(m interface{}) ([]string, error) {
	// slow but works.
	mapKeys := reflect.ValueOf(m).MapKeys()

	keys := make([]string, 0, len(mapKeys))
	for _, k := range mapKeys {
		keys = append(keys, fmt.Sprint(k))
	}
	if len(keys) != len(mapKeys) {
		panic("key count error")
	}
	return keys, nil
}

func GetMapStringKeys(m interface{}) []string {

	r, err := getMapStringKeys(m)
	if err != nil {
		return []string{}
	} else {
		return r
	}
}
