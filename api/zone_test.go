//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "crypto/sha1"
    "encoding/hex"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "testing"
)

// -----------------------------------------------------------------------------

func resetZoneTestEnv() {
    if hosts != nil {
        for _, hostsFile := range hosts.files {   // !!! avoid memory leaks
            hostsFile.file = nil
        }
        hosts = (*anchor)(nil)
    }
    Init()
}

// -----------------------------------------------------------------------------

func Test_LookupZone(t *testing.T) {
    var test string

    test = "found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.File = 42
        z.Name = "z"
        z.Notes = "..."
        z.fileZone = new(zoneObject)
        z.records = append(z.records, new(recordObject))
        addZone(z)

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)

        zone := LookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ LookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else {

            // --------------------

            if zone.id != 0 {
               t.Errorf("[ LookupZone(zQuery).id ] expected: %#v, actual: %#v", 0, zone)
            }

            // --------------------

            if zone.ID != int(z.id) {
               t.Errorf("[ LookupZone(zQuery).ID ] expected: %#v, actual: %#v", int(z.id), zone)
            }

            // --------------------

            if zone.File != z.File {
                t.Errorf("[ LookupZone(zQuery).File ] expected: %#v, actual: %#v", z.File, zone.File)
            }

            // --------------------

            if zone.Name != z.Name {
                t.Errorf("[ LookupZone(zQuery).Name ] expected: %#v, actual: %#v", z.Name, zone.Name)
            }

            // --------------------

            if zone.Notes != z.Notes {
                t.Errorf("[ LookupZone(zQuery).Notes ] expected: %#v, actual: %#v", z.Notes, zone.Notes)
            }

            // --------------------

            if zone.id != 0 {
               t.Errorf("[ LookupZone(zQuery).id ] expected: %#v, actual: %#v", 0, zone)
            }

            // --------------------

            if zone.fileZone != nil {
                t.Errorf("[ LookupZone(zQuery).fileZone ] expected: %#v, actual: %#v", nil, zone.fileZone)
            }

            // --------------------

            if zone.records != nil {
                t.Errorf("[ LookupZone(zQuery).records ] expected: %#v, actual: %#v", nil, zone.records)
            }
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = 42

        zone := LookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ LookupZone(zQuery).err ] expected: %#v, actual: %#v", nil, zone)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_CreateZone(t *testing.T) {
    var test string

    test = "missing-File"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.Name = "z"

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'zValues.File'") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'zValues.File'", err.Error())
        }
    })

    test = "missing-Name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 1

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'zValues.Name'") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'zValues.File'", err.Error())
        }

    })

    test = "illegal-Name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 1
        zValues.Name = "external"

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues.err) ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "illegal") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "illegal", err.Error())
        } else if !strings.Contains(err.Error(), "zValues.Name") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "zValues.Name", err.Error())
        }
    })

    test = "File-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 42
        zValues.Name = "z"

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'zValues.File' not found") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "'zValues.File' not found", err.Error())
        }
    })

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateZone() ] cannot create test-file")
        }
        f := LookupFile(fValues)

        z := new(Zone)
        z.File = f.ID
        z.Name = "z"
        addZone(z)

        // --------------------

        err = CreateZone(z)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(z).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "already exists") {
            t.Errorf("[ CreateZone(z).err.Error() ] expected: contains %#v, actual: %#v", "already exists", err.Error())
        }

        // --------------------

        os.Remove(fValues.Path)
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateZone() ] cannot create test-file")
        }
        f := LookupFile(fValues)

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "z"
        zValues.Notes = "..."

        err = CreateZone(zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createZone(zValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(fValues.Path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateZone() ] cannot create test-file")
        }
        f := LookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ CreateZone() ] cannot make test-directory")
        }

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "z"

        err = CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_createZone(t *testing.T) {
    var test string

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        fValues := new(File)
        fValues.Path = "_test-hosts.txt"
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ createZone() ] cannot create test-file")
        }

        f := LookupFile(fValues)

        expectedData := []byte(`##### Start Of Terraform Zone: z ###############################################
##### End Of Terraform Zone: z #################################################
`)

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "z"
        zValues.Notes = "..."

        err = createZone(zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createZone(zValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        z := lookupZone(zValues)
        if z == nil {
            t.Errorf("[ lookupZone(zValues) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if z.id == 0 {
                t.Errorf("[ lookupZone(zValues).id ] expected: not %#v, actual: %#v", 0, z.id)
            }

            // --------------------

            if z.ID != int(z.id) {
                t.Errorf("[ lookupZone(zValues).ID ] expected: %#v, actual: %#v", z.id, z.ID)
            }

            // --------------------

            if z.File != zValues.File {
                t.Errorf("[ lookupZone(zValues).File ] expected: %#v, actual: %#v", zValues.File, z.File)
            }

            // --------------------

            if z.Name != zValues.Name {
                t.Errorf("[ lookupZone(zValues).Name ] expected: %#v, actual: %#v", zValues.Name, z.Name)
            }

            // --------------------

            if z.Notes != zValues.Notes {
                t.Errorf("[ lookupZone(zValues).Notes ] expected: %#v, actual: %#v", zValues.Notes, z.Notes)
            }

            // --------------------

            if z.fileZone == nil {
                t.Errorf("[ lookupZone(zQuery).fileZone ] expected: not %#v, actual: %#v", nil, z.fileZone)
            } else {

                // --------------------

                if z.fileZone.zone != z {
                    t.Errorf("[ lookupZone(zQuery).fileZone.zone ] expected: %#v, actual: %#v", z, z.fileZone.zone)
                }

                // --------------------

                if len(z.fileZone.lines) != 2 {
                    t.Errorf("[ lookupZone(zQuery).fileZone.lines ] expected: %#v, actual: %#v", 2, len(z.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(expectedData)
                expected := hex.EncodeToString(checksum[:])
                if z.fileZone.checksum != expected {
                    t.Errorf("[ lookupZone(zQuery).fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
                }
            }

            // --------------------

            if len(z.records) != 0 {
                t.Errorf("[ lookupZone(zValues).records ] expected: %#v, actual: %#v", 0, len(z.records))
            }
        }

        // --------------------

        os.Remove(fValues.Path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ createFile() ] cannot create test-file")
        }
        f := LookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ createFile() ] cannot make test-directory")
        }

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "z"

        err = createZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ createZone(zValues).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        z := lookupZone(zValues)
        if z != nil {
            t.Errorf("[ lookupZone(zValues) ] expected: %#v, actual: %#v", nil, z)
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_zRead(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 0

        // --------------------

        _, err := z.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'z.ID'") {
            t.Errorf("[ z.Read().err.Error() ] expected: contains %#v, actual: %#v", "missing 'z.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 42

        // --------------------

        _, err := z.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.ID' not found") {
            t.Errorf("[ z.Read().err.Error() ] expected: contains %#v, actual: %#v", "'z.ID' not found", err.Error())
        }
    })

    test = "File-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        z := new(Zone)
        z.File = 42
        addZone(z)

        // --------------------

        _, err := z.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.File' not found") {
            t.Errorf("[ z.Read().err.Error() ] expected: contains %#v, actual: %#v", "'z.File' not found", err.Error())
        }
    })

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        path := "_test-hosts.txt"

        data := []byte(`
# some data

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        data = []byte(`
# some other data

`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot write test-file")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)
        z.Notes = "..."

        // --------------------

        zone, err := z.Read()

        // --------------------

        if err != nil {
            t.Errorf("[ z.Read().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if zone == nil {
            t.Errorf("[ z.Read().zone ] expected: not %#v, actual: %#v", nil, zone)
        } else {

            // --------------------

            if zone.id != 0 {
                t.Errorf("[ z.Read().zone.id ] expected: %#v, actual: %#v", 0, zone.id)
            }

            // --------------------

            if zone.ID != int(z.id) {
                t.Errorf("[ z.Read().zone.ID ] expected: %#v, actual: %#v", z.id, zone.ID)
            }

            // --------------------

            if zone.File != zQuery.File {
                t.Errorf("[ z.Read().zone.File ] expected: %#v, actual: %#v", zQuery.File, zone.File)
            }

            // --------------------

            if zone.Name != zQuery.Name {
                t.Errorf("[ z.Read().zone.Name ] expected: %#v, actual: %#v", zQuery.Name, zone.Name)
            }

            // --------------------

            if zone.Notes != z.Notes {
                t.Errorf("[ z.Read().zone.Notes ] expected: %#v, actual: %#v", z.Notes, zone.Notes)
            }

            // --------------------

            if zone.fileZone != nil {
                t.Errorf("[ z.Read().zone.fileZone ] expected: %#v, actual: %#v", nil, zone.fileZone)
            }

            // --------------------

            if zone.records != nil {
                t.Errorf("[ z.Read().zone.records ] expected: %#v, actual: %#v", nil, zone.records)
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        // --------------------

        _, err = z.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Read().err ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

func Test_readZone(t *testing.T) {
    var test string

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        data = []byte(`
# some other data

`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readZone() ] cannot write test-file")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)
        z.Notes = "..."

        // --------------------

        zone, err := readZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ readZone(z).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if zone == nil {
            t.Errorf("[ readZone(z).zone ] expected: not %#v, actual: %#v", nil, zone)
        } else {

            // --------------------

            if zone.id != 1 {
                t.Errorf("[ readZone(z).zone.id ] expected: %#v, actual: %#v", 1, zone.id)
            }

            // --------------------

            if zone.ID != 1 {
                t.Errorf("[ readZone(z).zone.ID ] expected: %#v, actual: %#v", 1, zone.ID)
            }

            // --------------------

            if zone.File != zQuery.File {
                t.Errorf("[ readZone(z).zone.File ] expected: %#v, actual: %#v", zQuery.File, zone.File)
            }

            // --------------------

            if zone.Name != zQuery.Name {
                t.Errorf("[ readZone(z).zone.Name ] expected: %#v, actual: %#v", zQuery.Name, zone.Name)
            }

            // --------------------

            if zone.Notes != z.Notes {
                t.Errorf("[ readZone(z).zone.Notes ] expected: %#v, actual: %#v", z.Notes, zone.Notes)
            }

            // --------------------

            if zone.fileZone == nil {
                t.Errorf("[ readZone(z).zone.fileZone ] expected: not %#v, actual: %#v", nil, zone.fileZone)
            } else {

                // --------------------

                if zone.fileZone.zone != z {
                    t.Errorf("[ readZone(z).zone.fileZone.zone ] expected: %#v, actual: %#v", z, zone.fileZone.zone)
                }

                // --------------------

                if len(zone.fileZone.lines) != 3 {
                    t.Errorf("[ readZone(z).zone.fileZone.lines ] expected: %#v, actual: %#v", 3, len(zone.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(data)
                expected := hex.EncodeToString(checksum[:])
                if zone.fileZone.checksum != expected {
                    t.Errorf("[ readZone(z).zone.fileZone.checksum ] expected: %#v, actual: %#v", expected, zone.fileZone.checksum)
                }
            }

            // --------------------

            if len(zone.records) != 3 {
                t.Errorf("[ readZone(z).zone.records ] expected: %#v, actual: %#v", 3, len(zone.records))
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        // --------------------

        _, err = readZone(z)

        // --------------------

        if err == nil {
            t.Errorf("[ readZone(z).err ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "zone-deleted"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        data = []byte(`
# some data

`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readZone() ] cannot write test-file")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)

        // --------------------

        zone, err := readZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ readZone(z).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if zone != nil {
            t.Errorf("[ readZone(z).zone ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_zUpdate(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 0

        // --------------------

        zValues := new(Zone)

        err := z.Update(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ z.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'z.ID'") {
            t.Errorf("[ z.Update().err.Error() ] expected: contains %#v, actual: %#v", "missing 'z.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 42

        // --------------------

        zValues := new(Zone)

        err := z.Update(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ z.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.ID' not found") {
            t.Errorf("[ z.Update().err.Error() ] expected: contains %#v, actual: %#v", "'z.ID' not found", err.Error())
        }
    })

    test = "File-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        z := new(Zone)
        z.File = 42
        addZone(z)

        // --------------------

        zValues := new(Zone)

        err := z.Update(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ z.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.File' not found") {
            t.Errorf("[ z.Update().err.Error() ] expected: contains %#v, actual: %#v", "'z.File' not found", err.Error())
        }
    })

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Update() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Update() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        z.Notes = "..."

        z.records[1].lines[0] = "# some updated data"
        checksum := sha1.Sum([]byte(z.records[1].lines[0]))
        z.records[1].checksum = hex.EncodeToString(checksum[:])

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        zValues.Notes = "...updated notes"

        err = updateZone(z, zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ z.Update().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Update() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Update() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ z.Update() ] cannot make test-directory")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)

        z.records[1].lines[0] = "# some updated data"
        checksum := sha1.Sum([]byte(z.records[1].lines[0]))
        z.records[1].checksum = hex.EncodeToString(checksum[:])

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"

        err = z.Update(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ z.Update().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_updateZone(t *testing.T) {
    var test string

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        zid := z.id
        z.Notes = "..."

        z.records[1].lines[0] = "# some updated data"
        checksum := sha1.Sum([]byte(z.records[1].lines[0]))
        z.records[1].checksum = hex.EncodeToString(checksum[:])

        expectedData := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################

# some updated data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        zValues.Notes = "...updated notes"

        err = updateZone(z, zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateZone(z).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        z = lookupZone(zValues)
        if z == nil {
            t.Errorf("[ lookupZone(zValues) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if z.id != zid {
                t.Errorf("[ lookupZone(zValues).id ] expected: %#v, actual: %#v", zid, z.id)
            }

            // --------------------

            if z.ID != int(z.id) {
                t.Errorf("[ lookupZone(zValues).ID ] expected: %#v, actual: %#v", int(z.id), z.ID)
            }

            // --------------------

            if z.File != zValues.File {
                t.Errorf("[ lookupZone(zValues).File ] expected: %#v, actual: %#v", zValues.File, z.File)
            }

            // --------------------

            if z.Name != zValues.Name {
                t.Errorf("[ lookupZone(zValues).Name ] expected: %#v, actual: %#v", zValues.Name, z.Name)
            }

            // --------------------

            if z.Notes != zValues.Notes {
                t.Errorf("[ lookupZone(zValues).Notes ] expected: %#v, actual: %#v", zValues.Notes, z.Notes)
            }

            // --------------------

            if z.fileZone == nil {
                t.Errorf("[ lookupZone(zValues).fileZone ] expected: not %#v, actual: %#v", nil, z.fileZone)
            } else {

                // --------------------

                if z.fileZone.zone != z {
                    t.Errorf("[ lookupZone(zValues).fileZone.zone ] expected: %#v, actual: %#v", z, z.fileZone.zone)
                }

                // --------------------

                if len(z.fileZone.lines) != 5 {
                    t.Errorf("[ lookupZone(zValues).fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(expectedData)
                expected := hex.EncodeToString(checksum[:])
                if z.fileZone.checksum != expected {
                    t.Errorf("[ lookupZone(zValues).fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
                }
            }

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ lookupZone(zValues).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "not-needed"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)
     
        info, err := os.Stat(path)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot stat test-file")
        }
        fileLastModified := info.ModTime()

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        zid := z.id
        z.Notes = "..."

        expectedData := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        zValues.Notes = "...updated notes"

        err = updateZone(z, zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateZone(z).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        z = lookupZone(zValues)
        if z == nil {
            t.Errorf("[ lookupZone(zValues) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if z.id != zid {
                t.Errorf("[ lookupZone(zValues).id ] expected: %#v, actual: %#v", zid, z.id)
            }

            // --------------------

            if z.ID != int(z.id) {
                t.Errorf("[ lookupZone(zValues).ID ] expected: %#v, actual: %#v", int(z.id), z.ID)
            }

            // --------------------

            if z.File != zValues.File {
                t.Errorf("[ lookupZone(zValues).File ] expected: %#v, actual: %#v", zValues.File, z.File)
            }

            // --------------------

            if z.Name != zValues.Name {
                t.Errorf("[ lookupZone(zValues).Name ] expected: %#v, actual: %#v", zValues.Name, z.Name)
            }

            // --------------------

            if z.Notes != zValues.Notes {
                t.Errorf("[ lookupZone(zValues).Notes ] expected: %#v, actual: %#v", zValues.Notes, z.Notes)
            }

            // --------------------

            if z.fileZone == nil {
                t.Errorf("[ lookupZone(zValues).fileZone ] expected: not %#v, actual: %#v", nil, z.fileZone)
            } else {

                // --------------------

                if z.fileZone.zone != z {
                    t.Errorf("[ lookupZone(zValues).fileZone.zone ] expected: %#v, actual: %#v", z, z.fileZone.zone)
                }

                // --------------------

                if len(z.fileZone.lines) != 5 {
                    t.Errorf("[ lookupZone(zValues).fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(expectedData)
                expected := hex.EncodeToString(checksum[:])
                if z.fileZone.checksum != expected {
                    t.Errorf("[ lookupZone(zValues).fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
                }
            }

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ lookupZone(zValues).records ] expected: %#v, actual: %#v", 3, len(z.records))
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

        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ updateZone() ] cannot make test-directory")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        z.Notes = "..."

        z.records[1].lines[0] = "# some updated data"
        checksum := sha1.Sum([]byte(z.records[1].lines[0]))
        z.records[1].checksum = hex.EncodeToString(checksum[:])

        expectedData := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        // --------------------

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        zValues.Notes = "...updated notes"

        err = updateZone(z, zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ updateZone().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        z = lookupZone(zValues)
        if z == nil {
            t.Errorf("[ lookupZone(zValues) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if z.File != zValues.File {
                t.Errorf("[ lookupZone(zValues).File ] expected: %#v, actual: %#v", zValues.File, z.File)
            }

            // --------------------

            if z.Name != zValues.Name {
                t.Errorf("[ lookupZone(zValues).Name ] expected: %#v, actual: %#v", zValues.Name, z.Name)
            }

            // --------------------

            if z.Notes != "..." {
                t.Errorf("[ lookupZone(zValues).Notes ] expected: %#v, actual: %#v", "...", z.Notes)
            }

            // --------------------

            if z.fileZone == nil {
                t.Errorf("[ lookupZone(zValues).fileZone ] expected: not %#v, actual: %#v", nil, z.fileZone)
            } else {

                // --------------------

                if z.fileZone.zone != z {
                    t.Errorf("[ lookupZone(zValues).fileZone.zone ] expected: %#v, actual: %#v", z, z.fileZone.zone)
                }

                // --------------------

                if len(z.fileZone.lines) != 5 {
                    t.Errorf("[ lookupZone(zValues).fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(expectedData)
                expected := hex.EncodeToString(checksum[:])
                if z.fileZone.checksum != expected {
                    t.Errorf("[ lookupZone(zValues).fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
                }
            }

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ lookupZone(zValues).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_zDelete(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 0

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'z.ID'") {
            t.Errorf("[ z.Delete().err.Error() ] expected: contains %#v, actual: %#v", "missing 'z.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 42

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.ID' not found") {
            t.Errorf("[ z.Delete().err.Error() ] expected: contains %#v, actual: %#v", "'z.ID' not found", err.Error())
        }
    })

    test = "cannot-delete-external"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 1
        z.Name = "external"
        addZone(z)

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues.err) ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "cannot delete") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "cannot delete", err.Error())
        } else if !strings.Contains(err.Error(), "external") {
            t.Errorf("[ CreateZone(zValues).err.Error() ] expected: contains %#v, actual: %#v", "external", err.Error())
        }
    })

    test = "File-not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        z := new(Zone)
        z.File = 42
        addZone(z)

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'z.File' not found") {
            t.Errorf("[ z.Delete().err.Error() ] expected: contains %#v, actual: %#v", "'z.File' not found", err.Error())
        }
    })

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Delete() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Delete() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)

        // --------------------

        err = z.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ z.Delete().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })


    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ z.Delete() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ z.Delete() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ z.Delete() ] cannot make test-directory")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)

        // --------------------

        err = z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_deleteZone(t *testing.T) {
    var test string

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ deleteZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        zid := z.ID
        z.Notes = "..."

        // --------------------

        err = deleteZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ deleteZone(z).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if z.id != 0 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 0, z.id)
        }

        // --------------------

        if z.ID != 0 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 0, z.ID)
        }

        // --------------------

        if z.File != 0 {
            t.Errorf("[ z.File ] expected: %#v, actual: %#v", 0, z.File)
        }

        // --------------------

        if z.Name != "" {
            t.Errorf("[ z.Name ] expected: %#v, actual: %#v", "", z.Name)
        }

        // --------------------

        if z.Notes != "" {
            t.Errorf("[ z.Notes ] expected: %#v, actual: %#v", "", z.Notes)
        }

        // --------------------

        if z.fileZone != nil {
            t.Errorf("[ z.fileZone ] expected: %#v, actual: %#v", nil, z.fileZone)
        }

        // --------------------

        if z.records != nil {
            t.Errorf("[ z.records ] expected: %#v, actual: %#v", nil, z.records)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.ID = zid
        z = lookupZone(zQuery)
        if z != nil {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", nil, z)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = "my-zone-1"
        z = lookupZone(zQuery)
        if z != nil {
            t.Errorf("[ lookupZone(zQuery.Name) ] expected: %#v, actual: %#v", nil, z)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = f.ID
        z = lookupZone(zQuery)
        if z == nil {   // still having the external zone
            t.Errorf("[ lookupZone(zQuery.File) ] expected: not %#v, actual: %#v", nil, z)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteZone() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ deleteZone() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ deleteZone() ] cannot make test-directory")
        }

        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "my-zone-1"
        z := lookupZone(zQuery)
        z.Notes = "..."

        expectedData := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`)

        // --------------------

        err = deleteZone(z)

        // --------------------

        if err == nil {
            t.Errorf("[ deleteZone(z).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        z = lookupZone(zQuery)
        if z == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: not %#v, actual: %#v", nil, z)
        } else {

            // --------------------

            if z.File != zQuery.File {
                t.Errorf("[ lookupZone(zQuery).File ] expected: %#v, actual: %#v", zQuery.File, z.File)
            }

            // --------------------

            if z.Name != zQuery.Name {
                t.Errorf("[ lookupZone(zQuery).Name ] expected: %#v, actual: %#v", zQuery.Name, z.Name)
            }

            // --------------------

            if z.Notes != "..." {
                t.Errorf("[ lookupZone(zQuery).Notes ] expected: %#v, actual: %#v", "...", z.Notes)
            }

            // --------------------

            if z.fileZone == nil {
                t.Errorf("[ lookupZone(zQuery).fileZone ] expected: not %#v, actual: %#v", nil, z.fileZone)
            } else {

                // --------------------

                if z.fileZone.zone != z {
                    t.Errorf("[ lookupZone(zQuery).fileZone.zone ] expected: %#v, actual: %#v", z, z.fileZone.zone)
                }

                // --------------------

                if len(z.fileZone.lines) != 5 {
                    t.Errorf("[ lookupZone(zQuery).fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
                }

                // --------------------

                checksum := sha1.Sum(expectedData)
                expected := hex.EncodeToString(checksum[:])
                if z.fileZone.checksum != expected {
                    t.Errorf("[ lookupZone(zQuery).fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
                }
            }

            // --------------------

            if len(z.records) != 3 {
                t.Errorf("[ lookupZone(zQuery).records ] expected: %#v, actual: %#v", 3, len(z.records))
            }
        }

        // --------------------

        os.Remove(path)
    })
}
 
//------------------------------------------------------------------------------

func Test_renderZone(t *testing.T) {
    var test string

    test = "short-name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.Name = "short-name"

        zo := new(zoneObject)
        z.fileZone = zo

        r1 := new(recordObject)
        r1.lines = append(r1.lines, "")
        addRecordObject(z, r1)
        r2 := new(recordObject)
        r2.lines = append(r2.lines, "# some data")
        addRecordObject(z, r2)
        r3 := new(recordObject)
        r3.lines = append(r3.lines, "")
        addRecordObject(z, r3)

        expectedData := `##### Start Of Terraform Zone: short-name ######################################

# some data

##### End Of Terraform Zone: short-name ########################################
`

        // --------------------

        renderZone(z)

        // --------------------

        if len(z.fileZone.lines) != 5 {
            t.Errorf("[ z.fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if z.fileZone.checksum != expected {
            t.Errorf("[ z.fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
        }
    })

    test = "long-name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.Name = "very-very-very-very-very-very-very-very-very-very-long-name"

        zo := new(zoneObject)
        z.fileZone = zo

        r1 := new(recordObject)
        r1.lines = append(r1.lines, "")
        addRecordObject(z, r1)
        r2 := new(recordObject)
        r2.lines = append(r2.lines, "# some data")
        addRecordObject(z, r2)
        r3 := new(recordObject)
        r3.lines = append(r3.lines, "")
        addRecordObject(z, r3)

        expectedData := `##### Start Of Terraform Zone: very-very-very-very-very-very-very-very-very-very-long-name #####

# some data

##### End Of Terraform Zone: very-very-very-very-very-very-very-very-very-very-long-name #####
`

        // --------------------

        renderZone(z)

        // --------------------

        if len(z.fileZone.lines) != 5 {
            t.Errorf("[ z.fileZone.lines ] expected: %#v, actual: %#v", 5, len(z.fileZone.lines))
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if z.fileZone.checksum != expected {
            t.Errorf("[ z.fileZone.checksum ] expected: %#v, actual: %#v", expected, z.fileZone.checksum)
        }
    })
}
 
//------------------------------------------------------------------------------

func Test_goScanZone(t *testing.T) {
    var test string

    test = "scanned/no-lines"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone != nil {
            t.Errorf("[ f.zones[0].zone ] expected: %#v, actual: %#v", nil, f.zones[0].zone)
        }
    })

    test = "scanned/new-external-zone"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "",
            "# some data",
            "",
        }

        expectedData := `
# some data

`

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "external" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "scanned/updated-external-zone"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "",
            "# some data",
            "",
        }

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        ls = []string{
            "",
            "# some updated data",
            "",
        }

        expectedData := `
# some updated data

`

        // --------------------

        f.zones = make([]*zoneObject, 0)
        zo = new(zoneObject)
        addZoneObject(f, zo)

        lines = make(chan string)
        done  = goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "external" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "scanned/rescan-no-change"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "",
            "# some data",
            "",
        }

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        expectedData := `
# some data

`

        // --------------------

        f.zones = make([]*zoneObject, 0)
        zo = new(zoneObject)
        addZoneObject(f, zo)

        lines = make(chan string)
        done  = goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "external" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "scanned/new-managed-zone"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "##### Start Of Terraform Zone: my-zone-1 #######################################",
            "",
            "# some data",
            "",
            "##### End Of Terraform Zone: my-zone-1 #########################################",
        }

        expectedData := `##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "my-zone-1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "scanned/updated-managed-zone"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "##### Start Of Terraform Zone: my-zone-1 #######################################",
            "",
            "# some data",
            "",
            "##### End Of Terraform Zone: my-zone-1 #########################################",
        }

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        ls = []string{
            "##### Start Of Terraform Zone: my-zone-1 #######################################",
            "",
            "# some updated data",
            "",
            "##### End Of Terraform Zone: my-zone-1 #########################################",
        }


        expectedData := `##### Start Of Terraform Zone: my-zone-1 #######################################

# some updated data

##### End Of Terraform Zone: my-zone-1 #########################################
`

        // --------------------

        f.zones = make([]*zoneObject, 0)
        zo = new(zoneObject)
        addZoneObject(f, zo)

        lines = make(chan string)
        done  = goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "my-zone-1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "anonymous-end-zone-marker/short-name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "##### Start Of Terraform Zone: my-zone-1 #######################################",
            "",
            "# some data",
            "",
            "##### End Of Terraform Zone: ",
        }

        expectedData := `##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-1 #########################################
`

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "my-zone-1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "my-zone-1", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "anonymous-end-zone-marker/long-name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "##### Start Of Terraform Zone: very-very-very-very-very-very-very-very-very-very-my-zone-1 #####",
            "",
            "# some data",
            "",
            "##### End Of Terraform Zone: ",
        }

        expectedData := `##### Start Of Terraform Zone: very-very-very-very-very-very-very-very-very-very-my-zone-1 #####

# some data

##### End Of Terraform Zone: very-very-very-very-very-very-very-very-very-very-my-zone-1 #####
`

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "very-very-very-very-very-very-very-very-very-very-my-zone-1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "very-very-very-very-very-very-very-very-very-very-my-zone-1", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "unexpected-end-zone-marker"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "##### Start Of Terraform Zone: my-zone-1 #######################################",
            "",
            "# some data",
            "",
            "##### End Of Terraform Zone: my-zone-2 #########################################",
        }

        expectedData := `##### Start Of Terraform Zone: my-zone-1 #######################################

# some data

##### End Of Terraform Zone: my-zone-2 #########################################
`

        // --------------------

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "my-zone-1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 3 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 3, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if f.zones[0].checksum != expected {
            t.Errorf("[ f.zones[0].checksum ] expected: %#v, actual: %#v", expected, f.zones[0].checksum)
        }
    })

    test = "cleanup-deleted-records"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        ls := []string{
            "1.1.1.1 my-host-1",
            "2.2.2.2 my-host-2",
        }

        f := new(File)
        zo := new(zoneObject)
        addZoneObject(f, zo)

        lines := make(chan string)
        done  := goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        ls = []string{
            "1.1.1.1 my-host-1",
        }

        // --------------------

        f.zones = make([]*zoneObject, 0)
        zo = new(zoneObject)
        addZoneObject(f, zo)

        lines = make(chan string)
        done  = goScanZone(f, zo, lines)

        for _, l := range ls {
            lines <- l
        }

        close(lines)
        _ = <-done

        // --------------------

        if f.zones[0].zone == nil {
            t.Errorf("[ f.zones[0].zone ] expected: not %#v, actual: %#v", nil, f.zones[0].zone)
        } else {

            // --------------------

            if f.zones[0].zone.Name != "external" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "external", f.zones[0].zone.Name)
            }

            // --------------------

            if f.zones[0].zone.fileZone.zone != f.zones[0].zone {
                t.Errorf("[ f.zones[0].zone.fileZone.lines ] expected: %#v, actual: %#v", f.zones[0].zone, f.zones[0].zone.fileZone.zone)
            }

            // --------------------

            if len(f.zones[0].zone.records) != 1 {
                t.Errorf("[ f.zones[0].zone.records ] expected: %#v, actual: %#v", 1, len(f.zones[0].zone.records))
            }
        }

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        if r == nil {
            t.Errorf("[ goScanZone() > lookupRecord(1.1.1.1) ] expected: not %#v, actual: %#v", nil, r)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = "2.2.2.2"
        r = lookupRecord(rQuery)

        if r != nil {
            t.Errorf("[ goScanZone() > lookupRecord(2.2.2.2) ] expected: %#v, actual: %#v", nil, r)
            if r.zoneRecord != nil {
                log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanZone()] r.Address: %q", r.Address)
                log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanZone()] r.zoneRecord.lines:")
                for i, line := range r.zoneRecord.lines {
                    log.Printf("[DEBUG][terraform-provider-hosts/api/testing goScanZone()] %d: %q", i, line)
                }
            }
        }
    })
}
 
//------------------------------------------------------------------------------

func Test_addRecordObject(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        z := new(Zone)
        r := new(recordObject)
        addRecordObject(z, r)

        // --------------------

        length := len(z.records)
        if length != 1 {
            t.Errorf("[ len(z.records) ] expected: %#v, actual: %#v", 1, length)
        }
    })
}

func Test_removeRecordObject(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        r := new(recordObject)

        removeRecordObject(z, r)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "removed"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        r := new(recordObject)
        addRecordObject(z, r)

        // --------------------

        removeRecordObject(z, r)

        // --------------------

        length := len(z.records)
        if length != 0 {
            t.Errorf("[ len(f.zones) ] expected: %#v, actual: %#v", 0, length)
        }
    })
}

func Test_deleteFromSliceOfRecordObjects(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        s1 := make([]*recordObject, 0)

        r := new(recordObject)

        // --------------------

        _ = deleteFromSliceOfRecordObjects(s1, r)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        s1 := make([]*recordObject, 0)

        r := new(recordObject)
        s1 = append(s1, r)

        // --------------------

        s2 := deleteFromSliceOfRecordObjects(s1, r)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }
    })

    test = "more-elements/first-element"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        s1 := make([]*recordObject, 0)

        r1 := new(recordObject)
        s1 = append(s1, r1)

        r2 := new(recordObject)
        s1 = append(s1, r2)

        r3 := new(recordObject)
        s1 = append(s1, r3)

        // --------------------

        s2 := deleteFromSliceOfRecordObjects(s1, r1)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == r1 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })

    test = "more-elements/middle-element"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        s1 := make([]*recordObject, 0)

        r1 := new(recordObject)
        s1 = append(s1, r1)

        r2 := new(recordObject)
        s1 = append(s1, r2)

        r3 := new(recordObject)
        s1 = append(s1, r3)

        // --------------------

        s2 := deleteFromSliceOfRecordObjects(s1, r2)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == r2 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })
    
    test = "more-elements/last-element"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        s1 := make([]*recordObject, 0)

        r1 := new(recordObject)
        s1 = append(s1, r1)

        r2 := new(recordObject)
        s1 = append(s1, r2)

        r3 := new(recordObject)
        s1 = append(s1, r3)

        // --------------------

        s2 := deleteFromSliceOfRecordObjects(s1, r3)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for i, element := range s2 {
            if element == r3 {
                t.Errorf("[ for s2[element].i ] expected: %s, actual: %#v", "<not found>", i)
            }
        }
    })
}
