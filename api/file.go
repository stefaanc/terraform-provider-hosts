package api

import (
    "crypto/sha1"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
)

// -----------------------------------------------------------------------------

type File struct {
    // readOnly
    ID       fileID   // indexed
    // read-writeOnce
    Path     string   // indexed
    // read-write
    Notes    string
    // private
    id       fileID
    data     []byte
    zones    map[string]string
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

    file = new(File)
    file.ID = f.id
    file.Path = f.Path
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
    path := fValues.Path
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
    f.data, err = ioutil.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            f.data = make([]byte, 0)
            err = ioutil.WriteFile(path, f.data, 0644)
        }

        if err != nil {
            hosts.removeFile(f)
            return err
        }
    }
    f.checksum = fmt.Sprintf("%x", sha1.Sum(f.data))

    scanFileZones(f)

    return nil
}

func readFile(f *File) (file *File, err error) {
    f.data, err = ioutil.ReadFile(f.Path)
    if err != nil {
        return nil, err
    }
    f.checksum = fmt.Sprintf("%x", sha1.Sum(f.data))

    scanFileZones(f)

    return f, nil
}

func updateFile(f *File, fValues *File) error {
    err := ioutil.WriteFile(f.Path, fValues.data, 0644)
    if err != nil {
        return err
    }
    f.checksum = fmt.Sprintf("%x", sha1.Sum(fValues.data))

    if f != fValues { 
        // copy slice so f doesn't change when caller changes fValues
        f.data = make([]byte, len(fValues.data))
        copy(f.data, fValues.data)
    }

    return nil
}

func deleteFile(f *File) error {
    path := f.Path

    // remove and zero file object
    f.Path = ""
    f.Notes = ""

    f.data = []byte(nil)
    f.checksum = ""

    hosts.removeFile(f)

    // delete physical file
    err := os.Remove(path)
    if err != nil {
        return err
    }

    return nil
}

func scanFileZones(f *File) {

}
