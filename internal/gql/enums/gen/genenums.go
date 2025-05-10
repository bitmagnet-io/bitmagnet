package main

import (
	"os"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/gql/enums"
)

func main() {
	gqlParts := make([]string, 0, len(enums.Enums))
	for _, e := range enums.Enums {
		gqlParts = append(gqlParts, genGql(e.Name, e.Values))
	}

	f, fErr := os.Create("./graphql/schema/enums.graphqls")
	checkErr(fErr)

	_, wErr := f.WriteString(strings.Join(gqlParts, "\n"))
	checkErr(wErr)
}

func genGql(name string, values []string) string {
	str := "enum " + name + " {\n"
	for _, value := range values {
		str += "  " + value + "\n"
	}

	str += "}\n"

	return str
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
