## Acceptance Testing

An overview of the CRUD methods called for the "terraform refresh" and "terraform apply" actions
<br/>

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
 &nbsp; &emsp; &emsp; &emsp;               > config diff from state      | &nbsp; |                      &nbsp; |                      &nbsp;
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed force-new fields |      8 | Read (refresh state)        | Delete, Create, Read
 &nbsp; &emsp; &emsp; &emsp; &emsp;          >> changed common fields    |      9 | Read (refresh state)        | Update, Read
