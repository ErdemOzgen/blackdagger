package executor

import (
	"os"
	"reflect"
)

// expandEnvHook is a mapstructure decode hook that expands environment variables in string fields.
func expandEnvHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String || t.Kind() != reflect.String {
		return data, nil
	}
	return os.ExpandEnv(data.(string)), nil
}
