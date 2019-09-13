package api

import (
)

// -----------------------------------------------------------------------------

type Zone struct {
    // readOnly
    ID    zoneID
    // read-writeOnce
    File  string
    Name  string
    // read-write
    Notes string
    // private
    id    zoneID   // copy to silently enforce 'readOnly'
    file  string   // copy to silently enforce 'read-writeOnce'
    name  string   // copy to silently enforce 'read-writeOnce'
}

// -----------------------------------------------------------------------------
