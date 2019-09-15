package api

import (
    "errors"
)

// -----------------------------------------------------------------------------

type Zone struct {
    // readOnly
    ID       zoneID   // indexed - read-write in zQuery
    Managed  bool
    File     string   // indexed - read-write in zQuery
    // read-writeOnce
    Name     string   // indexed
    // read-writeMany
    Notes    string
    // private
    id       zoneID
    lines    []string
    checksum string
}

func (f *File) CreateZone(zValues *Zone) error {
    z := GetZone(zValues)
    if z != nil && z.Managed {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone()] another zone with similar properties already exists")
    }
    if z.Name == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone()] missing 'zValues.Name'")
    }
    if z.Name == "external" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone()] illegal value \"external\" specified for 'zValues.Name'")
    }

    zValues.File    = f.Path
    zValues.Managed = true // take ownership

    return createZone(zValues)
}

func (z *Zone) Read() (zone *Zone, err error) {
    z, err = readZone(z)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    zone = new(Zone)
    zone.ID      = z.id
    zone.Managed = z.Managed
    zone.File    = z.File
    zone.Name    = z.Name
    zone.Notes   = z.Notes

    return zone, nil
}

func (z *Zone) Update(zValues *Zone) error {
    return updateZone(z, zValues)
}

func (z *Zone) Delete() error {
    return deleteZone(z)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (z *Zone)         the result of the create method and the GetZone method (hosts.go)
//                         this may not include the computed fields
//
//                     the input for the read/update/delete methods
//                         this must include the private 'id' field (meaning it is indexed)
//
// - (zQuery *Zone)    the input for the GetRecord method (hosts.go)
//                         this must include at least one of the indexed fields
//
//   (zValues *Zone)   the input for the create/update methods
//                         for a create method, this must include all writeMany and writeOnce fields
//                         for an update method, this must include all writeMany fields
//
// - (zone *Zone)      the result of the read method and the CreateZone method (file.go)
//                         this always includes all computed fields
//
// -----------------------------------------------------------------------------

func createZone(zValues *Zone) error {
    managed := zValues.Managed   // managed when called from CreateZone, unmanaged when called from goScanZones (file.go)
    file    := zValues.File
    name    := zValues.Name
    notes   := zValues.Notes

    // create zone
    z := new(Zone)
    z.Managed = managed
    z.File    = file
    z.Name    = name
    z.Notes   = notes

    err := hosts.addZone(z)
    if err != nil {
        return err
    }

    return nil
}

func readZone(z *Zone) (zone *Zone, err error) {
    return zone, nil
}

func updateZone(z *Zone, zValues *Zone) error {
    return nil
}

func deleteZone(z *Zone) error {
    return nil
}

// -----------------------------------------------------------------------------

func goScanLines(z *Zone, lines <-chan string) chan error {
    cerr := make(chan error)

    go func() {
        defer close(cerr)
        
        for _ = range lines {
//            ...
        }

        cerr <- nil
        return
    }()

    return cerr
}
