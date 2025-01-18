package cmd_test

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/repository"
	"bytes"
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
