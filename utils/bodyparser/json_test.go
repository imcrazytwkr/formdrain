package bodyparser_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
)

func TestJsonParser_Object(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"name":"Ada","n":3,"ok":true,"tags":["a","b"]}`))
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "Ada" {
		t.Fatalf("name = %#v", got["name"])
	}
	if got["n"] != int64(3) {
		t.Fatalf("n = %#v (%T)", got["n"], got["n"])
	}
	if got["ok"] != true {
		t.Fatalf("ok = %#v", got["ok"])
	}
}

func TestJsonParser_TrimStrings(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"name":"  Ada  ","tags":[" a "," b "]}`))
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "Ada" {
		t.Fatalf("name = %#v", got["name"])
	}
	tags, ok := got["tags"].([]any)
	if !ok || len(tags) != 2 || tags[0] != "a" || tags[1] != "b" {
		t.Fatalf("tags = %#v", got["tags"])
	}
}

func TestJsonParser_Numbers(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"i":42,"f":1.5}`))
	if err != nil {
		t.Fatal(err)
	}
	if got["i"] != int64(42) {
		t.Fatalf("i = %#v (%T)", got["i"], got["i"])
	}
	if got["f"] != 1.5 {
		t.Fatalf("f = %#v (%T)", got["f"], got["f"])
	}
}

func TestJsonParser_SkipNulls(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"a":null,"b":[null,"x",null],"c":[null,null]}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["a"]; ok {
		t.Fatalf("a should be omitted: %#v", got)
	}
	tags, ok := got["b"].([]any)
	if !ok || len(tags) != 1 || tags[0] != "x" {
		t.Fatalf("b = %#v", got["b"])
	}
	if _, ok := got["c"]; ok {
		t.Fatalf("c should be omitted: %#v", got)
	}
}

func TestJsonParser_Malformed(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	if _, err := p.Parse([]byte(`{`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestJsonParser_NestedArrayRejected(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	if _, err := p.Parse([]byte(`{"x":[[1]]}`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestJsonParser_CombinedTypesRejected(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	if _, err := p.Parse([]byte(`{"x":[1,"a"]}`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestJsonParser_EmptyObject(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}

func TestJsonParser_NestedObjectSkipped(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"name":"Ada","meta":{"x":1}}`))
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "Ada" {
		t.Fatalf("name = %#v", got["name"])
	}
	if _, ok := got["meta"]; ok {
		t.Fatalf("nested object should be skipped: %#v", got)
	}
}

func TestJsonParser_EmptyArrayOmitted(t *testing.T) {
	t.Parallel()

	p := bodyparser.NewJsonParser()
	got, err := p.Parse([]byte(`{"tags":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["tags"]; ok {
		t.Fatalf("empty array should be omitted: %#v", got)
	}
}
