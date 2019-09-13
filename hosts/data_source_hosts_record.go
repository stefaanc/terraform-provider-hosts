//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "log"
    "strings"

    "github.com/hashicorp/terraform/helper/schema"

//    "github.com/stefaanc/terraform-provider-hosts/api"
)

func dataSourceHostsRecord() *schema.Resource {
    return &schema.Resource {
        Read:   dataSourceHostsRecordRead,

        Schema: map[string]*schema.Schema {
            "query_name": &schema.Schema {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "address": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "name": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "aliases": &schema.Schema {
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

            "host": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "domain": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "rootdomain": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
            "fqdn": &schema.Schema {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func dataSourceHostsRecordRead(d *schema.ResourceData, m interface{}) error {
//    h := m.(*api.Hosts)
    query_name := strings.ToLower(d.Get("query_name").(string))

    log.Printf("[INFO][hosts] reading hosts-record for '%s'\n", query_name)

//    record := new(api.Record)
//    record.Name = query_name

//    r, err := api.ReadRecord(h, record)
//    if err != nil {
//        log.Printf("[ERROR][hosts] cannot read hosts-record, error:\n")
//        log.Printf("[WARNING][hosts]     Error: '%s'\n", err.Error())
//        return err
//    }

    // computed fields
//    d.Set("address", r.Address)
//    d.Set("name", r.Name)
//    d.Set("aliases", r.Aliases)
//    d.Set("comment", r.Comment)
//    d.Set("host", r.Host)
//    d.Set("domain", r.Domain)
//    d.Set("rootdomain", r.Rootdomain)
//    d.Set("fqdn", r.FQDN)

    // id
//    d.SetId(r.Name)

    log.Printf("[INFO][hosts] read hosts-record\n")

    return nil
//    return err
}
