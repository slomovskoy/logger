package logger

import "github.com/mitchellh/mapstructure"

// Hooks is a runtime config DAO
type Hooks interface{}

// Hook is an any hook struct for config file
type Hook struct {
	Name     string      `mapstructure:"name"`
	Settings interface{} `mapstructure:"settings"`
}

// HookKV is a key-value hook in config file
type HookKV map[string]interface{}

func listToHooks(rawHooks []interface{}) []Hook {
	var hookList []Hook
	for i := range rawHooks {
		if rawMapHooks, ok := rawHooks[i].(map[interface{}]interface{}); ok {
			// convert map[interface{}]interface{} => map[string]interface{}
			// and
			// make mapstructure.Decode
			mapHook := make(map[string]interface{})
			for k, v := range rawMapHooks {
				if name, strOk := k.(string); strOk {
					mapHook[name] = v
				}
			}
			var hook Hook
			if err := mapstructure.WeakDecode(mapHook, &hook); err == nil {
				hookList = append(hookList, hook)
			}
		}
	}
	return hookList
}
