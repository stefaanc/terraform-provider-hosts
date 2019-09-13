package api

import (
    "testing"
)

// -----------------------------------------------------------------------------

func Test_InitHosts(t *testing.T) {
    var test string

    test = "1st-init"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        // --------------------

        if hosts == nil {
            t.Errorf("[ hosts ] expected: not <nil>, actual: <nil>")
        }

        // --------------------

        fid := hosts.newFileID()
        if fid != 1 {
            t.Errorf("[ 1st hosts.newFileID() ] expected: 1, actual: %d", fid)
        }

        fid = hosts.newFileID()
        if fid != 2 {
            t.Errorf("[ 2nd hosts.newFileID() ] expected: 2, actual: %d", fid)
        }

        length := len(hosts.fileIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.fileIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.filePaths)
        if length != 0 {
            t.Errorf("[ len(hosts.filePaths) ] expected: 0, actual: %d", length)
        }

        // --------------------

        zid := hosts.newZoneID()
        if zid != 1 {
            t.Errorf("[ 1st hosts.newZoneID() ] expected: 1, actual: %d", zid)
        }

        zid = hosts.newZoneID()
        if zid != 2 {
            t.Errorf("[ 2nd hosts.newZoneID() ] expected: 2, actual: %d", zid)
        }

        length = len(hosts.zoneIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.zoneFiles)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneFiles) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.zoneNames)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneNames) ] expected: 0, actual: %d", length)
        }

        // --------------------

        rid := hosts.newRecordID()
        if rid != 1 {
            t.Errorf("[ 1st hosts.newRecordID() ] expected: 1, actual: %d", rid)
        }

        rid = hosts.newRecordID()
        if rid != 2 {
            t.Errorf("[ 2nd hosts.newRecordID() ] expected: 2, actual: %d", rid)
        }

        length = len(hosts.recordIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.recordIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.recordNames)
        if length != 0 {
            t.Errorf("[ len(hosts.recordNames) ] expected: 0, actual: %d", length)
        }
    })

    test = "2nd-init"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        hosts.fileIndex[fid] = f
        hosts.filePaths["p"] = append(hosts.filePaths["p"], f)

        z := new(Zone)
        zid := hosts.newZoneID()
        hosts.zoneIndex[zid] = z
        hosts.zoneFiles["f"] = append(hosts.zoneFiles["f"], z)
        hosts.zoneNames["n"] = append(hosts.zoneNames["n"], z)

        r := new(Record)
        rid := hosts.newRecordID()
        hosts.recordIndex[rid] = r
        hosts.recordNames["n"] = append(hosts.recordNames["n"], r)

        InitHosts()

        // --------------------

        if hosts == nil {
            t.Errorf("[ hosts ] expected: not <nil>, actual: <nil>")
        }

        // --------------------

        fid = hosts.newFileID()
        if fid != 1 {
            t.Errorf("[ hosts.newFileID() ] expected: 1, actual: %d", fid)
        }

        length := len(hosts.fileIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.fileIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.filePaths)
        if length != 0 {
            t.Errorf("[ len(hosts.filePaths) ] expected: 0, actual: %d", length)
        }

        // --------------------

        zid = hosts.newZoneID()
        if zid != 1 {
            t.Errorf("[ hosts.newZoneID() ] expected: 1, actual: %d", zid)
        }

        length = len(hosts.zoneIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.zoneFiles)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneFiles) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.zoneNames)
        if length != 0 {
            t.Errorf("[ len(hosts.zoneNames) ] expected: 0, actual: %d", length)
        }

        // --------------------

        rid = hosts.newRecordID()
        if rid != 1 {
            t.Errorf("[ hosts.newRecordID() ] expected: 1, actual: %d", rid)
        }

        length = len(hosts.recordIndex)
        if length != 0 {
            t.Errorf("[ len(hosts.recordIndex) ] expected: 0, actual: %d", length)
        }

        length = len(hosts.recordNames)
        if length != 0 {
            t.Errorf("[ len(hosts.recordNames) ] expected: 0, actual: %d", length)
        }
    })
}

// -----------------------------------------------------------------------------

func Test_getFile(t *testing.T) {
    var test string

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        f.id = hosts.newFileID()

        file := hosts.getFile(f)

        // --------------------

        if file == nil {
            t.Errorf("[ hosts.getFile(f) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ hosts.getFile(f) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/no-files"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        q := new(File)
        q.ID = fileID(42)

        file := hosts.getFile(q)

        // --------------------

        if file != nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.ID = fid

        file := hosts.getFile(q)

        // --------------------

        if file == nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.ID = fileID(42)

        file := hosts.getFile(q)

        // --------------------

        if file != nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.Path = "p"

        file := hosts.getFile(q)

        // --------------------

        if file == nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-Path/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.Path = "x"

        file := hosts.getFile(q)

        // --------------------

        if file != nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        f = new(File)
        fid = hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.Path = "p"

        file := hosts.getFile(q)

        // --------------------

        if file != nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/ID-and-Path/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.ID = fid
        q.Path = "p"

        file := hosts.getFile(q)

        // --------------------

        if file == nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/ID-and-Path/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        fid := hosts.newFileID()
        f.id = fid
        f.path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.path] = append(hosts.filePaths[f.path], f)

        q := new(File)
        q.ID = fid
        q.Path = "x"

        file := hosts.getFile(q)

        // --------------------

        if file != nil {
            t.Errorf("[ hosts.getFile(q) ] expected: %#v, actual: %#v", nil, file)
        }
    })
}

func Test_addFile(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        f.Path = "p"
        f.path = "p"

        err := hosts.addFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addFile(f) ] expected: %#v, actual: %#v", nil, err)
        }

        if f.ID != 1 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 1, f.ID)
        }

        if f.id != 1 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 1, f.id)
        }

        // --------------------

        q := new(File)
        q.ID = f.id

        file := hosts.getFile(q)
        if file == nil {
            t.Errorf("[ hosts.getFile(q.ID) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ hosts.getFile(q.ID) ] expected: %#v, actual: %#v", f, file)
        }

        // --------------------

        q = new(File)
        q.Path = f.path

        file = hosts.getFile(q)
        if file == nil {
            t.Errorf("[ hosts.getFile(q.Path) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ hosts.getFile(q.Path) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        f.Path = "p"
        f.path = "p"

        _ = hosts.addFile(f)
        err := hosts.addFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        f.Path = "p"
        f.path = "p"

        _ = hosts.addFile(f)

        f = new(File)
        f.Path = "p"
        f.path = "p"
        err := hosts.addFile(f)

        // --------------------

        if err == nil {
            t.Errorf("[ hosts.addFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })
}

func Test_removeFile(t *testing.T) {
    var test string

    test = "removed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)
        f.Path = "p"
        f.path = "p"

        _ = hosts.addFile(f)
        hosts.removeFile(f)

        // --------------------

        if f.ID != 0 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 0, f.ID)
        }

        if f.id != 0 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 0, f.id)
        }

        // --------------------

        q := new(File)
        q.ID = f.id

        file := hosts.getFile(q)
        if file != nil {
            t.Errorf("[ hosts.getFile(q.ID) ] expected: %#v, actual: %#v", nil, file)
        }

        // --------------------

        q = new(File)
        q.Path = f.path

        file = hosts.getFile(q)
        if file != nil {
            t.Errorf("[ hosts.getFile(q.Path) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        f := new(File)

        hosts.removeFile(f)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })
}

func Test_deleteFromSliceOfFiles(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*File, 0)

        f := new(File)
        f.id = hosts.newFileID()

        _ = deleteFromSliceOfFiles(s1, f)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*File, 0)

        f := new(File)
        f.id = hosts.newFileID()
        s1 = append(s1, f)

        s2 := deleteFromSliceOfFiles(s1, f)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }
    })

    test = "more-elements/first-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*File, 0)

        f1 := new(File)
        f1.id = hosts.newFileID()
        s1 = append(s1, f1)

        f2 := new(File)
        f2.id = hosts.newFileID()
        s1 = append(s1, f2)

        f3 := new(File)
        f3.id = hosts.newFileID()
        s1 = append(s1, f3)

        s2 := deleteFromSliceOfFiles(s1, f1)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == f1.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })

    test = "more-elements/middle-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*File, 0)

        f1 := new(File)
        f1.id = hosts.newFileID()
        s1 = append(s1, f1)

        f2 := new(File)
        f2.id = hosts.newFileID()
        s1 = append(s1, f2)

        f3 := new(File)
        f3.id = hosts.newFileID()
        s1 = append(s1, f3)

        s2 := deleteFromSliceOfFiles(s1, f2)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == f2.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
    
    test = "more-elements/last-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*File, 0)

        f1 := new(File)
        f1.id = hosts.newFileID()
        s1 = append(s1, f1)

        f2 := new(File)
        f2.id = hosts.newFileID()
        s1 = append(s1, f2)

        f3 := new(File)
        f3.id = hosts.newFileID()
        s1 = append(s1, f3)

        s2 := deleteFromSliceOfFiles(s1, f3)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == f3.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
}

// -----------------------------------------------------------------------------

func Test_getZone(t *testing.T) {
    var test string

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        z.id = hosts.newZoneID()

        zone := hosts.getZone(z)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(z) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(z) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/no-zones"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        q := new(Zone)
        q.ID = zoneID(42)

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zid

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zoneID(42)

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.File = "f"

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-File/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.File = "x"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z1"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        z = new(Zone)
        zid = hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z2"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.File = "f"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.Name = "z"

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-Name/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.Name = "x"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f1"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        z = new(Zone)
        zid = hosts.newZoneID()
        z.id = zid
        z.file = "f2"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.Name = "z"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-File/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zid
        q.File = "f"

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-File/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zid
        q.File = "x"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-Name/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zid
        q.Name = "z"

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-Name/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.ID = zid
        q.Name = "x"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/Name-and-File/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.File = "f"
        q.Name = "z"

        zone := hosts.getZone(q)

        // --------------------

        if zone == nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/Name-and-File/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.file = "f"
        z.name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.file] = append(hosts.zoneFiles[z.file], z)
        hosts.zoneNames[z.name] = append(hosts.zoneNames[z.name], z)

        q := new(Zone)
        q.File = "x"
        q.Name = "z"

        zone := hosts.getZone(q)

        // --------------------

        if zone != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, zone)
        }
    })
}

func Test_addZone(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        z.File = "f"
        z.file = "f"
        z.Name = "z"
        z.name = "z"

        err := hosts.addZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addZone(z) ] expected: %#v, actual: %#v", nil, err)
        }

        if z.ID != 1 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 1, z.ID)
        }

        if z.id != 1 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 1, z.id)
        }

        // --------------------

        q := new(Zone)
        q.ID = z.id

        zone := hosts.getZone(q)
        if zone == nil {
            t.Errorf("[ hosts.getZone(q.ID) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q.ID) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        q = new(Zone)
        q.File = z.file

        zone = hosts.getZone(q)
        if zone == nil {
            t.Errorf("[ hosts.getZone(q.File) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q.File) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        q = new(Zone)
        q.Name = z.name

        zone = hosts.getZone(q)
        if zone == nil {
            t.Errorf("[ hosts.getZone(q.Name) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ hosts.getZone(q.Name) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        z.File = "f"
        z.file = "f"
        z.Name = "z"
        z.name = "z"

        _ = hosts.addZone(z)
        err := hosts.addZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addZone(z) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        z.File = "f"
        z.file = "f"
        z.Name = "z"
        z.name = "z"

        _ = hosts.addZone(z)

        z = new(Zone)
        z.File = "f"
        z.file = "f"
        z.Name = "z"
        z.name = "z"
        err := hosts.addZone(z)

        // --------------------

        if err == nil {
            t.Errorf("[ hosts.addZone(z) ] expected: %#v, actual: %#v", "<error>", err)
        }
    })
}

func Test_removeZone(t *testing.T) {
    var test string

    test = "removed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)
        z.File = "f"
        z.file = "f"
        z.Name = "z"
        z.name = "z"

        _ = hosts.addZone(z)
        hosts.removeZone(z)

        // --------------------

        if z.ID != 0 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 0, z.ID)
        }

        if z.id != 0 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 0, z.id)
        }

        // --------------------

        q := new(Zone)
        q.ID = z.id

        zone := hosts.getZone(q)
        if zone != nil {
            t.Errorf("[ hosts.getZone(q.ID) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        q = new(Zone)
        q.File = z.file

        zone = hosts.getZone(q)
        if zone != nil {
            t.Errorf("[ hosts.getZone(q.File) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        q = new(Zone)
        q.Name = z.name

        zone = hosts.getZone(q)
        if zone != nil {
            t.Errorf("[ hosts.getZone(q.Name) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        z := new(Zone)

        hosts.removeZone(z)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })
}

func Test_deleteFromSliceOfZones(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Zone, 0)

        z := new(Zone)
        z.id = hosts.newZoneID()

        _ = deleteFromSliceOfZones(s1, z)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Zone, 0)

        z := new(Zone)
        z.id = hosts.newZoneID()
        s1 = append(s1, z)

        s2 := deleteFromSliceOfZones(s1, z)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }
    })

    test = "more-elements/first-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Zone, 0)

        z1 := new(Zone)
        z1.id = hosts.newZoneID()
        s1 = append(s1, z1)

        z2 := new(Zone)
        z2.id = hosts.newZoneID()
        s1 = append(s1, z2)

        z3 := new(Zone)
        z3.id = hosts.newZoneID()
        s1 = append(s1, z3)

        s2 := deleteFromSliceOfZones(s1, z1)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == z1.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })

    test = "more-elements/middle-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Zone, 0)

        z1 := new(Zone)
        z1.id = hosts.newZoneID()
        s1 = append(s1, z1)

        z2 := new(Zone)
        z2.id = hosts.newZoneID()
        s1 = append(s1, z2)

        z3 := new(Zone)
        z3.id = hosts.newZoneID()
        s1 = append(s1, z3)

        s2 := deleteFromSliceOfZones(s1, z2)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == z2.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
    
    test = "more-elements/last-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Zone, 0)

        z1 := new(Zone)
        z1.id = hosts.newZoneID()
        s1 = append(s1, z1)

        z2 := new(Zone)
        z2.id = hosts.newZoneID()
        s1 = append(s1, z2)

        z3 := new(Zone)
        z3.id = hosts.newZoneID()
        s1 = append(s1, z3)

        s2 := deleteFromSliceOfZones(s1, z3)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == z3.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
}

// -----------------------------------------------------------------------------

func Test_getRecord(t *testing.T) {
    var test string

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        r.id = hosts.newRecordID()

        record := hosts.getRecord(r)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(r) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(r) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/no-records"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        q := new(Record)
        q.ID = recordID(42)

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = rid

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = recordID(42)

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "a"

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Address/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "x"

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        r = new(Record)
        rid = hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n4", "n5", "n6" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "a"

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/1-name/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "n1" }

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Name/1-name/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "x" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Name/1-name/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        r = new(Record)
        rid = hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n4" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "n1" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/more-names/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "n1", "n3" }

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Name/more-names/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "n1", "x" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Name/more-names/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        r = new(Record)
        rid = hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n4" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Names = []string{ "n1", "n2" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getZone(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Address/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = rid
        q.Address = "a"

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Address/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = rid
        q.Address = "x"

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Names/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = rid
        q.Names = []string{ "n1", "n2" }

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Names/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.ID = rid
        q.Names = []string{ "x" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "a"
        q.Names = []string{ "n1", "n2" }

        record := hosts.getRecord(q)

        // --------------------

        if record == nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/Address-and-Names/not-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "a"
        q.Names = []string{ "x" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/multiple-found"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        rid := hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        r = new(Record)
        rid = hosts.newRecordID()
        r.id = rid
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n4" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        q := new(Record)
        q.Address = "a"
        q.Names = []string{ "n1", "n2" }

        record := hosts.getRecord(q)

        // --------------------

        if record != nil {
            t.Errorf("[ hosts.getRecord(q) ] expected: %#v, actual: %#v", nil, record)
        }
    })
}

func Test_addRecord(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        err := hosts.addRecord(r)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addRecord(r) ] expected: %#v, actual: %#v", nil, err)
        }

        if r.ID != 1 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 1, r.ID)
        }

        if r.id != 1 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 1, r.id)
        }

        // --------------------

        q := new(Record)
        q.ID = r.id

        record := hosts.getRecord(q)
        if record == nil {
            t.Errorf("[ hosts.getRecord(q.ID) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q.ID) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        q = new(Record)
        q.Address = r.Address

        record = hosts.getRecord(q)
        if record == nil {
            t.Errorf("[ hosts.getRecord(q.Address) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q.Address) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        q = new(Record)
        q.Names = r.Names

        record = hosts.getRecord(q)
        if record == nil {
            t.Errorf("[ hosts.getRecord(q.Names) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ hosts.getRecord(q.Names) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        _ = hosts.addRecord(r)
        err := hosts.addRecord(r)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addRecord(r) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        _ = hosts.addRecord(r)

        r = new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        err := hosts.addRecord(r)

        // --------------------

        if err == nil {
            t.Errorf("[ hosts.addRecord(r) ] expected: %#v, actual: %#v", "<error>", err)
        }
    })
}

func Test_removeRecord(t *testing.T) {
    var test string

    test = "removed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        _ = hosts.addRecord(r)
        hosts.removeRecord(r)

        // --------------------

        if r.ID != 0 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 0, r.ID)
        }

        if r.id != 0 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 0, r.id)
        }

        // --------------------

        q := new(Record)
        q.ID = r.id

        record := hosts.getRecord(q)
        if record != nil {
            t.Errorf("[ hosts.getRecord(q.ID) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        q = new(Record)
        q.Address = r.Address

        record = hosts.getRecord(q)
        if record != nil {
            t.Errorf("[ hosts.getRecord(q.Address) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        q = new(Record)
        q.Names = r.Names

        record = hosts.getRecord(q)
        if record != nil {
            t.Errorf("[ hosts.getRecord(q.Names) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        r := new(Record)

        hosts.removeRecord(r)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })
}

func Test_deleteFromSliceOfRecords(t *testing.T) {
    var test string

    test = "empty"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Record, 0)

        r := new(Record)
        r.id = hosts.newRecordID()

        _ = deleteFromSliceOfRecords(s1, r)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Record, 0)

        r := new(Record)
        r.id = hosts.newRecordID()
        s1 = append(s1, r)

        s2 := deleteFromSliceOfRecords(s1, r)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }
    })

    test = "more-elements/first-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Record, 0)

        r1 := new(Record)
        r1.id = hosts.newRecordID()
        s1 = append(s1, r1)

        r2 := new(Record)
        r2.id = hosts.newRecordID()
        s1 = append(s1, r2)

        r3 := new(Record)
        r3.id = hosts.newRecordID()
        s1 = append(s1, r3)

        s2 := deleteFromSliceOfRecords(s1, r1)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == r1.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })

    test = "more-elements/middle-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Record, 0)

        r1 := new(Record)
        r1.id = hosts.newRecordID()
        s1 = append(s1, r1)

        r2 := new(Record)
        r2.id = hosts.newRecordID()
        s1 = append(s1, r2)

        r3 := new(Record)
        r3.id = hosts.newRecordID()
        s1 = append(s1, r3)

        s2 := deleteFromSliceOfRecords(s1, r2)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == r2.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
    
    test = "more-elements/last-element"
    t.Run(test, func(t *testing.T) {

        InitHosts()

        s1 := make([]*Record, 0)

        r1 := new(Record)
        r1.id = hosts.newRecordID()
        s1 = append(s1, r1)

        r2 := new(Record)
        r2.id = hosts.newRecordID()
        s1 = append(s1, r2)

        r3 := new(Record)
        r3.id = hosts.newRecordID()
        s1 = append(s1, r3)

        s2 := deleteFromSliceOfRecords(s1, r3)

        // --------------------

        len1 := len(s1)
        len2 := len(s2)
        if len2 != len1-1 {
            t.Errorf("[ len(s2) ] expected: %#v, actual: %#v", len1-1, len2)
        }

        for _, element := range s2 {
            if element.id == r3.id {
                t.Errorf("[ for s2[element].id ] expected: %s, actual: %#v", "<not found>", element.id)
            }
        }
    })
}
