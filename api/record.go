package api

import (
    "errors"
)

// -----------------------------------------------------------------------------

type Record struct {
    // readOnly
    ID      int        // indexed   // read-write in a rQuery
    // read-writeOnce
    Zone    int        // indexed   // read-write in a rQuery
    // read-writeMany
    Address string     // indexed
    Names   []string   // indexed
    Comment string
    Notes   string
    // readOnly
    Managed bool
    // readOnly        //-computed
//    FQDN       string
//    Domain     string
//    RootDomain string
    // private
    id      recordID
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
    r.Managed   = rPrivate.Managed
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
    rValues.Managed = true

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
        return nil, errors.New("[ERROR][terraform-provider-hosts/api/r.Read()] record not found")
    }

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
    record.Managed = rPrivate.Managed
    // no computed fields

    return record, nil
}

func (r *Record) Update(rValues *Record) error {
    if r.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] missing 'r.ID'")
    }
    if rValues.Address == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] missing 'rValues.Address'")
    }
    if len(rValues.Names) == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] missing 'rValues.Names'")
    }

    // lookup the ID field only, ignore any other fields
    rQuery := new(Record)
    rQuery.ID = r.ID

    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Update(rValues)] record not found")
    }

    return updateRecord(rPrivate, rValues)   // rValues.ID and rValues.Zone will be ignored
}

func (r *Record) Delete() error {
    if r.ID == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete(rValues)] missing 'r.ID'")
    }

    // lookup the ID field only, ignore any other fields
    rQuery := new(Record)
    rQuery.ID = r.ID

    rPrivate := lookupRecord(rQuery)
    if rPrivate == nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/r.Delete()] record not found")
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
    r.Zone    = rValues.Zone
    r.Address = rValues.Address
    r.Names   = make([]string, len(rValues.Names))
    copy(r.Names, rValues.Names)
    r.Comment = rValues.Comment
    r.Notes   = rValues.Notes
    r.Managed = rValues.Managed

    addRecord(r)   // adds r.ID and r.id

//    managed := rValues.Managed   // managed when called from CreateRecord, unmanaged when called from goScanLines (zone.go)

    return nil
}

func readRecord(r *Record) (record *Record, err error) {
    return r, nil
}

func updateRecord(r *Record, rValues *Record) error {
    return nil
}

func deleteRecord(r *Record) error {
    // remove and zero file object
    r.Zone    = 0
    r.Address = ""
    r.Names   = make([]string, 0)
    r.Comment = ""
    r.Notes   = ""
    r.Managed = false

    removeRecord(r)   // zeroes r.ID and r.id

    return nil
}