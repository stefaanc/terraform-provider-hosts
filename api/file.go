//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "bufio"
    "bytes"
    "crypto/sha1"
    "errors"
    "encoding/hex"
    //"fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "strings"
)

// -----------------------------------------------------------------------------

type File struct {
    // readOnly
    ID        int      // indexed   // read-write in a fQuery
    // read-writeOnce
    Path      string   // indexed
    // read-writeMany
    Notes     string
    // private
    id        fileID
    hostsFile *fileObject
    zones     []*zoneObject   // !!! beware of memory leaks
}

func LookupFile(fQuery *File) (f *File) {
    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return nil
    }

    // make a copy without the private fields
    f = new(File)
    f.ID    = fPrivate.ID
    f.Path  = fPrivate.Path
    f.Notes = fPrivate.Notes
    // ignore computed fields

    return f
}

func CreateFile(fValues *File) error {
    if fValues.Path == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateFile(fValues)] missing 'fValues.Path'")
    }

    // lookup all indexed fields except ID
    fQuery := new(File)
    fQuery.Path = fValues.Path

    fPrivate := lookupFile(fQuery)
    if fPrivate != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateFile(fValues)] another file with similar properties already exists")
    }

    return createFile(fValues)   // fValues.ID will be ignored
}

func (f *File) Read() (file *File, err error) {
    if f.ID == 0 {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/f.Read()] missing 'f.ID'")
    }

    // lookup the ID field only, ignore any other fields
    fQuery := new(File)
    fQuery.ID = f.ID

    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/f.Read()] file not found")
    }

    // read file
    fPrivate, err = readFile(fPrivate)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    file = new(File)
    file.ID    = fPrivate.ID
    file.Path  = fPrivate.Path
    file.Notes = fPrivate.Notes
    // no computed fields

    return file, nil
}

func (f *File) Update(fValues *File) error {
    if f.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/f.Update(fValues)] missing 'f.ID'")
    }

    // lookup the ID field only, ignore any other fields
    fQuery := new(File)
    fQuery.ID = f.ID

    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/f.Update(fValues)] file not found")
    }

    return updateFile(fPrivate, fValues)   // fValues.ID and fValues.Path will be ignored
}

func (f *File) Delete() error {
    if f.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/f.Delete(fValues)] missing 'f.ID'")
    }

    // lookup the ID field only, ignore any other fields
    fQuery := new(File)
    fQuery.ID = f.ID

    fPrivate := lookupFile(fQuery)
    if fPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/f.Delete()] file not found")
    }

    return deleteFile(fPrivate)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (f *File)         the result of the public CreateFile method and LookupFile method
//                         this doesn't include the computed fields (always use a read method to get the computed fields)
//                         this doesn't include the private fields
//
//                     the result of the private createFile method and lookupFile method (hosts.go)
//                         this doesn't include the computed fields (always use a read method to get the computed fields)
//
//                     the anchor for the public Read/Update/Delete methods
//                         this must include the 'ID' field
//
//                     the input for the private readFile/updateFile/deleteFile methods
//                         this must include the private 'id' field
//
// - (fQuery *File)    the input for the public LookupFile method
//                     the input for the private lookupFile method (hosts.go)
//                         this should include at least one of the indexed fields
//
//   (fValues *File)   the input for the public CreateFile/Update methods
//                     the input for the private createFile/updateFile methods
//                         for a create method, this must include *all* writeMany and writeOnce fields
//                         for an update method, this must include *all* writeMany fields
//
// - (file *File)      the result of the public Read method
//                     the result of the private readFile method
//                         this does include all computed fields
//
// -----------------------------------------------------------------------------

func createFile(fValues *File) error {
    // create and initialize file object
    f := new(File)
    f.Path = fValues.Path
    f.Notes = fValues.Notes

    f.hostsFile = fValues.hostsFile  // requested by goScanFile()
    // f.zones                       // filled by goScanFile()

    addFile(f)   // adds f.ID and f.id

    if fValues.hostsFile == nil {   // if requested by CreateFile()
        // add the file to the hosts
        hostsFile := new(fileObject)
        hostsFile.file = f      // !!! beware of memory leaks
        f.hostsFile = hostsFile // !!! beware of memory leaks

        addFileObject(hostsFile)

        // read physical file, if it doesn't exist then create it
        data, err := ioutil.ReadFile(f.Path)
        if err == nil {
            log.Printf("[INFO][terraform-provider-hosts/api/readFile()] read physical file %d, path %q\n", f.ID, f.Path)
        } else {
            if os.IsNotExist(err) {
                data = []byte(nil)
                err = ioutil.WriteFile(f.Path, data, 0644)
                if err == nil {
                    log.Printf("[INFO][terraform-provider-hosts/api/createFile()] created physical file %d, path %q\n", f.ID, f.Path)
                }
            }
        }
        if err != nil {
            // restore consistent state
            removeFileObject(hostsFile)
            f.hostsFile = nil   // !!! avoid memory leaks
            removeFile(f)

            return err
        }

        checksum := sha1.Sum(data)
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])

        // process data
        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done
   }

    log.Printf("[INFO][terraform-provider-hosts/api/createFile()] created file %d, path %q\n", f.ID, f.Path)
    return nil
}

func readFile(f *File) (file *File, err error) {
    // read physical file
    data, err := ioutil.ReadFile(f.Path)
    if err != nil {
        return nil, err
    }
    log.Printf("[INFO][terraform-provider-hosts/api/readFile()] read physical file %d, path %q\n", f.ID, f.Path)

    checksum := sha1.Sum(data)
    newChecksum := hex.EncodeToString(checksum[:])

    if f.hostsFile.checksum != newChecksum {
        f.hostsFile.checksum = newChecksum

        // process data
        done := goScanFile(f.hostsFile, bytes.NewReader(data))
        _ = <-done
    }

    // no computed fields

    log.Printf("[INFO][terraform-provider-hosts/api/readFile()] read file %d, path %q\n", f.ID, f.Path)
    return f, nil
}

func updateFile(f *File, fValues *File) error {
    notes   := f.Notes     // save so we can restore if needed
    oldChecksum := f.hostsFile.checksum   // save to compare old with new

    // update file
    f.Notes    = fValues.Notes

    if fValues.hostsFile == nil || f == fValues {   // if requested by f.Update() or if forcing a render/write
        // render file to calculate new checksum
        done := goRenderFile(f)   // updates data & checksum
        _ = <-done
        
        if f.hostsFile.checksum != oldChecksum {
            // update physical file
            err := ioutil.WriteFile(f.Path, f.hostsFile.data, 0644)
            if err != nil {
                // restore consistent state
                f.Notes = notes
                f.hostsFile.data     = []byte(nil)
                f.hostsFile.checksum = oldChecksum

                return err
            }
            log.Printf("[INFO][terraform-provider-hosts/api/updateFile()] updated physical file %d, path %q\n", f.ID, f.Path)
        }

        // don't keep rendered data in memory
        f.hostsFile.data = []byte(nil)
    }

    log.Printf("[INFO][terraform-provider-hosts/api/updateFile()] updated file %d, path %q\n", f.ID, f.Path)
    return nil
}

func deleteFile(f *File) error {
    // remove the zone from the file
    if f.hostsFile != nil {   // if requested by f.Delete()
        removeFileObject(f.hostsFile)
        oldHostsFile := f.hostsFile   // save so we can restore if needed
        f.hostsFile = nil            // !!! avoid memory leaks

        if len(f.zones) == 0 {
            // delete physical file
            err := os.Remove(f.Path)
            if err != nil {
               // restore consistent state
               f.hostsFile = oldHostsFile   // !!! beware of memory leaks
               addFileObject(f.hostsFile)

               return err
            }
            log.Printf("[INFO][terraform-provider-hosts/api/deleteFile()] deleted physical file %d, path %q\n", f.ID, f.Path)
        }
    }

    // save for logging
    id := f.ID
    path := f.Path

    // remove and zero file object
    removeFile(f)   // zeroes f.ID and f.id

    f.Path = ""
    f.Notes = ""

    for _, zoneObject := range f.zones {   // !!! avoid memory leaks
        zoneObject.zone = nil
    }
    f.zones = []*zoneObject(nil)

    log.Printf("[INFO][terraform-provider-hosts/api/deleteFile()] deleted file %d, path %q\n", id, path)
    return nil
}

// -----------------------------------------------------------------------------

func goRenderFile(f *File) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // render bytes
        rendered := bytes.NewBuffer([]byte(nil))
        w := io.Writer(rendered)

        // render lines for the default external zone
        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)
        if z != nil {
            for _, line := range z.fileZone.lines {
                // update lines
                _, _ = io.WriteString(w, line)   // error cannot happen
                _, _ = io.WriteString(w, "\n")   // error cannot happen
            }
        }

        for _, zoneObject := range f.zones {
            if zoneObject.zone.Name == "external" {
                continue
            }

            for _, line := range zoneObject.lines {
                // update lines
                _, _ = io.WriteString(w, line)   // error cannot happen
                _, _ = io.WriteString(w, "\n")   // error cannot happen
            }
        }

        // calculate checksum for the lines
        data := rendered.Bytes()
        checksum := sha1.Sum(data)

        // update zoneObject
        f.hostsFile.data = data
        f.hostsFile.checksum = hex.EncodeToString(checksum[:])

        // finish goRenderZones()
        done <- true
        return
    }()

    return done
}

// -----------------------------------------------------------------------------

func goScanFile(hostsFile *fileObject, r io.Reader) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        f := hostsFile.file

        // keep the old slice of zoneObjects to cleanup old zones that aren't replaced
        oldZones := f.zones

        // update file with new slice of zoneObjects
        f.zones = make([]*zoneObject, 0)

        defer func() {
            // cleanup zones that aren't replaced
            for _, fileZone := range oldZones {
                z := fileZone.zone
                if z.fileZone == fileZone {   // if zone was deleted from the read file, fileZone was not replaced
                    // delete zone object
                    z.fileZone = nil   // !!! avoid memory leaks
                    _ = deleteZone(z)   // error cannot happen
                }
            }
        }()

        // create 'external' zoneObject
        fileZone := new(zoneObject)
        addZoneObject(f, fileZone)

        lines2 := make(chan string)
        done2  := goScanZone(f, fileZone, lines2)

        // save this channel for later use
        fileZoneExternal := fileZone
        linesExternal := lines2
        doneExternal := done2

        // start scanning
        scanner := bufio.NewScanner(r)
        for scanner.Scan() {
            line := scanner.Text()

            // create new zoneObject when startZoneMarker
            // complete old zoneObject if not external
            if strings.HasPrefix(line, startZoneMarker) {
                // line is a marker for the start of new zone
                if lines2 != linesExternal {
                    // unexpected startZoneMarker, probably an endZoneMarker missing => silently ignore
                    log.Printf("[WARNING][terraform-provider-hosts/api/goScanFile()] unexpected start-of-zone marker - missing end-of-zone marker: \n> %q", line)

                    // wait for goScanZone() of the current zone to finish
                    lines2 <- endZoneMarker   // insert an endZoneMarker
                    close(lines2)
                    _ = <-done2
                }

                // create new zone
                fileZone = new(zoneObject)
                addZoneObject(f, fileZone)

                lines2 = make(chan string)
                done2 = goScanZone(f, fileZone, lines2)
            }

            // complete old zoneObject if not external
            if strings.HasPrefix(line, endZoneMarker) {
                // line is a marker for the end of current zone
                if lines2 == linesExternal {
                    // unexpected endZoneMarker, probably a startZoneMarker missing => silently ignore
                    // all records from the zone with missing startZoneMarker will be in the external zone
                    log.Printf("[WARNING][terraform-provider-hosts/api/goScanFile()] unexpected end-of-zone marker, skipping line: \n> %q", line)
                    continue
                }

                // wait for goScanZone() of the current zone to finish
                fileZone.lines = append(fileZone.lines, line)
                lines2 <- line
                close(lines2)
                _ = <-done2

                // back to external zone
                fileZone = fileZoneExternal
                lines2 = linesExternal
                done2 = doneExternal

                continue
            }

            fileZone.lines = append(fileZone.lines, line)
            lines2 <- line
        }
        if err := scanner.Err(); err != nil {   // cannot happen at the moment - crash if code is modified
            // // scanner error
            log.Fatal(err)
        }

        if lines2 != linesExternal {
            // endZoneMarker missing => silently ignore
            log.Printf("[WARNING][terraform-provider-hosts/api/goScanFile()] missing end-of-zone marker")

            // wait for goScanZone() of the current zone to finish
            lines2 <- endZoneMarker   // insert an endZoneMarker
            close(lines2)
            _ = <-done2
        }

        // wait for goScanLines() of the external zone to finish
        close(linesExternal)
        _ = <-doneExternal
        if len(fileZoneExternal.lines) == 0 {
            removeZoneObject(f, fileZoneExternal)
        }

        // finish goScanZones()
        done <- true
        return
    }()

    return done
}

// -----------------------------------------------------------------------------

type zoneObject struct {
    lines    []string
    checksum string
    zone  *Zone   // !!! beware of memory leaks
}

func addZoneObject(f *File, z *zoneObject) {
    f.zones = append(f.zones, z)
    return
}

func removeZoneObject(f *File, z *zoneObject) {
    f.zones = deleteFromSliceOfZoneObjects(f.zones, z)
    return
}

func deleteFromSliceOfZoneObjects(zs []*zoneObject, z *zoneObject) []*zoneObject {
    if len(zs) == 0 {
        return []*zoneObject(nil)   // always return a copy
    }

    newZoneObjects := make([]*zoneObject, 0, len(zs) - 1)
    for _, zoneObject := range zs {
        if z == zoneObject {
            continue
        }
        newZoneObjects = append(newZoneObjects, zoneObject)
    }

    return newZoneObjects
}
