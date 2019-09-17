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
            removeFile(f)                                                       // need to remove zones/records too !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
            return err
        }
    }
    checksum := sha1.Sum(data)
    f.checksum = hex.EncodeToString(checksum[:])

    // process data
    cerr := goScanZones(f, bytes.NewReader(data))
    err = <-cerr
    if err != nil {
        removeFile(f)                                                           // need to remove zones/records too !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
        return err
    }

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
        f.checksum = newChecksum

        // process data
        cerr := goScanZones(f, bytes.NewReader(data))
        err = <-cerr
        if err != nil {
            return nil, err
        }
    }

    return f, nil
}

func updateFile(f *File, fValues *File) error {
    data:= []byte("# some updated data")                                        // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    checksum := sha1.Sum(data)                                                  // TBD
    f.checksum = hex.EncodeToString(checksum[:])                                // TBD

    err := ioutil.WriteFile(f.Path, data, 0644)
    if err != nil {
        return err
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

    // delete physical file
    err := os.Remove(path)
    if err != nil {
        return err
    }

    return nil
}

// -----------------------------------------------------------------------------

func goScanZones(f *File, r io.Reader) chan error {
    cerr := make(chan error)

    go func(cerr chan error) {
        defer close(cerr)

        startMarker := "##### Start Of Terraform Zone: "
        endMarker   := "##### End Of Terraform Zone: "

        // create a hash for the checksum of the file
//        hash := sha1.New()

        // create the default external zone
        zone := "external"

        zValues := new(Zone)
        zValues.File = f.ID
        zValues.Name = zone
        err := createZone(zValues)
        if err != nil {
            cerr <- err
            return
        }
        z := lookupZone(zValues)

        // create a channel to the external zone
        lines := make(chan string)
        cerr2 := goScanLines(z, lines)

        // save this channel for later use
        linesExternal := lines
        cerrExternal := cerr2

        // start scanning
//        firstLine := true
        scanner := bufio.NewScanner(r)
        for scanner.Scan() {
            line := scanner.Text()

            // update hash
            // in this case, it is safe ignore errors from io.WriteString
//            if firstLine {
//                firstLine = false
//                _, _ = io.WriteString(hash, line)
//            } else {
//                _, _ = io.WriteString(hash, "\n" + line)
//            }

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
                    if err != nil {
                        close(lines)
                        cerr <- err
                        return
                    }
                    z := lookupZone(zValues)

                    // create a channel to the new zone
                    lines = make(chan string)
                    cerr2 = goScanLines(z, lines)

                    continue
                } 
            } else {
                if strings.HasPrefix(line, endMarker + zone) {
                    // line is a marker for the end of current zone
                    // wait for goScanLines(z) of the current zone to finish
                    close(lines)
                    err = <-cerr2
                    if err != nil {
                        if lines != linesExternal {
                            close(linesExternal)
                        }
                        cerr <- err
                        return
                    }

                    // back to external zone
                    zone = "external"
                    lines = linesExternal
                    cerr2 = cerrExternal

                    continue
                }
            }

            // line is not a marker => send to scanLines(z) it
            lines <- line
        }
        if err := scanner.Err(); err != nil {
            // scanner error
            close(lines)
            if lines != linesExternal {
                close(linesExternal)
            }
            cerr <- err
            return
        }

        if lines != linesExternal {
            // missing end-marker for current zone => silently ignore
            // wait for goScanLines(z) of the current zone to finish
            close(lines)
            err = <-cerr2
            if err != nil {
                close(linesExternal)
                cerr <- err
                return
            }
        }

        // wait for goScanLines(z) of the external zone to finish
        close(linesExternal)
        err = <-cerrExternal
        if err != nil {
            cerr <- err
            return
        }

        // save checksum of the file
//        checksum := hash.Sum(nil)
//        f.checksum = hex.EncodeToString(checksum[:])

        // finish goScanZones(f, r)
        cerr <- nil
        return
    }(cerr)

    return cerr
}
