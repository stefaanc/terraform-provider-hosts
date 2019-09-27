//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "fmt"
    "log"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

    "github.com/stefaanc/terraform-provider-hosts/api"
)

func resourceHostsRecord() *schema.Resource {
    return &schema.Resource {
        Create: resourceHostsRecordCreate,
        Read:   resourceHostsRecordRead,
        Update: resourceHostsRecordUpdate,
        Delete: resourceHostsRecordDelete,

        Schema: map[string]*schema.Schema {
            "record_id": &schema.Schema {
                Type:     schema.TypeInt,
                Computed: true,   // this is a non-persistent id => computed
            },
            "address": &schema.Schema {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "names": &schema.Schema {
                Type:     schema.TypeList,
                Elem:     &schema.Schema {
                    Type: schema.TypeString,
                    StateFunc: func(val interface{}) string {
                        return strings.ToLower(val.(string))
                    },
                    DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
                        if old == strings.ToLower(new) {
                            return true 
                        }
                        return false
                    },
                },
                Required: true,
                ForceNew: true,
            },
            "comment": &schema.Schema {
                Type:     schema.TypeString,
                Optional: true,
                Default: "",
            },
            "notes": &schema.Schema {
                // remark that non-empty "notes" will always cause a 'diff' since it is not saved in the physical hosts-file
                Type:     schema.TypeString,
                Optional: true,
                Default: "",
            },
        },
    }
}

func resourceHostsRecordCreate(d *schema.ResourceData, m interface{}) error {
    zone := m.(*api.Zone)
    address := d.Get("address").(string)
    ns := d.Get("names").([]interface {})
    names := make([]string, len(ns))
    for i, _ := range ns {
        names[i] = ns[i].(string)
    }
    comment := d.Get("comment").(string)
    notes := d.Get("notes").(string)

    log.Printf(`[INFO][terraform-provider-hosts] creating hosts-record
                    [INFO][terraform-provider-hosts]     zone:    %#v
                    [INFO][terraform-provider-hosts]     address: %#v
                    [INFO][terraform-provider-hosts]     names:   %#v
                    [INFO][terraform-provider-hosts]     comment: %#v
                    [INFO][terraform-provider-hosts]     notes:   %#v
`   , zone.Name, address, names, comment, notes)

    rValues := new(api.Record)
    rValues.Zone    = zone.ID
    rValues.Address = address
    rValues.Names   = names
    rValues.Comment = comment
    rValues.Notes   = notes
    err := api.CreateRecord(rValues)
    if err != nil {
        // this is most probably because
        // - there is an error in the fields that wasn't checked by this provider
        // - the hosts-file cannot be read or created
        log.Printf("[ERROR][terraform-provider-hosts] cannot create hosts-record\n")
        return err
    }

    record := api.LookupRecord(rValues)
    if record == nil {
        // this is most probably because 
        // - the record was deleted out-of-band
        // - a record with the same indexed fields was added out-of-band
        log.Printf("[ERROR][terraform-provider-hosts] cannot find hosts-record\n")
        return fmt.Errorf("[ERROR][terraform-provider-hosts/hosts/resourceHostsRecordCreate] cannot find hosts-record")
    }

    // set identifying fields - need to set this in 'create' so we can lookup in 'read'
    _ = d.Set("record_id", record.ID)
    _ = d.Set("names", record.Names)

    // set id
    d.SetId(record.Names[0])

    log.Printf("[INFO][terraform-provider-hosts] created hosts-record %#v\n", record.ID)
    return resourceHostsRecordRead(d, m)
}

func resourceHostsRecordRead(d *schema.ResourceData, m interface{}) error {
    zone := m.(*api.Zone)
    recordID := d.Get("record_id").(int)

    // the recordID is usually sufficient to lookup a record
    // but names are needed to check the recordID hasn't changed since recordID is not persistent
    ns := d.Get("names").([]interface {})
    names := make([]string, len(ns))
    for i, _ := range ns {
        names[i] = ns[i].(string)
    }

    log.Printf(`[INFO][terraform-provider-hosts] reading hosts-record %#v
                    [INFO][terraform-provider-hosts]     zone:    %#v
                    [INFO][terraform-provider-hosts]     names:   %#v
`   , recordID, zone.Name, names)

    rQuery := new(api.Record)
    rQuery.ID    = recordID
    rQuery.Zone  = zone.ID
    rQuery.Names = names
    r := api.LookupRecord(rQuery)
    if r == nil {
        // this may be because the recordID changed vs previous terraform operation - retry without ID
        log.Printf("[ERROR][terraform-provider-hosts] try finding new record_id for hosts-record %#v\n", recordID)
        rQuery.ID = 0
        for _, n := range names {
            rQuery.Names = []string{ n }
            r = api.LookupRecord(rQuery)
            if r != nil {
                break
            }
        }
    }
    if r == nil {
        // this is most probably because 
        // - the record was deleted out-of-band
        // - a record with the same indexed fields was added out-of-band
        log.Printf("[WARNING][terraform-provider-hosts] cannot find hosts-record %#v\n", recordID)
        d.SetId("")

        log.Printf("[INFO][terraform-provider-hosts] deleted hosts-record %#v\n", recordID)
        return nil   // don't return an error to allow terraform refresh to update state
    }

    record, err := r.Read()
    if err != nil {
        // this is most probably because the hosts-file became inaccessible for reading
        log.Printf("[ERROR][terraform-provider-hosts] cannot read hosts-record %#v\n", recordID)
        return err
    }

    // set fields
    _ = d.Set("record_id", record.ID)
    _ = d.Set("address", record.Address)
    _ = d.Set("names", record.Names)
    _ = d.Set("comment", record.Comment)
    _ = d.Set("notes", record.Notes)

    if recordID != r.ID {
        log.Printf("[INFO][terraform-provider-hosts] read hosts-record %#v - found new record_id: %#v\n", recordID, r.ID)
    } else {
        log.Printf("[INFO][terraform-provider-hosts] read hosts-record %#v\n", recordID)
    }
    return nil
}

func resourceHostsRecordUpdate(d *schema.ResourceData, m interface{}) error {
    zone := m.(*api.Zone)
    recordID := d.Get("record_id").(int)
    comment := d.Get("comment").(string)
    notes := d.Get("notes").(string)

    log.Printf(`[INFO][terraform-provider-hosts] updating hosts-record %#v
                    [INFO][terraform-provider-hosts]     zone:    %#v
                    [INFO][terraform-provider-hosts]     comment: %#v
                    [INFO][terraform-provider-hosts]     notes:   %#v
`   , recordID, zone.Name, comment, notes)

    rQuery := new(api.Record)
    rQuery.ID = recordID
    rQuery.Zone = zone.ID
    r := api.LookupRecord(rQuery)
    if r == nil {
        // this is most probably because 
        // - the record was deleted out-of-band
        log.Printf("[ERROR][terraform-provider-hosts] cannot find hosts-record %#v\n", recordID)
        return fmt.Errorf("[ERROR][terraform-provider-hosts/hosts/resourceHostsRecordUpdate] cannot find hosts-record [id=%s]", d.Id())
    }

    rValues := new(api.Record)
    rValues.Comment = comment
    rValues.Notes   = notes
    err := r.Update(rValues)
    if err != nil {
        // this is most probably because the hosts-file became inaccessible for writing - perhaps reading still possible
        log.Printf("[ERROR][terraform-provider-hosts] cannot update hosts-record %#v\n", recordID)
        return err
    }

    log.Printf("[INFO][terraform-provider-hosts] updated hosts-record %#v\n", recordID)
    return resourceHostsRecordRead(d, m)
}

func resourceHostsRecordDelete(d *schema.ResourceData, m interface{}) error {
    zone := m.(*api.Zone)
    recordID := d.Get("record_id").(int)

    log.Printf(`[INFO][terraform-provider-hosts] deleting hosts-record %#v
                    [INFO][terraform-provider-hosts]     zone: %#v
`   , recordID, zone.Name)

    rQuery := new(api.Record)
    rQuery.ID = recordID
    rQuery.Zone = zone.ID
    r := api.LookupRecord(rQuery)
    if r == nil {
        // this is most probably because 
        // - the record was deleted out-of-band
        log.Printf("[WARNING][terraform-provider-hosts] cannot find hosts-record %#v\n", recordID)
        d.SetId("")

        log.Printf("[INFO][terraform-provider-hosts] deleted hosts-record %#v\n", recordID)
        return nil
    }

    err := r.Delete()
    if err != nil {
        // this is most probably because the hosts-file became inaccessible for writing - perhaps reading still possible
        log.Printf("[ERROR][terraform-provider-hosts] cannot delete hosts-record %#v\n", recordID)
        return err
    }

    // set id
    d.SetId("")

    log.Printf("[INFO][terraform-provider-hosts] deleted hosts-record %#v\n", recordID)
    return nil
}
