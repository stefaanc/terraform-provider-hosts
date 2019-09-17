//
// Copyright (c) 2019 Stefaan Coussement
// MIT License
//
// more info: https://github.com/stefaanc/terraform-provider-hosts
//
package api

import (
)

// -----------------------------------------------------------------------------

func Init() {
    initHosts()
    return
}

// -----------------------------------------------------------------------------

type anchor struct {
    newFileID func () fileID
    fileIndex map[fileID]*File
    filePaths map[string][]*File

    newZoneID func () zoneID
    zoneIndex map[zoneID]*Zone
    zoneFiles map[fileID][]*Zone
    zoneNames map[string][]*Zone

    newRecordID func () recordID
    recordIndex map[recordID]*Record
    recordZones map[zoneID][]*Record
    recordAddresses map[string][]*Record
    recordNames map[string][]*Record
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
    hosts.fileIndex = make(map[fileID]*File)
    hosts.filePaths = make(map[string][]*File)

    lastZoneID := zoneID(0)
    hosts.newZoneID = func() zoneID {
        lastZoneID += 1
        return lastZoneID
    }
    hosts.zoneIndex = make(map[zoneID]*Zone)
    hosts.zoneFiles = make(map[fileID][]*Zone)
    hosts.zoneNames = make(map[string][]*Zone)

    lastRecordID := recordID(0)
    hosts.newRecordID = func() recordID {
        lastRecordID += 1
        return lastRecordID
    }
    hosts.recordIndex = make(map[recordID]*Record)
    hosts.recordZones = make(map[zoneID][]*Record)
    hosts.recordAddresses = make(map[string][]*Record)
    hosts.recordNames = make(map[string][]*Record)
}

// -----------------------------------------------------------------------------

type fileID int

func lookupFile(fQuery *File) (f *File) {
    if fQuery.id != 0 {
        return fQuery
    }

    if fQuery.ID != 0 {
        f := hosts.fileIndex[fileID(fQuery.ID)]
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
        fs := hosts.filePaths[fQuery.Path]

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

    hosts.fileIndex[id] = f
    hosts.filePaths[path] = append(hosts.filePaths[path], f)

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

    delete(hosts.fileIndex, id)
    hosts.filePaths[path] = deleteFromSliceOfFiles(hosts.filePaths[path], f)

    f.ID = 0
    f.id = fileID(0)

    return
}

func deleteFromSliceOfFiles(fs []*File, f *File) []*File {
    if len(fs) == 0 {
        return []*File(nil)   // always return a copy
    }

    newFiles := make([]*File, len(fs) - 1)
    decr := 0
    for i, file := range fs {
        if file.id == f.id {
            decr = 1
        } else {
            newFiles[i-decr] = fs[i]
        }
    }

    return newFiles
}

// -----------------------------------------------------------------------------

type zoneID int

func lookupZone(zQuery *Zone) (z *Zone) {
    if zQuery.id != 0 {
        return zQuery
    }

    if zQuery.ID != 0 {
        z := hosts.zoneIndex[zoneID(zQuery.ID)]
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
        zs := hosts.zoneNames[zQuery.Name]
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
        zs := hosts.zoneFiles[fileID(zQuery.File)]

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

    hosts.zoneIndex[id] = z
    hosts.zoneFiles[file] = append(hosts.zoneFiles[file], z)
    hosts.zoneNames[name] = append(hosts.zoneNames[name], z)

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

    delete(hosts.zoneIndex, id)
    hosts.zoneFiles[file] = deleteFromSliceOfZones(hosts.zoneFiles[file], z)
    hosts.zoneNames[name] = deleteFromSliceOfZones(hosts.zoneNames[name], z)

    z.ID = 0
    z.id = zoneID(0)

    return
}

func deleteFromSliceOfZones(zs []*Zone, z *Zone) []*Zone {
    if len(zs) == 0 {
        return []*Zone(nil)   // always return a copy
    }

    newZones := make([]*Zone, len(zs) - 1)
    decr := 0
    for i, zone := range zs {
        if zone.id == z.id {
            decr = 1
        } else {
            newZones[i-decr] = zs[i]
        }
    }

    return newZones
}

// -----------------------------------------------------------------------------

type recordID int

func lookupRecord(rQuery *Record) (r *Record) {
    if rQuery.id != 0 {
        return rQuery
    }

    if rQuery.ID != 0 {
        r := hosts.recordIndex[recordID(rQuery.ID)]
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
        rs := hosts.recordNames[rQuery.Names[0]]
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
        rs := hosts.recordAddresses[rQuery.Address]
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
        rs := hosts.recordZones[zoneID(rQuery.Zone)]

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

    hosts.recordIndex[id] = r
    hosts.recordZones[zone] = append(hosts.recordZones[zone], r)
    hosts.recordAddresses[address] = append(hosts.recordAddresses[address], r)
    for _, n := range names {
        hosts.recordNames[n] = append(hosts.recordNames[n], r)
    }

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

    delete(hosts.recordIndex, id)
    hosts.recordZones[zone] = deleteFromSliceOfRecords(hosts.recordZones[zone], r)
    hosts.recordAddresses[address] = deleteFromSliceOfRecords(hosts.recordAddresses[address], r)
    for _, n := range names {
       hosts.recordNames[n] = deleteFromSliceOfRecords(hosts.recordNames[n], r)
    }

    r.ID = 0
    r.id = recordID(0)

    return
}

func deleteFromSliceOfRecords(rs []*Record, r *Record) []*Record {
    if len(rs) == 0 {
        return []*Record(nil)   // always return a copy
    }

    newRecords := make([]*Record, len(rs) - 1)
    decr := 0
    for i, record := range rs {
        if record.id == r.id {
            decr = 1
        } else {
            newRecords[i-decr] = rs[i]
        }
    }

    return newRecords
}

// -----------------------------------------------------------------------------
