package env

import (
	"os"

	"github.com/joho/godotenv"
)

func providerValues(values ...string) VarsProvider {
	return func() map[string]string {
		env := make(map[string]string, len(values))

		for _, s := range values {
			for j := range len(s) {
				if s[j] == '=' {
					key := s[:j]
					if _, ok := env[key]; !ok {
						var value string

						if len(s) > j+1 {
							value = s[j+1:]
						}

						env[key] = value
					}

					break
				}
			}
		}

		return env
	}
}

func providerOSEnviron() VarsProvider {
	return providerValues(os.Environ()...)
}

func providerDotEnv(filenames ...string) VarsProvider {
	return func() map[string]string {
		env, err := godotenv.Read(filenames...)
		if err != nil {
			return nil
		}

		return env
	}
}
