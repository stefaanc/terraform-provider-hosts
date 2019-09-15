package api

import (
    "errors"
)

// -----------------------------------------------------------------------------

type Record struct {
    // readOnly
    ID      recordID   // indexed - read-write in zQuery
    // read-writeMany
    Address string     // indexed
    Names   []string   // indexed
    Comment string
    Notes   string
    // readOnly        //-computed
//    FQDN       string
//    Domain     string
//    RootDomain string
    // private
    id      recordID
}

func (z *Zone) CreateRecord(rValues *Record) error {
    r := GetRecord(rValues)
    if r != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord()] another record with similar properties already exists")
    }
    if r.Address == "" {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord()] missing 'rValues.Address'")
    }
    if len(r.Names) == 0 {
        return errors.New("[ERROR][terraform-provider-hosts/api/CreateRecord()] missing 'rValues.Names'")
    }

    return createRecord(rValues)
}

func (r *Record) Read() (record *Record, err error) {
    r, err = readRecord(r)
    if err != nil {
        return nil, err
    }

    // make a copy without the private fields
    record = new(Record)
    record.ID      = r.id
    record.Address = r.Address
    record.Names   = make([]string, len(r.Names))
    copy(record.Names, r.Names)
    record.Comment = r.Comment
    record.Notes   = r.Notes

    return record, nil
}

func (r *Record) Update(rValues *Record) error {
    return updateRecord(r, rValues)
}

func (r *Record) Delete() error {
    return deleteRecord(r)
}

// -----------------------------------------------------------------------------
//
// naming guidelines:
//
// - (r *Record)         the result of the create method and the GetRecord method (hosts.go)
//                           this may not include the computed fields
//
//                       the input for the read/update/delete methods
//                           this must include the private 'id' field (meaning it is indexed)
//
// - (rQuery *Record)    the input for the GetRecord method (hosts.go)
//                           this must include at least one of the indexed fields
//
//   (rValues *Record)   the input for the create/update methods
//                           for a create method, this must include all writeMany and writeOnce fields
//                           for an update method, this must include all writeMany fields
//
// - (record *Record)    the result of the read method and the CreateRecord method (zone.go)
//                           this always includes all computed fields
//
// -----------------------------------------------------------------------------

func createRecord(rValues *Record) error {
    return nil
}

func readRecord(r *Record) (record *Record, err error) {
    return record, nil
}

func updateRecord(r *Record, rValues *Record) error {
    return nil
}

func deleteRecord(r *Record) error {
    return nil
}