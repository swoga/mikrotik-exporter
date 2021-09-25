package utils

func SetDefaultString(v *string, d string) {
	if *v == "" {
		*v = d
	}
}
