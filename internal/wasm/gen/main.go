package gen

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen/jen"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/gen/spec"
	_ "github.com/planetscale/vtprotobuf/features/marshal"
	_ "github.com/planetscale/vtprotobuf/features/size"
	_ "github.com/planetscale/vtprotobuf/features/unmarshal"
	vtgenerator "github.com/planetscale/vtprotobuf/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var vtFeatureNames = []string{"marshal", "unmarshal", "size"}

func GeneratePlugin(plugin *protogen.Plugin) error {
	vtgen, err := vtgenerator.NewGenerator(plugin.Files, vtFeatureNames, &vtgenerator.Extensions{})
	if err != nil {
		return err
	}
	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		genFile := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_vtproto.pb.go", file.GoImportPath)
		vtgen.GenerateFile(genFile, file)
	}

	s, err := spec.New(plugin)
	if err != nil {
		plugin.Error(err)
		return err
	}

	specJson, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal spec to JSON: %w", err)
	}
	specFile := plugin.NewGeneratedFile("proto-spec.json", "")
	specFile.P(string(specJson))

	for _, f := range s.Files {
		if !f.Generate {
			continue
		}
		if err := genPBFile(plugin, f); err != nil {
			return err
		}
		if err := genHostFile(plugin, f); err != nil {
			return err
		}
		if err := genPluginFile(plugin, f); err != nil {
			return err
		}
	}

	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	return nil
}

func genPBFile(p *protogen.Plugin, f spec.FileInfo) error {
	return renderFile(p, jen.GenTypes(f), f.GeneratedFilenamePrefix+".pb.go")
}

func genHostFile(p *protogen.Plugin, f spec.FileInfo) error {
	if len(f.PluginServices) == 0 && f.HostService == nil {
		return nil
	}

	return renderFile(p, jen.GenHostFile(f), f.GeneratedFilenamePrefix+"_host.pb.go")
}

func genPluginFile(p *protogen.Plugin, f spec.FileInfo) error {
	if len(f.PluginServices) == 0 && f.HostService == nil {
		return nil
	}
	return renderFile(p, jen.GenPluginFile(f), f.GeneratedFilenamePrefix+"_plugin.pb.go")
}

func renderFile(p *protogen.Plugin, f *jen.File, path string) error {
	g := p.NewGeneratedFile(path, protogen.GoImportPath(path))
	buf := &strings.Builder{}
	if err := f.Render(buf); err != nil {
		err = fmt.Errorf("failed to render file %s: %w", path, err)
		p.Error(err)
		return err
	}
	g.P(buf.String())
	return nil
}
