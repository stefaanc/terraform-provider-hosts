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
    "log"
    "strings"
)

// -----------------------------------------------------------------------------

type Record struct {
    // readOnly
    ID         int        // indexed   // read-write in a rQuery
    // read-writeOnce
    Zone       int        // indexed
    Address    string     // indexed
    Names      []string   // indexed
    // read-writeMany
    Comment    string
    Notes      string
    // readOnly        //-computed
//    FQDN         string
//    Domain       string
//    RootDomain   string
    // private
    id         recordID
    managed    bool
    zoneRecord *recordObject   // !!! beware of memory leaks
}

func LookupRecord(rQuery *Record) (r *Record) {
    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return nil
    }

    // make a copy without the private fields
    r = new(Record)
    r.ID      = rPrivate.ID
    r.Zone    = rPrivate.Zone
    r.Address = rPrivate.Address
    r.Names   = make([]string, len(rPrivate.Names))
    copy(r.Names, rPrivate.Names)
    r.Comment = rPrivate.Comment
    r.Notes   = rPrivate.Notes
    // ignore computed fields

    return r
}

func CreateRecord(rValues *Record) error {
    if rValues.Zone == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord(rValues)] missing 'rValues.Zone'")
    }
    if rValues.Address == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord(rValues)] missing 'rValues.Address'")
    }
    if len(rValues.Names) == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord(rValues)] missing 'rValues.Names'")
    }

    // check zone
    zQuery := new(Zone)
    zQuery.ID = rValues.Zone
    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] zone 'rValues.Zone' not found")
    }
    if zPrivate.Name == "external" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord(rValues)] cannot create records in the \"external\" zone")
    }

    // lookup all indexed fields except ID
    rQuery := new(Record)
    rQuery.Zone    = rValues.Zone
    rQuery.Address = rValues.Address
    rQuery.Names   = make([]string, len(rValues.Names))
    copy(rQuery.Names, rValues.Names)
    rPrivate := lookupRecord(rQuery)
    if rPrivate != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord(rValues)] another record with similar properties already exists")
    }

    // take ownership
    rValues.managed = true

    return createRecord(rValues)   // rValues.ID will be ignored
}

func (r *Record) Read() (record *Record, err error) {
    if r.ID == 0 {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/r.Read()] missing 'r.ID'")
    }

    // lookup the ID field only, ignore any other fields
    rQuery := new(Record)
    rQuery.ID = r.ID
    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/r.Read()] record 'r.ID' not found")
    }

    // check zone
    zQuery := new(Zone)
    zQuery.ID = rPrivate.Zone
    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] zone 'r.Zone' not found")
    }

    // read record
    rPrivate, err = readRecord(rPrivate)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    record = new(Record)
    record.ID      = rPrivate.ID
    record.Zone    = rPrivate.Zone
    record.Address = rPrivate.Address
    record.Names   = make([]string, len(rPrivate.Names))
    copy(record.Names, rPrivate.Names)
    record.Comment = rPrivate.Comment
    record.Notes   = rPrivate.Notes
    // no computed fields

    return record, nil
}

func (r *Record) Update(rValues *Record) error {
    if r.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] missing 'r.ID'")
    }

    // lookup the ID field only, ignore any other fields
    rQuery := new(Record)
    rQuery.ID = r.ID
    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] record 'r.ID' not found")
    }

    // check zone
    zQuery := new(Zone)
    zQuery.ID = rPrivate.Zone
    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] zone 'r.Zone' not found")
    }
    if zPrivate.Name == "external" {
        if rValues.Comment != rPrivate.Comment {
            return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] cannot update 'r.Comment' for records in the \"external\" zone")
        }
    }

    return updateRecord(rPrivate, rValues)   // rValues.ID and rValues.Zone will be ignored
}

func (r *Record) Delete() error {
    if r.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete()] missing 'r.ID'")
    }

    // lookup the ID field only, ignore any other fields
    rQuery := new(Record)
    rQuery.ID = r.ID
    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete()] record 'r.ID' not found")
    }

    // check zone
    zQuery := new(Zone)
    zQuery.ID = rPrivate.Zone
    zPrivate := lookupZone(zQuery)
    if zPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete()] zone 'r.Zone' not found")
    }
    if zPrivate.Name == "external" {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete()] cannot delete records in the \"external\" zone")
    }

    return deleteRecord(rPrivate)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (r *Record)         the result of the public CreateRecord method and LookupRecord method
//                           this doesn't include the computed fields (always use a read method to get the computed fields)
//                           this doesn't include private fields
//
//                       the result of the private createRecord method and lookupRecord method (hosts.go)
//                           this doesn't include the computed fields (always use a read method to get the computed fields)
//
//                       the anchor for the public Read/Update/Delete methods
//                           this must include the 'ID' field
//
//                       the input for the private readRecord/updateRecord/deleteRecord methods
//                           this must include the private 'id' field
//
// - (rQuery *Record)    the input for the public LookupRecord method
//                       the input for the private lookupRecord method (hosts.go)
//                           this should include at least one of the indexed fields
//
//   (rValues *Record)   the input for the public CreateRecord/Update methods
//                       the input for the private createRecord/updateRecord methods
//                           for a create method, this must include *all* writeMany and writeOnce fields
//                           for an update method, this must include *all* writeMany fields
//
// - (record *Record)    the result of the public Read method
//                       the result of the private readRecord method
//                           this does include all computed fields
//
// -----------------------------------------------------------------------------

func createRecord(rValues *Record) error {
    // create record
    r := new(Record)
    r.Zone       = rValues.Zone
    r.Address    = rValues.Address
    r.Names      = make([]string, len(rValues.Names))
    copy(r.Names, rValues.Names)
    r.Comment    = rValues.Comment
    r.Notes      = rValues.Notes

    r.managed    = rValues.managed      // requested by CreateRecord()

    addRecord(r)   // updates r.ID and r.id

    if rValues.zoneRecord == nil {   // if requested by CreateRecord()
        // add the record to the zone
        zoneRecord := new(recordObject)
        zoneRecord.record = r       // !!! beware of memory leaks
        r.zoneRecord = zoneRecord   // !!! beware of memory leaks

        zQuery := new(Zone)
        zQuery.ID = r.Zone
        z := lookupZone(zQuery)
        addRecordObject(z, zoneRecord)
    
        // render record
        renderRecord(r)   // updates lines & checksum

        // update zone
        err := updateZone(z, z)
        if err != nil {
            // restore consistent state
            removeRecordObject(z, zoneRecord)
            r.zoneRecord = nil   // !!! avoid memory leaks
            removeRecord(r)

            return err
        }
    } else {                         // requested by goScanRecord()
        // update record & recordObject
        r.zoneRecord = rValues.zoneRecord   // !!! beware of memory leaks
        r.zoneRecord.record = r             // !!! beware of memory leaks
    }

    log.Printf("[INFO][terraform-provider-hosts/api/createRecord()] created zone %d, record %q - %#v\n", r.Zone, r.Address, r.Names)
    return nil
}

func readRecord(r *Record) (record *Record, err error) {
    // read zone
    zQuery := new(Zone)
    zQuery.ID = r.Zone
    z := lookupZone(zQuery)
    _, err = readZone(z)
    if err != nil {
        return nil, err
    }

    // don't return r, instead lookup record
    // - to cover case where record was deleted by external programs
    rQuery := new(Record)
    rQuery.ID = r.ID
    record = lookupRecord(rQuery)

    // no computed fields

    if record != nil {
        log.Printf("[INFO][terraform-provider-hosts/api/readRecord()] read zone %d, record %q - %#v\n", record.Zone, record.Address, record.Names)
    }
    return record, nil
}

func updateRecord(r *Record, rValues *Record) error {
    comment := r.Comment   // save so we can restore if needed
    notes   := r.Notes     // save so we can restore if needed
    oldChecksum := r.zoneRecord.checksum   // save to compare old with new

    // update record
    r.Comment  = rValues.Comment
    r.Notes    = rValues.Notes

    if rValues.zoneRecord == nil || r == rValues {   // if requested by r.Update() or if forcing a render/write
        // render record to calculate new checksum
        renderRecord(r)   // updates lines & checksum
        
        if r.zoneRecord.checksum != oldChecksum {
            // update zone
            zQuery := new(Zone)
            zQuery.ID = r.Zone
            z := lookupZone(zQuery)
            err := updateZone(z, z)
            if err != nil {
                // restore consistent state
                r.Comment = comment
                r.Notes   = notes
                renderRecord(r)

                return err
            }
        }
    } else {                         // requested by goScanRecord()
        // update record & recordObject
        r.zoneRecord = rValues.zoneRecord   // !!! beware of memory leaks
        r.zoneRecord.record = r             // !!! beware of memory leaks
    }

    log.Printf("[INFO][terraform-provider-hosts/api/updateRecord()] updated zone %d, record %q - %#v\n", r.Zone, r.Address, r.Names)
    return nil
}

func deleteRecord(r *Record) error {
    // remove the record from the zone
    if r.zoneRecord != nil {   // if requested by r.Delete()
        zQuery := new(Zone)
        zQuery.ID = r.Zone
        z := lookupZone(zQuery)

        removeRecordObject(z, r.zoneRecord)
        oldZoneRecord := r.zoneRecord   // save so we can restore if needed
        r.zoneRecord = nil              // !!! avoid memory leaks

        err := updateZone(z, z)
        if err != nil {
            // restore consistent state
            r.zoneRecord = oldZoneRecord   // !!! beware of memory leaks
            addRecordObject(z, r.zoneRecord)

            return err
        }
    }

    // save for logging
    zone := r.Zone
    address := r.Address
    names := r.Names

    // remove and zero record
    removeRecord(r)   // zeroes r.ID and r.id

    r.Zone       = 0
    r.Address    = ""
    r.Names      = []string(nil)
    r.Comment    = ""
    r.Notes      = ""

    r.managed    = false

    log.Printf("[INFO][terraform-provider-hosts/api/deleteRecord()] deleted zone %d, record %q - %#v\n", zone, address, names)
    return nil
}

// -----------------------------------------------------------------------------

func renderRecord(r *Record) {
    // render strings
    rendered := make([]string, 0, 1)                                            // at this moment we support only single-line records

    line := r.Address

    for _, name := range r.Names {
        line += " " + name
    }

    if r.Comment != "" {
        line += " #" + r.Comment
    }

    rendered = append(rendered, line)

    // calculate checksum for the line
    checksum := sha1.Sum([]byte(line))

    // update recordObject
    r.zoneRecord.lines = rendered
    r.zoneRecord.checksum = hex.EncodeToString(checksum[:])

    return
}

// -----------------------------------------------------------------------------

func goScanRecord(z *Zone, zoneRecord *recordObject, lines <-chan string) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // create a hash for the checksum of the record
        hash := sha1.New()
        
        // collect lines
        collected := make([]string, 0)
        for line := range lines {
            // update hash
            _, _ = io.WriteString(hash, line)

            // update lines
            collected = append(collected, line)
        }
        if len(collected) == 0 {
            done <- true
            return
        }

        // calculate checksum for the lines
        checksum := hash.Sum(nil)
        zoneRecord.checksum = hex.EncodeToString(checksum[:])   // we cannot check if checksum changed because we need to parse the lines to know which record this is

        // update recordObject
        zoneRecord.lines = collected

        // process lines                                                        // at this moment we support only single-line records

        // split the line in an information-part and a comment-part
        parts := strings.SplitN(zoneRecord.lines[0], "#", 2)

        if parts[0] == "" {
            // the line doesn't have an information-part
            done <- true
            return
        }

        comment := ""
        if len(parts) > 1 {
            comment = strings.TrimRight(parts[1], " \t")
        }

        // split the information-part
        parts = strings.Fields(parts[0])

        if len(parts) < 2 {
            // the information-part doesn't have both an address and a name
            log.Printf("[WARNING][terraform-provider-hosts/api/goScanRecord()] information-part doesn't have both an address and a name, skipping line: \n> %q", zoneRecord.lines[0])

            done <- true
            return
        }

        // create a new record if it doesn't exist, otherwise update it
        rQuery := new(Record)
        rQuery.Zone = z.ID
        rQuery.Address = parts[0]
        rQuery.Names = parts[1:]
        r := lookupRecord(rQuery)

        if r == nil {
            // create record
            rQuery.Comment = comment
            // rQuery.Notes   = ""   // notes are not saved in file

            rQuery.zoneRecord = zoneRecord
        
            _ = createRecord(rQuery)   // error cannot happen
        } else {
            // update record
            rQuery.Comment = comment
            rQuery.Notes   = r.Notes   // notes are not saved in file, need to pick up from old record

            rQuery.zoneRecord = zoneRecord
        
            _ = updateRecord(r, rQuery)   // error cannot happen
        }

        done <- true
        return
    }()

    return done
}
