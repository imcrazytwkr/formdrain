package bodyparser_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
)

func TestFormParser(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewFormParser()
	got, err := p.Parse([]byte("name=Ada&tags=a&tags=b&empty=null&skip="))
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "Ada" {
		t.Fatalf("name = %#v", got["name"])
	}
	tags, ok := got["tags"].([]string)
	if !ok || len(tags) != 2 {
		t.Fatalf("tags = %#v", got["tags"])
	}
	if _, ok := got["empty"]; ok {
		t.Fatalf("empty should be omitted: %#v", got)
	}
}
