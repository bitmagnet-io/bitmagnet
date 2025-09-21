package batch

import (
	"io"
	"mime"
	"net/http"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Transport provides a transport compatible with apollo-link-batch-http requests.
// The code was borrowed from https://github.com/skaji/gqlgen-apollo-batch.
type Transport struct{}

var _ graphql.Transport = Transport{}

func (Transport) Supports(r *http.Request) bool {
	if r.Header.Get("Upgrade") != "" {
		return false
	}
	if r.Method != http.MethodPost {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return false
	}
	return mediaType == "application/json"
}

func (Transport) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	start := graphql.Now()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res := graphqlErrorResponse("failed to read request body: " + err.Error())
		writeJSON(w, http.StatusBadRequest, res)
		return
	}

	isBatch := isBatchRequestBody(body)
	var paramsList []*graphql.RawParams
	if isBatch {
		if err := jsonUnmarshal(body, &paramsList); err != nil {
			res := graphqlErrorResponse("json body could not be decoded: " + err.Error())
			writeJSON(w, http.StatusBadRequest, res)
			return
		}
		if len(paramsList) == 0 {
			res := graphqlErrorResponse("json body must not be an empty array")
			writeJSON(w, http.StatusBadRequest, res)
			return
		}
	} else {
		var params *graphql.RawParams
		if err := jsonUnmarshal(body, &params); err != nil {
			res := graphqlErrorResponse("json body could not be decoded: " + err.Error())
			writeJSON(w, http.StatusBadRequest, res)
			return
		}
		paramsList = append(paramsList, params)
	}

	end := graphql.Now()
	for _, params := range paramsList {
		params.ReadTime = graphql.TraceTiming{
			Start: start,
			End:   end,
		}
	}

	out := make([]*graphql.Response, len(paramsList))
	errs := make([]gqlerror.List, len(paramsList))
	if len(paramsList) == 1 {
		ctx := r.Context()
		params := paramsList[0]
		out[0], errs[0] = doExec(ctx, exec, params)
	} else {
		var wg sync.WaitGroup
		wg.Add(len(paramsList))
		for index, params := range paramsList {
			ctx := r.Context()
			index := index
			params := params
			go func() {
				defer wg.Done()
				out[index], errs[index] = doExec(ctx, exec, params)
			}()
		}
		wg.Wait()
	}

	if isBatch {
		writeJSON(w, http.StatusOK, out) // XXX always http.StatusOK
	} else {
		writeJSON(w, statusFor(errs[0]), out[0])
	}
}
