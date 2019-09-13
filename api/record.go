package api

import (
)

// -----------------------------------------------------------------------------

type Record struct {
    // read-write
    Address    string
    Names      []string
    Comment    string
    // read-only
    ID         recordID
    FQDN       string
    Domain     string
    RootDomain string
    // private
    id         recordID
    isManaged  bool
}

// -----------------------------------------------------------------------------
