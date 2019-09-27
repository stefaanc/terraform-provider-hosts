//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "log"

    "github.com/stefaanc/terraform-provider-hosts/api"
)

type Config struct {
    file string
    zone string
}

func (c *Config) Client() (interface{}, error) {
    log.Printf(`[INFO][terraform-provider-hosts] configuring hosts-provider
                    [INFO][terraform-provider-hosts]     file: %q
                    [INFO][terraform-provider-hosts]     zone: %q
`   , c.file, c.zone)

    api.Init()

    fValues := new(api.File)
    fValues.Path = c.file
    f := api.LookupFile(fValues)
    if f == nil {
        err := api.CreateFile(fValues)
        if err != nil {
            return nil, err
        }
        f = api.LookupFile(fValues)
    }

    zValues := new(api.Zone)
    zValues.File = f.ID
    zValues.Name = c.zone
    z := api.LookupZone(zValues)
    if z == nil {
        err := api.CreateZone(zValues)
        if err != nil {
            return nil, err
        }
        z = api.LookupZone(zValues)
    }

    log.Printf("[INFO][terraform-provider-hosts] configured hosts-provider\n")
    return z, nil
}