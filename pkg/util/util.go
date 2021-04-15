/* Util package will have the utility functions required for the project to run.

 */

package util

import (
	"os"
)

// GetEnv read the value from the environment using a key and if no key present uses the fallback default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
