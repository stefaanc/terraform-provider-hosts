#
# Copyright (c) 2019 Stefaan Coussement
# MIT License
#
# more info: https://github.com/stefaanc/kluster
#

terraform {
    required_version = ">= 0.12"

    backend "local" {
        path = "./state/terraform.tfstate"
    }
}