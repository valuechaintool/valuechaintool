package utils

import "fmt"

func MapOnlyContains(m map[string]interface{}, eks []string) error {
	for k := range m {
		var ok bool
		for _, ek := range eks {
			if ek == k {
				ok = true
			}
		}
		if !ok {
			return fmt.Errorf("found field %v that was not expected", k)
		}
	}
	return nil
}
