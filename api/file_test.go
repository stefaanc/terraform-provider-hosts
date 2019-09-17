package api

import (
    "crypto/sha1"
    "encoding/hex"
    "io/ioutil"
    "os"
    "testing"
)

// -----------------------------------------------------------------------------

func resetFileTestEnv() {
    hosts = (*anchor)(nil)
    Init()
}

// -----------------------------------------------------------------------------

func Test_LookupFile(t *testing.T) {
    var test string

    test = "found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "f"
        f.Notes = "..."
        addFile(f)

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := LookupFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ LookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else {

            // --------------------

            if file.ID != int(f.id) {
               t.Errorf("[ LookupFile(fQuery).ID ] expected: %#v, actual: %#v", int(f.id), file.ID)
            }

            // --------------------

            if file.Path != f.Path {
                t.Errorf("[ LookupFile(fQuery).Path ] expected: %#v, actual: %#v", f.Path, file.Path)
            }

            // --------------------

            if file.Notes != f.Notes {
                t.Errorf("[ LookupFile(fQuery).Notes ] expected: %#v, actual: %#v", f.Notes, file.Notes)
            }

            // --------------------

            if file.id != 0 {
               t.Errorf("[ LookupFile(fQuery).id ] expected: %#v, actual: %#v", 0, file.id)
            }

            // --------------------

            if file.checksum != "" {
                t.Errorf("[ LookupFile(fQuery).checksum ] expected: %#v, actual: %#v", "", file.checksum)
            }

        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fQuery := new(File)
        fQuery.ID = 42

        file := LookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ LookupFile(fQuery) ] expected: %s, actual: %#v", "<error>", file)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_CreateFile(t *testing.T) {
    var test string

    test = "missing-Path"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        fValues := new(File)

        err := CreateFile(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateFile(fValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "f"
        addFile(f)

        // --------------------

        err := CreateFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"

        err := CreateFile(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createFile(fValues) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(fValues.Path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ createFile() ] cannot make test-directory")
        }

        // --------------------

        fValues := new(File)
        fValues.Path = path

        err = CreateFile(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateFile(fValues) ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_createFile(t *testing.T) {
    var test string

    test = "created-new-file"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"
        fValues.Notes = "..."

        err := createFile(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createFile(fValues) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        f := lookupFile(fValues)
        if f == nil {
            t.Errorf("[ lookupFile(fValues) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {

            // --------------------

            if f.id == 0 {
                t.Errorf("[ lookupFile(fValues).id ] expected: not %#v, actual: %#v", 0, f.id)
            }

            // --------------------

            if f.ID != int(f.id) {
                t.Errorf("[ lookupFile(fValues).ID ] expected: %#v, actual: %#v", f.id, f.ID)
            }

            // --------------------

            if f.Path != fValues.Path {
                t.Errorf("[ lookupFile(fValues).Path ] expected: %#v, actual: %#v", fValues.Path, f.Path)
            }

            // --------------------

            if f.Notes != fValues.Notes {
                t.Errorf("[ lookupFile(fValues).Notes ] expected: %#v, actual: %#v", fValues.Notes, f.Notes)
            }

            // --------------------

            checksum := sha1.Sum(nil)
            expected := hex.EncodeToString(checksum[:])
            if f.checksum != expected {
                t.Errorf("[ lookupFile(fValues).checksum ] expected: %#v, actual: %#v", expected, f.checksum)
            }

            // --------------------

            data, err := ioutil.ReadFile(fValues.Path)
            if err != nil {
                t.Errorf("[ ioutil.ReadFile(fValues.Path) ] expected: %s, actual: %#v", "not <nil>", err)
            } else if len(data) != 0 {
                t.Errorf("[ len(ioutil.ReadFile(fValues.Path)) ] expected: %#v, actual: %#v", 0, len(data))
            }
        }

        // --------------------

        os.Remove(fValues.Path)
    })

    test = "read-existing-file"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ createFile() ] cannot write test-file")
        }

        // --------------------

        fValues := new(File)
        fValues.Path = path

        err = createFile(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createFile(fValues) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        f := lookupFile(fValues)
        if f == nil {
            t.Errorf("[ lookupFile(fValues) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {

            // --------------------

            if f.id == 0 {
                t.Errorf("[ lookupFile(fValues).id ] expected: %#v, actual: %#v", 0, f.id)
            }

            // --------------------

            if f.ID != int(f.id) {
                t.Errorf("[ lookupFile(fValues).ID ] expected: %#v, actual: %#v", f.id, f.ID)
            }

            // --------------------

            if f.Path != fValues.Path {
                t.Errorf("[ lookupFile(fValues).Path ] expected: %#v, actual: %#v", fValues.Path, f.Path)
            }

            // --------------------

            if f.Notes != fValues.Notes {
                t.Errorf("[ lookupFile(fValues).Notes ] expected: %#v, actual: %#v", fValues.Notes, f.Notes)
            }

            // --------------------

            checksum := sha1.Sum(data)
            expected := hex.EncodeToString(checksum[:])
            if f.checksum != expected {
                t.Errorf("[ lookupFile(fValues).checksum ] expected: %#v, actual: %#v", expected, f.checksum)
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ createFile() ] cannot make test-directory")
        }

        // --------------------

        fValues := new(File)
        fValues.Path = path

        err = createFile(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ createFile(fValues) ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        f := lookupFile(fValues)
        if f != nil {
            t.Errorf("[ lookupFile(fValues) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_fRead(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)
        f.ID = 0

        // --------------------

        _, err = f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        f.ID = 1

        // --------------------

        _, err = f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        _, err = f.Read()

        // --------------------

        if err != nil {
            t.Errorf("[ f.Read() ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_doesnt-exist.txt"
        addFile(f)

        // --------------------

        _, err := f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

func Test_readFile(t *testing.T) {
    var test string

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        file, err := readFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ readFile(f) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if file == nil {
            t.Errorf("[ readFile(f) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {


            // --------------------

            checksum := sha1.Sum(data)
            expected := hex.EncodeToString(checksum[:])
            if file.checksum != expected {
                t.Errorf("[ file.checksum ] expected: %#v, actual: %#v", expected, file.checksum)
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ readFile() ] cannot make test-directory")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        _, err = readFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ readFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_fUpdate(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)
        f.ID = 0

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "not found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        f.ID = 1

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Update() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = f.Update(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ f.Update(fValues) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ f.Update() ] cannot make test-directory")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Update(fValues) ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_updateFile(t *testing.T) {
    var test string

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = updateFile(f, fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateFile(f) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if f == nil {
            t.Errorf("[ updateFile(f) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {

            // --------------------

            checksum := sha1.Sum([]byte("# some updated data"))
            expected := hex.EncodeToString(checksum[:])
            if f.checksum != expected {
                t.Errorf("[ f.checksum ] expected: %#v, actual: %#v", expected, f.checksum)
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot make test-directory")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("# some updated data")

        err = updateFile(f, fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ updateFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_fDelete(t *testing.T) {
    var test string

    test = "no-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Delete() ] cannot write test-file")
        }

        f := new(File)
        addFile(f)
        f.ID = 0

        // --------------------

        err = f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Delete() ] cannot write test-file")
        }

        f := new(File)
        f.ID = 1

        // --------------------

        err = f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Delete() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        addFile(f)

        // --------------------

        err = f.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ f.Delete() ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_doesnt-exist.txt"
        addFile(f)

        // --------------------

        err := f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

func Test_deleteFile(t *testing.T) {
    var test string

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        f.Notes = "..."
        addFile(f)

        // --------------------

        err = deleteFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ deleteFile(f) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if f.id != 0 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 0, f.id)
        }

        // --------------------

        if f.ID != 0 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 0, f.ID)
        }

        // --------------------

        if f.Path != "" {
            t.Errorf("[ f.Path ] expected: %#v, actual: %#v", "", f.Path)
        }

        // --------------------

        if f.Notes != "" {
            t.Errorf("[ f.Notes ] expected: %#v, actual: %#v", "", f.Notes)
        }

        // --------------------

        if f.checksum != "" {
            t.Errorf("[ f.checksum ] expected: %#v, actual: %#v", "", f.checksum)
        }

        // --------------------

        f = lookupFile(f)
        if f != nil {
            t.Errorf("[ lookupFile(f) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_doesnt-exist.txt"
        addFile(f)

        // --------------------

        err := deleteFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ deleteFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}
