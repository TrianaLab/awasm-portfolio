package cmd_test

import (
	"bytes"
	"github.com/TrianaLab/awasm-portfolio/cmd"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupTestCommand(t *testing.T, command *cobra.Command, args []string) (string, error) {
	buf := new(bytes.Buffer)
	command.SetOut(buf)
	command.SetErr(buf)
	command.SetArgs(args)

	err := command.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, err := setupTestCommand(t, rootCmd, []string{"help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "A CLI tool for managing resources") {
		t.Errorf("expected help output, got: %s", output)
	}
}

func TestCreateCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// Valid creation
	output, err := setupTestCommand(t, rootCmd, []string{"create", "namespace", "testns"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "namespace/testns created") {
		t.Errorf("expected namespace creation success, got: %s", output)
	}

	// Invalid creation (duplicate)
	output, err = setupTestCommand(t, rootCmd, []string{"create", "namespace", "testns"})
	if err == nil && !strings.Contains(output, "already exists") {
		t.Errorf("expected duplicate creation error, got: %s", output)
	}
}

func TestDeleteCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// Precreate a namespace
	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "testns"})

	// Valid deletion
	output, err := setupTestCommand(t, rootCmd, []string{"delete", "namespace", "testns"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "namespace/testns in namespace '' deleted") {
		t.Errorf("expected namespace deletion success, got: %s", output)
	}

	// Invalid deletion
	output, err = setupTestCommand(t, rootCmd, []string{"delete", "namespace", "nonexistent"})
	if err == nil && !strings.Contains(output, "not found") {
		t.Errorf("expected not found error, got: %s", output)
	}
}

func TestGetCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// Precreate a namespace
	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "testns"})

	// Valid get
	output, err := setupTestCommand(t, rootCmd, []string{"get", "namespace", "testns"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "testns") {
		t.Errorf("expected namespace in output, got: %s", output)
	}

	// Invalid kind
	output, err = setupTestCommand(t, rootCmd, []string{"get", "invalidkind"})
	if err == nil && !strings.Contains(output, "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got: %s", output)
	}
}

func TestDescribeCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// Precreate a namespace
	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "testns"})

	// Valid describe
	output, err := setupTestCommand(t, rootCmd, []string{"describe", "namespace", "testns"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "testns") {
		t.Errorf("expected namespace in describe output, got: %s", output)
	}

	// Nonexistent describe
	output, err = setupTestCommand(t, rootCmd, []string{"describe", "namespace", "nonexistent"})
	if err == nil && !strings.Contains(output, "not found") {
		t.Errorf("expected not found error, got: %s", output)
	}
}

func TestVersionCommand(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, err := setupTestCommand(t, rootCmd, []string{"version"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "App version:") {
		t.Errorf("expected 'App version:' in output, got: %s", output)
	}
}

func TestCreateCommand_OutputFlagRejected(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, _ := setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns1", "-o", "yaml"})
	if !strings.Contains(output, "output flag is not supported") {
		t.Errorf("expected output flag rejection, got: %s", output)
	}
}

func TestCreateCommand_AllNamespacesFlag(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, err := setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns-a", "-A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "namespace/ns-a created") {
		t.Errorf("expected namespace creation with -A flag, got: %s", output)
	}
}

func TestDeleteCommand_OutputFlagRejected(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns1"})
	output, _ := setupTestCommand(t, rootCmd, []string{"delete", "namespace", "ns1", "-o", "yaml"})
	if !strings.Contains(output, "output flag is not supported") {
		t.Errorf("expected output flag rejection, got: %s", output)
	}
}

func TestDeleteCommand_AllNamespacesFlag(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns-del"})
	output, err := setupTestCommand(t, rootCmd, []string{"delete", "namespace", "ns-del", "-A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "deleted") {
		t.Errorf("expected deletion with -A flag, got: %s", output)
	}
}

func TestDescribeCommand_OutputFlagRejected(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, _ := setupTestCommand(t, rootCmd, []string{"describe", "namespace", "ns1", "-o", "yaml"})
	if !strings.Contains(output, "output flag is not supported") {
		t.Errorf("expected output flag rejection, got: %s", output)
	}
}

func TestDescribeCommand_NoArgs(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// No args triggers the Help path
	output, err := setupTestCommand(t, rootCmd, []string{"describe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Describe") {
		t.Errorf("expected help output, got: %s", output)
	}
}

func TestDescribeCommand_AllNamespacesFlag(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns-desc"})
	// -A with kind-only (no name) lists all namespaces across the cluster.
	output, err := setupTestCommand(t, rootCmd, []string{"describe", "namespace", "-A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "ns-desc") {
		t.Errorf("expected namespace in output with -A flag, got: %s", output)
	}
}

func TestGetCommand_WrongOutputFormat(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, _ := setupTestCommand(t, rootCmd, []string{"get", "namespace", "-o", "xml"})
	if !strings.Contains(output, "wrong format") {
		t.Errorf("expected wrong-format error, got: %s", output)
	}
}

func TestGetCommand_AllNamespacesFlag(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	_, _ = setupTestCommand(t, rootCmd, []string{"create", "namespace", "ns-get"})
	output, err := setupTestCommand(t, rootCmd, []string{"get", "namespace", "-A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "ns-get") {
		t.Errorf("expected namespace listed with -A flag, got: %s", output)
	}
}

func TestGetCommand_NoArgs(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	output, err := setupTestCommand(t, rootCmd, []string{"get"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Get resources") {
		t.Errorf("expected help output, got: %s", output)
	}
}

func TestGetCommand_ServiceError(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// Asking for a name in an empty repo triggers the service error path.
	output, _ := setupTestCommand(t, rootCmd, []string{"get", "namespace", "nonexistent"})
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error output, got: %s", output)
	}
}

func TestRootCommand_HelpFlag(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	rootCmd := cmd.NewRootCommand(repo)

	// --help triggers the custom help function which lists resources & aliases.
	output, err := setupTestCommand(t, rootCmd, []string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(output, "Available Commands") {
		t.Errorf("expected help output, got: %s", output)
	}
	if !strings.Contains(output, "namespace") {
		t.Errorf("expected resource listing in help, got: %s", output)
	}
}
