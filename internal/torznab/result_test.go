package torznab

// yep this was a bad idea
//
//import (
//	_ "embed"
//)
//
////go:embed examples/*
//var testExamples embed.FS
//
//func TestResult(t *testing.T) {
//	example, readErr := testExamples.ReadFile("examples/result1.xml")
//	assert.NoError(t, readErr)
//	result := &SearchResult{}
//	unmarshalErr := xml.Unmarshal(example, &result)
//	assert.NoError(t, unmarshalErr)
//	marshaled, marshalErr := result.Xml()
//	assert.NoError(t, marshalErr)
//	t.Logf("marshaled: %s", marshaled)
//	result2 := &SearchResult{}
//	unmarshalErr = xml.Unmarshal(marshaled, &result2)
//	assert.NoError(t, unmarshalErr)
//	assert.Equal(t, result, result2)
//}
