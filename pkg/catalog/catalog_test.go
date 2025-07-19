package catalog

import (
    "os"
    "testing"
)

func TestParseCatalogLine(t *testing.T) {
    l := "2007-07-11_Jakob.pdf*jakob*kindergeld"
    entry := ParseCatalogLine(l)
    if entry.Name != "2007-07-11_Jakob.pdf" {
        t.Error("Fail name parse")
    }
    if len(entry.Tags) != 2 || entry.Tags[0] != "jakob" || entry.Tags[1] != "kindergeld" {
        t.Error("Fail tags parse")
    }
    if entry.Type != "file" {
        t.Error("Fail type file parse")
    }
    u := ParseCatalogLine("https://foo.bar:any*other")
    if u.Type != "url" {
        t.Error("Fail detect url")
    }
}

func TestCatalogSaveLoadRoundtrip(t *testing.T) {
    tmpfile := ".cat_test"
    os.Remove(tmpfile)
    orig := []CatEntry{
        {Name: "foo.pdf", Type: "file", Tags: []string{"x", "y"}},
        {Name: "http://test.com", Type: "url", Tags: []string{"z"}},
    }
    // Temporarily monkeypatch CatalogFilename
    old := CatalogFilename
    CatalogFilename = tmpfile
    defer func() { CatalogFilename = old; os.Remove(tmpfile) }()
    if err := SaveCatalog(orig); err != nil {
        t.Fatalf("SaveCatalog: %v", err)
    }
    loaded, err := LoadCatalog()
    if err != nil {
        t.Fatalf("LoadCatalog: %v", err)
    }
    if len(loaded) != 2 || loaded[0].Name != "foo.pdf" || loaded[1].Type != "url" {
        t.Error("Roundtrip failed")
    }
    if len(loaded[0].Tags) != 2 || loaded[0].Tags[1] != "y" {
        t.Error("Tags not roundtripped")
    }
}
