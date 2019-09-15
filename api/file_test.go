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

func Test_CreateFile(t *testing.T) {
    var test string

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "f"
        _ = hosts.addFile(f)

        err := CreateFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "no-path"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fValues := new(File)

        err := CreateFile(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateFile(fValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "create"
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

    test = "create-new-file"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"

        err := createFile(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createFile(fValues) ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        f := GetFile(fValues)
        if f == nil {
            t.Errorf("[ GetFile(fValues) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {

            // --------------------

            if f.id == 0 {
                t.Errorf("[ GetFile(fValues).id ] expected: not %#v, actual: %#v", 0, f.id)
            }

            // --------------------

            if f.ID != f.id {
                t.Errorf("[ GetFile(fValues).ID ] expected: %#v, actual: %#v", f.id, f.ID)
            }

            // --------------------

            if f.Path != fValues.Path {
                t.Errorf("[ GetFile(fValues).Path ] expected: %#v, actual: %#v", fValues.Path, f.Path)
            }


            // --------------------

            checksum := sha1.Sum(nil)
            expected := hex.EncodeToString(checksum[:])
            if f.checksum != expected {
                t.Errorf("[ GetFile(fValues).checksum ] expected: %#v, actual: %#v", expected, f.checksum)
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
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ hosts.createFile() ] cannot write test-file")
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

        f := GetFile(fValues)
        if f == nil {
            t.Errorf("[ GetFile(fValues) ] expected: %s, actual: %#v", "not <nil>", f)
        } else {

            // --------------------

            if f.id == 0 {
                t.Errorf("[ GetFile(fValues).id ] expected: %#v, actual: %#v", 0, f.id)
            }

            // --------------------

            if f.ID != f.id {
                t.Errorf("[ GetFile(fValues).ID ] expected: %#v, actual: %#v", f.id, f.ID)
            }

            // --------------------

            if f.Path != fValues.Path {
                t.Errorf("[ GetFile(fValues).Path ] expected: %#v, actual: %#v", fValues.Path, f.Path)
            }

            // --------------------

            checksum := sha1.Sum(data)
            expected := hex.EncodeToString(checksum[:])
            if f.checksum != expected {
                t.Errorf("[ GetFile(fValues).checksum ] expected: %#v, actual: %#v", expected, f.checksum)
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

        f := GetFile(fValues)
        if f != nil {
            t.Errorf("[ GetFile(fValues) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_fRead(t *testing.T) {
    var test string

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

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

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ f.Read() ] cannot make test-directory")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

        // --------------------

        _, err = f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read() ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_readFile(t *testing.T) {
    var test string

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

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
        _ = hosts.addFile(f)

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

    test = "update"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ f.Update() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("some updated data")

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
        _ = hosts.addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("some updated data")

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

    test = "update"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("some updated data")

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

            checksum := sha1.Sum([]byte("some updated data"))
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
        _ = hosts.addFile(f)

        // --------------------

        fValues := new(File)
//        fValues.data = []byte("some updated data")

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


    test = "delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

        // --------------------

        err = deleteFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ deleteFile(f) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_test-hosts.txt"
        _ = hosts.addFile(f)

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

    test = "delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("some data")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteFile() ] cannot write test-file")
        }

        f := new(File)
        f.Path = path
        _ = hosts.addFile(f)

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

        if f.checksum != "" {
            t.Errorf("[ f.checksum ] expected: %#v, actual: %#v", "", f.checksum)
        }

        // --------------------

        f = GetFile(f)
        if f != nil {
            t.Errorf("[ GetFile(f) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_test-hosts.txt"
        _ = hosts.addFile(f)

        // --------------------

        err := deleteFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ deleteFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}
