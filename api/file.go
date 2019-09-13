package api

import (
)

// -----------------------------------------------------------------------------

type File struct {
    // readOnly
    ID       fileID
    // read-writeOnce
    Path     string
    // read-write
    Notes    string
    // private
    id       fileID   // copy to silently enforce 'readOnly'
    path     string   // copy to silently enforce 'read-writeOnce'
    data     []byte
    checksum string
}

// -----------------------------------------------------------------------------
