package zaprec

import (
	"go.uber.org/zap/zapcore"
	"sync"
)

// Record holds all data that is recorded by recorder in recorder.Write-calls.
// Keep in mind that everything is kept as direct reference.
type Record struct {
	// Entry is the zapcore.Entry to be logged.
	Entry zapcore.Entry
	// Fields are all passed fields in the Write-call.
	Fields []zapcore.Field
}

// RecordStore holds a Record list allowing centralized collection of records.
type RecordStore struct {
	// records is a Record list with all events that have been logged.
	records []Record
	// recordsMutex locks records.
	recordsMutex sync.RWMutex
}

// add the given Record to the store.
func (rs *RecordStore) add(record Record) {
	rs.recordsMutex.Lock()
	defer rs.recordsMutex.Unlock()
	rs.records = append(rs.records, record)
}

// Records returns a copy of the list of records.
func (rs *RecordStore) Records() []Record {
	rs.recordsMutex.RLock()
	defer rs.recordsMutex.RUnlock()
	records := make([]Record, 0, len(rs.records))
	for _, event := range rs.records {
		records = append(records, event)
	}
	return records
}

// RecordsByLevel returns all records that match the given zapcore.Level.
func (rs *RecordStore) RecordsByLevel(level zapcore.Level) []Record {
	rs.recordsMutex.RLock()
	defer rs.recordsMutex.RUnlock()
	found := make([]Record, 0)
	for _, record := range rs.records {
		if record.Entry.Level == level {
			found = append(found, record)
		}
	}
	return found
}
