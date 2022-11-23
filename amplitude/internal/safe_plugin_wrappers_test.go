package internal_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/internal"
	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestSafePluginWrappers(t *testing.T) {
	suite.Run(t, new(SafePluginWrappersSuite))
}

type SafePluginWrappersSuite struct {
	suite.Suite
}

func (t *SafePluginWrappersSuite) TestSafeBeforePluginWrapper() {
	plugin := &testBeforePlugin{}
	logger := &mockLogger{}
	wrapper := internal.SafeBeforePluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	require := t.Require()
	require.Equal(plugin.Name(), wrapper.Name())
	require.Equal(plugin.Type(), wrapper.Type())

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Return(event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeBeforePluginWrapper_PanicOnSetup() {
	plugin := &testBeforePlugin{raisePanicOnSetup: true}
	logger := &mockLogger{}
	wrapper := internal.SafeBeforePluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Setup: %s", []interface{}{"test-before-plugin", "panic in test-before-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeBeforePluginWrapper_PanicOnExecute() {
	plugin := &testBeforePlugin{raisePanicOnExecute: true}
	logger := &mockLogger{}
	wrapper := internal.SafeBeforePluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Execute: %s", []interface{}{"test-before-plugin", "panic in test-before-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeEnrichmentPluginWrapper() {
	plugin := &testEnrichmentPlugin{}
	logger := &mockLogger{}
	wrapper := internal.SafeEnrichmentPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	require := t.Require()
	require.Equal(plugin.Name(), wrapper.Name())
	require.Equal(plugin.Type(), wrapper.Type())

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Return(event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeEnrichmentPluginWrapper_PanicOnSetup() {
	plugin := &testEnrichmentPlugin{raisePanicOnSetup: true}
	logger := &mockLogger{}
	wrapper := internal.SafeEnrichmentPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Setup: %s", []interface{}{"test-enrichment-plugin", "panic in test-enrichment-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeEnrichmentPluginWrapper_PanicOnExecute() {
	plugin := &testEnrichmentPlugin{raisePanicOnExecute: true}
	logger := &mockLogger{}
	wrapper := internal.SafeEnrichmentPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Execute: %s", []interface{}{"test-enrichment-plugin", "panic in test-enrichment-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeDestinationPluginWrapper() {
	plugin := &testDestinationPlugin{}
	logger := &mockLogger{}
	wrapper := internal.SafeDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	require := t.Require()
	require.Equal(plugin.Name(), wrapper.Name())
	require.Equal(plugin.Type(), wrapper.Type())

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeDestinationPluginWrapper_PanicOnSetup() {
	plugin := &testDestinationPlugin{raisePanicOnSetup: true}
	logger := &mockLogger{}
	wrapper := internal.SafeDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Setup: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeDestinationPluginWrapper_PanicOnExecute() {
	plugin := &testDestinationPlugin{raisePanicOnExecute: true}
	logger := &mockLogger{}
	wrapper := internal.SafeDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Execute: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeExtendedDestinationPluginWrapper() {
	plugin := &testDestinationPlugin{}
	logger := &mockLogger{}
	wrapper := internal.SafeExtendedDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	require := t.Require()
	require.Equal(plugin.Name(), wrapper.Name())
	require.Equal(plugin.Type(), wrapper.Type())

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeExtendedDestinationPluginWrapper_PanicOnSetup() {
	plugin := &testDestinationPlugin{raisePanicOnSetup: true}
	logger := &mockLogger{}
	wrapper := internal.SafeExtendedDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Setup: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	wrapper.Execute(event)

	wrapper.Flush()

	wrapper.Shutdown()

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeExtendedDestinationPluginWrapper_PanicOnExecute() {
	plugin := &testDestinationPlugin{raisePanicOnExecute: true}
	logger := &mockLogger{}
	wrapper := internal.SafeExtendedDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Execute: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	event := &types.Event{}
	plugin.On("Execute", event).Once()
	wrapper.Execute(event)

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeExtendedDestinationPluginWrapper_PanicOnFlush() {
	plugin := &testDestinationPlugin{raisePanicOnFlush: true}
	logger := &mockLogger{}
	wrapper := internal.SafeExtendedDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Flush: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	plugin.On("Flush")
	wrapper.Flush()

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

func (t *SafePluginWrappersSuite) TestSafeExtendedDestinationPluginWrapper_PanicOnShutdown() {
	plugin := &testDestinationPlugin{raisePanicOnShutdown: true}
	logger := &mockLogger{}
	wrapper := internal.SafeExtendedDestinationPluginWrapper{
		Plugin: plugin,
		Logger: logger,
	}

	logger.On("Errorf", "Panic in plugin %s.Shutdown: %s", []interface{}{"test-destination-plugin", "panic in test-destination-plugin"}).Return().Once()

	config := types.Config{}
	plugin.On("Setup", config).Once()
	wrapper.Setup(config)

	plugin.On("Shutdown")
	wrapper.Shutdown()

	plugin.AssertExpectations(t.T())
	logger.AssertExpectations(t.T())
}

type testBeforePlugin struct {
	mock.Mock
	raisePanicOnSetup   bool
	raisePanicOnExecute bool
}

func (p *testBeforePlugin) Name() string {
	return "test-before-plugin"
}

func (p *testBeforePlugin) Type() types.PluginType {
	return types.PluginTypeBefore
}

func (p *testBeforePlugin) Setup(config types.Config) {
	p.Called(config)

	if p.raisePanicOnSetup {
		panic("panic in test-before-plugin")
	}
}

func (p *testBeforePlugin) Execute(event *types.Event) *types.Event {
	args := p.Called(event)

	if p.raisePanicOnExecute {
		panic("panic in test-before-plugin")
	}

	return args[0].(*types.Event)
}

type testEnrichmentPlugin struct {
	mock.Mock
	raisePanicOnSetup   bool
	raisePanicOnExecute bool
}

func (p *testEnrichmentPlugin) Name() string {
	return "test-enrichment-plugin"
}

func (p *testEnrichmentPlugin) Type() types.PluginType {
	return types.PluginTypeEnrichment
}

func (p *testEnrichmentPlugin) Setup(config types.Config) {
	p.Called(config)

	if p.raisePanicOnSetup {
		panic("panic in test-enrichment-plugin")
	}
}

func (p *testEnrichmentPlugin) Execute(event *types.Event) *types.Event {
	args := p.Called(event)

	if p.raisePanicOnExecute {
		panic("panic in test-enrichment-plugin")
	}

	return args[0].(*types.Event)
}

type testDestinationPlugin struct {
	mock.Mock
	raisePanicOnSetup    bool
	raisePanicOnExecute  bool
	raisePanicOnFlush    bool
	raisePanicOnShutdown bool
}

func (p *testDestinationPlugin) Name() string {
	return "test-destination-plugin"
}

func (p *testDestinationPlugin) Type() types.PluginType {
	return types.PluginTypeDestination
}

func (p *testDestinationPlugin) Setup(config types.Config) {
	p.Called(config)

	if p.raisePanicOnSetup {
		panic("panic in test-destination-plugin")
	}
}

func (p *testDestinationPlugin) Execute(event *types.Event) {
	p.Called(event)

	if p.raisePanicOnExecute {
		panic("panic in test-destination-plugin")
	}
}

func (p *testDestinationPlugin) Flush() {
	p.Called()

	if p.raisePanicOnFlush {
		panic("panic in test-destination-plugin")
	}
}

func (p *testDestinationPlugin) Shutdown() {
	p.Called()

	if p.raisePanicOnShutdown {
		panic("panic in test-destination-plugin")
	}
}

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Debugf(message string, args ...interface{}) {
	l.Called(message, args)
}

func (l *mockLogger) Infof(message string, args ...interface{}) {
	l.Called(message, args)
}

func (l *mockLogger) Warnf(message string, args ...interface{}) {
	l.Called(message, args)
}

func (l *mockLogger) Errorf(message string, args ...interface{}) {
	l.Called(message, args)
}
