package utils

func CopyStringStringMap(s map[string]string) map[string]string {
	d := make(map[string]string)
	for key, value := range s {
		d[key] = value
	}
	return d
}
