//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "testing"
)

// -----------------------------------------------------------------------------

func resetRecordTestEnv() {
    hosts = (*anchor)(nil)
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
        r.Managed = true
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

            if record.Managed != r.Managed {
                t.Errorf("[ LookupRecord(rQuery).Managed ] expected: %#v, actual: %#v", r.Managed, record.Managed)
            }

            // --------------------

            if record.id != 0 {
               t.Errorf("[ LookupRecord(rQuery).id ] expected: %#v, actual: %#v", 0, record.id)
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
            t.Errorf("[ LookupRecord(rQuery) ] expected: %s, actual: %#v", "<error>", record)
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
            t.Errorf("[ CreateRecord(rValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "missing-Address"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 42
        rValues.Names = []string{ "n1", "n2", "n3" }

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "missing-Names"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 42
        rValues.Address = "a"

        err := CreateRecord(rValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(rValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        r := new(Record)
        r.Zone = 42
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        addRecord(r)

        // --------------------

        err := CreateRecord(r)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateRecord(r) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Zone = 42
        rValues.Address = "a"
        rValues.Names = []string{ "n1", "n2", "n3" }

        err := CreateRecord(rValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createRecord(rValues) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        // --------------------

        rValues := new(Record)
        rValues.Address = "a"

//        err := CreateRecord(rValues)                                          // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

        // --------------------

//        if err == nil {
//            t.Errorf("[ CreateRecord(rValues) ] expected: %s, actual: %#v", "<error>", err)
//        }
    })
}

func Test_rDelete(t *testing.T) {
    var test string

    test = "no-ID"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        addRecord(r)
        r.ID = 0

        // --------------------

        err := r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        r.ID = 1

        // --------------------

        err := r.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ r.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        addRecord(r)

        // --------------------

        err := r.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ r.Delete() ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetRecordTestEnv()

        r := new(Record)
        addRecord(r)

        // --------------------

//        err := r.Delete()                                                     // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

        // --------------------

//        if err == nil {
//            t.Errorf("[ r.Delete() ] expected: %s, actual: %#v", "<error>", err)
//        }
    })
}
