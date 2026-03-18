package metadata

// The first key is the object key, whose attribute path is dot separated.
// Especially, for function return value, the key is an empty string.
// The second key is the field name and the value is the field description.
type ObjectDescription map[string]map[string]string
