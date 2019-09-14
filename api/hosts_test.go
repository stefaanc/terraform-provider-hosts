package api

import (
    "testing"
)

// -----------------------------------------------------------------------------

func resetHostsTestEnv() {
    hosts = (*anchor)(nil)
    Init()
}

// -----------------------------------------------------------------------------

func Test_init(t *testing.T) {
    var test string

    test = "1st-init"
    t.Run(test, func(t *testing.T) {

        hosts = (*anchor)(nil)

        // --------------------

        Init()

        // --------------------

        if hosts == nil {
            t.Errorf("[ hosts ] expected: not <nil>, actual: <nil>")
        } else {

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
        }
    })

    test = "2nd-init"
    t.Run(test, func(t *testing.T) {

        hosts = (*anchor)(nil)
        Init()

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

        // --------------------

        Init()

        // --------------------

        if hosts == nil {
            t.Errorf("[ hosts ] expected: not <nil>, actual: <nil>")
        } else {

            // --------------------

            fid = hosts.newFileID()
            if fid != 2 {
                t.Errorf("[ hosts.newFileID() ] expected: 2, actual: %d", fid)
            }

            length := len(hosts.fileIndex)
            if length != 1 {
                t.Errorf("[ len(hosts.fileIndex) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.filePaths)
            if length != 1 {
                t.Errorf("[ len(hosts.filePaths) ] expected: 1, actual: %d", length)
            }

            // --------------------

            zid = hosts.newZoneID()
            if zid != 2 {
                t.Errorf("[ hosts.newZoneID() ] expected: 2, actual: %d", zid)
            }

            length = len(hosts.zoneIndex)
            if length != 1 {
                t.Errorf("[ len(hosts.zoneIndex) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.zoneFiles)
            if length != 1 {
                t.Errorf("[ len(hosts.zoneFiles) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.zoneNames)
            if length != 1 {
                t.Errorf("[ len(hosts.zoneNames) ] expected: 1, actual: %d", length)
            }

            // --------------------

            rid = hosts.newRecordID()
            if rid != 2 {
                t.Errorf("[ hosts.newRecordID() ] expected: 2, actual: %d", rid)
            }

            length = len(hosts.recordIndex)
            if length != 1 {
                t.Errorf("[ len(hosts.recordIndex) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.recordNames)
            if length != 1 {
                t.Errorf("[ len(hosts.recordNames) ] expected: 1, actual: %d", length)
            }
        }
    })
}

// -----------------------------------------------------------------------------

func Test_getFile(t *testing.T) {
    var test string

    setupIndex1 := func () (f *File) {
        f = new(File)
        f.id = hosts.newFileID()
        f.Path = "p"
        hosts.fileIndex[f.id] = f
        hosts.filePaths[f.Path] = append(hosts.filePaths[f.Path], f)

        return f
    }

    setupIndex2 := func () (f1 *File, f2 *File) {
        f1 = new(File)
        f1.id = hosts.newFileID()
        f1.Path = "p"
        hosts.fileIndex[f1.id] = f1
        hosts.filePaths[f1.Path] = append(hosts.filePaths[f1.Path], f1)

        f2 = new(File)
        f2.id = hosts.newFileID()
        f2.Path = "p"
        hosts.fileIndex[f2.id] = f2
        hosts.filePaths[f2.Path] = append(hosts.filePaths[f2.Path], f2)

        return f1, f2
    }

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        file := GetFile(f)

        // --------------------

        if file == nil {
            t.Errorf("[ GetFile(f) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ GetFile(f) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/no-files"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        fQuery := new(File)
        fQuery.ID = fileID(42)

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = f.id

        file := GetFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = fileID(42)

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "p"

        file := GetFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-Path/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "x"

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "p"

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/ID-and-Path/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = f.id
        fQuery.Path = "p"

        file := GetFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/ID-and-Path/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = f.id
        fQuery.Path = "x"

        file := GetFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ GetFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })
}

func Test_addFile(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        f := new(File)
        f.Path = "p"

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

        fQuery := new(File)
        fQuery.ID = f.id

        file := GetFile(fQuery)
        if file == nil {
            t.Errorf("[ GetFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ GetFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = f.Path

        file = GetFile(fQuery)
        if file == nil {
            t.Errorf("[ GetFile(fQuery.Path) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ GetFile(fQuery.Path) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        f := new(File)
        f.Path = "p"

        _ = hosts.addFile(f)

        // --------------------

        err := hosts.addFile(f)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addFile(f) ] expected: %s, actual: %#v", "<error>", err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        f := new(File)
        f.Path = "p"

        _ = hosts.addFile(f)

        // --------------------

        f = new(File)
        f.Path = "p"
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

        resetHostsTestEnv()

        f := new(File)
        f.Path = "p"

        _ = hosts.addFile(f)

        // --------------------

        hosts.removeFile(f)

        // --------------------

        if f.ID != 0 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 0, f.ID)
        }

        if f.id != 0 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 0, f.id)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = f.id

        file := GetFile(fQuery)
        if file != nil {
            t.Errorf("[ GetFile(fQuery.ID) ] expected: %#v, actual: %#v", nil, file)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = f.Path

        file = GetFile(fQuery)
        if file != nil {
            t.Errorf("[ GetFile(fQuery.Path) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

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

        resetHostsTestEnv()

        s1 := make([]*File, 0)

        f := new(File)
        f.id = hosts.newFileID()

        // --------------------

        _ = deleteFromSliceOfFiles(s1, f)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        s1 := make([]*File, 0)

        f := new(File)
        f.id = hosts.newFileID()
        s1 = append(s1, f)

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

    setupIndex1 := func () (z *Zone) {
        z = new(Zone)
        zid := hosts.newZoneID()
        z.id = zid
        z.File = "f"
        z.Name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[z.File] = append(hosts.zoneFiles[z.File], z)
        hosts.zoneNames[z.Name] = append(hosts.zoneNames[z.Name], z)

        return z
    }

    setupIndex2z := func () (z1 *Zone, z2 *Zone) {
        z1 = new(Zone)
        z1.id = hosts.newZoneID()
        z1.File = "f"
        z1.Name = "z1"
        hosts.zoneIndex[z1.id] = z1
        hosts.zoneFiles[z1.File] = append(hosts.zoneFiles[z1.File], z1)
        hosts.zoneNames[z1.Name] = append(hosts.zoneNames[z1.Name], z1)

        z2 = new(Zone)
        z2.id = hosts.newZoneID()
        z2.File = "f"
        z2.Name = "z2"
        hosts.zoneIndex[z2.id] = z2
        hosts.zoneFiles[z2.File] = append(hosts.zoneFiles[z2.File], z2)
        hosts.zoneNames[z2.Name] = append(hosts.zoneNames[z2.Name], z2)

        return z1, z2
    }

    setupIndex2f := func () (z1 *Zone, z2 *Zone) {
        z1 = new(Zone)
        z1.id = hosts.newZoneID()
        z1.File = "f1"
        z1.Name = "z"
        hosts.zoneIndex[z1.id] = z1
        hosts.zoneFiles[z1.File] = append(hosts.zoneFiles[z1.File], z1)
        hosts.zoneNames[z1.Name] = append(hosts.zoneNames[z1.Name], z1)

        z2 = new(Zone)
        z2.id = hosts.newZoneID()
        z2.File = "f2"
        z2.Name = "z"
        hosts.zoneIndex[z2.id] = z2
        hosts.zoneFiles[z2.File] = append(hosts.zoneFiles[z2.File], z2)
        hosts.zoneNames[z2.Name] = append(hosts.zoneNames[z2.Name], z2)

        return z1, z2
    }

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zone := GetZone(z)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(z) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(z) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/no-zones"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = zoneID(42)

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = zoneID(42)

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = "f"

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-File/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = "x"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2z()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = "f"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "z"

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-Name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "x"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2f()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "z"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id
        zQuery.File = "f"

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-File/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id
        zQuery.File = "x"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-Name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id
        zQuery.Name = "z"

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-Name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id
        zQuery.Name = "x"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/Name-and-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = "f"
        zQuery.Name = "z"

        zone := GetZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/Name-and-File/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = "x"
        zQuery.Name = "z"

        zone := GetZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ GetZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })
}

func Test_addZone(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        z := new(Zone)
        z.File = "f"
        z.Name = "z"

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

        zQuery := new(Zone)
        zQuery.ID = z.id

        zone := GetZone(zQuery)
        if zone == nil {
            t.Errorf("[ GetZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = z.File

        zone = GetZone(zQuery)
        if zone == nil {
            t.Errorf("[ GetZone(zQuery.File) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery.File) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = z.Name

        zone = GetZone(zQuery)
        if zone == nil {
            t.Errorf("[ GetZone(zQuery.Name) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ GetZone(zQuery.Name) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        z := new(Zone)
        z.File = "f"
        z.Name = "z"
        _ = hosts.addZone(z)

        // --------------------

        err := hosts.addZone(z)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addZone(z) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        z := new(Zone)
        z.File = "f"
        z.Name = "z"
        _ = hosts.addZone(z)

        // --------------------

        z = new(Zone)
        z.File = "f"
        z.Name = "z"

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

        resetHostsTestEnv()

        z := new(Zone)
        z.File = "f"
        z.Name = "z"
        _ = hosts.addZone(z)

        // --------------------

        hosts.removeZone(z)

        // --------------------

        if z.ID != 0 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 0, z.ID)
        }

        if z.id != 0 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 0, z.id)
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = z.id

        zone := GetZone(zQuery)
        if zone != nil {
            t.Errorf("[ GetZone(zQuery.ID) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = z.File

        zone = GetZone(zQuery)
        if zone != nil {
            t.Errorf("[ GetZone(zQuery.File) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = z.Name

        zone = GetZone(zQuery)
        if zone != nil {
            t.Errorf("[ GetZone(zQuery.Name) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

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

        resetHostsTestEnv()

        s1 := make([]*Zone, 0)

        // --------------------

        z := new(Zone)
        z.id = hosts.newZoneID()

        _ = deleteFromSliceOfZones(s1, z)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        s1 := make([]*Zone, 0)

        z := new(Zone)
        z.id = hosts.newZoneID()
        s1 = append(s1, z)

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

    setupIndex1 := func() (r *Record) {
        r = new(Record)
        r.id = hosts.newRecordID()
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        return r
    }

    setupIndex2 := func() (r1 *Record, r2 *Record) {
        r1 = new(Record)
        r1.id = hosts.newRecordID()
        r1.Address = "a"
        r1.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r1.id] = r1
        hosts.recordAddresses[r1.Address] = append(hosts.recordAddresses[r1.Address], r1)
        for _, n := range r1.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r1)
        }

        r2 = new(Record)
        r2.id = hosts.newRecordID()
        r2.Address = "a"
        r2.Names = []string{ "n1", "n2", "n4" }
        hosts.recordIndex[r2.id] = r2
        hosts.recordAddresses[r2.Address] = append(hosts.recordAddresses[r2.Address], r2)
        for _, n := range r2.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r2)
        }

        return r1, r2
    }

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        record := GetRecord(r)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(r) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(r) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/no-records"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = recordID(42)

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = recordID(42)

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Address/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "x"

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/1-name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1" }

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Name/1-name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "x" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Name/1-name/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/more-names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "n3" }

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Name/more-names/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "x" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Name/more-names/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "n2" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Address/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id
        rQuery.Address = "a"

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Address/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id
        rQuery.Address = "x"

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id
        rQuery.Names = []string{ "n1", "n2" }

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Names/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id
        rQuery.Names = []string{ "x" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"
        rQuery.Names = []string{ "n1", "n2" }

        record := GetRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/Address-and-Names/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"
        rQuery.Names = []string{ "x" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"
        rQuery.Names = []string{ "n1", "n2" }

        record := GetRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ GetRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })
}

func Test_addRecord(t *testing.T) {
    var test string

    test = "added"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

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

        rQuery := new(Record)
        rQuery.ID = r.id

        record := GetRecord(rQuery)
        if record == nil {
            t.Errorf("[ GetRecord(rQuery.ID) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery.ID) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = r.Address

        record = GetRecord(rQuery)
        if record == nil {
            t.Errorf("[ GetRecord(rQuery.Address) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery.Address) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = r.Names

        record = GetRecord(rQuery)
        if record == nil {
            t.Errorf("[ GetRecord(rQuery.Names) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ GetRecord(rQuery.Names) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        _ = hosts.addRecord(r)

        // --------------------

        err := hosts.addRecord(r)

        // --------------------

        if err != nil {
            t.Errorf("[ hosts.addRecord(r) ] expected: %#v, actual: %#v", nil, err)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        _ = hosts.addRecord(r)

        // --------------------

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

        resetHostsTestEnv()

        r := new(Record)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        _ = hosts.addRecord(r)

        // --------------------

        hosts.removeRecord(r)

        // --------------------

        if r.ID != 0 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 0, r.ID)
        }

        if r.id != 0 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 0, r.id)
        }

        // --------------------

        rQuery := new(Record)
        rQuery.ID = r.id

        record := GetRecord(rQuery)
        if record != nil {
            t.Errorf("[ GetRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = r.Address

        record = GetRecord(rQuery)
        if record != nil {
            t.Errorf("[ GetRecord(rQuery.Address) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = r.Names

        record = GetRecord(rQuery)
        if record != nil {
            t.Errorf("[ GetRecord(rQuery.Names) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

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

        resetHostsTestEnv()

        s1 := make([]*Record, 0)

        r := new(Record)
        r.id = hosts.newRecordID()

        // --------------------

        _ = deleteFromSliceOfRecords(s1, r)

        // --------------------

        // nothing to test - making sure this doesn't throw a Fatal error

    })

    test = "1-element"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        s1 := make([]*Record, 0)

        r := new(Record)
        r.id = hosts.newRecordID()
        s1 = append(s1, r)

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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

        resetHostsTestEnv()

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

        // --------------------

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
