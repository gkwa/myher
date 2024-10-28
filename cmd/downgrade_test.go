package cmd

import (
   "bytes"
   "fmt"
   "io"
   "os"
   "strings"
   "testing"

   "github.com/gkwa/myher/core/gomod"
   "github.com/gkwa/myher/internal/logger"
)

type mockGoModService struct {
   err error
}

func (m *mockGoModService) GetModuleInfo() (*gomod.ModuleInfo, error) {
   return &gomod.ModuleInfo{}, nil
}

func (m *mockGoModService) PrettyPrint(info *gomod.ModuleInfo) {
}

func (m *mockGoModService) GenerateDowngradeCommands(concurrency int, alternating bool) ([]string, error) {
   if m.err != nil {
   	return nil, m.err
   }
   if alternating {
   	return []string{
   		"go get github.com/pkg/errors@v0.9.0",
   		"# go get github.com/stretchr/testify@v1.8.3",
   	}, nil
   }
   return []string{
   	"go get github.com/pkg/errors@v0.9.0",
   	"go get github.com/stretchr/testify@v1.8.3",
   }, nil
}

func TestDowngradeCommand(t *testing.T) {
   // Save the original service constructor
   originalNewService := gomod.NewService
   defer func() {
   	gomod.NewService = originalNewService
   }()

   tests := []struct {
   	name     string
   	args     []string
   	wantErr  bool
   	contains []string
   }{
   	{
   		name: "basic downgrade",
   		args: []string{"downgrade"},
   		contains: []string{
   			"github.com",
   			"@v",
   		},
   	},
   	{
   		name: "with concurrency",
   		args: []string{"downgrade", "--concurrent", "3"},
   		contains: []string{
   			"github.com",
   			"@v",
   		},
   	},
   	{
   		name: "with alternating comments",
   		args: []string{"downgrade", "--enable-alternating-comments"},
   		contains: []string{
   			"github.com",
   			"@v",
   			"# go get",
   		},
   	},
   }

   for _, tt := range tests {
   	t.Run(tt.name, func(t *testing.T) {
   		// Override the service constructor for each test
   		mockSvc := &mockGoModService{}
   		gomod.NewService = func(logger gomod.Logger) gomod.Service {
   			return mockSvc
   		}

   		oldStdout := os.Stdout
   		r, w, _ := os.Pipe()
   		os.Stdout = w

   		oldVerbose, oldLogFormat := verbose, logFormat
   		verbose, logFormat = 1, "text"
   		defer func() {
   			verbose, logFormat = oldVerbose, oldLogFormat
   		}()

   		customLogger := logger.NewConsoleLogger(verbose, logFormat == "json")
   		cliLogger = customLogger

   		rootCmd.SetArgs(tt.args)
   		execErr := rootCmd.Execute()

   		w.Close()
   		os.Stdout = oldStdout

   		var buf bytes.Buffer
   		if _, err := io.Copy(&buf, r); err != nil {
   			t.Fatalf("Failed to copy output: %v", err)
   		}
   		output := buf.String()

   		if (execErr != nil) != tt.wantErr {
   			t.Errorf("Execute() error = %v, wantErr %v", execErr, tt.wantErr)
   			return
   		}

   		for _, s := range tt.contains {
   			if !strings.Contains(output, s) {
   				t.Errorf("Output should contain %q but got:\n%s", s, output)
   			}
   		}

   		if tt.name == "with alternating comments" {
   			lines := strings.Split(strings.TrimSpace(output), "\n")
   			var hasCommented bool
   			for _, line := range lines {
   				if strings.HasPrefix(line, "# go get") {
   					hasCommented = true
   					break
   				}
   			}
   			if !hasCommented {
   				t.Error("Expected at least one commented line with --enable-alternating-comments")
   			}
   		}

   		t.Logf("Command output:\n%s", output)
   	})
   }
}

func TestDowngradeFlags(t *testing.T) {
   tests := []struct {
   	name        string
   	args        []string
   	wantErr     bool
   	checkValues func() error
   }{
   	{
   		name: "default values",
   		args: []string{"downgrade"},
   		checkValues: func() error {
   			if concurrent != 5 {
   				return fmt.Errorf("expected concurrent=5, got %d", concurrent)
   			}
   			if alternatingComments != false {
   				return fmt.Errorf("expected alternatingComments=false, got %v", alternatingComments)
   			}
   			return nil
   		},
   	},
   	{
   		name: "set concurrent",
   		args: []string{"downgrade", "--concurrent", "10"},
   		checkValues: func() error {
   			if concurrent != 10 {
   				return fmt.Errorf("expected concurrent=10, got %d", concurrent)
   			}
   			return nil
   		},
   	},
   	{
   		name: "set alternating comments",
   		args: []string{"downgrade", "--enable-alternating-comments"},
   		checkValues: func() error {
   			if !alternatingComments {
   				return fmt.Errorf("expected alternatingComments=true, got false")
   			}
   			return nil
   		},
   	},
   	{
   		name:    "invalid concurrent value",
   		args:    []string{"downgrade", "--concurrent", "invalid"},
   		wantErr: true,
   	},
   }

   for _, tt := range tests {
   	t.Run(tt.name, func(t *testing.T) {
   		// Reset flags before each test
   		concurrent = 5
   		alternatingComments = false

   		rootCmd.SetArgs(tt.args)
   		err := rootCmd.Execute()

   		if (err != nil) != tt.wantErr {
   			t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
   			return
   		}

   		if tt.checkValues != nil {
   			if err := tt.checkValues(); err != nil {
   				t.Error(err)
   			}
   		}
   	})
   }
}
