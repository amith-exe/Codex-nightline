package main

import "testing"

func TestCommandLineHasFlagValueRequiresExactTokenBoundary(t *testing.T) {
	commandLine := "flowsyncd -root /tmp/flowsync-other -socket=/tmp/flowsync/daemon.sock-extra"

	if commandLineHasFlagValue(commandLine, "-root", "/tmp/flowsync") {
		t.Fatal("matched root value prefix from longer flag value")
	}
	if commandLineHasFlagValue(commandLine, "-socket", "/tmp/flowsync/daemon.sock") {
		t.Fatal("matched socket value prefix from longer flag value")
	}
}

func TestCommandLineHasFlagValueMatchesSeparateAndEqualsForms(t *testing.T) {
	commandLine := "flowsyncd -root /tmp/flowsync -socket=/tmp/flowsync/daemon.sock"

	if !commandLineHasFlagValue(commandLine, "-root", "/tmp/flowsync") {
		t.Fatal("did not match separate flag value")
	}
	if !commandLineHasFlagValue(commandLine, "-socket", "/tmp/flowsync/daemon.sock") {
		t.Fatal("did not match equals flag value")
	}
}
