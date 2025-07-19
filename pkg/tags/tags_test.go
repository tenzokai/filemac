package tags

import (
    "os"
    "testing"
    "github.com/tenzokai/filemac/pkg/catalog"
)

func withTempCatalog(entries []catalog.CatEntry, testfunc func()) {
    tmp := ".cat_test_tags"
    old := catalog.CatalogFilename
    catalog.CatalogFilename = tmp
    defer func() { catalog.CatalogFilename = old; os.Remove(tmp) }()
    catalog.SaveCatalog(entries)
    testfunc()
}

func TestCmdAddRemoveReplaceTag(t *testing.T) {
    withTempCatalog([]catalog.CatEntry{
        {Name: "foo.txt", Tags: []string{"a", "b"}},
        {Name: "bar.txt", Tags: []string{"c"}},
    }, func() {
        CmdAddTag("1", "c") // add new tag to entry 1
        newCat, _ := catalog.LoadCatalog()
        found := false
        for _, tag := range newCat[0].Tags {
            if tag == "c" {
                found = true
            }
        }
        if !found {
            t.Error("Tag not added")
        }
        CmdRemoveTag("1", "a")
        newCat, _ = catalog.LoadCatalog()
        for _, tag := range newCat[0].Tags {
            if tag == "a" {
                t.Error("Tag not removed")
            }
        }
        CmdReplaceTag("2", "c", "z")
        newCat, _ = catalog.LoadCatalog()
        tagOk := false
        for _, tag := range newCat[1].Tags {
            if tag == "z" {
                tagOk = true
            }
        }
        if !tagOk {
            t.Error("Tag not replaced")
        }
    })
}
