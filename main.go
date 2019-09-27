//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/plugin"
    "github.com/stefaanc/terraform-provider-hosts/hosts"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: hosts.Provider,
    })
}