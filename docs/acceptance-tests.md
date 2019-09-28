## Acceptance Testing

An overview of the CRUD methods called for the "terraform refresh" and "terraform apply" actions
<br/>

resource                                                                                 | &nbsp; | terraform refresh       | terraform apply
:----------------------------------------------------------------------------------------|--------|:------------------------|:---------------------
                                                  - not in infrastructure                | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp;                                             -- not in terraform state            | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp;                                  --- not in terraform config       |      - | ---                     | ---
 &nbsp; &nbsp; &nbsp;                                  --- in terraform config           |      1 | ---                     | Create, Read
 &nbsp;                                             -- in terraform state                | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp;                                  --- not in terraform config       |      2 | Read (clear state)      |                      &nbsp;
 &nbsp; &nbsp; &nbsp;                                  --- in terraform config           |      3 | Read (clear state)      | Create, Read
                                                  - in infrastructure                    | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp;                                             -- not in terraform state            | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp;                                  --- not in terraform config       |      - | ---                     | ---
 &nbsp; &nbsp; &nbsp;                                  --- in terraform config           |      4 | ---                     | Create (error)
 &nbsp;                                             -- in terraform state                | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp;                                  --- not in terraform config       |      5 | Read (refresh state)    | Delete
 &nbsp; &nbsp; &nbsp;                                  --- in terraform config           | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;                        > config same as state        |      6 | Read (refresh state)    | ---
 &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;            >> changed computed fields  |      7 | Read (refresh state)    | ---
 &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;                        > config diff from state      | &nbsp; |                  &nbsp; |                      &nbsp;
 &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;            >> changed force-new fields |      8 | Read (refresh state)    | Delete, Create, Read
 &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;            >> changed common fields    |      9 | Read (refresh state)    | Update, Read
