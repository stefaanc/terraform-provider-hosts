#
# Copyright (c) 2019 Stefaan Coussement
# MIT License
provider "hosts" {
    version = "~> 0.0.0"

    file = "playground/hosts-test.txt"
}

data "hosts_record" "myhost1" {
    query_name = "myhost6"
}
