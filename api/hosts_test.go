//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
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

func Test_Init(t *testing.T) {
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

            length = len(hosts.recordZones)
            if length != 0 {
                t.Errorf("[ len(hosts.recordZones) ] expected: 0, actual: %d", length)
            }

            length = len(hosts.recordAddresses)
            if length != 0 {
                t.Errorf("[ len(hosts.recordAddresses) ] expected: 0, actual: %d", length)
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
        hosts.zoneFiles[fid] = append(hosts.zoneFiles[fid], z)
        hosts.zoneNames["n"] = append(hosts.zoneNames["n"], z)

        r := new(Record)
        rid := hosts.newRecordID()
        hosts.recordIndex[rid] = r
        hosts.recordZones[zid] = append(hosts.recordZones[zid], r)
        hosts.recordAddresses["a"] = append(hosts.recordAddresses["a"], r)
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

            length = len(hosts.recordZones)
            if length != 1 {
                t.Errorf("[ len(hosts.recordZones) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.recordAddresses)
            if length != 1 {
                t.Errorf("[ len(hosts.recordAddresses) ] expected: 1, actual: %d", length)
            }

            length = len(hosts.recordNames)
            if length != 1 {
                t.Errorf("[ len(hosts.recordNames) ] expected: 1, actual: %d", length)
            }
        }
    })
}

// -----------------------------------------------------------------------------

func Test_lookupFile(t *testing.T) {
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

        file := lookupFile(f)

        // --------------------

        if file == nil {
            t.Errorf("[ lookupFile(f) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ lookupFile(f) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/no-files"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        fQuery := new(File)
        fQuery.ID = 42

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := lookupFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = 42

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "p"

        file := lookupFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else  if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/only-Path/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "x"

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/only-Path/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2()

        // --------------------

        fQuery := new(File)
        fQuery.Path = "p"

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "by-query/ID-and-Path/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)
        fQuery.Path = "p"

        file := lookupFile(fQuery)

        // --------------------

        if file == nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "by-query/ID-and-Path/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        f := setupIndex1()

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)
        fQuery.Path = "x"

        file := lookupFile(fQuery)

        // --------------------

        if file != nil {
            t.Errorf("[ lookupFile(fQuery) ] expected: %#v, actual: %#v", nil, file)
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

        addFile(f)

        // --------------------

        if f.ID != 1 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 1, f.ID)
        }

        if f.id != 1 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 1, f.id)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := lookupFile(fQuery)
        if file == nil {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = f.Path

        file = lookupFile(fQuery)
        if file == nil {
            t.Errorf("[ lookupFile(fQuery.Path) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery.Path) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        f := new(File)

        addFile(f)

        // --------------------

        addFile(f)

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := lookupFile(fQuery)
        if file == nil {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        } else if file.id != f.id {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f, file)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        f1 := new(File)
        f1.Path = "p"

        addFile(f1)

        // --------------------

        f2 := new(File)
        f2.Path = "p"
        addFile(f2)

        // --------------------

        if f2.ID != 2 {
            t.Errorf("[ f2.ID ] expected: %#v, actual: %#v", 2, f2.ID)
        }

        if f2.id != 2 {
            t.Errorf("[ f2.id ] expected: %#v, actual: %#v", 2, f2.id)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f2.id)

        file := lookupFile(fQuery)
        if file == nil {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f2, file)
        } else if file.id != f2.id {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", f2, file)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = f2.Path

        file = lookupFile(fQuery)
        if file != nil {
            t.Errorf("[ lookupFile(fQuery.Path) ] expected: %#v, actual: %#v", nil, file)
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

        addFile(f)

        // --------------------

        removeFile(f)

        // --------------------

        if f.ID != 0 {
            t.Errorf("[ f.ID ] expected: %#v, actual: %#v", 0, f.ID)
        }

        if f.id != 0 {
            t.Errorf("[ f.id ] expected: %#v, actual: %#v", 0, f.id)
        }

        // --------------------

        fQuery := new(File)
        fQuery.ID = int(f.id)

        file := lookupFile(fQuery)
        if file != nil {
            t.Errorf("[ lookupFile(fQuery.ID) ] expected: %#v, actual: %#v", nil, file)
        }

        // --------------------

        fQuery = new(File)
        fQuery.Path = f.Path

        file = lookupFile(fQuery)
        if file != nil {
            t.Errorf("[ lookupFile(fQuery.Path) ] expected: %#v, actual: %#v", nil, file)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        f := new(File)

        removeFile(f)

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

func Test_lookupZone(t *testing.T) {
    var test string

    setupIndex1 := func () (z *Zone) {
        fileID := hosts.newFileID()

        z = new(Zone)
        z.id = hosts.newZoneID()
        z.File = int(fileID)
        z.Name = "z"
        hosts.zoneIndex[z.id] = z
        hosts.zoneFiles[fileID] = append(hosts.zoneFiles[fileID], z)
        hosts.zoneNames[z.Name] = append(hosts.zoneNames[z.Name], z)

        return z
    }

    setupIndex2n := func () (z1 *Zone, z2 *Zone) {
        fileID := hosts.newFileID()

        z1 = new(Zone)
        z1.id = hosts.newZoneID()
        z1.File = int(fileID)
        z1.Name = "z1"
        hosts.zoneIndex[z1.id] = z1
        hosts.zoneFiles[fileID] = append(hosts.zoneFiles[fileID], z1)
        hosts.zoneNames[z1.Name] = append(hosts.zoneNames[z1.Name], z1)

        z2 = new(Zone)
        z2.id = hosts.newZoneID()
        z2.File = int(fileID)
        z2.Name = "z2"
        hosts.zoneIndex[z2.id] = z2
        hosts.zoneFiles[fileID] = append(hosts.zoneFiles[fileID], z2)
        hosts.zoneNames[z2.Name] = append(hosts.zoneNames[z2.Name], z2)

        return z1, z2
    }

    setupIndex2f := func () (z1 *Zone, z2 *Zone) {
        fileID1 := hosts.newFileID()
        fileID2 := hosts.newFileID()

        z1 = new(Zone)
        z1.id = hosts.newZoneID()
        z1.File = int(fileID1)
        z1.Name = "z"
        hosts.zoneIndex[z1.id] = z1
        hosts.zoneFiles[fileID1] = append(hosts.zoneFiles[fileID1], z1)
        hosts.zoneNames[z1.Name] = append(hosts.zoneNames[z1.Name], z1)

        z2 = new(Zone)
        z2.id = hosts.newZoneID()
        z2.File = int(fileID2)
        z2.Name = "z"
        hosts.zoneIndex[z2.id] = z2
        hosts.zoneFiles[fileID2] = append(hosts.zoneFiles[fileID2], z2)
        hosts.zoneNames[z2.Name] = append(hosts.zoneNames[z2.Name], z2)

        return z1, z2
    }

    test = "by-id"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zone := lookupZone(z)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(z) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(z) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/no-zones"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = 42

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = 42

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = z.File

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-File/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = 42

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-File/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z1, _ := setupIndex2n()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = z1.File

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "z"

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else  if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/only-Name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "x"

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/only-Name/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2f()

        // --------------------

        zQuery := new(Zone)
        zQuery.Name = "z"

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)
        zQuery.File = z.File

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-File/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)
        zQuery.File = 42

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/ID-and-Name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)
        zQuery.Name = "z"

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/ID-and-Name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)
        zQuery.Name = "x"

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/Name-and-File/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = z.File
        zQuery.Name = "z"

        zone := lookupZone(zQuery)

        // --------------------

        if zone == nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "by-query/Name-and-File/name-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z := setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = z.File
        zQuery.Name = "x"

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "by-query/Name-and-File/file-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        zQuery := new(Zone)
        zQuery.File = 42
        zQuery.Name = "z"

        zone := lookupZone(zQuery)

        // --------------------

        if zone != nil {
            t.Errorf("[ lookupZone(zQuery) ] expected: %#v, actual: %#v", nil, zone)
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
        z.File = int(hosts.newFileID())
        z.Name = "z"

        addZone(z)

        // --------------------

        if z.ID != 1 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 1, z.ID)
        }

        if z.id != 1 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 1, z.id)
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)

        zone := lookupZone(zQuery)
        if zone == nil {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = z.File

        zone = lookupZone(zQuery)
        if zone == nil {
            t.Errorf("[ lookupZone(zQuery.File) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery.File) ] expected: %#v, actual: %#v", z, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = z.Name

        zone = lookupZone(zQuery)
        if zone == nil {
            t.Errorf("[ lookupZone(zQuery.Name) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery.Name) ] expected: %#v, actual: %#v", z, zone)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        z := new(Zone)
        z.File = int(hosts.newFileID())
        z.Name = "z"
        addZone(z)

        // --------------------

        addZone(z)

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)

        zone := lookupZone(zQuery)
        if zone == nil {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        } else if zone.id != z.id {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z, zone)
        }
   })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        fileID := int(hosts.newFileID())

        z1 := new(Zone)
        z1.File = fileID
        z1.Name = "z"
        addZone(z1)

        // --------------------

        z2 := new(Zone)
        z2.File = fileID
        z2.Name = "z"

        addZone(z2)

        // --------------------

        if z2.ID != 2 {
            t.Errorf("[ z2.ID ] expected: %#v, actual: %#v", 2, z2.ID)
        }

        if z2.id != 2 {
            t.Errorf("[ z2.id ] expected: %#v, actual: %#v", 2, z2.id)
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z2.id)

        zone := lookupZone(zQuery)
        if zone == nil {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z2, zone)
        } else if zone.id != z2.id {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", z2, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = z2.File

        zone = lookupZone(zQuery)
        if zone != nil {
            t.Errorf("[ lookupZone(zQuery.File) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = z2.Name

        zone = lookupZone(zQuery)
        if zone != nil {
            t.Errorf("[ lookupZone(zQuery.Name) ] expected: %#v, actual: %#v", nil, zone)
        }
    })
}

func Test_removeZone(t *testing.T) {
    var test string

    test = "removed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        z := new(Zone)
        z.File = int(hosts.newFileID())
        z.Name = "z"
        addZone(z)

        // --------------------

        removeZone(z)

        // --------------------

        if z.ID != 0 {
            t.Errorf("[ z.ID ] expected: %#v, actual: %#v", 0, z.ID)
        }

        if z.id != 0 {
            t.Errorf("[ z.id ] expected: %#v, actual: %#v", 0, z.id)
        }

        // --------------------

        zQuery := new(Zone)
        zQuery.ID = int(z.id)

        zone := lookupZone(zQuery)
        if zone != nil {
            t.Errorf("[ lookupZone(zQuery.ID) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.File = z.File

        zone = lookupZone(zQuery)
        if zone != nil {
            t.Errorf("[ lookupZone(zQuery.File) ] expected: %#v, actual: %#v", nil, zone)
        }

        // --------------------

        zQuery = new(Zone)
        zQuery.Name = z.Name

        zone = lookupZone(zQuery)
        if zone != nil {
            t.Errorf("[ lookupZone(zQuery.Name) ] expected: %#v, actual: %#v", nil, zone)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        z := new(Zone)

        removeZone(z)

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

func Test_lookupRecord(t *testing.T) {
    var test string

    setupIndex1 := func() (r *Record) {
        zoneID := hosts.newZoneID()

        r = new(Record)
        r.id = hosts.newRecordID()
        r.Zone = int(zoneID)
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r.id] = r
        hosts.recordZones[zoneID] = append(hosts.recordZones[zoneID], r)
        hosts.recordAddresses[r.Address] = append(hosts.recordAddresses[r.Address], r)
        for _, n := range r.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r)
        }

        return r
    }

    setupIndex2a := func() (r1 *Record, r2 *Record) {
        zoneID := hosts.newZoneID()

        r1 = new(Record)
        r1.id = hosts.newRecordID()
        r1.Zone = int(zoneID)
        r1.Address = "a1"
        r1.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r1.id] = r1
        hosts.recordZones[zoneID] = append(hosts.recordZones[zoneID], r1)
        hosts.recordAddresses[r1.Address] = append(hosts.recordAddresses[r1.Address], r1)
        for _, n := range r1.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r1)
        }

        r2 = new(Record)
        r2.id = hosts.newRecordID()
        r2.Zone = int(zoneID)
        r2.Address = "a2"
        r2.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r2.id] = r2
        hosts.recordZones[zoneID] = append(hosts.recordZones[zoneID], r2)
        hosts.recordAddresses[r2.Address] = append(hosts.recordAddresses[r2.Address], r2)
        for _, n := range r2.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r2)
        }

        return r1, r2
    }

    setupIndex2n := func() (r1 *Record, r2 *Record) {
        zoneID := hosts.newZoneID()

        r1 = new(Record)
        r1.id = hosts.newRecordID()
        r1.Zone = int(zoneID)
        r1.Address = "a"
        r1.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r1.id] = r1
        hosts.recordZones[zoneID] = append(hosts.recordZones[zoneID], r1)
        hosts.recordAddresses[r1.Address] = append(hosts.recordAddresses[r1.Address], r1)
        for _, n := range r1.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r1)
        }

        r2 = new(Record)
        r2.id = hosts.newRecordID()
        r2.Zone = int(zoneID)
        r2.Address = "a"
        r2.Names = []string{ "n1", "n2", "n4" }
        hosts.recordIndex[r2.id] = r2
        hosts.recordZones[zoneID] = append(hosts.recordZones[zoneID], r2)
        hosts.recordAddresses[r2.Address] = append(hosts.recordAddresses[r2.Address], r2)
        for _, n := range r2.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r2)
        }

        return r1, r2
    }

    setupIndex2z := func() (r1 *Record, r2 *Record) {
        zoneID1 := hosts.newZoneID()
        zoneID2 := hosts.newZoneID()

        r1 = new(Record)
        r1.id = hosts.newRecordID()
        r1.Zone = int(zoneID1)
        r1.Address = "a"
        r1.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r1.id] = r1
        hosts.recordZones[zoneID1] = append(hosts.recordZones[zoneID1], r1)
        hosts.recordAddresses[r1.Address] = append(hosts.recordAddresses[r1.Address], r1)
        for _, n := range r1.Names {
            hosts.recordNames[n] = append(hosts.recordNames[n], r1)
        }

        r2 = new(Record)
        r2.id = hosts.newRecordID()
        r2.Zone = int(zoneID2)
        r2.Address = "a"
        r2.Names = []string{ "n1", "n2", "n3" }
        hosts.recordIndex[r2.id] = r2
        hosts.recordZones[zoneID2] = append(hosts.recordZones[zoneID2], r2)
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

        record := lookupRecord(r)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(r) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(r) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/no-records"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = 42

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/empty-query"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-ID/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-ID/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = 42

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Zone/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r.Zone

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Zone/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = 42

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Zone/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        z1, _ := setupIndex2a()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = z1.Zone

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Address/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "x"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Address/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2z()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/1-name/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1" }

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Names/1-name/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "x" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/1-name/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2z()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/more-names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else  if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/only-Names/more-names/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "x" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/only-Names/more-names/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2z()

        // --------------------

        rQuery := new(Record)
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupZone(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Zone/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Zone = r.Zone

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Zone/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Zone = 42

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Address/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Address/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Address = "x"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/ID-and-Names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/ID-and-Names/not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)
        rQuery.Names = []string{ "x" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Address/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r.Zone
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/Zone-and-Address/zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = 42
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Address/address-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r.Zone
        rQuery.Address = "x"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Address/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r1, _ := setupIndex2n()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r1.Zone
        rQuery.Address = "a"

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Names/found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r.Zone
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/Zone-and-Names/zone-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = 42
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Names/names-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r := setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r.Zone
        rQuery.Names = []string{ "x" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Zone-and-Names/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        r1, _ := setupIndex2a()

        // --------------------

        rQuery := new(Record)
        rQuery.Zone = r1.Zone
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
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

        record := lookupRecord(rQuery)

        // --------------------

        if record == nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "by-query/Address-and-Names/address-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "x"
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/names-not-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex1()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"
        rQuery.Names = []string{ "x" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "by-query/Address-and-Names/multiple-found"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()
        setupIndex2z()

        // --------------------

        rQuery := new(Record)
        rQuery.Address = "a"
        rQuery.Names = []string{ "n1", "n2" }

        record := lookupRecord(rQuery)

        // --------------------

        if record != nil {
            t.Errorf("[ lookupRecord(rQuery) ] expected: %#v, actual: %#v", nil, record)
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
        r.Zone = int(hosts.newZoneID())
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }

        addRecord(r)

        // --------------------

        if r.ID != 1 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 1, r.ID)
        }

        if r.id != 1 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 1, r.id)
        }

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)

        record := lookupRecord(rQuery)
        if record == nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Zone = r.Zone

        record = lookupRecord(rQuery)
        if record == nil {
            t.Errorf("[ lookupRecord(rQuery.Zone) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery.Zone) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = r.Address

        record = lookupRecord(rQuery)
        if record == nil {
            t.Errorf("[ lookupRecord(rQuery.Address) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery.Address) ] expected: %#v, actual: %#v", r, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = r.Names

        record = lookupRecord(rQuery)
        if record == nil {
            t.Errorf("[ lookupRecord(rQuery.Names) ] expected: %#v, actual: %#v", r, record)
        } else if record.id != r.id {
            t.Errorf("[ lookupRecord(rQuery.Names) ] expected: %#v, actual: %#v", r, record)
        }
    })

    test = "already-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        r := new(Record)
        r.Zone = int(hosts.newZoneID())
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        addRecord(r)

        // --------------------

        addRecord(r)

        // --------------------

        if r.ID != 1 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 1, r.ID)
        }

        if r.id != 1 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 1, r.id)
        }
    })

    test = "duplicate"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        zoneID := int(hosts.newZoneID())

        r1 := new(Record)
        r1.Zone = zoneID
        r1.Address = "a"
        r1.Names = []string{ "n1", "n2", "n3" }
        addRecord(r1)

        // --------------------

        r2 := new(Record)
        r2.Zone = zoneID
        r2.Address = "a"
        r2.Names = []string{ "n1", "n2", "n3" }

        addRecord(r2)

        // --------------------

        if r2.ID != 2 {
            t.Errorf("[ r2.ID ] expected: %#v, actual: %#v", 2, r2.ID)
        }

        if r2.id != 2 {
            t.Errorf("[ r2.id ] expected: %#v, actual: %#v", 2, r2.id)
        }

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r2.id)

        record := lookupRecord(rQuery)
        if record == nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", r2, record)
        } else if record.id != r2.id {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", r2, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Zone = r2.Zone

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Zone) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = r2.Address

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Address) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = r2.Names

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Names) ] expected: %#v, actual: %#v", nil, record)
        }
    })
}

func Test_removeRecord(t *testing.T) {
    var test string

    test = "removed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        r := new(Record)
        r.Zone = int(hosts.newZoneID())
        r.Address = "a"
        r.Names = []string{ "n1", "n2", "n3" }
        addRecord(r)

        // --------------------

        removeRecord(r)

        // --------------------

        if r.ID != 0 {
            t.Errorf("[ r.ID ] expected: %#v, actual: %#v", 0, r.ID)
        }

        if r.id != 0 {
            t.Errorf("[ r.id ] expected: %#v, actual: %#v", 0, r.id)
        }

        // --------------------

        rQuery := new(Record)
        rQuery.ID = int(r.id)

        record := lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.ID) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Zone = r.Zone

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Zone) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Address = r.Address

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Address) ] expected: %#v, actual: %#v", nil, record)
        }

        // --------------------

        rQuery = new(Record)
        rQuery.Names = r.Names

        record = lookupRecord(rQuery)
        if record != nil {
            t.Errorf("[ lookupRecord(rQuery.Names) ] expected: %#v, actual: %#v", nil, record)
        }
    })

    test = "not-indexed"
    t.Run(test, func(t *testing.T) {

        resetHostsTestEnv()

        r := new(Record)

        removeRecord(r)

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
