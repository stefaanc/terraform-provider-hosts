//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "crypto/sha1"
    "errors"
    "encoding/hex"
    "io"
)

// -----------------------------------------------------------------------------

type Zone struct {
    // readOnly
    ID       int      // indexed   // read-write in a zQuery
    // read-writeOnce
    File     int      // indexed   // read-write in a zQuery
    Name     string   // indexed   // read-write in a zQuery
    // read-writeMany
    Notes    string
    // private
    id       zoneID
    lines    []string
    checksum string
}

func LookupZone(zQuery *Zone) (z *Zone) {
    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return nil
    }

    // make a copy without the private fields
    z = new(Zone)
    z.ID    = zPrivate.ID
    z.File  = zPrivate.File
    z.Name  = zPrivate.Name
    z.Notes = zPrivate.Notes
    // ignore computed fields

    return z
}

func CreateZone(zValues *Zone) error {
    if zValues.File == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] missing 'zValues.File'")
    }
    if zValues.Name == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] missing 'zValues.Name'")
    }
    if zValues.Name == "external" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] illegal value \"external\" specified for 'zValues.Name'")
    }

    // lookup all indexed fields except ID
    zQuery := new(Zone)
    zQuery.File = zValues.File
    zQuery.Name = zValues.Name

    zPrivate := lookupZone(zQuery)
    if zPrivate != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] another zone with similar properties already exists")
    }

    return createZone(zValues)   // zValues.ID will be ignored
}

func (z *Zone) Read() (zone *Zone, err error) {
    if z.ID == 0 {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] missing 'z.ID'")
    }

    // lookup the ID field only, ignore any other fields
    zQuery := new(Zone)
    zQuery.ID = z.ID

    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] zone not found")
    }

    zPrivate, err = readZone(zPrivate)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    zone = new(Zone)
    zone.ID      = zPrivate.ID
    zone.File    = zPrivate.File
    zone.Name    = zPrivate.Name
    zone.Notes   = zPrivate.Notes
    // no computed fields

    return zone, nil
}

func (z *Zone) Update(zValues *Zone) error {
    if z.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Update(zValues)] missing 'z.ID'")
    }

    // lookup the ID field only, ignore any other fields
    zQuery := new(Zone)
    zQuery.ID = z.ID

    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Update(zValues)] zone not found")
    }

    return updateZone(zPrivate, zValues)   // zValues.ID, zValues.Name and zValues.File will be ignored
}

func (z *Zone) Delete() error {
    if z.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Delete(zValues)] missing 'z.ID'")
    }

    // lookup the ID field only, ignore any other fields
    zQuery := new(Zone)
    zQuery.ID = z.ID

    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Delete()] zone not found")
    }

    return deleteZone(zPrivate)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (z *Zone)         the result of the public CreateFile method and LookupFile method
//                         this doesn't include the computed fields (always use a read method to get the computed fields)
//                         this doesn't include the private fields
//
//                     the result of the private createFile method and lookupFile method (hosts.go)
//                         this doesn't include the computed fields (always use a read method to get the computed fields)
//
//                     the anchor for the public Read/Update/Delete methods
//                         this must include the 'ID' field
//
//                     the input for the private read/update/delete methods
//                         this must include the private 'id' field
//
// - (zQuery *Zone)    the input for the public LookupFile method
//                     the input for the private lookupFile method (hosts.go)
//                         this should include at least one of the indexed fields
//
//   (zValues *Zone)   the input for the public CreateFile/Update methods
//                     the input for the private createFile/update methods
//                         for a create method, this must include *all* writeMany and writeOnce fields
//                         for an update method, this must include *all* writeMany fields
//
// - (zone *Zone)      the result of the public Read method
//                     the result of the private read method
//                         this does include all computed fields
//
// -----------------------------------------------------------------------------

func createZone(zValues *Zone) error {
    // create zone
    z := new(Zone)
    z.File  = zValues.File
    z.Name  = zValues.Name
    z.Notes = zValues.Notes

//    z.lines = make([]string, 0)                                               // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

    addZone(z)   // adds z.ID and z.id

    return nil
}

func readZone(z *Zone) (zone *Zone, err error) {
    return z, nil
}

func updateZone(z *Zone, zValues *Zone) error {
    return nil
}

func deleteZone(z *Zone) error {
    // remove and zero file object
    z.File    = 0
    z.Name    = ""
    z.Notes   = ""

    removeZone(z)   // zeroes z.ID and z.id

    return nil
}

// -----------------------------------------------------------------------------

func goRenderLines(z *Zone) chan string {
    lines := make(chan string)

    go func() {
        defer close(lines)

        for _, line := range z.lines {
            lines <- line
        }
    }()

    return lines
}

// -----------------------------------------------------------------------------

func goScanLines(z *Zone, lines <-chan string) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // create a hash for the checksum of the zone
        hash := sha1.New()
        
        for line := range lines {
            // update hash
            _, _ = io.WriteString(hash, line)

            // update zone
            z.lines = append(z.lines, line)
        }

        // save checksum of the zone
        checksum := hash.Sum(nil)
        z.checksum = hex.EncodeToString(checksum[:])

        done <- true
        return
    }()

    return done
}
