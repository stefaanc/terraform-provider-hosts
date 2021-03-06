//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package hosts

import (
    "runtime"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
    return &schema.Provider {
        Schema: map[string]*schema.Schema {
            "file": {
                Description: "The path to the hosts-file",
                Type:        schema.TypeString,
                Optional:    true,
                DefaultFunc: func() (interface{}, error) {
                    if runtime.GOOS == "windows" {
                        return "C:\\Windows\\System32\\drivers\\etc\\hosts", nil
                    } else {
                        return "/etc/hosts", nil
                    }
                },
            },
            "zone": {
                Description: "The zone in the hosts-file",
                Type:        schema.TypeString,
                Optional:    true,
                Default:     "external",
            },
        },

        DataSourcesMap: map[string]*schema.Resource {
            "hosts_record": dataSourceHostsRecord(),
        },

        ResourcesMap: map[string]*schema.Resource {
            "hosts_record": resourceHostsRecord(),
        },

        ConfigureFunc: providerConfigure,
    }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
    config := Config{
        file: d.Get("file").(string),
        zone: d.Get("zone").(string),
    }

    return config.Client()
}
