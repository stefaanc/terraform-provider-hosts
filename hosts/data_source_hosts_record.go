//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "errors"
    "log"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

    "github.com/stefaanc/terraform-provider-hosts/api"
)

func dataSourceHostsRecord() *schema.Resource {
    return &schema.Resource {
        Read:   dataSourceHostsRecordRead,

        Schema: map[string]*schema.Schema {
            "name": &schema.Schema {
                Type:     schema.TypeString,
                StateFunc: func(val interface{}) string {
                    return strings.ToLower(val.(string))
                },
                DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
                    if old == strings.ToLower(new) {
                        return true 
                    }
                    return false
                },
                Required: true,
                ForceNew: true,
            },
            
            "record_id": &schema.Schema {
                Type:     schema.TypeInt,
                Computed: true,
            },
            "address": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "names": &schema.Schema {
                Type:     schema.TypeList,
                Elem:     &schema.Schema {
                    Type: schema.TypeString,
                },
                Computed: true,
            },
            "comment": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "notes": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceHostsRecordRead(d *schema.ResourceData, m interface{}) error {
    zone := m.(*api.Zone)
    name := d.Get("name").(string)

    log.Printf(`[INFO][terraform-provider-hosts] reading hosts-record %#v
                    [INFO][terraform-provider-hosts]     zone: %#v
`   , name, zone.Name)

    rQuery := new(api.Record)
    rQuery.Zone = zone.ID
    rQuery.Names = []string{ name }
    r := api.LookupRecord(rQuery)
    if r == nil {
        d.SetId("")
        log.Printf("[ERROR][terraform-provider-hosts] cannot find hosts-record %#v\n", name)
        return errors.New("[ERROR][terraform-provider-hosts/hosts/dataSourceHostsRecordRead] cannot find hosts-record")
    }

    record, err := r.Read()
    if err != nil {
        log.Printf("[ERROR][terraform-provider-hosts] cannot read hosts-record %#v\n", name)
        return err
    }

    // set computed fields
    _ = d.Set("record_id", record.ID)
    _ = d.Set("address", record.Address)
    _ = d.Set("names", record.Names)
    _ = d.Set("comment", record.Comment)
    _ = d.Set("notes", record.Notes)

    // set id
    d.SetId(name)

    log.Printf("[INFO][terraform-provider-hosts] read hosts-record %#v\n", name)
    return nil
}
