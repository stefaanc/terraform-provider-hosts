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

type Zone struct {
    // readOnly
    ID       int      // indexed   // read-write in a zQuery
    // read-writeOnce
    File     int      // indexed
    Name     string   // indexed
    // read-writeMany
    Notes    string
    // private
    id       zoneID
    managed  bool
    fileZone *zoneObject       // !!! beware of memory leaks
    records  []*recordObject   // !!! beware of memory leaks
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

    // check file
    fQuery := new(File)
    fQuery.ID = zValues.File
    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] file 'zValues.File' not found")
    }

    // lookup all indexed fields except ID
    zQuery := new(Zone)
    zQuery.File = zValues.File
    zQuery.Name = zValues.Name
    zPrivate := lookupZone(zQuery)
    if zPrivate != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateZone(zValues)] another zone with similar properties already exists")
    }

    // take ownership
    zValues.managed = true

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
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] zone 'z.ID' not found")
    }

    // check file
    fQuery := new(File)
    fQuery.ID = zPrivate.File
    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] file 'z.File' not found")
    }

    // read zone
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
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Update(zValues)] zone 'z.ID' not found")
    }

    // check file
    fQuery := new(File)
    fQuery.ID = zPrivate.File
    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] file 'z.File' not found")
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
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Delete()] zone 'z.ID' not found")
    }
    if zPrivate.Name == "external" {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Delete()] cannot delete zone \"external\"")
    }

    // check file
    fQuery := new(File)
    fQuery.ID = zPrivate.File
    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/z.Read()] file 'z.File' not found")
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
    z.File     = zValues.File
    z.Name     = zValues.Name
    z.Notes    = zValues.Notes

    z.managed  = zValues.managed   // requested by CreateZone()

    addZone(z)   // adds z.ID and z.id

    if zValues.fileZone == nil {   // if requested by CreateZone()
        // add the zone to the file
        fileZone := new(zoneObject)
        fileZone.zone = z       // !!! beware of memory leaks
        z.fileZone = fileZone   // !!! beware of memory leaks

        fQuery := new(File)
        fQuery.ID = z.File
        f := lookupFile(fQuery)
        addZoneObject(f, fileZone)
    
        // render zone
        renderZone(z)   // updates lines & checksum

        // update file
        err := updateFile(f, f)
        if err != nil {
            // restore consistent state
            removeZoneObject(f, fileZone)
            z.fileZone = nil   // !!! avoid memory leaks
            removeZone(z)

            return err
        }
    } else {                       // requested by goScanZone()
        // update zone & zoneObject
        z.fileZone = zValues.fileZone   // !!! beware of memory leaks
        z.fileZone.zone = z             // !!! beware of memory leaks
    }

    log.Printf("[INFO][terraform-provider-hosts/api/createZone()] created file %d, zone %q\n", z.File, z.Name)
    return nil
}

func readZone(z *Zone) (zone *Zone, err error) {
    // read file
    fQuery := new(File)
    fQuery.ID = z.File
    f := lookupFile(fQuery)
    _, err = readFile(f)
    if err != nil {
        return nil, err
    }

    // don't return z, instead lookup zone
    // - to cover case where zone was deleted by external programs
    zQuery := new(Zone)
    zQuery.ID = z.ID
    zone = lookupZone(zQuery)

    // no computed fields
    if zone != nil {
        log.Printf("[INFO][terraform-provider-hosts/api/readZone()] read file %d, zone %q\n", zone.File, zone.Name)
    }
    return zone, nil
}

func updateZone(z *Zone, zValues *Zone) error {
    notes   := z.Notes     // save so we can restore if needed
    oldLines    := z.fileZone.lines      // save so we can restore if needed
    oldChecksum := z.fileZone.checksum   // save to compare old with new

    // update zone
    z.Notes    = zValues.Notes

    if zValues.fileZone == nil || z == zValues {   // if requested by z.Update() or if forcing a render/write
        // render zone to calculate new checksum
        renderZone(z)   // updates lines & checksum
        
        if z.fileZone.checksum != oldChecksum {
            // update file
            fQuery := new(File)
            fQuery.ID = z.File
            f := lookupFile(fQuery)
            err := updateFile(f, f)
            if err != nil {
                // restore consistent state
                z.Notes    = notes
                z.fileZone.lines    = oldLines
                z.fileZone.checksum = oldChecksum

                return err
            }
        }
    } else {                       // requested by goScanZone()
        // update zone & zoneObject
        z.fileZone = zValues.fileZone   // !!! beware of memory leaks
        z.fileZone.zone = z             // !!! beware of memory leaks
    }

    log.Printf("[INFO][terraform-provider-hosts/api/updateZone()] updated file %d, zone %q\n", z.File, z.Name)
    return nil
}

func deleteZone(z *Zone) error {
    // remove the zone from the file
    if z.fileZone != nil {   // if requested by z.Delete()
        fQuery := new(File)
        fQuery.ID = z.File
        f := lookupFile(fQuery)

        removeZoneObject(f, z.fileZone)
        oldFileZone := z.fileZone   // save so we can restore if needed
        z.fileZone = nil            // !!! avoid memory leaks

        err := updateFile(f, f)
        if err != nil {
            // restore consistent state
            z.fileZone = oldFileZone   // !!! beware of memory leaks
            addZoneObject(f, z.fileZone)

            return err
        }
    }

    // save for logging
    file := z.File
    name := z.Name

    // remove and zero zone object
    removeZone(z)   // zeroes z.ID and z.id

    z.File     = 0
    z.Name     = ""
    z.Notes    = ""

    z.managed  = false
    for _, recordObject := range z.records {   // !!! avoid memory leaks
        recordObject.record = nil
    }
    z.records  = []*recordObject(nil)

    log.Printf("[INFO][terraform-provider-hosts/api/deleteZone()] deleted file %d, zone %q\n", file, name)
    return nil
}

// -----------------------------------------------------------------------------

var startZoneMarker string = "##### Start Of Terraform Zone: "
var endZoneMarker string   = "##### End Of Terraform Zone: "

// -----------------------------------------------------------------------------

func renderZone(z *Zone) {
    // create a hash for the checksum of the zone
    hash := sha1.New()

    // render strings
    rendered := make([]string, 0)

    // render marker
    line := startZoneMarker + z.Name + " #####"
    padding := 80 - len(line)
    if padding < 0 { padding = 0 }
    line += strings.Repeat("#", padding)

    // update hash
    _, _ = io.WriteString(hash, line)   // error cannot happen
    _, _ = io.WriteString(hash, "\n")   // error cannot happen

    // update lines
    rendered = append(rendered, line)

    for _, recordObject := range z.records {
        // update hash
        _, _ = io.WriteString(hash, recordObject.lines[0])   // error cannot happen   // at this moment we support only single-line records
        _, _ = io.WriteString(hash, "\n")                    // error cannot happen

        // update lines
        rendered = append(rendered, recordObject.lines[0])                      // at this moment we support only single-line records
    }

    // render marker
    line = endZoneMarker + z.Name + " #####"
    padding = 80 - len(line)
    if padding < 0 { padding = 0 }
    line += strings.Repeat("#", padding)

    // update hash
    _, _ = io.WriteString(hash, line)   // error cannot happen
    _, _ = io.WriteString(hash, "\n")   // error cannot happen

    // update lines
    rendered = append(rendered, line)

    // calculate checksum for the lines
    checksum := hash.Sum(nil)

    // update zoneObject
    z.fileZone.lines = rendered
    z.fileZone.checksum = hex.EncodeToString(checksum[:])

    return
}

// -----------------------------------------------------------------------------

func goScanZone(f *File, fileZone *zoneObject, lines <-chan string) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // create a hash for the checksum of the zone
        hash := sha1.New()

        // scan first line, possibly a start-zone-marker
        var line string
        var receivedLine bool
        for line = range lines {   // using for to cover case where lines channel is closed before a line is sent - applicable to files without external records
            receivedLine = true
            break
        }
        if !receivedLine {
            done <- true
            return
        }

        // get zone name
        var zone string
        isStartMarker := strings.HasPrefix(line, startZoneMarker)
        if isStartMarker {
            zone = strings.Trim(line[len(startZoneMarker):], " #")
        } else {
            zone = "external"
        }

        // update hash
        _, _ = io.WriteString(hash, line)   // error cannot happen
        _, _ = io.WriteString(hash, "\n")   // error cannot happen

        // create/update zone
        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = zone
        z := lookupZone(zQuery)

        var oldChecksum string
        if z == nil {
            // create zone
            // zQuery.Notes   = ""   // notes are not saved in file

            zQuery.fileZone = fileZone

            _ = createZone(zQuery)   // error cannot happen
            z = lookupZone(zQuery)
        } else {
            // pickup old checksum
            oldChecksum = z.fileZone.checksum

            // update zone
            zQuery.Notes   = z.Notes   // notes are not saved in file, need to pick up from old zone

            zQuery.fileZone = fileZone
        
            _ = updateZone(z, zQuery)   // error cannot happen
        }

        // keep the old slice of recordObjects to cleanup old records that aren't replaced
        oldRecords := z.records

        // update zone with new slice of recordObjects
        z.records = make([]*recordObject, 0)

        defer func() {
            // cleanup records that aren't replaced
            for _, zoneRecord := range oldRecords {
                r := zoneRecord.record
                if r != nil {   // if zoneRecord is a record, not a comment/blank-line
                    if r.zoneRecord == zoneRecord {   // if record was deleted from the read zone, zoneRecord was not replaced
                        // delete record object
                        r.zoneRecord = nil   // !!! avoid memory leaks
                        _ = deleteRecord(r)   // error cannot happen
                    }
                }
            }
        }()

        // update zone
        if !isStartMarker {
            zoneRecord := new(recordObject)
            zoneRecord.lines = append(zoneRecord.lines, line)                   // at this moment we support only single-line records
            addRecordObject(z, zoneRecord)
        }

        // collect lines
        for line := range lines {

            if strings.HasPrefix(line, endZoneMarker) {
                if len(line) == len(endZoneMarker) {
                    // no zone name in end-marker - goScanFile probably inserted anonymous endZoneMarker
                    // render missing marker
                    line = endZoneMarker + zone + " #####"
                    padding := 80 - len(line)
                    if padding < 0 { padding = 0 }
                    line += strings.Repeat("#", padding)
                    fileZone.lines = append(fileZone.lines, line)
                }

                // update hash
                _, _ = io.WriteString(hash, line)   // error cannot happen
                _, _ = io.WriteString(hash, "\n")   // error cannot happen

                // end of zone
                if !strings.HasPrefix(line, endZoneMarker + zone) {
                    // unexpected endZoneMarker, probably an endZone- and startZone-Marker missing => silently ignore
                    // all records from the zone with missing startZoneMarker will be in the current zone
                    log.Printf("[WARNING][terraform-provider-hosts/api/goScanZone()] unexpected end-of-zone marker - missing end-of-zone and start-of-zone marker: \n> %q\n", line)
                }

                break
            } else {
                // update hash
                _, _ = io.WriteString(hash, line)   // error cannot happen
                _, _ = io.WriteString(hash, "\n")   // error cannot happen
            }

            zoneRecord := new(recordObject)
            zoneRecord.lines = append(zoneRecord.lines, line)               // at this moment we support only single-line records
            addRecordObject(z, zoneRecord)
        }

        // calculate checksum for the lines
        checksum := hash.Sum(nil)
        fileZone.checksum = hex.EncodeToString(checksum[:])

        if fileZone.checksum == oldChecksum {
            done <- true
            return
        }

        // process lines
        for _, zoneRecord := range z.records {
            lines2 := make(chan string)
            done2 := goScanRecord(z, zoneRecord, lines2)

            lines2 <- zoneRecord.lines[0]                                        // at this moment we support only single-line records

            close(lines2)
            _ = <-done2
        }

        done <- true
        return
    }()

    return done
}

// -----------------------------------------------------------------------------

type recordObject struct {
    // remark that a record can be split over multiple lines
    // - f.i. a comment-line before to information-line
    // - f.i. names split over several lines
    lines    []string
    checksum string
    // remark that a recordObject may not have an associated record
    // - f.i. a comment-line in the external zone of the hosts-file
    record   *Record   // !!! beware of memory leaks
}

func addRecordObject(z *Zone, r *recordObject) {
    z.records = append(z.records, r)
    return
}

func removeRecordObject(z *Zone, r *recordObject) {
    z.records = deleteFromSliceOfRecordObjects(z.records, r)
    return
}

func deleteFromSliceOfRecordObjects(rs []*recordObject, r *recordObject) []*recordObject {
    if len(rs) == 0 {
        return []*recordObject(nil)   // always return a copy
    }

    newRecordObjects := make([]*recordObject, 0, len(rs) - 1)
    for _, recordObject := range rs {
        if r == recordObject {
            continue
        }
        newRecordObjects = append(newRecordObjects, recordObject)
    }

    return newRecordObjects
}
