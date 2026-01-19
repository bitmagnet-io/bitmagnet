# Fake sample plugins used for testing
PLUGINS := fake_search fake_search2

all: $(PLUGINS:%=%/plugin.wasm)

clean:
	rm -f $(PLUGINS:%=%/plugin.wasm)

%/plugin.wasm: %/plugin.go
	GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o $@ ./$*

%/plugin.wasm: %/plugin.ts
	cd $(@D) && npm run build