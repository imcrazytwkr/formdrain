package form_response_test

import (
	"testing"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_response"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestListFormResponsesByField_String(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	testutil.SeedForm(t, db, 11, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	idD := "00000000-0000-4000-8000-00000000000d"
	idE := "00000000-0000-4000-8000-00000000000e"
	idF := "00000000-0000-4000-8000-00000000000f"
	idOther := "00000000-0000-4000-8000-0000000000ff"

	mustSave := func(id string, formID int64, payload map[string]any) {
		t.Helper()
		if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
			Id: id, FormId: formID, SchemaVersion: 1, Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(idA, 10, map[string]any{"email": "b"})
	mustSave(idB, 10, map[string]any{"email": "a"})
	mustSave(idC, 10, map[string]any{})
	mustSave(idD, 10, map[string]any{"email": nil})
	mustSave(idE, 10, map[string]any{"email": 1})
	mustSave(idF, 10, map[string]any{"email": "a"})
	mustSave(idOther, 11, map[string]any{"email": "aaa"})

	asc, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, false, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idB, idF, idA, idC, idD, idE}; !sameResponseIDs(asc, want) {
		t.Fatalf("asc = %#v want %#v", responseIDs(asc), want)
	}

	desc, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, true, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idB, idF, idC, idD, idE}; !sameResponseIDs(desc, want) {
		t.Fatalf("desc = %#v want %#v", responseIDs(desc), want)
	}
}

func TestListFormResponsesByField_Number(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	idD := "00000000-0000-4000-8000-00000000000d"

	mustSave := func(id string, payload map[string]any) {
		t.Helper()
		if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
			Id: id, FormId: 10, SchemaVersion: 1, Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(idA, map[string]any{"age": 10})
	mustSave(idB, map[string]any{"age": 2})
	mustSave(idC, map[string]any{})
	mustSave(idD, map[string]any{"age": 2})

	asc, err := repo.ListFormResponsesByField(ctx, 10, "age", fc.FieldTypeNumber, false, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idB, idD, idA, idC}; !sameResponseIDs(asc, want) {
		t.Fatalf("asc = %#v want %#v", responseIDs(asc), want)
	}

	desc, err := repo.ListFormResponsesByField(ctx, 10, "age", fc.FieldTypeNumber, true, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idB, idD, idC}; !sameResponseIDs(desc, want) {
		t.Fatalf("desc = %#v want %#v", responseIDs(desc), want)
	}

	page, err := repo.ListFormResponsesByField(ctx, 10, "age", fc.FieldTypeNumber, false, false, 2, idD, 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idC}; !sameResponseIDs(page, want) {
		t.Fatalf("after 2/D = %#v want %#v", responseIDs(page), want)
	}
}

func TestListFormResponsesByField_Boolean(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	idD := "00000000-0000-4000-8000-00000000000d"

	mustSave := func(id string, payload map[string]any) {
		t.Helper()
		if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
			Id: id, FormId: 10, SchemaVersion: 1, Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(idA, map[string]any{"ok": true})
	mustSave(idB, map[string]any{"ok": false})
	mustSave(idC, map[string]any{})
	mustSave(idD, map[string]any{"ok": false})

	asc, err := repo.ListFormResponsesByField(ctx, 10, "ok", fc.FieldTypeBoolean, false, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idB, idD, idA, idC}; !sameResponseIDs(asc, want) {
		t.Fatalf("asc = %#v want %#v", responseIDs(asc), want)
	}

	desc, err := repo.ListFormResponsesByField(ctx, 10, "ok", fc.FieldTypeBoolean, true, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idB, idD, idC}; !sameResponseIDs(desc, want) {
		t.Fatalf("desc = %#v want %#v", responseIDs(desc), want)
	}
}

func TestListFormResponsesByField_ExactKey(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	field := `user.email'; DROP TABLE forms;--`

	if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
		Id: idA, FormId: 10, SchemaVersion: 1,
		Payload: map[string]any{field: "b"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
		Id: idB, FormId: 10, SchemaVersion: 1,
		Payload: map[string]any{field: "a"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
		Id: idC, FormId: 10, SchemaVersion: 1,
		Payload: map[string]any{"user": map[string]any{"email": "zzz"}},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := repo.ListFormResponsesByField(ctx, 10, field, fc.FieldTypeString, false, false, nil, "", 10)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idB, idA, idC}; !sameResponseIDs(got, want) {
		t.Fatalf("got %#v want %#v", responseIDs(got), want)
	}
}

func TestListFormResponsesByField_AllNull(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"

	for _, id := range []string{idB, idA} {
		if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
			Id: id, FormId: 10, SchemaVersion: 1, Payload: map[string]any{},
		}); err != nil {
			t.Fatal(err)
		}
	}

	for _, desc := range []bool{false, true} {
		got, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, desc, false, nil, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{idA, idB}; !sameResponseIDs(got, want) {
			t.Fatalf("desc=%v got %#v want %#v", desc, responseIDs(got), want)
		}
	}
}

func TestListFormResponsesByField_Keyset(t *testing.T) {
	db := testutil.OpenSqlite(t)
	ctx := t.Context()
	testutil.SeedSite(t, db, 1, "example.com")
	testutil.SeedForm(t, db, 10, 1)
	repo := frr.NewSqliteFormResponseRepository(db)

	idA := "00000000-0000-4000-8000-00000000000a"
	idB := "00000000-0000-4000-8000-00000000000b"
	idC := "00000000-0000-4000-8000-00000000000c"
	idD := "00000000-0000-4000-8000-00000000000d"
	idE := "00000000-0000-4000-8000-00000000000e"
	idF := "00000000-0000-4000-8000-00000000000f"

	mustSave := func(id string, payload map[string]any) {
		t.Helper()
		if _, err := repo.SaveFormResponse(ctx, &form_response.FormResponse{
			Id: id, FormId: 10, SchemaVersion: 1, Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(idA, map[string]any{"email": "b"})
	mustSave(idB, map[string]any{"email": "a"})
	mustSave(idC, map[string]any{})
	mustSave(idD, map[string]any{"email": nil})
	mustSave(idE, map[string]any{"email": 1})
	mustSave(idF, map[string]any{"email": "a"})

	page, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, false, false, "a", idF, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idA, idC}; !sameResponseIDs(page, want) {
		t.Fatalf("after a/F = %#v want %#v", responseIDs(page), want)
	}

	rest, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, false, true, nil, idC, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idD, idE}; !sameResponseIDs(rest, want) {
		t.Fatalf("after null/C = %#v want %#v", responseIDs(rest), want)
	}

	descPage, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, true, false, "a", idB, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idF, idC}; !sameResponseIDs(descPage, want) {
		t.Fatalf("desc after a/B = %#v want %#v", responseIDs(descPage), want)
	}

	nullPage, err := repo.ListFormResponsesByField(ctx, 10, "email", fc.FieldTypeString, true, true, nil, idC, 1)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{idD}; !sameResponseIDs(nullPage, want) {
		t.Fatalf("desc null group = %#v want %#v", responseIDs(nullPage), want)
	}
}
