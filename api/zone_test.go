package api

import (
    "testing"
)

// -----------------------------------------------------------------------------

func resetZoneTestEnv() {
    hosts = (*anchor)(nil)
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
            t.Errorf("[ LookupZone(zQuery) ] expected: %s, actual: %#v", "<error>", zone)
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
            t.Errorf("[ CreateZone(zValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "missing-Name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 42

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "invalid-Name"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 42
        zValues.Name = "external"

        err := CreateZone(zValues)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(zValues) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "already-exists"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        z := new(Zone)
        z.File = 42
        z.Name = "z"
        addZone(z)

        // --------------------

        err := CreateZone(z)

        // --------------------

        if err == nil {
            t.Errorf("[ CreateZone(z) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "created"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 42
        zValues.Name = "z"

        err := CreateZone(zValues)

        // --------------------

        if err != nil {
            t.Errorf("[ createZone(zValues) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "cannot-create"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        // --------------------

        zValues := new(Zone)
        zValues.File = 42
        zValues.Name = "z"

//        err = CreateZone(zValues)                                             // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

        // --------------------

//        if err == nil {
//            t.Errorf("[ CreateZone(zValues) ] expected: %s, actual: %#v", "<error>", err)
//        }
    })
}

func Test_zDelete(t *testing.T) {
    var test string

    test = "no-ID"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        addZone(z)
        z.ID = 0

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "not-found"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        z.ID = 1

        // --------------------

        err := z.Delete()

        // --------------------

        if err == nil {
            t.Errorf("[ z.Delete() ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "deleted"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        addZone(z)

        // --------------------

        err := z.Delete()

        // --------------------

        if err != nil {
            t.Errorf("[ z.Delete() ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "cannot-delete"
    t.Run(test, func(t *testing.T) {

        resetZoneTestEnv()

        z := new(Zone)
        addZone(z)

        // --------------------

//        err := z.Delete()                                                     // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

        // --------------------

//        if err == nil {
//            t.Errorf("[ z.Delete() ] expected: %s, actual: %#v", "<error>", err)
//        }
    })
}
