package api

import (
    "errors"
)

// -----------------------------------------------------------------------------

func Init() {
    hosts.init()
    return
}

func GetFile(fQuery *File) (f *File) {
    return hosts.getFile(fQuery)
}

func GetZone(zQuery *Zone) (z *Zone) {
    return hosts.getZone(zQuery)
}

func GetRecord(rQuery *Record) (r *Record) {
    return hosts.getRecord(rQuery)
}

// -----------------------------------------------------------------------------

type anchor struct {
    newFileID func () fileID
    fileIndex map[fileID]*File
    filePaths map[string][]*File

    newZoneID func () zoneID
    zoneIndex map[zoneID]*Zone
    zoneFiles map[string][]*Zone
    zoneNames map[string][]*Zone

    newRecordID func () recordID
    recordIndex map[recordID]*Record
    recordAddresses map[string][]*Record
    recordNames map[string][]*Record
}

var hosts *anchor

func (h *anchor) init() {
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
    hosts.zoneFiles = make(map[string][]*Zone)
    hosts.zoneNames = make(map[string][]*Zone)

    lastRecordID := recordID(0)
    hosts.newRecordID = func() recordID {
        lastRecordID += 1
        return lastRecordID
    }
    hosts.recordIndex = make(map[recordID]*Record)
    hosts.recordAddresses = make(map[string][]*Record)
    hosts.recordNames = make(map[string][]*Record)
}

// -----------------------------------------------------------------------------

type fileID int

func (h *anchor) getFile(fQuery *File) (f *File) {
    if fQuery.id != 0 {
        return fQuery
    }

    if fQuery.ID != 0 {
        f := hosts.fileIndex[fQuery.ID]
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

func (h *anchor) addFile(f *File) error {
    if f.id != 0 {
        // file already indexed
        return nil
    }

    if GetFile(f) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addFile()] another file with these indexing properties is already indexed")
    }

    id := hosts.newFileID()
    path := f.Path

    hosts.fileIndex[id] = f
    hosts.filePaths[path] = append(hosts.filePaths[path], f)

    f.ID = id
    f.id = id

    return nil
}

func (h *anchor) removeFile(f *File) {
    if f.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := f.id
    path := f.Path

    delete(hosts.fileIndex, id)
    hosts.filePaths[path] = deleteFromSliceOfFiles(hosts.filePaths[path], f)

    f.ID = fileID(0)
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

func (h *anchor) getZone(zQuery *Zone) (z *Zone) {
    if zQuery.id != 0 {
        return zQuery
    }

    if zQuery.ID != 0 {
        z := hosts.zoneIndex[zQuery.ID]
        if z == nil {
            return nil
        }

        // check other identifying properties
        if zQuery.File != "" && zQuery.File != z.File {
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
        z := (*Zone)(nil)

        // check other identifying properties
        if zQuery.File == "" {
            if len(zs) == 1 {
                z = zs[0]
            }
        } else {
            for _, candidate := range zs {
                if zQuery.File == candidate.File {
                    z = candidate
                    break
                }
            }
        }

        return z
    }

    if zQuery.File != "" {
        zs := hosts.zoneFiles[zQuery.File]
        if len(zs) != 1 {
            // if more than 1 valid zone, 'get' cannot decide which one to return
            return nil
        }

        return zs[0]
    }

    return nil
}

func (h *anchor) addZone(z *Zone) error {
    if z.id != 0 {
        // zone already indexed
        return nil
    }

    if GetZone(z) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addZone()] another zone with these indexing properties is already indexed")
    }

    id := hosts.newZoneID()
    file := z.File
    name := z.Name

    hosts.zoneIndex[id] = z
    hosts.zoneFiles[file] = append(hosts.zoneFiles[file], z)
    hosts.zoneNames[name] = append(hosts.zoneNames[name], z)

    z.ID = id
    z.id = id

    return nil
}

func (h *anchor) removeZone(z *Zone) {
    if z.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := z.id
    file := z.File
    name := z.Name

    delete(hosts.zoneIndex, id)
    hosts.zoneFiles[file] = deleteFromSliceOfZones(hosts.zoneFiles[file], z)
    hosts.zoneNames[name] = deleteFromSliceOfZones(hosts.zoneNames[name], z)

    z.ID = zoneID(0)
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

func (h *anchor) getRecord(rQuery *Record) (r *Record) {
    if rQuery.id != 0 {
        return rQuery
    }

    if rQuery.ID != 0 {
        r := hosts.recordIndex[rQuery.ID]
        if r == nil {
            return nil
        }

        // check other identifying properties
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

    if rQuery.Address != "" {
        rs := hosts.recordAddresses[rQuery.Address]
        if len(rs) == 0 {
            return nil
        }
        r := (*Record)(nil)

        if len(rQuery.Names) == 0 {
            if len(rs) == 1 {
                r = rs[0]
            }
        } else {
            // check other identifying properties
            count := 0
            for _, candidate := range rs {
                // a valid candidate has all names (or more) that are found in r.Names
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
                    count += 1
                    if count > 1 {
                        // if more than 1 valid record, 'get' cannot decide which one to return
                        return nil
                    }
                    r = candidate
                }
            }
        }

        return r
    }

    if len(rQuery.Names) > 0 {
        rs := hosts.recordNames[rQuery.Names[0]]
        if len(rs) == 0 {
            return nil
        }
        r := (*Record)(nil)

        if len(rQuery.Names) == 1 {
            if len(rs) == 1 {
                r = rs[0]
            }
        } else {
            // check other identifying names
            count := 0
            for _, candidate := range rs {
                // a valid candidate has all names (or more) that are found in r.Names
                valid := true
                for i, n := range rQuery.Names {
                    if i == 0 { continue }   // rs is based on r.Names[0]

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
                    count += 1
                    if count > 1 {
                        // if more than 1 valid record, 'get' cannot decide which one to return
                        return nil
                    }
                    r = candidate
                }
            }
        }

        return r
    }

    return nil
}

func (h *anchor) addRecord(r *Record) error {
    if r.id != 0 {
        // record already indexed
        return nil
    }

    if GetRecord(r) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addRecord()] another record with these indexing properties is already indexed")
    }

    id := hosts.newRecordID()
    address := r.Address
    names := r.Names

    hosts.recordIndex[id] = r
    hosts.recordAddresses[address] = append(hosts.recordAddresses[address], r)
    for _, n := range names {
        hosts.recordNames[n] = append(hosts.recordNames[n], r)
    }

    r.ID = id
    r.id = id

    return nil
}

func (h *anchor) removeRecord(r *Record) {
    if r.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := r.id
    address := r.Address
    names := r.Names

    delete(hosts.recordIndex, id)
    hosts.recordAddresses[address] = deleteFromSliceOfRecords(hosts.recordAddresses[address], r)
    for _, n := range names {
       hosts.recordNames[n] = deleteFromSliceOfRecords(hosts.recordNames[n], r)
    }

    r.ID = recordID(0)
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
