# Terraform Provider Hosts

**a terraform provider to work with a hosts-file**

#!!! UNDER CONSTRUCTION !!!

## Prerequisites

To build:
- [GNU make](https://www.gnu.org/software/make/)
- [Go](https://golang.org/) >= v1.12

To use:
- [Terraform](https://terraform.io) >= v0.12.7

Optional:
- [PSConsole](https://github.com/stefaanc/psconsole) for PowerShell


<br>

## Building The Provider

1. Clone the git-repository on your machine

   ```shell
   mkdir -p $my_repositories
   cd $my_repositories
   git clone git@github.com:stefaanc/terraform-provider-hosts
   ```

2. Build the provider

   ```shell
   cd $my_repositories/terraform-provider/hosts
   make build
   ```

   This will build the provider and put it in 
   - `%AppData%\terraform.d\plugins` on Windows
   - `$HOME\.terraform.d\plugins` on Linux
    


<br>

## Installing The Provider

1. Download the provider to your machine

   - go to [the releases tab on github](https://github.com/stefaanc/terraform-provider-hosts/releases)
   - download the file that is appropriate for your machine

2. Move the provider from your `Downloads` folder to

   - `%AppData%\terraform.d\plugins` on Windows
   - `$HOME\.terraform.d\plugins` on Linux



<br>

## Using The Provider

### provider "hosts"

```terraform
provider "hosts" {
    file  = "my-path/hosts"
}
```

arguments | optional / required | description
----------|:--------:|------------
`file`    | optional | The path to the `hosts`-file <br/>- defaults to `C:\Windows\System32\drivers\etc\hosts` on Windows or `/etc/hosts` on Linux

<br>

### Data-sources

#### data "hosts_record"

Reads a record from the hosts-file.  Records in a hosts-file look something like `192.168.0.42   myhost myhost.example.com myalias1 myalias2 # mycomment`

```terraform
data "hosts_record" "myhost" {
    query {
        hostname = "myhost"
    }
}
```

arguments  | optional / required | description
-----------|:--------:|------------
`query`    | required | The properties of the record that will be read

* The data-source will return an error when no record found.  

query keys | optional / required | description
-----------|:--------:|------------
`address`  | optional | The `address` of the record that will be read. 
`hostname` | optional | The `hostname` of the record that will be read. 
`name`     | optional | The `hostname` or `alias` of the record that will be read. 
`domain`   | optional | The `domain` of the record that will be read. 

* the keys `hostname` and `name` are mutually exclusive.
* other keys are not allowed. 

> :bulb:  
> Remark that it is perfectly legal to have multiple records with the same `address`, but it is illegal to have multiple records with the same name (`hostname` or `alias`).  The terraform `"hosts_record"`-resource doesn't allow to create records with such conflicting names, but an externally managed hosts-file may allow it.  
> 
> When the hosts-file contains multiple matching records, the first that is found will be returned.  
> - The hosts-file is traversed line-by-line from top to bottom.   
> - When a name (`query["name"]`) is provided, the first name on every line (the `hostname`) is checked in a first pass.  In a second pass, the other names on every line are checked (the `aliases`), from the second name to the last name.  

exports    | description
-----------|------------
`address`  | The IP address of the host
`hostname` | The hostname of the host
`domain`   | The domain of the host
`aliases`  | A list of aliases for the host, including `"hostname.domain"` when the hostname is not an FQDN
`comment`  | The comment for the record

> :bulb:  
> When the matching record has no aliases of the form `"hostname.domain"`and when the `"hostname"` doesn't contain a `"."`, the domain will be set to `"localdomain"`  
> When the matching record has no aliases of the form `"hostname.domain"`and when the `"hostname"` contains a `"."`, the domain will be set to substring starting from the first character after the first `"."`  
> When the matching record has an alias of the form `"hostname.domain"`, the domain will be set to `"domain"`  
> When the matching record has multiple aliases of the form `"hostname.domain"`, the domain will be set to `"domain"` of the first such alias.  Aliases are checked from the second name on the hosts-line to the last name.  

<br>

### Resources

#### resource "hosts_record"   

Manages a record in the hosts-file.  Records in a hosts-file look something like `192.168.0.42   my-host myhost.example.com myalias1 myalias2 # mycomment`

```terraform
resource "hosts_record" "myhost" {
    address = "192.168.0.42"
    hostname = "myhost"
    domain = "example.com"
    aliases = [ "myalias1", "myalias2" ]
    comment = "mycomment"
}
```

arguments  | optional / required | description
-----------|:--------:|------------
`address`  | required | The IP address of the host to create
`hostname` | required | The hostname of the host to create
`domain`   | optional | The domain of the host to create <br>- defaults to `"localdomain"` when the `"hostname"` doesn't contain a `"."`, otherwise defaults to the substring starting from the first character after the first `"."`
`aliases`  | optional | A list of aliases for the host to create <br/>- defaults to an empty list <br/>- when the `hostname` is not an FQDN,  `"hostname.domain"` is automatically added to the list
`comment`  | optional | The comment for the host to create



<br>

## For Further Investigation

- working with `hosts`-files on remote servers