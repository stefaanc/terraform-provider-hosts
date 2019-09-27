###
### myhost1
###

provider "hosts" {
    version = "~> 0.0.0"
    alias = "external"

    file = "./hosts-test.txt"
    // zone = "external"
}

data "hosts_record" "myhost1" {
    provider = hosts.external

    name = "myhost1"
}

###
### myhost111
###

provider "hosts" {
    version = "~> 0.0.0"
    alias = "myzone"

    file = "./hosts-test.txt"
    zone = "myzone"
}

resource "hosts_record" "myhost111" {
    provider = hosts.myzone

    address = "111.111.111.111"
    names   = [ "myHost111", "myHost111.local" ]
    comment = "server myhost111"
    notes   = "a first test-server"
}
