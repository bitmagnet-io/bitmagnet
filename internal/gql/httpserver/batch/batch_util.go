package batch

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"unicode"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func isBatchRequestBody(body []byte) bool {
	b := bytes.TrimLeftFunc(body, unicode.IsSpace)
	if len(b) == 0 {
		return false
	}

	return b[0] == '['
}

func doExec(
	ctx context.Context,
	exec graphql.GraphExecutor,
	params *graphql.RawParams,
) (*graphql.Response, gqlerror.List) {
	rc, err := exec.CreateOperationContext(ctx, params)
	if err != nil {
		return exec.DispatchError(graphql.WithOperationContext(ctx, rc), err), err
	}

	ctx2 := graphql.WithOperationContext(ctx, rc)
	responses, ctx3 := exec.DispatchOperation(ctx2, rc)

	return responses(ctx3), nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		b = []byte(`{"errors":[{"message":"failed to marshal data to json"}],"data":null}`)
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func jsonUnmarshal(body []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	return dec.Decode(v)
}

func graphqlErrorResponse(msg string) *graphql.Response {
	return &graphql.Response{Errors: gqlerror.List{{Message: msg}}}
}

func statusFor(err gqlerror.List) int {
	if err != nil && errcode.GetErrorKind(err) == errcode.KindProtocol {
		return http.StatusUnprocessableEntity
	}

	return http.StatusOK
}
