package zaprec

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

// TestNewRecorder tests NewRecorder.
func TestNewRecorder(t *testing.T) {
	levelEnabler := zapcore.LevelEnabler(zapcore.InfoLevel)
	logger, rec := NewRecorder(levelEnabler)
	recorderLogger, ok := logger.Core().(*recorder)
	require.True(t, ok, "core should be of correct time")
	assert.Same(t, rec, recorderLogger.store, "should set correct store")
	assert.Equal(t, levelEnabler, recorderLogger.levelEnabler, "should set correct level enabler")
}

// recorderEnabledSuite tests recorder.Enabled.
type recorderEnabledSuite struct {
	suite.Suite
}

func (suite *recorderEnabledSuite) TestNotEnabled() {
	r := &recorder{levelEnabler: zapcore.InfoLevel}
	suite.False(r.Enabled(zapcore.DebugLevel), "should not be enabled")
}

func (suite *recorderEnabledSuite) TestEnabled1() {
	r := &recorder{levelEnabler: zapcore.InfoLevel}
	suite.True(r.Enabled(zapcore.InfoLevel), "should be enabled")
}

func (suite *recorderEnabledSuite) TestEnabled2() {
	r := &recorder{levelEnabler: zapcore.InfoLevel}
	suite.True(r.Enabled(zapcore.FatalLevel), "should be enabled")
}

func TestRecorder_Enabled(t *testing.T) {
	suite.Run(t, new(recorderEnabledSuite))
}

func TestRecorder_With(t *testing.T) {
	levelEnabler := zapcore.LevelEnabler(zapcore.FatalLevel)
	logger, rec := NewRecorder(levelEnabler)
	with := logger.With(zap.String("meow", "woof"))
	recorderWith, ok := with.Core().(*recorder)
	require.True(t, ok, "with should return correct type")
	assert.Equal(t, levelEnabler, recorderWith.levelEnabler, "should keep level recorder")
	assert.Same(t, rec, recorderWith.store, "should keep record store")
	assert.Equal(t, []zapcore.Field{
		zap.String("meow", "woof"),
	}, recorderWith.fields, "should add new fields")
}

// recorderSuite for general testing with zap.Logger stuff.
type recorderSuite struct {
	suite.Suite
}

func (suite *recorderSuite) TestLevelNotEnabled() {
	logger, rec := NewRecorder(zapcore.InfoLevel)
	logger.Debug("meow")
	suite.Len(rec.Records(), 0, "should not have recorded")
}

func (suite *recorderSuite) TestOK() {
	logger, rec := NewRecorder(nil)
	logger.Info("woof", zap.String("meow", "cluck"))
	suite.Len(rec.Records(), 1, "should have recorded")
	suite.Equal("woof", rec.Records()[0].Entry.Message, "should have recorded correct message")
	suite.Equal(zapcore.InfoLevel, rec.Records()[0].Entry.Level, "should have recorded correct level")
	suite.Contains(rec.Records()[0].Fields, zap.String("meow", "cluck"), "should contain field")
}

func (suite *recorderSuite) TestWith() {
	logger, rec := NewRecorder(nil)
	logger = logger.With(zap.String("hello", "world"))
	logger.Debug("cluck", zap.String("cluck", "cluck"))
	suite.Len(rec.Records(), 1, "should have recorded")
	suite.Contains(rec.Records()[0].Fields, zap.String("hello", "world"), "should contain logger field")
	suite.Contains(rec.Records()[0].Fields, zap.String("cluck", "cluck"), "should contain entry field")
}

func TestRecorder(t *testing.T) {
	suite.Run(t, new(recorderSuite))
}
