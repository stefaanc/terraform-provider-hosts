//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "log"

//    "github.com/stefaanc/terraform-provider-hosts/api"
)

type Config struct {
    File string
}

func (c *Config) Client() (interface{}, error) {
    log.Printf("[INFO][hosts] configuring hosts-provider with:\n")
    log.Printf("[INFO][hosts]     File: %s\n", c.File)

    // file := new(api.File)
    // file.Path = c.File
    
    // hosts := new(api.Hosts)
    // hosts.File = file

    // log.Printf("[INFO][hosts] reading hosts-file '%s'\n", hosts.File)
    // h, err := api.ReadHosts(hosts)
    // if err != nil {
    //     log.Printf("[WARNING][hosts] cannot read hosts-file, error:\n")
    //     log.Printf("[WARNING][hosts]     Error: '%s'\n", err.Error())
    //     log.Printf("[INFO][hosts] creating an empty hosts-object\n")
    //     h, err = api.CreateHosts(hosts)
    // }

    log.Printf("[INFO][hosts] configured hosts-provider\n")

    // return h, err
    return nil, nil
}