terraform {
    required_version = ">= 0.12.9"

    backend "local" {
        path = "./_terraform.tfstate"
    }
}