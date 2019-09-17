//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "log"

    hosts "github.com/stefaanc/terraform-provider-hosts/api"
)

type Config struct {
    File string
    Zone string
}

func (c *Config) Client() (interface{}, error) {
    log.Printf("[INFO][terraform-provider-hosts] configuring hosts-provider for:\n")
    log.Printf("[INFO][terraform-provider-hosts]     File: %q\n", c.File)
    log.Printf("[INFO][terraform-provider-hosts]     Zone: %q\n", c.Zone)

    file := c.File
    zone := c.Zone

    hosts.Init()

    fValues := new(hosts.File)
    fValues.Path = file
    f := hosts.LookupFile(fValues)
    if f == nil {
        err := hosts.CreateFile(fValues)
        if err != nil {
            return nil, err
        }
        f = hosts.LookupFile(fValues)
    }

    zValues := new(hosts.Zone)
    zValues.File = f.ID
    zValues.Name = zone
    z := hosts.LookupZone(zValues)
    if z == nil {
        err := hosts.CreateZone(zValues)
        if err != nil {
            return nil, err
        }
        z = hosts.LookupZone(zValues)
    }

    log.Printf("[INFO][terraform-provider-hosts] configured hosts-provider\n")

    return z, nil
}