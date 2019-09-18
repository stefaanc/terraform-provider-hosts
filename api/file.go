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
    ID       int      // indexed   // read-write in a fQuery
    // read-writeOnce
    Path     string   // indexed   // read-write in a fQuery
    // read-writeMany
    Notes    string
    // private
    id       fileID
    checksum string
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

    addFile(f)   // adds f.ID and f.id

    // read physical file, if it doesn't exist then create it
    data, err := ioutil.ReadFile(f.Path)
    if err != nil {
        if os.IsNotExist(err) {
            data = []byte(nil)
            err = ioutil.WriteFile(f.Path, data, 0644)
        }

        if err != nil {
            removeFile(f)
            return err
        }
    }
    checksum := sha1.Sum(data)
    f.checksum = hex.EncodeToString(checksum[:])

    // process data
    done := goScanZones(f, bytes.NewReader(data))
    _ = <-done

    return nil
}

func readFile(f *File) (file *File, err error) {
    // read physical file
    data, err := ioutil.ReadFile(f.Path)
    if err != nil {
        return nil, err
    }

    checksum := sha1.Sum(data)
    newChecksum := hex.EncodeToString(checksum[:])

    if f.checksum != newChecksum {
        // process data
        b := bytes.NewBuffer(data)
        done := goScanZones(f, io.Reader(b))
        _ = <-done

        f.checksum = newChecksum
    }

    return f, nil
}

func updateFile(f *File, fValues *File) error {
    // collect data
    b := bytes.NewBuffer([]byte(nil))
    done := goRenderZones(f, io.Writer(b))
    _ = <-done
    data := b.Bytes()

    checksum := sha1.Sum(data)
    newChecksum := hex.EncodeToString(checksum[:])

    if f.checksum !=  newChecksum {
        // write physical file
        err := ioutil.WriteFile(f.Path, data, 0644)
        if err != nil {
            return err
        }

        f.checksum = newChecksum
    }

    return nil
}

func deleteFile(f *File) error {
    path := f.Path

    // remove and zero file object
    f.Path = ""
    f.Notes = ""

    f.checksum = ""

    removeFile(f)   // zeroes f.ID and f.id

    if len(hosts.zoneIndex) == 0 {
        // delete physical file
        err := os.Remove(path)
        if err != nil {
            return err
        }
    }

    return nil
}

// -----------------------------------------------------------------------------

var startMarker string = "##### Start Of Terraform Zone: "
var endMarker string   = "##### End Of Terraform Zone: "

// -----------------------------------------------------------------------------

func goRenderZones(f *File, w io.Writer) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // Remark:
        // io.WriteString error cannot happen with the current implementation of updateFile
        // this would need an io.Writer writing to file instead of to buffer
        // we leave this for later development
        // for now we throw a fatal error in case my assumptions are/become wrong

        // render lines for the default external zone
        zQuery := new(Zone)
        zQuery.File = f.ID
        zQuery.Name = "external"
        z := lookupZone(zQuery)

        lines := goRenderLines(z)
        for line := range lines {
            _, err := io.WriteString(w, line + "\n")
            if err != nil {   // cannot happen
                log.Fatal(err)
                // done <- true
                // return
            }
        }

        for _, zone := range hosts.zoneFiles[f.id] {
            if zone.Name == "external" {
                continue
            }

            line := startMarker + zone.Name + " "
            line = line + strings.Repeat("#", 80 - len(line)) + "\n"
            _, err := io.WriteString(w, line)
            if err != nil {   // cannot happen
                log.Fatal(err)
                // done <- true
                // return
            }

            lines := goRenderLines(zone)
            for line := range lines {
                _, err = io.WriteString(w, line + "\n")
                if err != nil {   // cannot happen
                    log.Fatal(err)
                    // done <- true
                    // return
                }
            }

            line = endMarker + zone.Name + " "
            line = line + strings.Repeat("#", 80 - len(line)) + "\n"
            _, err = io.WriteString(w, line)
            if err != nil {   // cannot happen
                log.Fatal(err)
                // done <- true
                // return
            }
        }

        // finish goRenderZones()
        done <- true
        return
    }()

    return done
}

// -----------------------------------------------------------------------------

func goScanZones(f *File, r io.Reader) chan bool {
    done := make(chan bool)

    go func() {
        defer close(done)

        // Remark:
        // Scanner error cannot happen with the current implementation of createFile/readFile
        // this would need an io.Reader reading from file instead of from buffer
        // we leave this for later development
        // for now we throw a fatal error in case my assumptions are/become wrong

        // create the default external zone
        zone := "external"

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = zone
        err := createZone(zValues)
        if err != nil {   // cannot happen
            log.Fatal(err)
            //done <- true
            //return
        }
        z := lookupZone(zValues)

        // create a channel to the external zone
        lines := make(chan string)
        done2 := goScanLines(z, lines)

        // save this channel for later use
        linesExternal := lines
        doneExternal := done2

        // start scanning
        scanner := bufio.NewScanner(r)
        for scanner.Scan() {
            line := scanner.Text()

            // scan line
            if lines == linesExternal {
                if strings.HasPrefix(line, startMarker) {
                    // line is a marker for the start of new zone
                    // create new zone
                    zone = strings.Trim(line[len(startMarker):], " #")

                    zValues = new(Zone)
                    zValues.File = f.ID
                    zValues.Name = zone
                    err = createZone(zValues)
                    if err != nil {   // cannot happen
                        close(lines)
                        log.Fatal(err)
                        //done <- true
                        //return
                    }
                    z := lookupZone(zValues)

                    // create a channel to the new zone
                    lines = make(chan string)
                    done2 = goScanLines(z, lines)

                    continue
                } 
            } else {
                if strings.HasPrefix(line, endMarker + zone) {
                    // line is a marker for the end of current zone
                    // wait for goScanLines() of the current zone to finish
                    close(lines)
                    _ = <-done2

                    // back to external zone
                    zone = "external"
                    lines = linesExternal
                    done2 = doneExternal

                    continue
                }
            }

            // line is not a marker => send to scanLines(z) it
            lines <- line
        }
        if err := scanner.Err(); err != nil {   // cannot happen
            // scanner error
            if lines != linesExternal {
                close(lines)
            }
            close(linesExternal)
            log.Fatal(err)
            //done <- true
            //return
        }

        if lines != linesExternal {
            // missing end-marker for current zone => silently ignore
            // wait for goScanLines(z) of the current zone to finish
            close(lines)
            _ = <-done2
        }

        // wait for goScanLines(z) of the external zone to finish
        close(linesExternal)
        _ = <-doneExternal

        // finish goScanZones()
        done <- true
        return
    }()

    return done
}
