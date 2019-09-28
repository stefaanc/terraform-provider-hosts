## Resource Lifecycle Methods

An overview of a resource's lifecycle methods called for the "terraform refresh" and "terraform apply" actions

<br/>

### The Create, Read & Delete Methods

resource                                                                 | &nbsp; | terraform refresh           | terraform apply
:------------------------------------------------------------------------|--------|:----------------------------|:---------------------
 &nbsp;                           - not in infrastructure                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                         | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      1 | ---                         | Create, Read
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      2 | Read (clear state)          |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      3 | Read (clear state)          | Create, Read
 &nbsp;                           - in infrastructure                    | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                         | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      4 | ---                         | Create (error)
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      5 | Read (refresh state)        | Delete
 &nbsp; &emsp; &emsp;                  --- in terraform config           | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp;               > config same as state        |      6 | Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed computed fields  |      7 | Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp;               > config different from state | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed force-new fields |      8 | Read (refresh state)        | Delete, Create, Read
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed common fields    |      - | n.a.                        | n.a.

> :information_source:  
> `ForceNew: true,` **must** be set for all fields in the schema, except the `Computed` fields

<br/>

### Adding The Update Method

> :bulb:  
> `terraform-provider-hosts` **DOES** support the Update method

resource                                                                 | &nbsp; | terraform refresh           | terraform apply
:------------------------------------------------------------------------|--------|:----------------------------|:---------------------
 &nbsp;                           - not in infrastructure                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                         | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      1 | ---                         | Create, Read
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      2 | Read (clear state)          |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      3 | Read (clear state)          | Create, Read
 &nbsp;                           - in infrastructure                    | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                         | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      4 | ---                         | Create (error)
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      5 | Read (refresh state)        | Delete
 &nbsp; &emsp; &emsp;                  --- in terraform config           | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp;               > config same as state        |      6 | Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed computed fields  |      7 | Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp;               > config different from state | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed force-new fields |      8 | Read (refresh state)        | Delete, Create, Read
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed common fields    |      9 | **Read (refresh state)**    | **Update, Read**

<br/>

### Adding The Exists Method

> :bulb:  
> `terraform-provider-hosts` **DOES NOT** supports the Update method

resource                                                                 | &nbsp; | terraform refresh                       | terraform apply
:------------------------------------------------------------------------|--------|:----------------------------------------|:---------------------
 &nbsp;                           - not in infrastructure                | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                                     | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      1 | ---                                     | Create, Read
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      2 | **Exists** (clear state)                |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      3 | **Exists** (clear state)                | Create, Read
 &nbsp;                           - in infrastructure                    | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp;                      -- not in terraform state            | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      - | ---                                     | ---
 &nbsp; &emsp; &emsp;                  --- in terraform config           |      4 | ---                                     | Create (error)
 &nbsp; &emsp;                      -- in terraform state                | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp;                  --- not in terraform config       |      5 | **Exists**, Read (refresh state)        | Delete
 &nbsp; &emsp; &emsp;                  --- in terraform config           | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp;               > config same as state        |      6 | **Exists**, Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed computed fields  |      7 | **Exists**, Read (refresh state)        | ---
 &nbsp; &emsp; &emsp; &emsp;               > config different from state | &nbsp; |                                  &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed force-new fields |      8 | **Exists**, Read (refresh state)        | Delete, Create, Read
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed common fields    |      9 | **Exists**, Read (refresh state)        | Update, Read
