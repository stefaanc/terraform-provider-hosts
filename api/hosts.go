//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
    "sync"
)

// -----------------------------------------------------------------------------

func Init() {
    initHosts()
    return
}

// -----------------------------------------------------------------------------

type anchor struct {
    files []*fileObject   // !!! beware of memory leaks

    newFileID func () fileID
    fileIndex *fileIndex

    newZoneID func () zoneID
    zoneIndex *zoneIndex

    newRecordID func () recordID
    recordIndex *recordIndex
}

var hosts *anchor

func initHosts() {
    if hosts != nil {
        // already initialized
        return
    }

    hosts = new(anchor)

    lastFileID := fileID(0)
    hosts.newFileID = func() fileID {
        lastFileID += 1
        return lastFileID
    }
    hosts.fileIndex = new(fileIndex)
    hosts.fileIndex.index = make(map[fileID]*File)
    hosts.fileIndex.paths = make(map[string][]*File)

    lastZoneID := zoneID(0)
    hosts.newZoneID = func() zoneID {
        lastZoneID += 1
        return lastZoneID
    }
    hosts.zoneIndex = new(zoneIndex)
    hosts.zoneIndex.index = make(map[zoneID]*Zone)
    hosts.zoneIndex.files = make(map[fileID][]*Zone)
    hosts.zoneIndex.names = make(map[string][]*Zone)

    lastRecordID := recordID(0)
    hosts.newRecordID = func() recordID {
        lastRecordID += 1
        return lastRecordID
    }
    hosts.recordIndex = new(recordIndex)
    hosts.recordIndex.index = make(map[recordID]*Record)
    hosts.recordIndex.zones = make(map[zoneID][]*Record)
    hosts.recordIndex.addresses = make(map[string][]*Record)
    hosts.recordIndex.names = make(map[string][]*Record)
}

// -----------------------------------------------------------------------------

type fileID int

type fileIndex struct {
    sync.RWMutex
    index map[fileID]*File
    paths map[string][]*File
}

func lookupFile(fQuery *File) (f *File) {
    if fQuery.id != 0 {
        return fQuery
    }

    if fQuery.ID != 0 {
        hosts.fileIndex.RLock()
        f := hosts.fileIndex.index[fileID(fQuery.ID)]
        hosts.fileIndex.RUnlock()

        if f == nil {
            return nil
        }

        // check other identifying properties
        if fQuery.Path != "" && fQuery.Path != f.Path {
            return nil
        }

        return f
    }

    if fQuery.Path != "" {
        hosts.fileIndex.RLock()
        fs := hosts.fileIndex.paths[fQuery.Path]
        hosts.fileIndex.RUnlock()

        if len(fs) != 1 {
            // if more than 1 valid record, 'get' cannot decide which one to return
            return nil
        }

        return fs[0]
    }

    return nil
}

func addFile(f *File) {
    if f.id != 0 {
        // file already indexed
        return
    }

    id := hosts.newFileID()
    path := f.Path

    hosts.fileIndex.Lock()
    hosts.fileIndex.index[id] = f
    hosts.fileIndex.paths[path] = append(hosts.fileIndex.paths[path], f)
    hosts.fileIndex.Unlock()

    f.ID = int(id)
    f.id = id

    return
}

func removeFile(f *File) {
    if f.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := f.id
    path := f.Path

    hosts.fileIndex.Lock()
    delete(hosts.fileIndex.index, id)
    hosts.fileIndex.paths[path] = deleteFromSliceOfFiles(hosts.fileIndex.paths[path], f)
    hosts.fileIndex.Unlock()

    f.ID = 0
    f.id = fileID(0)

    return
}

func deleteFromSliceOfFiles(fs []*File, f *File) []*File {
    if len(fs) == 0 {
        return []*File(nil)   // always return a copy
    }

    newFiles := make([]*File, 0, len(fs) - 1)
    for _, file := range fs {
        if file == f {
            continue
        }

        newFiles = append(newFiles, file)
    }

    return newFiles
}

// -----------------------------------------------------------------------------

type zoneID int

type zoneIndex struct {
    sync.RWMutex
    index map[zoneID]*Zone
    files map[fileID][]*Zone
    names map[string][]*Zone
}

func lookupZone(zQuery *Zone) (z *Zone) {
    if zQuery.id != 0 {
        return zQuery
    }

    if zQuery.ID != 0 {
        hosts.zoneIndex.RLock()
        z := hosts.zoneIndex.index[zoneID(zQuery.ID)]
        hosts.zoneIndex.RUnlock()

        if z == nil {
            return nil
        }

        // check other identifying properties
        if zQuery.File != 0 && zQuery.File != z.File {
            return nil
        }
        if zQuery.Name != "" && zQuery.Name != z.Name {
            return nil
        }

        return z
    }

    if zQuery.Name != "" {
        hosts.zoneIndex.RLock()
        zs := hosts.zoneIndex.names[zQuery.Name]
        hosts.zoneIndex.RUnlock()
        if len(zs) == 0 {
            return nil
        }

        // check other identifying properties
        if zQuery.File != 0 {
            zsReduced := make([]*Zone, 0)
            for _, candidate := range zs {
                // a valid candidate has a file equal to zQuery.File
                if candidate.File == zQuery.File {
                    zsReduced = append(zsReduced, candidate)
                }
            }
            zs = zsReduced

            if len(zs) == 0 {
                return nil
            }
        }
        
        if len(zs) != 1 {
            // if more than 1 valid zone, 'get' cannot decide which one to return
            return nil
        }

        return zs[0]
    }

    if zQuery.File != 0 {
        hosts.zoneIndex.RLock()
        zs := hosts.zoneIndex.files[fileID(zQuery.File)]
        hosts.zoneIndex.RUnlock()

        if len(zs) != 1 {
            // if more than 1 valid zone, 'get' cannot decide which one to return
            return nil
        }

        return zs[0]
    }

    return nil
}

func addZone(z *Zone) {
    if z.id != 0 {
        // zone already indexed
        return
    }

    id := hosts.newZoneID()
    file := fileID(z.File)
    name := z.Name

    hosts.zoneIndex.Lock()
    hosts.zoneIndex.index[id] = z
    hosts.zoneIndex.files[file] = append(hosts.zoneIndex.files[file], z)
    hosts.zoneIndex.names[name] = append(hosts.zoneIndex.names[name], z)
    hosts.zoneIndex.Unlock()

    z.ID = int(id)
    z.id = id

    return
}

func removeZone(z *Zone) {
    if z.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := z.id
    file := fileID(z.File)
    name := z.Name

    hosts.zoneIndex.Lock()
    delete(hosts.zoneIndex.index, id)
    hosts.zoneIndex.files[file] = deleteFromSliceOfZones(hosts.zoneIndex.files[file], z)
    hosts.zoneIndex.names[name] = deleteFromSliceOfZones(hosts.zoneIndex.names[name], z)
    hosts.zoneIndex.Unlock()

    z.ID = 0
    z.id = zoneID(0)

    return
}

func deleteFromSliceOfZones(zs []*Zone, z *Zone) []*Zone {
    if len(zs) == 0 {
        return []*Zone(nil)   // always return a copy
    }

    newZones := make([]*Zone, 0, len(zs) - 1)
    for _, zone := range zs {
        if zone == z {
            continue
        }
        newZones = append(newZones, zone)
    }

    return newZones
}

// -----------------------------------------------------------------------------

type recordID int

type recordIndex struct {
    sync.RWMutex
    index map[recordID]*Record
    zones map[zoneID][]*Record
    addresses map[string][]*Record
    names map[string][]*Record
}

func lookupRecord(rQuery *Record) (r *Record) {
    if rQuery.id != 0 {
        return rQuery
    }

    if rQuery.ID != 0 {
        hosts.recordIndex.RLock()
        r := hosts.recordIndex.index[recordID(rQuery.ID)]
        hosts.recordIndex.RUnlock()

        if r == nil {
            return nil
        }

        // check other identifying properties
        if rQuery.Zone != 0 && rQuery.Zone != r.Zone {
            return nil
        }
        if rQuery.Address != "" && rQuery.Address != r.Address {
            return nil
        }
        if len(rQuery.Names) > 0 {
            // a valid record has all names (or more) that are found in r.Names
            valid := true
            for _, n := range rQuery.Names {
                found := false
                for _, name := range r.Names {
                    if n == name {
                        found = true
                        break
                    }
                }
                if !found {
                    valid = false
                    break
                }
            }
            if !valid {
                return nil
            }
        }

        return r
    }

    if len(rQuery.Names) > 0 {
        hosts.recordIndex.RLock()
        rs := hosts.recordIndex.names[rQuery.Names[0]]
        hosts.recordIndex.RUnlock()

        if len(rs) == 0 {
            return nil
        }

        // check other identifying names
        if len(rQuery.Names) > 0 {
            rsReduced := make([]*Record, 0)
            for _, candidate := range rs {
                // a valid candidate has all names (or more) that are found in rQuery.Names
                valid := true
                for _, n := range rQuery.Names {
                    found := false
                    for _, name := range candidate.Names {
                        if n == name {
                            found = true
                            break
                        }
                    }
                    if !found {
                        valid = false
                        break
                    }
                }
                if valid {
                    rsReduced = append(rsReduced, candidate)
                }
            }
            rs = rsReduced

            if len(rs) == 0 {
                return nil
            }
        }

        // check other identifying properties
        if rQuery.Address != "" {
            rsReduced := make([]*Record, 0)
            for _, candidate := range rs {
                // a valid candidate has an address equal to rQuery.Address
                if candidate.Address == rQuery.Address {
                    rsReduced = append(rsReduced, candidate)
                }
            }
            rs = rsReduced

            if len(rs) == 0 {
                return nil
            }
        }
        if rQuery.Zone != 0 {
            rsReduced := make([]*Record, 0)
            for _, candidate := range rs {
                // a valid candidate has a zone equal to rQuery.Zone
                if candidate.Zone == rQuery.Zone {
                    rsReduced = append(rsReduced, candidate)
                }
            }
            rs = rsReduced

            if len(rs) == 0 {
                return nil
            }
        }

        if len(rs) != 1 {
            // if more than 1 valid record, 'get' cannot decide which one to return
            return nil
        }

        return rs[0]
    }

    if rQuery.Address != "" {
        hosts.recordIndex.RLock()
        rs := hosts.recordIndex.addresses[rQuery.Address]
        hosts.recordIndex.RUnlock()

        if len(rs) == 0 {
            return nil
        }

        // check other identifying properties
        if rQuery.Zone != 0 {
            rsReduced := make([]*Record, 0)
            for _, candidate := range rs {
                // a valid candidate has a zone equal to rQuery.Zone
                if candidate.Zone == rQuery.Zone {
                    rsReduced = append(rsReduced, candidate)
                }
            }
            rs = rsReduced

            if len(rs) == 0 {
                return nil
            }
        }

        if len(rs) != 1 {
            // if more than 1 valid record, 'get' cannot decide which one to return
            return nil
        }

        return rs[0]
    }

    if rQuery.Zone != 0 {
        hosts.recordIndex.RLock()
        rs := hosts.recordIndex.zones[zoneID(rQuery.Zone)]
        hosts.recordIndex.RUnlock()

        if len(rs) != 1 {
            // if more than 1 valid record, 'get' cannot decide which one to return
            return nil
        }

        return rs[0]
    }

    return nil
}

func addRecord(r *Record) {
    if r.id != 0 {
        // record already indexed
        return
    }

    id := hosts.newRecordID()
    zone := zoneID(r.Zone)
    address := r.Address
    names := r.Names

    hosts.recordIndex.Lock()
    hosts.recordIndex.index[id] = r
    hosts.recordIndex.zones[zone] = append(hosts.recordIndex.zones[zone], r)
    hosts.recordIndex.addresses[address] = append(hosts.recordIndex.addresses[address], r)
    for _, n := range names {
        hosts.recordIndex.names[n] = append(hosts.recordIndex.names[n], r)
    }
    hosts.recordIndex.Unlock()

    r.ID = int(id)
    r.id = id

    return
}

func removeRecord(r *Record) {
    if r.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := r.id
    zone := zoneID(r.Zone)
    address := r.Address
    names := r.Names

    hosts.recordIndex.Lock()
    delete(hosts.recordIndex.index, id)
    hosts.recordIndex.zones[zone] = deleteFromSliceOfRecords(hosts.recordIndex.zones[zone], r)
    hosts.recordIndex.addresses[address] = deleteFromSliceOfRecords(hosts.recordIndex.addresses[address], r)
    for _, n := range names {
       hosts.recordIndex.names[n] = deleteFromSliceOfRecords(hosts.recordIndex.names[n], r)
    }
    hosts.recordIndex.Unlock()

    r.ID = 0
    r.id = recordID(0)

    return
}

func deleteFromSliceOfRecords(rs []*Record, r *Record) []*Record {
    if len(rs) == 0 {
        return []*Record(nil)   // always return a copy
    }

    newRecords := make([]*Record, 0, len(rs) - 1)
    for _, record := range rs {
        if record == r {
            continue
        }

        newRecords = append(newRecords, record)
    }

    return newRecords
}

// -----------------------------------------------------------------------------

type fileObject struct {
    data     []byte   // filled by goRenderFile(), cleared by goScanFile()
    checksum string
    file     *File    // !!! beware of memory leaks
}

func addFileObject(f *fileObject) {
    hosts.files = append(hosts.files, f)
    return
}

func removeFileObject(f *fileObject) {
    hosts.files = deleteFromSliceOfFileObjects(hosts.files, f)
    return
}

func deleteFromSliceOfFileObjects(fs []*fileObject, f *fileObject) []*fileObject {
    if len(fs) == 0 {
        return []*fileObject(nil)   // always return a copy
    }

    newFileObjects := make([]*fileObject, 0, len(fs) - 1)
    for _, fileObject := range fs {
        if f == fileObject {
            continue
        }
        newFileObjects = append(newFileObjects, fileObject)
    }

    return newFileObjects
}
