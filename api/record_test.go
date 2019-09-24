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
    "os"
    "strings"
    "testing"
)

// -----------------------------------------------------------------------------

func resetRecordTestEnv() {
    if hosts != nil {
        for _, hostsFile := range hosts.files {   // !!! avoid memory leaks
            hostsFile.file = nil
        }
        hosts = (*anchor)(nil)
    }
    Init()
}

// -----------------------------------------------------------------------------

func Test_LookupRecord(t *testing.T) {
    var test string

    test = "found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Zone = 42
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        r.Comment = " ccc"
        r.Notes = "..."
        r.managed = true
        r.zoneRecord = new(recordObject)
        addRecord(r)

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)

        record := LookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ LookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else {

            // --------------------

            if record.ID != int(r.id) {
               t.Errorf("[ LookupRecord(rQuery).ID ] expected: %#v, actual: %#v", int(r.id), record.ID)
            }

            // --------------------

            if record.Address != r.Address {
                t.Errorf("[ LookupRecord(rQuery).Address ] expected: %#v, actual: %#v", r.Address, record.Address)
            }

            // --------------------

            if len(record.Names) != len(r.Names) {
                t.Errorf("[ LookupRecord(rQuery).Names ] expected: %#v, actual: %#v", r.Names, record.Names)
            }

            // --------------------

            if record.Comment != r.Comment {
                t.Errorf("[ LookupRecord(rQuery).Comment ] expected: %#v, actual: %#v", r.Comment, record.Comment)
            }

            // --------------------

            if record.Notes != r.Notes {
                t.Errorf("[ LookupRecord(rQuery).Notes ] expected: %#v, actual: %#v", r.Notes, record.Notes)
            }

            // --------------------

            if record.id != 0 {
               t.Errorf("[ LookupRecord(rQuery).id ] expected: %#v, actual: %#v", 0, record.id)
            }

            // --------------------

            if record.managed != false {
                t.Errorf("[ LookupRecord(rQuery).managed ] expected: %#v, actual: %#v", false, record.managed)
            }

            // --------------------

            if record.zoneRecord != nil {
                t.Errorf("[ LookupRecord(rQuery).zoneRecord ] expected: %#v, actual: %#v", nil, record.zoneRecord)
            }
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = 42

        record := LookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ LookupRecord(rQuery).err ] expected: %s, actual: %#v", "<error>", record)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_CreateRecord(t *testing.T) {
    var test string

    test = "missing-Zone"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'rValues.Zone'") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'rValues.Zone'", err.Error())
        }
    })

    test = "missing-Address"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 1
        rValues.Names = []string{ "n1", "n2", "n3" }

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'rValues.Address'") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'rValues.Address'", err.Error())
        }
    })

    test = "missing-Names"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 1
        rValues.Address = "a"

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'rValues.Names'") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "missing 'rValues.Names'", err.Error())
        }
    })

    test = "Zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 42
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'rValues.Zone' not found") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "'rValues.Zone' not found", err.Error())
        }
    })

    test = "cannot-create-external-records"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
# some data

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }

        zQuery := new(Zone)
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err = CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "cannot create records") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "cannot create records", err.Error())
        } else if !strings.Contains(err.Error(), "external") {
            t.Errorf("[ CreateRecord(rValues).err.Error() ] expected: contains %#v, actual: %#v", "external", err.Error())
        }

        // --------------------

        os.Remove(path)
    })

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        _ = CreateZone(zValues)
        z := lookupZone(zValues)

        r := new(Record)
        r.Zone = z.ID
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        addRecord(r)

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1" }

        err = CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(r).err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "already exists") {
            t.Errorf("[ CreateRecord(r).err.Error() ] expected: contains %#v, actual: %#v", "already exists", err.Error())
        }

        // --------------------

        os.Remove(path)
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        _ = CreateZone(zValues)
        z := lookupZone(zValues)

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err = CreateRecord(rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createRecord(rValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        _ = CreateZone(zValues)
        z := lookupZone(zValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot make test-directory")
        }

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err = CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_createRecord(t *testing.T) {
    var test string

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        _ = CreateZone(zValues)
        z := lookupZone(zValues)

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        rValues.managed = true

        err = createRecord(rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        r := lookupRecord(rValues)
        if r == nil {
            t.Errorf("[ lookupRecord(rValues) ] expected: not %#v, actual: %#v", nil, r)
        } else {

            // --------------------

            if r.id == 0 {
                t.Errorf("[ lookupRecord(rValues).id ] expected: not %#v, actual: %#v", 0, r.id)
            }

            // --------------------

            if r.ID != int(r.id) {
                t.Errorf("[ lookupRecord(rValues).ID ] expected: %#v, actual: %#v", r.id, r.ID)
            }

            // --------------------

            if r.Address != rValues.Address {
                t.Errorf("[ lookupRecord(rValues).Address ] expected: %#v, actual: %#v", rValues.Address, r.Address)
            }

            // --------------------

            if len(r.Names) != len(rValues.Names) {
                t.Errorf("[ lookupRecord(rValues).Names ] expected: %#v, actual: %#v", rValues.Names, r.Names)
            }

            // --------------------

            if r.Comment != rValues.Comment {
                t.Errorf("[ lookupRecord(rValues).Comment ] expected: %#v, actual: %#v", rValues.Comment, r.Comment)
            }

            // --------------------

            if r.Notes != rValues.Notes {
                t.Errorf("[ lookupRecord(rValues).Notes ] expected: %#v, actual: %#v", rValues.Notes, r.Notes)
            }

            // --------------------

            if r.managed != rValues.managed {
                t.Errorf("[ lookupRecord(rValues).managed ] expected: %#v, actual: %#v", rValues.managed, r.managed)
            }

            // --------------------

            if r.zoneRecord == nil {
                t.Errorf("[ lookupRecord(rValues).zoneRecord ] expected: not %#v, actual: %#v", nil, r.zoneRecord)
            } else {

                // --------------------

                if r.zoneRecord.record != r {
                   t.Errorf("[ lookupRecord(rValues).zoneRecord.record ] expected: not %#v, actual: %#v", r, r.zoneRecord.record)
                }

                // --------------------

                if len(r.zoneRecord.lines) != 1 {
                    t.Errorf("[ lookupRecord(rValues).zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(r.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(r.zoneRecord.lines[0]))
                    expected := hex.EncodeToString(checksum[:])
                    if r.zoneRecord.checksum != expected {
                        t.Errorf("[ lookupRecord(rValues).zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
                    }
                }
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        fValues := new(File)
        fValues.Path = path
        err := CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }
        f := lookupFile(fValues)

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = "my-zone-1"
        _ = CreateZone(zValues)
        z := lookupZone(zValues)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot make test-directory")
        }

        // --------------------

        rValues := new(Record)
        rValues.Zone = z.ID
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err = createRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        r := lookupRecord(rValues)
        if r != nil {
            t.Errorf("[ lookupRecord(rValues) ] expected: %#v, actual: %#v", nil, r)
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_rRead(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 0

        // --------------------

        _, err := r.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'r.ID'") {
            t.Errorf("[ r.Read().err.Error() ] expected: contains %#v, actual: %#v", "missing 'r.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 42

        // --------------------

        _, err := r.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.ID' not found") {
            t.Errorf("[ r.Read().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
    })

    test = "Zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Zone = 42
        addRecord(r)

        // --------------------

        _, err := r.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Read().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.Zone' not found") {
            t.Errorf("[ r.Read().err.Error() ] expected: contains %#v, actual: %#v", "'r.Zone' not found", err.Error())
        }
    })

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Read() ] cannot create test-file")
        }

        data = []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some other comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Read() ] cannot write test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        // --------------------

        record, err := r.Read()

        // --------------------

        if err != nil {
            t.Errorf("[ r.Read().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if record == nil {
            t.Errorf("[ r.Read().record ] expected: not %#v, actual: %#v", nil, record)
        } else {

            // --------------------

            if record.id != 0 {
                t.Errorf("[ r.Read().record.id ] expected: %#v, actual: %#v", 0, record.id)
            }

            // --------------------

            if record.ID != 1 {
                t.Errorf("[ r.Read().record.ID ] expected: %#v, actual: %#v", 1, record.ID)
            }

            // --------------------

            if record.Address != "1.1.1.1" {
                t.Errorf("[ r.Read().record.Address ] expected: %#v, actual: %#v", "1.1.1.1", record.Address)
            }

            // --------------------

            if len(record.Names) != 1 {
                t.Errorf("[ r.Read().record.Names ] expected: %#v, actual: %#v", 1, record.Names)
            }

            // --------------------

            if record.Comment != " some other comment" {
                t.Errorf("[ r.Read().record.Comment ] expected: %#v, actual: %#v", " some other comment", record.Comment)
            }

            // --------------------

            if record.Notes != r.Notes {
                t.Errorf("[ r.Read().record.Notes ] expected: %#v, actual: %#v", r.Notes, record.Notes)
            }

            // --------------------

            if record.managed != false {
                t.Errorf("[ r.Read().record.managed ] expected: %#v, actual: %#v", false, record.managed)
            }

            // --------------------

            if record.zoneRecord != nil {
                t.Errorf("[ r.Read().record.zoneRecord ] expected: %#v, actual: %#v", nil, record.zoneRecord)
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Read() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Read() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        os.Remove(path)

        // --------------------

        _, err = r.Read()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Read().err ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

func Test_readRecord(t *testing.T) {
    var test string

    test = "read"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot create test-file")
        }

        data = []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some other comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot write test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        expectedData := "1.1.1.1 my-host-1 # some other comment"

        // --------------------

        record, err := readRecord(r)

        // --------------------

        if err != nil {
            t.Errorf("[ readRecord(r).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if record == nil {
            t.Errorf("[ readRecord(r).record ] expected: not %#v, actual: %#v", nil, record)
        } else {

            // --------------------

            if record.id == 0 {
                t.Errorf("[ readRecord(r).record.id ] expected: not %#v, actual: %#v", 0, record.id)
            }

            // --------------------

            if record.ID != int(record.id) {
                t.Errorf("[ readRecord(r).record.ID ] expected: %#v, actual: %#v", record.id, record.ID)
            }

            // --------------------

            if record.Address != "1.1.1.1" {
                t.Errorf("[ readRecord(r).record.Address ] expected: %#v, actual: %#v", "1.1.1.1", record.Address)
            }

            // --------------------

            if len(record.Names) != 1 {
                t.Errorf("[ readRecord(r).record.Names ] expected: %#v, actual: %#v", 1, record.Names)
            }

            // --------------------

            if record.Comment != " some other comment" {
                t.Errorf("[ readRecord(r).record.Comment ] expected: %#v, actual: %#v", " some other comment", record.Comment)
            }

            // --------------------

            if record.Notes != r.Notes {
                t.Errorf("[ readRecord(r).record.Notes ] expected: %#v, actual: %#v", r.Notes, record.Notes)
            }

            // --------------------

            if record.managed != true {
                t.Errorf("[ readRecord(r).record.managed ] expected: %#v, actual: %#v", true, record.managed)
            }

            // --------------------

            if record.zoneRecord == nil {
                t.Errorf("[ readRecord(r).record.zoneRecord ] expected: not %#v, actual: %#v", nil, record.zoneRecord)
            } else {

                // --------------------

                if record.zoneRecord.record != r {
                   t.Errorf("[ readRecord(r).record.zoneRecord.record ] expected: not %#v, actual: %#v", r, record.zoneRecord.record)
                }

                // --------------------

                if len(record.zoneRecord.lines) != 1 {
                    t.Errorf("[ readRecord(r).record.zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(record.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(expectedData))
                    expected := hex.EncodeToString(checksum[:])
                    if record.zoneRecord.checksum != expected {
                        t.Errorf("[ readRecord(r).record.zoneRecord.checksum ] expected: %#v, actual: %#v", expected, record.zoneRecord.checksum)
                    }
                }
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-read"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        os.Remove(path)

        // --------------------

        _, err = readRecord(r)

        // --------------------

        if err == nil {
            t.Errorf("[ readRecord(r).err ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "record-deleted"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot create test-file")
        }

        data = []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err = ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ readRecord() ] cannot write test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        os.Remove(path)

        // --------------------

        record, err := readRecord(r)

        // --------------------

        if err == nil {
            t.Errorf("[ readRecord(r).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        if record != nil {
            t.Errorf("[ readRecord(r).record ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_rUpdate(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 0

        // --------------------

        rValues := new(Record)
        
        err := r.Update(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ r.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'r.ID'") {
            t.Errorf("[ r.Update().err.Error() ] expected: contains %#v, actual: %#v", "missing 'r.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 42

        // --------------------

        rValues := new(Record)
        
        err := r.Update(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ r.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.ID' not found") {
            t.Errorf("[ r.Update().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
    })

    test = "Zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Zone = 42
        addRecord(r)

        // --------------------

        rValues := new(Record)

        err := r.Update(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ r.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.Zone' not found") {
            t.Errorf("[ r.Update().err.Error() ] expected: contains %#v, actual: %#v", "'r.Zone' not found", err.Error())
        }
    })

    test = "cannot-update-external-records"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
1.1.1.1 my-host-1

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some comment"

        err = r.Update(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ r.Update().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "cannot update 'r.Comment'") {
            t.Errorf("[ r.Update().err.Error() ] expected: contains %#v, actual: %#v", "cannot update 'r.Comment'", err.Error())
        } else if !strings.Contains(err.Error(), "external") {
            t.Errorf("[ r.Update().err.Error() ] expected: contains %#v, actual: %#v", "external", err.Error())
        }

        // --------------------

        os.Remove(path)
    })

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Update() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Update() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some updated comment"

        err = r.Update(rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ r.Update().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-update"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Update() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Update() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ r.Update() ] cannot make test-directory")
        }

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some updated comment"

        err = r.Update(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ r.Update().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}


func Test_updateRecord(t *testing.T) {
    var test string

    test = "updated"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        expectedData := "1.1.1.1 my-host-1 # some updated comment"

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some updated comment"
        rValues.Notes = "...updated notes"

        err = updateRecord(r, rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateRecord(r).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        r = lookupRecord(rQuery)
        if r == nil {
            t.Errorf("[ lookupRecord(rValues) ] expected: not %#v, actual: %#v", nil, r)
        } else {

            // --------------------

            if r.id == 0 {
                t.Errorf("[ lookupRecord(rValues).id ] expected: not %#v, actual: %#v", 0, r.id)
            }

            // --------------------

            if r.ID != int(r.id) {
                t.Errorf("[ lookupRecord(rValues).ID ] expected: %#v, actual: %#v", r.id, r.ID)
            }

            // --------------------

            if r.Address != "1.1.1.1" {
                t.Errorf("[ lookupRecord(rValues).Address ] expected: %#v, actual: %#v", "1.1.1.1", r.Address)
            }

            // --------------------

            if len(r.Names) != 1 {
                t.Errorf("[ lookupRecord(rValues).Names ] expected: %#v, actual: %#v", 1, r.Names)
            }

            // --------------------

            if r.Comment != rValues.Comment {
                t.Errorf("[ lookupRecord(rValues).Comment ] expected: %#v, actual: %#v", rValues.Comment, r.Comment)
            }

            // --------------------

            if r.Notes != rValues.Notes {
                t.Errorf("[ lookupRecord(rValues).Notes ] expected: %#v, actual: %#v", rValues.Notes, r.Notes)
            }

            // --------------------

            if r.managed != true {
                t.Errorf("[ lookupRecord(rValues).managed ] expected: %#v, actual: %#v", true, r.managed)
            }

            // --------------------

            if r.zoneRecord == nil {
                t.Errorf("[ lookupRecord(rValues).zoneRecord ] expected: not %#v, actual: %#v", nil, r.zoneRecord)
            } else {

                // --------------------

                if r.zoneRecord.record != r {
                   t.Errorf("[ lookupRecord(rValues).zoneRecord.record ] expected: not %#v, actual: %#v", r, r.zoneRecord.record)
                }

                // --------------------

                if len(r.zoneRecord.lines) != 1 {
                    t.Errorf("[ lookupRecord(rValues).zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(r.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(expectedData))
                    expected := hex.EncodeToString(checksum[:])
                    if r.zoneRecord.checksum != expected {
                        t.Errorf("[ lookupRecord(rValues).zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
                    }
                }
            }
        }

        // --------------------

        os.Remove(path)
    })

    test = "not-needed"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot create test-file")
        }
     
        info, err := os.Stat(path)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot stat test-file")
        }
        fileLastModified := info.ModTime()

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        expectedData := "1.1.1.1 my-host-1 # some comment"

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some comment"
        rValues.Notes = "..."

        err = updateRecord(r, rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ updateRecord(r).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        r = lookupRecord(rQuery)
        if r == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: not %#v, actual: %#v", nil, r)
        } else {

            // --------------------

            if r.id == 0 {
                t.Errorf("[ lookupRecord(rQuery).id ] expected: not %#v, actual: %#v", 0, r.id)
            }

            // --------------------

            if r.ID != int(r.id) {
                t.Errorf("[ lookupRecord(rQuery).ID ] expected: %#v, actual: %#v", r.id, r.ID)
            }

            // --------------------

            if r.Address != "1.1.1.1" {
                t.Errorf("[ lookupRecord(rQuery).Address ] expected: %#v, actual: %#v", "1.1.1.1", r.Address)
            }

            // --------------------

            if len(r.Names) != 1 {
                t.Errorf("[ lookupRecord(rQuery).Names ] expected: %#v, actual: %#v", 1, r.Names)
            }

            // --------------------

            if r.Comment != rValues.Comment {
                t.Errorf("[ lookupRecord(rQuery).Comment ] expected: %#v, actual: %#v", rValues.Comment, r.Comment)
            }

            // --------------------

            if r.Notes != rValues.Notes {
                t.Errorf("[ lookupRecord(rQuery).Notes ] expected: %#v, actual: %#v", rValues.Notes, r.Notes)
            }

            // --------------------

            if r.managed != true {
                t.Errorf("[ lookupRecord(rQuery).managed ] expected: %#v, actual: %#v", true, r.managed)
            }

            // --------------------

            if r.zoneRecord == nil {
                t.Errorf("[ lookupRecord(rQuery).zoneRecord ] expected: not %#v, actual: %#v", nil, r.zoneRecord)
            } else {

                // --------------------

                if r.zoneRecord.record != r {
                   t.Errorf("[ lookupRecord(rQuery).zoneRecord.record ] expected: not %#v, actual: %#v", r, r.zoneRecord.record)
                }

                // --------------------

                if len(r.zoneRecord.lines) != 1 {
                    t.Errorf("[ lookupRecord(rQuery).zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(r.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(expectedData))
                    expected := hex.EncodeToString(checksum[:])
                    if r.zoneRecord.checksum != expected {
                        t.Errorf("[ lookupRecord(rQuery).zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
                    }
                }
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

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ updateRecord() ] cannot make test-directory")
        }

        expectedData := "1.1.1.1 my-host-1 # some comment"

        // --------------------

        rValues := new(Record)
        rValues.Comment = " some updated comment"
        rValues.Notes = "...updated notes"

        err = updateRecord(r, rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ updateRecord(r).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        r = lookupRecord(rQuery)
        if r == nil {
            t.Errorf("[ lookupRecord(rValues) ] expected: not %#v, actual: %#v", nil, r)
        } else {

            // --------------------

            if r.id == 0 {
                t.Errorf("[ lookupRecord(rValues).id ] expected: not %#v, actual: %#v", 0, r.id)
            }

            // --------------------

            if r.ID != int(r.id) {
                t.Errorf("[ lookupRecord(rValues).ID ] expected: %#v, actual: %#v", r.id, r.ID)
            }

            // --------------------

            if r.Address != "1.1.1.1" {
                t.Errorf("[ lookupRecord(rValues).Address ] expected: %#v, actual: %#v", "1.1.1.1", r.Address)
            }

            // --------------------

            if len(r.Names) != 1 {
                t.Errorf("[ lookupRecord(rValues).Names ] expected: %#v, actual: %#v", 1, r.Names)
            }

            // --------------------

            if r.Comment != " some comment" {
                t.Errorf("[ lookupRecord(rValues).Comment ] expected: %#v, actual: %#v", " some comment", r.Comment)
            }

            // --------------------

            if r.Notes != "..." {
                t.Errorf("[ lookupRecord(rValues).Notes ] expected: %#v, actual: %#v", "...", r.Notes)
            }

            // --------------------

            if r.managed != true {
                t.Errorf("[ lookupRecord(rValues).managed ] expected: %#v, actual: %#v", true, r.managed)
            }

            // --------------------

            if r.zoneRecord == nil {
                t.Errorf("[ lookupRecord(rValues).zoneRecord ] expected: not %#v, actual: %#v", nil, r.zoneRecord)
            } else {

                // --------------------

                if r.zoneRecord.record != r {
                   t.Errorf("[ lookupRecord(rValues).zoneRecord.record ] expected: not %#v, actual: %#v", r, r.zoneRecord.record)
                }

                // --------------------

                if len(r.zoneRecord.lines) != 1 {
                    t.Errorf("[ lookupRecord(rValues).zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(r.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(expectedData))
                    expected := hex.EncodeToString(checksum[:])
                    if r.zoneRecord.checksum != expected {
                        t.Errorf("[ lookupRecord(rValues).zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
                    }
                }
            }
        }

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_rDelete(t *testing.T) {
    var test string

    test = "missing-ID"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 0

        // --------------------

        err := r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "missing 'r.ID'") {
            t.Errorf("[ r.Delete().err.Error() ] expected: contains %#v, actual: %#v", "missing 'r.ID'", err.Error())
        }
    })

    test = "ID-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 42

        // --------------------

        err := r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.ID' not found") {
            t.Errorf("[ r.Delete().err.Error() ] expected: contains %#v, actual: %#v", "not found", err.Error())
        }
    })

    test = "Zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Zone = 42
        addRecord(r)

        // --------------------

        err := r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "'r.Zone' not found") {
            t.Errorf("[ r.Delete().err.Error() ] expected: contains %#v, actual: %#v", "'r.Zone' not found", err.Error())
        }
    })

    test = "cannot-delete-external-records"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`
1.1.1.1 my-host-1

`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ CreateRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        // --------------------

        err = r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        } else if !strings.Contains(err.Error(), "cannot delete records") {
            t.Errorf("[ r.Delete().err.Error() ] expected: contains %#v, actual: %#v", "cannot delete records", err.Error())
        } else if !strings.Contains(err.Error(), "external") {
            t.Errorf("[ r.Delete().err.Error() ] expected: contains %#v, actual: %#v", "external", err.Error())
        }

        // --------------------

        os.Remove(path)
    })

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Delete() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Delete() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        // --------------------

        err = r.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ r.Delete().err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ r.Delete() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ r.Delete() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ r.Delete() ] cannot make test-directory")
        }

        // --------------------

        err = r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete().err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        os.Remove(path)
    })
}

func Test_deleteRecord(t *testing.T) {
    var test string

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ deleteRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        rid := r.ID
        zid := r.Zone

        // --------------------

        err = deleteRecord(r)

        // --------------------

        if err != nil {
            t.Errorf("[ deleteRecord(r).err ] expected: %#v, actual: %#v", nil, err)
        }

        // --------------------

        if r.id != 0 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 0, r.id)
        }

        // --------------------

        if r.ID != 0 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 0, r.ID)
        }

        // --------------------

        if r.Zone != 0 {
            t.Errorf("[ r.Zone ] expected: %#v, actual: %#v", 0, r.Zone)
        }

        // --------------------

        if r.Address != "" {
            t.Errorf("[ r.Address ] expected: %#v, actual: %#v", "", r.Address)
        }

        // --------------------

        if r.Names != nil {
            t.Errorf("[ r.Names ] expected: %#v, actual: %#v", nil, r.Names)
        }

        // --------------------

        if r.Notes != "" {
            t.Errorf("[ r.Notes ] expected: %#v, actual: %#v", "", r.Notes)
        }

        // --------------------

        if r.managed != false {
            t.Errorf("[ r.managed ] expected: %#v, actual: %#v", false, r.managed)
        }

        // --------------------

        if r.zoneRecord != nil {
            t.Errorf("[ r.zoneRecord ] expected: %#v, actual: %#v", nil, r.zoneRecord)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.ID = rid
        r = lookupRecord(rQuery)
        if r != nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, r)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Zone = zid
        r = lookupRecord(rQuery)
        if r != nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, r)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = "1.1.1.1"
        r = lookupRecord(rQuery)
        if r != nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, r)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = []string{ "my-host-1"}
        r = lookupRecord(rQuery)
        if r != nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, r)
        }

        // --------------------

        os.Remove(path)
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        path := "_test-hosts.txt"

        data := []byte(`##### Start Of Terraform Zone: my-zone-1 #######################################
1.1.1.1 my-host-1 # some comment
##### End Of Terraform Zone: my-zone-1 #########################################
`)
        err := ioutil.WriteFile(path, data, 0644)
        if err != nil {
            t.Errorf("[ deleteRecord() ] cannot write test-file")
        }

        fValues := new(File)
        fValues.Path = path
        err = CreateFile(fValues)
        if err != nil {
            t.Errorf("[ deleteRecord() ] cannot create test-file")
        }

        rQuery := new(Record)
        rQuery.Address = "1.1.1.1"
        r := lookupRecord(rQuery)
        r.Notes = "..."
        r.managed = true

        os.Remove(path)
        err = os.Mkdir(path, 0644)
        if err != nil {
            t.Errorf("[ deleteRecord() ] cannot make test-directory")
        }

        // --------------------

        err = deleteRecord(r)

        // --------------------

        if err == nil {
            t.Errorf("[ deleteRecord(r).err ] expected: %s, actual: %#v", "<error>", err)
        }

        // --------------------

        r = lookupRecord(rQuery)
        if r == nil {
            t.Errorf("[ lookupRecord(rValues) ] expected: not %#v, actual: %#v", nil, r)
        } else {

            // --------------------

            if r.id == 0 {
                t.Errorf("[ lookupRecord(rValues).id ] expected: not %#v, actual: %#v", 0, r.id)
            }

            // --------------------

            if r.ID != int(r.id) {
                t.Errorf("[ lookupRecord(rValues).ID ] expected: %#v, actual: %#v", r.id, r.ID)
            }

            // --------------------

            if r.Address != "1.1.1.1" {
                t.Errorf("[ lookupRecord(rValues).Address ] expected: %#v, actual: %#v", "1.1.1.1", r.Address)
            }

            // --------------------

            if len(r.Names) != 1 {
                t.Errorf("[ lookupRecord(rValues).Names ] expected: %#v, actual: %#v", 1, r.Names)
            }

            // --------------------

            if r.Comment != " some comment" {
                t.Errorf("[ lookupRecord(rValues).Comment ] expected: %#v, actual: %#v", " some comment", r.Comment)
            }

            // --------------------

            if r.Notes != "..." {
                t.Errorf("[ lookupRecord(rValues).Notes ] expected: %#v, actual: %#v", "...", r.Notes)
            }

            // --------------------

            if r.managed != true {
                t.Errorf("[ lookupRecord(rValues).managed ] expected: %#v, actual: %#v", true, r.managed)
            }

            // --------------------

            if r.zoneRecord == nil {
                t.Errorf("[ lookupRecord(rValues).zoneRecord ] expected: not %#v, actual: %#v", nil, r.zoneRecord)
            } else {

                // --------------------

                if r.zoneRecord.record != r {
                   t.Errorf("[ lookupRecord(rValues).zoneRecord.record ] expected: not %#v, actual: %#v", r, r.zoneRecord.record)
                }

                // --------------------

                if len(r.zoneRecord.lines) != 1 {
                    t.Errorf("[ lookupRecord(rValues).zoneRecord.lines ] expected: %#v, actual: %#v", 2, len(r.zoneRecord.lines))
                } else {

                    // --------------------

                    checksum := sha1.Sum([]byte(r.zoneRecord.lines[0]))
                    expected := hex.EncodeToString(checksum[:])
                    if r.zoneRecord.checksum != expected {
                        t.Errorf("[ lookupRecord(rValues).zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
                    }
                }
            }
        }

        // --------------------

        os.Remove(path)
    })
}

// -----------------------------------------------------------------------------

func Test_renderRecord(t *testing.T) {
    var test string

    test = "without-comment"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Address = "1.1.1.1"
        r.Names = []string{
            "my-host-1",
            "my-host-2",
            "my-host-3",
        }

        ro := new(recordObject)
        r.zoneRecord = ro

        expectedData := "1.1.1.1 my-host-1 my-host-2 my-host-3"

        // --------------------

        renderRecord(r)

        // --------------------

        if len(r.zoneRecord.lines) != 1 {
            t.Errorf("[ r.zoneRecord.lines ] expected: %#v, actual: %#v", 1, len(r.zoneRecord.lines))
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if r.zoneRecord.checksum != expected {
            t.Errorf("[ r.zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
        }
    })

    test = "with-comment"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.Address = "1.1.1.1"
        r.Names = []string{
            "my-host-1",
            "my-host-2",
            "my-host-3",
        }
        r.Comment = " some comment"

        ro := new(recordObject)
        r.zoneRecord = ro

        expectedData := "1.1.1.1 my-host-1 my-host-2 my-host-3 # some comment"

        // --------------------

        renderRecord(r)

        // --------------------

        if len(r.zoneRecord.lines) != 1 {
            t.Errorf("[ r.zoneRecord.lines ] expected: %#v, actual: %#v", 1, len(r.zoneRecord.lines))
        }

        // --------------------

        checksum := sha1.Sum([]byte(expectedData))
        expected := hex.EncodeToString(checksum[:])
        if r.zoneRecord.checksum != expected {
            t.Errorf("[ r.zoneRecord.checksum ] expected: %#v, actual: %#v", expected, r.zoneRecord.checksum)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_goScanRecord(t *testing.T) {
    var test string

    test = "scanned/no-lines"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record != nil {
            t.Errorf("[ z.records[0].record ] expected: %#v, actual: %#v", nil, z.records[0].record)
        }
    })

    test = "scanned/new-record"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := "1.1.1.1   my-host-1 my-host-2   # some comment"

        // --------------------

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record == nil {
            t.Errorf("[ z.records[0].record ] expected: not %#v, actual: %#v", nil, z.records[0].record)
        } else {

            // --------------------

            if z.records[0].record.Address != "1.1.1.1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "1.1.1.1", z.records[0].record.Address)
            }

            // --------------------

            if len(z.records[0].record.Names) != 2 {
                t.Errorf("[ z.records[0].record.Names ] expected: not %#v, actual: %#v", 2, z.records[0].record.Names)
            } else {

                // --------------------

                if z.records[0].record.Names[0] != "my-host-1" {
                    t.Errorf("[ z.records[0].record.Names[0] ] expected: not %#v, actual: %#v", "my-host-1", z.records[0].record.Names[0])
                }

                // --------------------

                if z.records[0].record.Names[1] != "my-host-2" {
                    t.Errorf("[ z.records[0].record.Names[1] ] expected: not %#v, actual: %#v", "my-host-2", z.records[0].record.Names[1])
                }
            }

            // --------------------

            if z.records[0].record.Comment != " some comment" {
                t.Errorf("[ z.records[0].record.Comment ] expected: not %#v, actual: %#v", " some comment", z.records[0].record.Comment)
            }

            // --------------------

            if z.records[0].record.zoneRecord.record != z.records[0].record {
                t.Errorf("[ z.records[0].record.zoneRecord.record ] expected: %#v, actual: %#v", z.records[0].record, z.records[0].record.zoneRecord.record)
            }
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })

    test = "scanned/updated-record"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := "1.1.1.1   my-host-1 my-host-2   # some comment"

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        l = "1.1.1.1   my-host-1 my-host-2   # some other comment"

        // --------------------

        z.records = make([]*recordObject, 0)
        ro = new(recordObject)
        addRecordObject(z, ro)

        lines = make(chan string)
        done  = goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record == nil {
            t.Errorf("[ z.records[0].record ] expected: not %#v, actual: %#v", nil, z.records[0].record)
        } else {

            // --------------------

            if z.records[0].record.Address != "1.1.1.1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "1.1.1.1", z.records[0].record.Address)
            }

            // --------------------

            if len(z.records[0].record.Names) != 2 {
                t.Errorf("[ z.records[0].record.Names ] expected: not %#v, actual: %#v", 2, z.records[0].record.Names)
            } else {

                // --------------------

                if z.records[0].record.Names[0] != "my-host-1" {
                    t.Errorf("[ z.records[0].record.Names[0] ] expected: not %#v, actual: %#v", "my-host-1", z.records[0].record.Names[0])
                }

                // --------------------

                if z.records[0].record.Names[1] != "my-host-2" {
                    t.Errorf("[ z.records[0].record.Names[1] ] expected: not %#v, actual: %#v", "my-host-2", z.records[0].record.Names[1])
                }
            }

            // --------------------

            if z.records[0].record.Comment != " some other comment" {
                t.Errorf("[ z.records[0].record.Comment ] expected: not %#v, actual: %#v", " some other comment", z.records[0].record.Comment)
            }

            // --------------------

            if z.records[0].record.zoneRecord.record != z.records[0].record {
                t.Errorf("[ z.records[0].record.zoneRecord.record ] expected: %#v, actual: %#v", z.records[0].record, z.records[0].record.zoneRecord.record)
            }
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })

    test = "scanned/rescan-no-change"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := "1.1.1.1   my-host-1 my-host-2   # some comment"

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        z.records = make([]*recordObject, 0)
        ro = new(recordObject)
        addRecordObject(z, ro)

        lines = make(chan string)
        done  = goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record == nil {
            t.Errorf("[ z.records[0].record ] expected: not %#v, actual: %#v", nil, z.records[0].record)
        } else {

            // --------------------

            if z.records[0].record.Address != "1.1.1.1" {
                t.Errorf("[ f.zones[0].zone.Name ] expected: not %#v, actual: %#v", "1.1.1.1", z.records[0].record.Address)
            }

            // --------------------

            if len(z.records[0].record.Names) != 2 {
                t.Errorf("[ z.records[0].record.Names ] expected: not %#v, actual: %#v", 2, z.records[0].record.Names)
            } else {

                // --------------------

                if z.records[0].record.Names[0] != "my-host-1" {
                    t.Errorf("[ z.records[0].record.Names[0] ] expected: not %#v, actual: %#v", "my-host-1", z.records[0].record.Names[0])
                }

                // --------------------

                if z.records[0].record.Names[1] != "my-host-2" {
                    t.Errorf("[ z.records[0].record.Names[1] ] expected: not %#v, actual: %#v", "my-host-2", z.records[0].record.Names[1])
                }
            }

            // --------------------

            if z.records[0].record.Comment != " some comment" {
                t.Errorf("[ z.records[0].record.Comment ] expected: not %#v, actual: %#v", " some comment", z.records[0].record.Comment)
            }

            // --------------------

            if z.records[0].record.zoneRecord.record != z.records[0].record {
                t.Errorf("[ z.records[0].record.zoneRecord.record ] expected: %#v, actual: %#v", z.records[0].record, z.records[0].record.zoneRecord.record)
            }
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })

    test = "scanned/blank-record"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := ""

        // --------------------

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record != nil {
            t.Errorf("[ z.records[0].record ] expected: %#v, actual: %#v", nil, z.records[0].record)
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })

    test = "scanned/record-without-information"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := "# some comment"

        // --------------------

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record != nil {
            t.Errorf("[ z.records[0].record ] expected: %#v, actual: %#v", nil, z.records[0].record)
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })

    test = "scanned/invalid-record"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        l := "1.1.1.1   # some comment"

        // --------------------

        z := new(Zone)
        ro := new(recordObject)
        addRecordObject(z, ro)

        lines := make(chan string)
        done  := goScanRecord(z, ro, lines)

        lines <- l

        close(lines)
        _ = <-done

        // --------------------

        if z.records[0].record != nil {
            t.Errorf("[ z.records[0].record ] expected: %#v, actual: %#v", nil, z.records[0].record)
        }

        // --------------------

        if len(z.records[0].lines) != 1 {
            t.Errorf("[ z.records[0].lines ] expected: %#v, actual: %#v", 1, z.records[0].lines)
        }

        // --------------------

        checksum := sha1.Sum([]byte(l))
        expected := hex.EncodeToString(checksum[:])
        if z.records[0].checksum != expected {
            t.Errorf("[ z.records[0].checksum ] expected: %#v, actual: %#v", expected, z.records[0].checksum)
        }
    })
}
