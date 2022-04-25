package zaprec

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// recorder is a Logger that records logged entries for testing. Remember that
// all With-calls return a logger with the same original RecordStore.
type recorder struct {
	// store is where records are added in Write-calls.
	store *RecordStore
	// levelEnabler is the zapcore.LevelEnabler to use in Enabled.
	levelEnabler zapcore.LevelEnabler
	// fields is structured context for the logger.
	fields []zapcore.Field
	// fieldsMutex locks fields.
	fieldsMutex sync.RWMutex
	// writes holds all active write calls.
	writes sync.WaitGroup
}

// NewRecorder creates and initializes a new recorder with the given
// zapcore.LevelEnabler. The returned zap.Logger records all entries to the
// returned RecordStore.
func NewRecorder(levelEnabler zapcore.LevelEnabler) (*zap.Logger, *RecordStore) {
	rs := &RecordStore{}
	return zap.New(&recorder{
		store:        rs,
		levelEnabler: levelEnabler,
	}), rs
}

// Enabled uses the internal level enabler to check if the given level is
// enabled.
func (r *recorder) Enabled(level zapcore.Level) bool {
	if r.levelEnabler == nil {
		return true
	}
	return r.levelEnabler.Enabled(level)
}

// With creates a new zapcore.Core with the given added fields. Keep in mind
// that the same RecordStore is used.
func (r *recorder) With(fields []zapcore.Field) zapcore.Core {
	r.fieldsMutex.RLock()
	defer r.fieldsMutex.RUnlock()
	return &recorder{
		store:        r.store,
		levelEnabler: r.levelEnabler,
		fields:       append(r.fields, fields...),
	}
}

// Check performs the Enabled check and accepts if ok.
func (r *recorder) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !r.Enabled(entry.Level) {
		return ce
	}
	return ce.AddCore(entry, r)
}

// Write creates a Record from the passed data and adds them to the RecordStore.
func (r *recorder) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	r.writes.Add(1)
	defer r.writes.Done()
	r.fieldsMutex.RLock()
	defer r.fieldsMutex.RUnlock()
	r.store.add(Record{
		Entry:  entry,
		Fields: append(r.fields, fields...),
	})
	return nil
}

// Sync waits on all Write calls to finish.
func (r *recorder) Sync() error {
	r.writes.Wait()
	return nil
}
