package plugins

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	"github.com/dop251/goja"
// 	"github.com/os-webui/os-webui/internal/utils"
// )

// var DefaultPluginsManager = NewPluginsManager()

// type PluginsManager struct {
// 	plugins map[string]*Plugin // Updated to store pointers for safe mutations
// }

// func NewPluginsManager() *PluginsManager {
// 	return &PluginsManager{
// 		plugins: map[string]*Plugin{},
// 	}
// }

// // Info maps the JavaScript self-describing manifest object schema
// type Info struct {
// 	ID          string `json:"id"`
// 	Name        string `json:"name"`
// 	Version     string `json:"version"`
// 	Description string `json:"description"`
// 	Author      string `json:"author"`
// }

// type Plugin struct {
// 	ID      string
// 	Install string
// 	Data    string
// 	Config  string
// 	Info    Info          // Cached structural metadata
// 	Program *goja.Program // Compiled pre-parsed bytecode ready for VM execution
// }

// // Load compiles the single-file plugin script, verifies metadata, and stores it in the runtime map
// func (m *PluginsManager) Load(id, install, data, config string) (*Plugin, error) {
// 	scriptBytes, err := os.ReadFile(install)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read plugin script at %s: %w", install, err)
// 	}

// 	// Performance Optimization: Convert slice via zero-allocation and parse straight into compiled Goja bytecode
// 	scriptSource := utils.BytesToString(scriptBytes)
// 	program, err := goja.Compile(install, scriptSource, false)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to compile JavaScript syntax for plugin %s: %w", id, err)
// 	}

// 	// Initialize a temporary sandboxed VM context strictly to extract metadata
// 	vm := goja.New()

// 	// Stub global helper mocks to prevent structural reference crashes during manifest extraction
// 	_ = vm.Set("log", func(msg string) {})

// 	// Execute the compiled program inside the scratchpad VM to register the global variables
// 	_, err = vm.RunProgram(program)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize script execution context: %w", err)
// 	}

// 	// Extract the global self-describing 'pluginInfo' object
// 	jsInfo := vm.Get("pluginInfo")
// 	if jsInfo == nil || goja.IsUndefined(jsInfo) {
// 		return nil, fmt.Errorf("missing mandatory global 'pluginInfo' object in %s", install)
// 	}

// 	// Safely transform the JavaScript object via native JSON stringification to cross type boundaries safely
// 	jsJSON := vm.Get("JSON")
// 	if jsJSON == nil {
// 		return nil, fmt.Errorf("failed to locate global JSON engine")
// 	}
// 	jsonStringify, ok := goja.AssertFunction(jsJSON.(*goja.Object).Get("stringify"))
// 	if !ok {
// 		return nil, fmt.Errorf("failed to assert global JSON.stringify function")
// 	}

// 	jsonStrVal, err := jsonStringify(goja.Undefined(), jsInfo)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to stringify pluginInfo object: %w", err)
// 	}

// 	var info Info
// 	if err = json.Unmarshal(utils.StringToBytes(jsonStrVal.String()), &info); err != nil {
// 		return nil, fmt.Errorf("failed to parse schema mapping for pluginInfo: %w", err)
// 	}

// 	// Instantiate operational plugin runtime container
// 	p := &Plugin{
// 		ID:      id,
// 		Install: install,
// 		Data:    data,
// 		Config:  config,
// 		Info:    info,
// 		Program: program,
// 	}

// 	// Register the loaded instance into the state map
// 	m.plugins[id] = p

// 	return p, nil
// }
