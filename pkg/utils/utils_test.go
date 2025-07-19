package utils

import (
    "os"
    "testing"
)

func TestDirExistsAndMkdir(t *testing.T) {
    testDir := "test_tmp_dir"
    defer os.RemoveAll(testDir)
    if DirExists(testDir) {
        t.Fatalf("DirExists: Should not yet exist: %s", testDir)
    }
    if err := MkdirIfMissing(testDir); err != nil {
        t.Fatalf("MkdirIfMissing: %v", err)
    }
    if !DirExists(testDir) {
        t.Fatalf("DirExists: Should exist after MkdirIfMissing: %s", testDir)
    }
}

func TestListVisibleFiles(t *testing.T) {
    testDir := "test_tmp_dir2"
    os.RemoveAll(testDir)
    if err := MkdirIfMissing(testDir); err != nil {
        t.Fatalf("MkdirIfMissing: %v", err)
    }
    f1 := testDir + "/foo.txt"
    f2 := testDir + "/.hidden"
    f, _ := os.Create(f1); f.Close()
    fh, _ := os.Create(f2); fh.Close()
    files, err := ListVisibleFiles(testDir)
    if err != nil {
        t.Fatalf("ListVisibleFiles: %v", err)
    }
    want := "foo.txt"
    found := false
    for _, file := range files {
        if file == want {
            found = true
        }
        if len(file) > 0 && file[0] == '.' {
            t.Errorf("Should not list hidden file: %s", file)
        }
    }
    if !found {
        t.Errorf("Expected to find %s", want)
    }
    os.RemoveAll(testDir)
}
