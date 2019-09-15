package api

import (
    "bufio"
    "bytes"
    "crypto/sha1"
    "errors"
    "encoding/hex"
    "io"
    "io/ioutil"
    "os"
    "strings"
)

// -----------------------------------------------------------------------------

type File struct {
    // readOnly
    ID       fileID   // indexed - read-write in zQuery
    // read-writeOnce
    Path     string   // indexed
    // read-writeMany
    Notes    string
    // private
    id       fileID
    checksum string
}

func CreateFile(fValues *File) error {
    f := GetFile(fValues)
    if f != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateFile()] another file with similar properties already exists")
    }

    if fValues.Path == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateFile()] missing 'fValues.Path'")
    }

    return createFile(fValues)
}

func (f *File) Read() (file *File, err error) {
    f, err = readFile(f)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    file = new(File)
    file.ID    = f.id
    file.Path  = f.Path
    file.Notes = f.Notes

    return file, nil
}

func (f *File) Update(fValues *File) error {
    return updateFile(f, fValues)
}

func (f *File) Delete() error {
    return deleteFile(f)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (f *File)         the result of the create method and the GetFile method (hosts.go)
//                         this may not include the computed fields
//
//                     the input for the read/update/delete methods
//                         this must include the private 'id' field (meaning it is indexed)
//
// - (fQuery *File)    the input for the GetFile method (hosts.go)
//                         this must include at least one of the indexed fields
//
//   (fValues *File)   the input for the create/update methods
//                         for a create method, this must include all writeMany and writeOnce fields
//                         for an update method, this must include all writeMany fields
//
// - (file *File)      the result of the read method and the CreateFile method (hosts.go)
//                         this always includes all computed fields
//
// -----------------------------------------------------------------------------

func createFile(fValues *File) error {
    path  := fValues.Path
    notes := fValues.Notes

    // create and initialize file object
    f := new(File)
    f.Path = path
    f.Notes = notes

    err := hosts.addFile(f)
    if err != nil {
        return err
    }

    // read physical file, if it doesn't exist then create it
    data, err := ioutil.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            data = []byte(nil)
            err = ioutil.WriteFile(path, data, 0644)
        }

        if err != nil {
            hosts.removeFile(f)
            return err
        }
    }

    cerr := goScanZones(f, bytes.NewReader(data))
    err = <-cerr
    if err != nil {
        return err
    }

    return nil
}

func readFile(f *File) (file *File, err error) {
    data, err := ioutil.ReadFile(f.Path)
    if err != nil {
        return nil, err
    }
    checksum := sha1.Sum(data)
    f.checksum = hex.EncodeToString(checksum[:])

    cerr := goScanZones(f, bytes.NewReader(data))
    err = <-cerr
    if err != nil {
        return nil, err
    }

    return f, nil
}

func updateFile(f *File, fValues *File) error {
    data:= []byte("some updated data")                                          // TBD !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

    err := ioutil.WriteFile(f.Path, data, 0644)
    if err != nil {
        return err
    }

    cerr := goScanZones(f, bytes.NewReader(data))
    err = <-cerr
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

    hosts.removeFile(f)

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
        hash := sha1.New()

        // create the default external zone
        zone := "external"

        zValues := new(Zone)
        zValues.Managed = false
        zValues.File = f.Path
        zValues.Name = zone
        err := createZone(zValues)
        if err != nil {
            cerr <- err
            return
        }
        z := hosts.getZone(zValues)

        // create a channel to the external zone
        lines := make(chan string)
        cerr2 := goScanLines(z, lines)

        // save this channel for later use
        linesExternal := lines
        cerrExternal := cerr2

        // start scanning
        firstLine := true
        scanner := bufio.NewScanner(r)
        for scanner.Scan() {
            line := scanner.Text()

            // update hash
            // in this case, it is safe ignore errors from io.WriteString
            if firstLine {
                firstLine = false
                _, _ = io.WriteString(hash, line)
            } else {
                _, _ = io.WriteString(hash, "\n" + line)
            }

            // scan line
            if lines == linesExternal {
                if strings.HasPrefix(line, startMarker) {
                    // line is a marker for the start of new zone
                    // create new zone
                    zone = strings.Trim(line[len(startMarker):], " #")

                    zValues = new(Zone)
                    zValues.Managed = false
                    zValues.File = f.Path
                    zValues.Name = zone
                    err = createZone(zValues)
                    if err != nil {
                        close(lines)
                        cerr <- err
                        return
                    }
                    z := hosts.getZone(zValues)

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
        checksum := hash.Sum(nil)
        f.checksum = hex.EncodeToString(checksum[:])

        // finish goScanZones(f, r)
        cerr <- nil
        return
    }(cerr)

    return cerr
}
