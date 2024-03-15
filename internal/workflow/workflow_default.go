package workflow

import (
	_ "embed"
)

//go:embed workflow.default.yaml
var workflowDefaultYaml string

//var workflowDefault actionDefinition

//func init() {
//  var rawWorkflow map[string]interface{}
//  parseErr := yaml.Unmarshal([]byte(workflowDefaultYaml), &rawWorkflow)
//  if parseErr != nil {
//    panic(parseErr)
//  }
//
//}
