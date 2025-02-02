package metrics

import (
	"errors"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestInitializeStatsdCollector(t *testing.T) {
	// Test Statter client
	err := InitializeStatsdCollector(&StatsdCollectorConfig{FlushBytes: 0})
	assert.Nil(t, err)
	assert.NotNil(t, instance)
	assert.NotNil(t, instance.client)

	// Test sampleRate
	err = InitializeStatsdCollector(&StatsdCollectorConfig{SampleRate: 5})
	assert.Nil(t, err)
	assert.Equal(t, float32(5), instance.sampleRate)

	err = InitializeStatsdCollector(&StatsdCollectorConfig{})
	assert.Nil(t, err)
	assert.Equal(t, float32(1), instance.sampleRate)
}

func TestUpdate(t *testing.T) {
	// Test when instance is nil
	instance = nil
	Update(UpdateOption{})

	_ = InitializeStatsdCollector(&StatsdCollectorConfig{})

	mc := &mockStatsdClient{}
	instance.client = mc

	now := time.Now()
	mc.On("TimingDuration",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
		mock.AnythingOfType("float32")).Return(nil)
	Update(UpdateOption{Type: Duration, Duration: time.Since(now)})

	mc.On("Inc",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int64"),
		mock.AnythingOfType("float32")).Return(nil)
	Update(UpdateOption{Type: Count})

	// error case
	mc.On("Gauge",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int64"),
		mock.AnythingOfType("float32")).Return(errors.New("error"))
	Update(UpdateOption{Type: Gauge, NumValue: -500})

	mc.AssertExpectations(t)
}

func TestFormatter(t *testing.T) {
	cases := []struct {
		updateOption   UpdateOption
		expectedOutput string
	}{
		{
			updateOption: UpdateOption{
				Scope: "process_image",
				Name:  "duration",
				Type:  Duration,
			},
			expectedOutput: "process_image.duration",
		},
		{
			updateOption: UpdateOption{
				Name: "duration",
				Type: Duration,
			},
			expectedOutput: fmt.Sprintf("%s.%s", DefaultScope, "duration"),
		},
		{
			updateOption: UpdateOption{
				Scope: ".scope",
				Name:  "duration",
				Type:  Duration,
			},
			expectedOutput: "scope.duration",
		},
		{
			updateOption: UpdateOption{
				Scope: ".scope.",
				Name:  "duration",
				Type:  Duration,
			},
			expectedOutput: "scope.duration",
		},
		{
			updateOption: UpdateOption{
				Scope: "scope.name",
				Name:  "duration",
				Type:  Duration,
			},
			expectedOutput: "scope.name.duration",
		},
		{
			updateOption: UpdateOption{
				Scope: ".scope.name.",
				Name:  ".duration.time.",
				Type:  Duration,
			},
			expectedOutput: "scope.name.duration.time",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expectedOutput, formatter(c.updateOption))
	}
}

type mockStatsdClient struct {
	mock.Mock
}

func (msc *mockStatsdClient) Inc(str string, i int64, sr float32) error {
	args := msc.Called(str, i, sr)
	return args.Error(0)
}

func (msc *mockStatsdClient) Dec(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Gauge(str string, i int64, sr float32) error {
	args := msc.Called(str, i, sr)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return args.Error(0)
}

func (msc *mockStatsdClient) GaugeDelta(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Timing(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) TimingDuration(str string, t time.Duration, sr float32) error {
	args := msc.Called(str, t, sr)
	return args.Error(0)
}

func (msc *mockStatsdClient) Set(string, string, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) SetInt(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Raw(string, string, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) NewSubStatter(string) statsd.SubStatter {
	panic("implement me")
}

func (msc *mockStatsdClient) SetPrefix(string) {
	panic("implement me")
}

func (msc *mockStatsdClient) Close() error {
	panic("implement me")
}
