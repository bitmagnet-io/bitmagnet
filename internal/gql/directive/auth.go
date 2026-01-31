package directive

type (
	AuthDirective struct {
		Object string
		Action string
	}

	AuthDirectives map[AuthDirective]struct{}
)

func ExtractAuthDirectives(directives SchemaDirectives) AuthDirectives {
	result := make(AuthDirectives)

	for _, mp := range directives {
		for _, mp := range mp {
			for _, mp := range mp {
				result[AuthDirective{
					Action: mp["action"],
					Object: mp["object"],
				}] = struct{}{}
			}
		}
	}

	return result
}
