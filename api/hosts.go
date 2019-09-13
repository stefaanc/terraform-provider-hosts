package api

import (
    "errors"
)

// -----------------------------------------------------------------------------

func InitHosts() {
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

// -----------------------------------------------------------------------------

type fileID int

func (hosts *anchor) getFile(f *File) (file *File) {
    if f.id != 0 {
        return f
    }

    if f.ID != 0 {
        file := hosts.fileIndex[f.ID]
        if file == nil {
            return nil
        }

        // check other identifying properties
        if f.Path != "" && f.Path != file.path {
            return nil
        }

        return file
    }

    if f.Path != "" {
        fs := hosts.filePaths[f.Path]
        if len(fs) != 1 {
            // if more than 1 valid record, 'get' cannot decide which one to return
            return nil
        }

        return fs[0]
    }

    return nil
}

func (hosts *anchor) addFile(file *File) error {
    if file.id != 0 {
        // file already indexed
        return nil
    }

    if hosts.getFile(file) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addFile()] another file with these indexing properties is already indexed")
    }

    id := hosts.newFileID()
    path := file.path

    hosts.fileIndex[id] = file
    hosts.filePaths[path] = append(hosts.filePaths[path], file)

    file.ID = id
    file.id = id

    return nil
}

func (hosts *anchor) removeFile(file *File) {
    if file.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := file.id
    path := file.path

    delete(hosts.fileIndex, id)
    hosts.filePaths[path] = deleteFromSliceOfFiles(hosts.filePaths[path], file)

    file.ID = fileID(0)
    file.id = fileID(0)

    return
}

func deleteFromSliceOfFiles(fs []*File, file *File) []*File {
    if len(fs) == 0 {
        return make([]*File, 0)   // always return a copy
    }

    newFiles := make([]*File, len(fs) - 1)
    decr := 0
    for i, f := range fs {
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

func (hosts *anchor) getZone(z *Zone) *Zone {
    if z.id != 0 {
        return z
    }

    if z.ID != 0 {
        zone := hosts.zoneIndex[z.ID]
        if zone == nil {
            return nil
        }

        // check other identifying properties
        if z.File != "" && z.File != zone.file {
            return nil
        }
        if z.Name != "" && z.Name != zone.name {
            return nil
        }

        return zone
    }

    if z.Name != "" {
        zs := hosts.zoneNames[z.Name]
        if len(zs) == 0 {
            return nil
        }
        zone := (*Zone)(nil)

        // check other identifying properties
        if z.File == "" {
            if len(zs) == 1 {
                zone = zs[0]
            }
        } else {
            for _, candidate := range zs {
                if z.File == candidate.file {
                    zone = candidate
                    break
                }
            }
        }

        return zone
    }

    if z.File != "" {
        zs := hosts.zoneFiles[z.File]
        if len(zs) != 1 {
            // if more than 1 valid zone, 'get' cannot decide which one to return
            return nil
        }

        return zs[0]
    }

    return nil
}

func (hosts *anchor) addZone(zone *Zone) error {
    if zone.id != 0 {
        // zone already indexed
        return nil
    }

    if hosts.getZone(zone) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addZone()] another zone with these indexing properties is already indexed")
    }

    id := hosts.newZoneID()
    file := zone.file
    name := zone.name

    hosts.zoneIndex[id] = zone
    hosts.zoneFiles[file] = append(hosts.zoneFiles[file], zone)
    hosts.zoneNames[name] = append(hosts.zoneNames[name], zone)

    zone.ID = id
    zone.id = id

    return nil
}

func (hosts *anchor) removeZone(zone *Zone) {
    if zone.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := zone.id
    file := zone.file
    name := zone.name

    delete(hosts.zoneIndex, id)
    hosts.zoneFiles[file] = deleteFromSliceOfZones(hosts.zoneFiles[file], zone)
    hosts.zoneNames[name] = deleteFromSliceOfZones(hosts.zoneNames[name], zone)

    zone.ID = zoneID(0)
    zone.id = zoneID(0)

    return
}

func deleteFromSliceOfZones(zs []*Zone, zone *Zone) []*Zone {
    if len(zs) == 0 {
        return make([]*Zone, 0)   // always return a copy
    }

    newZones := make([]*Zone, len(zs) - 1)
    decr := 0
    for i, z := range zs {
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

func (hosts *anchor) getRecord(r *Record) *Record {
    if r.id != 0 {
        return r
    }

    if r.ID != 0 {
        record := hosts.recordIndex[r.ID]
        if record == nil {
            return nil
        }

        // check other identifying properties
        if r.Address != "" && r.Address != record.Address {
            return nil
        }
        if len(r.Names) > 0 {
            // a valid record has all names (or more) that are found in r.Names
            valid := true
            for _, n := range r.Names {
                found := false
                for _, name := range record.Names {
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

        return record
    }

    if r.Address != "" {
        rs := hosts.recordAddresses[r.Address]
        if len(rs) == 0 {
            return nil
        }
        record := (*Record)(nil)

        if len(r.Names) == 0 {
            if len(rs) == 1 {
                record = rs[0]
            }
        } else {
            // check other identifying properties
            count := 0
            for _, candidate := range rs {
                // a valid candidate has all names (or more) that are found in r.Names
                valid := true
                for _, n := range r.Names {
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
                    record = candidate
                }
            }
        }

        return record
    }

    if len(r.Names) > 0 {
        rs := hosts.recordNames[r.Names[0]]
        if len(rs) == 0 {
            return nil
        }
        record := (*Record)(nil)

        if len(r.Names) == 1 {
            if len(rs) == 1 {
                record = rs[0]
            }
        } else {
            // check other identifying names
            count := 0
            for _, candidate := range rs {
                // a valid candidate has all names (or more) that are found in r.Names
                valid := true
                for i, n := range r.Names {
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
                    record = candidate
                }
            }
        }

        return record
    }

    return nil
}

func (hosts *anchor) addRecord(record *Record) error {
    if record.id != 0 {
        // record already indexed
        return nil
    }

    if hosts.getRecord(record) != nil {
        return errors.New("[ERROR][terraform-provider-hosts/api/hosts.addRecord()] another record with these indexing properties is already indexed")
    }

    id := hosts.newRecordID()
    address := record.Address
    names := record.Names

    hosts.recordIndex[id] = record
    hosts.recordAddresses[address] = append(hosts.recordAddresses[address], record)
    for _, n := range names {
        hosts.recordNames[n] = append(hosts.recordNames[n], record)
    }

    record.ID = id
    record.id = id

    return nil
}

func (hosts *anchor) removeRecord(record *Record) {
    if record.id == 0 {
        // file not indexed yet/anymore
        return
    }

    id := record.id
    address := record.Address
    names := record.Names

    delete(hosts.recordIndex, id)
    hosts.recordAddresses[address] = deleteFromSliceOfRecords(hosts.recordAddresses[address], record)
    for _, n := range names {
       hosts.recordNames[n] = deleteFromSliceOfRecords(hosts.recordNames[n], record)
    }

    record.ID = recordID(0)
    record.id = recordID(0)

    return
}

func deleteFromSliceOfRecords(rs []*Record, record *Record) []*Record {
    if len(rs) == 0 {
        return make([]*Record, 0)   // always return a copy
    }

    newRecords := make([]*Record, len(rs) - 1)
    decr := 0
    for i, r := range rs {
        if record.id == r.id {
            decr = 1
        } else {
            newRecords[i-decr] = rs[i]
        }
    }

    return newRecords
}

// -----------------------------------------------------------------------------
