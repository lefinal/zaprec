package zaprec

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"sync"
	"testing"
)

func TestRecordStore_add(t *testing.T) {
	rs := &RecordStore{}
	records := make([]Record, 128)
	var wg sync.WaitGroup
	for _, record := range records {
		wg.Add(1)
		go func(record Record) {
			defer wg.Done()
			rs.add(record)
		}(record)
	}
	wg.Wait()
	assert.Equal(t, records, rs.records)
}

// RecordStoreRecords tests RecordStore.Records.
type RecordStoreRecords struct {
	suite.Suite
}

func (suite *RecordStoreRecords) TestCorrectOrder() {
	logger, rec := NewRecorder(nil)
	entries := 32
	for i := 0; i < entries; i++ {
		logger.Info(fmt.Sprintf("entry: %d", i))
	}
	records := rec.Records()
	suite.Require().Equal(entries, len(records), "should return correct amount of records")
	for i := 0; i < entries; i++ {
		suite.Equal(fmt.Sprintf("entry: %d", i), records[i].Entry.Message, "should be expected message")
	}
}

// TestSliceCopy assures that the returned records slice is a copy so that thread
// safety is provided.
func (suite *RecordStoreRecords) TestSliceCopy() {
	logger, rec := NewRecorder(nil)
	// Record multiple entries.
	for i := 0; i < 32; i++ {
		logger.Info("Hello World!")
	}
	// Alter returned records.
	rec.Records()[16].Entry.Message = "I am a changed message"
	// Record more entries.
	for i := 0; i < 32; i++ {
		logger.Info("More!")
	}
	// Assure not changed.
	records := rec.Records()
	suite.Require().Len(records, 64, "should contain expected record amount")
	for i, record := range records {
		if i < 32 {
			suite.Equal("Hello World!", record.Entry.Message, "should hold correct message")
		} else {
			suite.Equal("More!", record.Entry.Message, "should hold correct message")
		}
	}
}

func TestRecorder_Records(t *testing.T) {
	suite.Run(t, new(RecordStoreRecords))
}

// TestRecorder_RecordsByLevel tests RecordStore.RecordsByLevel.
func TestRecorder_RecordsByLevel(t *testing.T) {
	logger, recorder := NewRecorder(nil)
	logger.Info("info")
	logger.Error("error")
	logger.Info("info")
	records := recorder.RecordsByLevel(zap.InfoLevel)
	require.Len(t, records, 2, "should return correct amount of records")
	for _, record := range records {
		assert.Equal(t, zap.InfoLevel, record.Entry.Level, "record should have correct level")
	}
}
