//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "bytes"
    "crypto/sha1"
    "encoding/hex"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "testing"
)

// -----------------------------------------------------------------------------

func resetFileTestEnv() {
    if hosts != nil {
        for _, hostsFile := range hosts.files {   // !!! avoid memory leaks
            hostsFile.file = nil
        }
        hosts = (*anchor)(nil)
    }
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
        f.hostsFile = new(fileObject)
        f.zones = append(f.zones, new(zoneObject))
        addFile(f)

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := LookupFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ LookupFile(fQuery) ] expected: not %#v, actual: %#v", nil, file)
        } else {

            // --------------------

            if file.id != 0 {
               t.Errorf("[ LookupFile(fQuery).id ] expected: %#v, actual: %#v", 0, file.id)
            }

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

            if file.hostsFile != nil {
                t.Errorf("[ lookupFile(fValues).hostsFile ] expected: not %#v, actual: %#v", nil, file.hostsFile)
            }

            // --------------------

            if file.zones != nil {
                t.Errorf("[ lookupFile(fValues).zones ] expected: %#v, actual: %#v", nil, file.zones)
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
            t.Errorf("[ LookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
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
            t.Errorf("[ CreateFile(fValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'fValues.Path'") {
            t.Errorf("[ CreateFile(fValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'fValues.Path'", err.Error())
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
            t.Errorf("[ CreateFile(f).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "already exists") {
            t.Errorf("[ CreateFile(f).err.Error() ] expected: contains %#v, actual: %#v", "already exists", err.Error())
        }
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"
        fValues.Notes = "..."

        err := CreateFile(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createFile(fValues).err ] expected: %#v, actual: %#v", nil, err)
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
            t.Errorf("[ CreateFile(fValues).err ] expected: %s, actual: %#v", "<error>", err)
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
            t.Errorf("[ createFile(fValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        f := lookupFile(fValues)
        if f == nil {
            t.Errorf("[ lookupFile(fValues) ] expected: not %#v, actual: %#v", nil, f)
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

            if f.hostsFile == nil {
                t.Errorf("[ lookupFile(fValues).hostsFile ] expected: not %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ lookupFile(fValues).hostsFile.file ] expected: %#v, actual: %#v", nil, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ lookupFile(fValues).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum(nil)
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ lookupFile(fValues).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if len(f.zones) != 0 {
                t.Errorf("[ lookupFile(fValues).zones ] expected: %#v, actual: %#v", 0, len(f.zones))
            }

            // --------------------

            data, err := ioutil.ReadFile(fValues.Path)
            if err != nil {
                t.Errorf("[ ioutil.ReadFile(fValues.Path).err ] expected: %#v, actual: %#v", nil, err)
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
            t.Errorf("[ createFile(fValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        f := lookupFile(fValues)
        if f == nil {
            t.Errorf("[ lookupFile(fValues) ] expected: not %#v, actual: %#v", nil, f)
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

            if f.hostsFile == nil {
                t.Errorf("[ lookupFile(fValues).hostsFile ] expected: %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ lookupFile(fValues).hostsFile.file ] expected: %#v, actual: %#v", nil, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ lookupFile(fValues).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum(data)
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ lookupFile(fValues).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if len(f.zones) != 1 {
                t.Errorf("[ lookupFile(fValues).zones ] expected: %#v, actual: %#v", 1, len(f.zones))
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
            t.Errorf("[ createFile(fValues).err ] expected: %s, actual: %#v", "<error>", err)
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

// -----------------------------------------------------------------------------

func Test_fRead(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 0

        // --------------------

        _, err := f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'f.ID'") {
            t.Errorf("[ f.Read().err.Error() ] expected: contains %#v, actual: %#v", "missing 'f.ID'", err.Error())
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 42

        // --------------------

        _, err := f.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "not found") {
            t.Errorf("[ f.Read().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
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
        f.Notes = "..."
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)
        
        f.zones = make([]*zoneObject, 0)

        // --------------------

        file, err := f.Read()

        // --------------------

        if err != nil {
            t.Errorf("[ f.Read().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if file == nil {
            t.Errorf("[ f.Read().file ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if file.id != 0 {
                t.Errorf("[ f.Read().file.id ] expected: %#v, actual: %#v", 0, file.id)
            }

            // --------------------

            if file.ID != 1 {
                t.Errorf("[ f.Read().file.ID ] expected: %#v, actual: %#v", 1, file.ID)
            }

            // --------------------

            if file.Path != path {
                t.Errorf("[ f.Read().file.Path ] expected: %#v, actual: %#v", path, file.Path)
            }

            // --------------------

            if file.Notes != "..." {
                t.Errorf("[ f.Read().file.Notes ] expected: %#v, actual: %#v", "...", file.Notes)
            }

            // --------------------

            if file.hostsFile != nil {
                t.Errorf("[ f.Read().file.hostsFile ] expected: %#v, actual: %#v", nil, file.hostsFile)
            }

            // --------------------

            if file.zones != nil {
                t.Errorf("[ f.Read().file.zones ] expected: %#v, actual: %#v", 0, file.zones)
            }
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
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
            t.Errorf("[ f.Read().err ] expected: %s, actual: %#v", "<error>", err)
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

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        f.zones = make([]*zoneObject, 0)

        // --------------------

        file, err := readFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ readFile(f).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if file == nil {
            t.Errorf("[ readFile(f) ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if file.id == 0 {
                t.Errorf("[ readFile(f).id ] expected: not %#v, actual: %#v", 0, file.id)
            }

            // --------------------

            if file.ID != int(f.id) {
                t.Errorf("[ readFile(f).ID ] expected: %#v, actual: %#v", f.id, file.ID)
            }

            // --------------------

            if file.Path != f.Path {
                t.Errorf("[ readFile(f).Path ] expected: %#v, actual: %#v", f.Path, file.Path)
            }

            // --------------------

            if file.Notes != f.Notes {
                t.Errorf("[ readFile(f).Notes ] expected: %#v, actual: %#v", f.Notes, file.Notes)
            }

            // --------------------

            if file.hostsFile == nil {
                t.Errorf("[ readFile(f).hostsFile ] expected: not %#v, actual: %#v", nil, file.hostsFile)
            } else {

                // --------------------

                if file.hostsFile.file != f {
                    t.Errorf("[ readFile(f).hostsFile.file ] expected: %#v, actual: %#v", f, file.hostsFile.file)
                }

                // --------------------

                if file.hostsFile.data != nil {
                    t.Errorf("[ lookupFile(fValues).hostsFile.data ] expected: %#v, actual: %#v", nil, file.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum(data)
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ readFile(f).hostsFile.checksum ] expected: %#v, actual: %#v", expected, file.hostsFile.checksum)
                }
            }

            // --------------------

            if len(file.zones) != 1 {
                t.Errorf("[ readFile(f).zones ] expected: %#v, actual: %#v", 1, len(file.zones))
            }
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_doesnt-exist.txt"
        addFile(f)

        // --------------------

        _, err := readFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ readFile(f).err ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_fUpdate(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 0

        // --------------------

        fValues := new(File)

        err := f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'f.ID'") {
            t.Errorf("[ f.Update().err.Error() ] expected: contains %#v, actual: %#v", "missing 'f.ID'", err.Error())
        }
    })

    test = "not found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 42

        // --------------------

        fValues := new(File)

        err := f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "not found") {
            t.Errorf("[ f.Update().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
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

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = int(f.id)
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "# some updated data")
        f.zones = append(f.zones, zo)

        // --------------------

        fValues := new(File)

        err = f.Update(fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ f.Update(fValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
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

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = int(f.id)
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "# some updated data")
        f.zones = append(f.zones, zo)

        // --------------------

        fValues := new(File)

        err = f.Update(fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ f.Update(fValues).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
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
        f.Notes = "..."
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = int(f.id)
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "# some updated data")
        f.zones = append(f.zones, zo)

        // --------------------

        fValues := new(File)
        fValues.Notes = "...updated notes"

        err = updateFile(f, fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateFile(f).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if f == nil {
            t.Errorf("[ updateFile(f) ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if f.Path != path {
                t.Errorf("[ updateFile(f).Path ] expected: %#v, actual: %#v", path, f.Path)
            }

            // --------------------

            if f.Notes != fValues.Notes {
                t.Errorf("[ updateFile(f).Notes ] expected: %#v, actual: %#v", fValues.Notes, f.Notes)
            }

            // --------------------

            if f.hostsFile == nil {
                t.Errorf("[ updateFile(f).hostsFile ] expected: not %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ updateFile(f).hostsFile.file ] expected: %#v, actual: %#v", f, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ updateFile(f).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum([]byte(zo.lines[0] + "\n"))
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ updateFile(f).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if len(f.zones) != 1 {
                t.Errorf("[ updateFile(f).zones ] expected: %#v, actual: %#v", 1, len(f.zones))
            }
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "not-needed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte("# some data\n")
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot write test-file")
        }
     
        info, err := os.Stat(path)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot stat test-file")
        }
        fileLastModified := info.ModTime()

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = int(f.id)
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "# some data")
        f.zones = append(f.zones, zo)

        // --------------------

        fValues := new(File)
        fValues.Notes = "...updated notes"

        err = updateFile(f, fValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateFile(f).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if f == nil {
            t.Errorf("[ updateFile(f) ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if f.Path != path {
                t.Errorf("[ updateFile(f).Path ] expected: %#v, actual: %#v", path, f.Path)
            }

            // --------------------

            if f.Notes != fValues.Notes {
                t.Errorf("[ updateFile(f).Notes ] expected: %#v, actual: %#v", fValues.Notes, f.Notes)
            }

            // --------------------

            if f.hostsFile == nil {
                t.Errorf("[ updateFile(f).hostsFile ] expected: not %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ updateFile(f).hostsFile.file ] expected: %#v, actual: %#v", f, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ updateFile(f).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum([]byte(zo.lines[0] + "\n"))
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ updateFile(f).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if len(f.zones) != 1 {
                t.Errorf("[ updateFile(f).zones ] expected: %#v, actual: %#v", 1, len(f.zones))
            }

            // --------------------
     
            info, err := os.Stat(path)
            if err != nil {
                t.Errorf("[ updateFile() ] cannot stat written-file")
            }
            fLastModified := info.ModTime()

            if fLastModified != fileLastModified {
                t.Errorf("[ f.lastModified ] expected: %#v, actual: %#v", fileLastModified, fLastModified)
            }
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        data := []byte("# some data")

        path := "_test-hosts"
        err := os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ updateFile() ] cannot make test-directory")
        }

        f := new(File)
        f.Path = path
        f.Notes = "..."
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = int(f.id)
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "# some updated data")
        f.zones = append(f.zones, zo)

        // --------------------

        fValues := new(File)
        fValues.Notes = "...updated notes"

        err = updateFile(f, fValues)

        // --------------------

        if err == nil {
            t.Errorf("[ updateFile(f).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        if f == nil {
            t.Errorf("[ updateFile(f) ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if f.Path != path {
                t.Errorf("[ updateFile(f).Path ] expected: %#v, actual: %#v", path, f.Path)
            }

            // --------------------

            if f.Notes != "..." {
                t.Errorf("[ updateFile(f).Notes ] expected: %#v, actual: %#v", "...", f.Notes)
            }

            // --------------------

            if f.hostsFile == nil {
                t.Errorf("[ updateFile(f).hostsFile ] expected: not %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ updateFile(f).hostsFile.file ] expected: %#v, actual: %#v", f, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ updateFile(f).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum(data)
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ updateFile(f).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if len(f.zones) != 1 {
                t.Errorf("[ updateFile(f).zones ] expected: %#v, actual: %#v", 1, len(f.zones))
            }
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_fDelete(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 0

        // --------------------

        err := f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'f.ID'") {
            t.Errorf("[ f.Delete().err.Error() ] expected: contains %#v, actual: %#v", "missing 'f.ID'", err.Error())
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.ID = 42

        // --------------------

        err := f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "not found") {
            t.Errorf("[ f.Delete().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
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

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        // --------------------

        err = f.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ f.Delete().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_doesnt-exist.txt"
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        err := f.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ f.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
    })
}

func Test_deleteFile(t *testing.T) {
    var test string

    test = "deleted/no-zones"
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
        fid := f.ID

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        f.zones = make([]*zoneObject, 0)

        // --------------------

        err = deleteFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ deleteFile(f).err ] expected: %#v, actual: %#v", nil, err)
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

        if f.hostsFile != nil {
            t.Errorf("[ f.hostsFile ] expected: %#v, actual: %#v", nil, f.hostsFile)
        }

        // --------------------

        if f.zones != nil {
            t.Errorf("[ f.zones ] expected: %#v, actual: %#v", nil, f.zones)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = fid
        f = lookupFile(fQuery)
        if f != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = path
        f = lookupFile(fQuery)
        if f != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        _, err = ioutil.ReadFile(path)
        if err == nil {
                t.Errorf("[ ioutil.ReadFile(f.Path).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "deleted/existing-zones"
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
        fid := f.ID

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        f.zones = append(f.zones, new(zoneObject))

        // --------------------

        err = f.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ deleteFile(f).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = fid
        f = lookupFile(fQuery)
        if f != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = path
        f = lookupFile(fQuery)
        if f != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, f)
        }

        // --------------------

        readData, err := ioutil.ReadFile(path)
        if err != nil {
            t.Errorf("[ ioutil.ReadFile(fValues.Path).err ] expected: %#v, actual: %#v", nil, err)
        } else if len(readData) != len(data) {
            t.Errorf("[ len(ioutil.ReadFile(fValues.Path)).data ] expected: %#v, actual: %#v", len(data), len(readData))
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        data := []byte("# some data")
        path := "_doesnt-exist.txt"

        f := new(File)
        f.Path = path
        f.Notes = "..."
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])
        addFileObject(f.hostsFile)

        f.zones = make([]*zoneObject, 0)

        // --------------------

        err := deleteFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ deleteFile(f).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        if f == nil {
            t.Errorf("[ updateFile(f) ] expected: not %#v, actual: %#v", nil, f)
        } else {

            // --------------------

            if f.Path != path {
                t.Errorf("[ updateFile(f).Path ] expected: %#v, actual: %#v", path, f.Path)
            }

            // --------------------

            if f.Notes != "..." {
                t.Errorf("[ updateFile(f).Notes ] expected: %#v, actual: %#v", "...", f.Notes)
            }

            // --------------------

            if f.hostsFile == nil {
                t.Errorf("[ updateFile(f).hostsFile ] expected: not %#v, actual: %#v", nil, f.hostsFile)
            } else {

                // --------------------

                if f.hostsFile.file != f {
                    t.Errorf("[ updateFile(f).hostsFile.file ] expected: %#v, actual: %#v", f, f.hostsFile.file)
                }

                // --------------------

                if f.hostsFile.data != nil {
                    t.Errorf("[ updateFile(f).hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
                }

                // --------------------

                checksum := sha1.Sum(data)
                expected := hex.EncodeToString(checksum[:])
                if f.hostsFile.checksum != expected {
                    t.Errorf("[ updateFile(f).hostsFile.checksum ] expected: %#v, actual: %#v", expected, f.hostsFile.checksum)
                }
            }

            // --------------------

            if f.zones == nil {
                t.Errorf("[ updateFile(f).zones ] expected: not %#v, actual: %#v", nil, f.zones)
            }
        }

        // --------------------

        fo.file = nil   // !!! avoid memory leaks
    })
}

// -----------------------------------------------------------------------------

func Test_goRenderFile(t *testing.T) {
    var test string

    test = "rendered/only-external"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_test-hosts.txt"
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        z := new(Zone)
        z.File = f.ID
        z.Name = "external"
        addZone(z)

        zo := new(zoneObject)
        z.fileZone = zo   // !!! beware of memory leaks
        zo.zone = z       // !!! beware of memory leaks
        zo.lines = append(zo.lines, "")
        zo.lines = append(zo.lines, "# some data")
        zo.lines = append(zo.lines, "")
        f.zones = append(f.zones, zo)

        expectedData := []byte(`
# some data

`)

        // --------------------

        done := goRenderFile(f)
        _ = <-done

        // --------------------

        hash := sha1.Sum(expectedData)
        expectedChecksum := hex.EncodeToString(hash[:])

        hash = sha1.Sum(f.hostsFile.data)
        actualChecksum := hex.EncodeToString(hash[:])
        if actualChecksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.data ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        if f.hostsFile.checksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.checksum ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        zo.zone = nil   // !!! avoid memory leaks
        fo.file = nil   // !!! avoid memory leaks
    })

    test = "rendered/1-managed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_test-hosts.txt"
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        z1 := new(Zone)
        z1.File = f.ID
        z1.Name = "external"
        addZone(z1)

        zo1 := new(zoneObject)
        z1.fileZone = zo1   // !!! beware of memory leaks
        zo1.zone = z1       // !!! beware of memory leaks
        zo1.lines = append(zo1.lines, "")
        zo1.lines = append(zo1.lines, "# some data")
        zo1.lines = append(zo1.lines, "")
        f.zones = append(f.zones, zo1)

        z2 := new(Zone)
        z2.File = f.ID
        z2.Name = "my-zone-1"
        addZone(z2)

        zo2 := new(zoneObject)
        z2.fileZone = zo2   // !!! beware of memory leaks
        zo2.zone = z2       // !!! beware of memory leaks
        zo2.lines = append(zo2.lines, "##### Start Of Terraform Zone: my-zone-1 #######################################")
        zo2.lines = append(zo2.lines, "")
        zo2.lines = append(zo2.lines, "# some data")
        zo2.lines = append(zo2.lines, "")
        zo2.lines = append(zo2.lines, "##### End Of Terraform Zone: my-zone-1 #########################################")
        f.zones = append(f.zones, zo2)

        expectedData := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        // --------------------

        done := goRenderFile(f)
        _ = <-done

        // --------------------

        hash := sha1.Sum(expectedData)
        expectedChecksum := hex.EncodeToString(hash[:])

        hash = sha1.Sum(f.hostsFile.data)
        actualChecksum := hex.EncodeToString(hash[:])
        if actualChecksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.data ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        if f.hostsFile.checksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.checksum ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        zo1.zone = nil   // !!! avoid memory leaks
        zo2.zone = nil   // !!! avoid memory leaks
        fo.file = nil    // !!! avoid memory leaks
    })

    test = "rendered/more-managed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        f.Path = "_test-hosts.txt"
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        z1 := new(Zone)
        z1.File = f.ID
        z1.Name = "external"
        addZone(z1)

        zo1 := new(zoneObject)
        z1.fileZone = zo1   // !!! beware of memory leaks
        zo1.zone = z1       // !!! beware of memory leaks
        zo1.lines = append(zo1.lines, "")
        zo1.lines = append(zo1.lines, "# some data")
        zo1.lines = append(zo1.lines, "")
        f.zones = append(f.zones, zo1)

        z2 := new(Zone)
        z2.File = f.ID
        z2.Name = "my-zone-1"
        addZone(z2)

        zo2 := new(zoneObject)
        z2.fileZone = zo2   // !!! beware of memory leaks
        zo2.zone = z2       // !!! beware of memory leaks
        zo2.lines = append(zo2.lines, "##### Start Of Terraform Zone: my-zone-1 #######################################")
        zo2.lines = append(zo2.lines, "")
        zo2.lines = append(zo2.lines, "# some data")
        zo2.lines = append(zo2.lines, "")
        zo2.lines = append(zo2.lines, "##### End Of Terraform Zone: my-zone-1 #########################################")
        f.zones = append(f.zones, zo2)

        z3 := new(Zone)
        z3.File = f.ID
        z3.Name = "my-zone-2"
        addZone(z3)

        zo3 := new(zoneObject)
        z3.fileZone = zo3   // !!! beware of memory leaks
        zo3.zone = z3       // !!! beware of memory leaks
        zo2.lines = append(zo2.lines, "##### Start Of Terraform Zone: my-zone-2 #######################################")
        zo3.lines = append(zo3.lines, "")
        zo3.lines = append(zo3.lines, "# some data")
        zo3.lines = append(zo3.lines, "")
        zo3.lines = append(zo3.lines, "##### End Of Terraform Zone: my-zone-2 #########################################")
        f.zones = append(f.zones, zo3)

        expectedData := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
##### Start Of Terraform Zone: my-zone-2 #######################################

# some data

##### End Of Terraform Zone: my-zone-2 #########################################
`)

        // --------------------

        done := goRenderFile(f)
        _ = <-done

        // --------------------

        hash := sha1.Sum(expectedData)
        expectedChecksum := hex.EncodeToString(hash[:])

        hash = sha1.Sum(f.hostsFile.data)
        actualChecksum := hex.EncodeToString(hash[:])
        if actualChecksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.data ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        if f.hostsFile.checksum != expectedChecksum {
            t.Errorf("[ goRenderFile() > f.hostsFile.checksum ] expected: %#v, actual: %#v", expectedData, f.hostsFile.data)
        }

        // --------------------

        zo1.zone = nil   // !!! avoid memory leaks
        zo2.zone = nil   // !!! avoid memory leaks
        zo3.zone = nil   // !!! avoid memory leaks
        fo.file = nil    // !!! avoid memory leaks
    })
}

// -----------------------------------------------------------------------------

func Test_goScanFile(t *testing.T) {
    var test string

    test = "scanned/only-external"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 1 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 1, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "scanned/1-managed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################

# some final data

`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 2 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 2, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 6 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 6, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-1) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-1).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "scanned/more-managed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################

# some other data

##### Start Of Terraform Zone: my-zone-2 #######################################

# some data

##### End Of Terraform Zone: my-zone-2 #########################################

# some final data

`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 3 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 3, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 9 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 9, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-1) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-1).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-2"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-2) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-2).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "unexpected-start-zone-marker"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### Start Of Terraform Zone: my-zone-2 #######################################

# some data

##### End Of Terraform Zone: my-zone-2 #########################################
`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 3 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 3, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-1) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-1).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-2"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-2) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-2).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "unexpected-end-zone-marker"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

##### End Of Terraform Zone: my-zone-1 #########################################

# some final data

`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 1 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 1, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 6 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 6, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "missing-end-zone-marker"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        // --------------------

        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 2 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 2, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-1) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(my-zone-1).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })

    test = "cleanup-deleted-zones"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        path := "_test-hosts.txt"
        data1 := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        f := new(File)
        f.Path = path
        addFile(f)

        fo := new(fileObject)
        f.hostsFile = fo   // !!! beware of memory leaks
        fo.file = f        // !!! beware of memory leaks
        addFileObject(f.hostsFile)

        done := goScanFile(f.hostsFile, bytes.NewReader(data1))
        _ = <-done

        // --------------------

        data2 := []byte(`
# some other data

`)

        done = goScanFile(f.hostsFile, bytes.NewReader(data2))
        _ = <-done

        // --------------------

        if f.hostsFile.data != nil {
            t.Errorf("[ goScanFile() > f.hostsFile.data ] expected: %#v, actual: %#v", nil, f.hostsFile.data)
        }

        if len(f.zones) != 1 {
            t.Errorf("[ goScanFile() > f.zones ] expected: %#v, actual: %#v", 1, len(f.zones))
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        if z == nil {
            t.Errorf("[ goScanFile() > lookupZone(external) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ goScanFile() > lookupZone(external).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }

        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)

        if z != nil {
            t.Errorf("[ goScanFile() > lookupZone(my-zone-1) ] expected: %#v, actual: %#v", nil, z)
            if z.fileZone != nil {
                log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanFile()] z.Name: %q", z.Name)
                log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanFile()] z.fileZone.lines:")
                for i, line := range z.fileZone.lines {
                    log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanFile()] %d: %q", i, line)
                }
            }
        }

        // --------------------

        fo.file = nil    // !!! avoid memory leaks
    })
}

//------------------------------------------------------------------------------

func Test_addZoneObject(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        // --------------------

        f := new(File)
        z := new(zoneObject)
        addZoneObject(f, z)

        // --------------------

        length := len(f.zones)
        if length != 1 {
            t.Errorf("[ len(f.zones) ] expected: %#v, actual: %#v", 1, length)
        }
    })
}

func Test_removeZoneObject(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        z := new(zoneObject)

        removeZoneObject(f, z)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "removed"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        f := new(File)
        z := new(zoneObject)
        addZoneObject(f, z)

        // --------------------

        removeZoneObject(f, z)

        // --------------------

        length := len(f.zones)
        if length != 0 {
            t.Errorf("[ len(f.zones) ] expected: %#v, actual: %#v", 0, length)
        }
    })
}

func Test_deleteFromSliceOfZoneObjects(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        s1 := make([]*zoneObject, 0)

        z := new(zoneObject)

        // --------------------

        _ = deleteFromSliceOfZoneObjects(s1, z)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        s1 := make([]*zoneObject, 0)

        z := new(zoneObject)
        s1 = append(s1, z)

        // --------------------

        s2 := deleteFromSliceOfZoneObjects(s1, z)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }
    })

    test = "more-elements/first-element"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        s1 := make([]*zoneObject, 0)

        z1 := new(zoneObject)
        s1 = append(s1, z1)

        z2 := new(zoneObject)
        s1 = append(s1, z2)

        z3 := new(zoneObject)
        s1 = append(s1, z3)

        // --------------------

        s2 := deleteFromSliceOfZoneObjects(s1, z1)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == z1 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })

    test = "more-elements/middle-element"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        s1 := make([]*zoneObject, 0)

        z1 := new(zoneObject)
        s1 = append(s1, z1)

        z2 := new(zoneObject)
        s1 = append(s1, z2)

        z3 := new(zoneObject)
        s1 = append(s1, z3)

        // --------------------

        s2 := deleteFromSliceOfZoneObjects(s1, z2)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == z2 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })
    
    test = "more-elements/last-element"
    t.Run(test, func(t *testing.T) {

        resetFileTestEnv()

        s1 := make([]*zoneObject, 0)

        z1 := new(zoneObject)
        s1 = append(s1, z1)

        z2 := new(zoneObject)
        s1 = append(s1, z2)

        z3 := new(zoneObject)
        s1 = append(s1, z3)

        // --------------------

        s2 := deleteFromSliceOfZoneObjects(s1, z3)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == z3 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })
}
