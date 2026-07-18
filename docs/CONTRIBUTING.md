# Contributing Guide

Thank you for your interest in contributing to FlowSync! We welcome bug reports, feature suggestions, documentation updates, and pull requests.

---

## 🤝 Code of Conduct

By participating in this project, you agree to maintain a respectful, welcoming, and inclusive community environment.

---

## 🛠️ How to Contribute

### 1. Reporting Bugs
*   Ensure the bug is reproducible with the latest version of the main branch.
*   Check existing issues to ensure it hasn't already been reported.
*   Open a new issue with a clear description, reproduction steps, expected behavior, and relevant logs from `flowsync status` or the daemon output.

### 2. Suggesting Enhancements
*   Open an issue explaining the proposed feature.
*   Describe the use case, why this enhancement is valuable to users, and how it fits into the existing client-daemon architecture.

### 3. Submitting Pull Requests
*   **Fork the repo** and create a branch from the `main` branch.
*   Ensure all code conforms to the standard Go style and format guidelines.
*   Include unit tests for any new features or bug fixes.
*   Verify that all tests pass: `go test ./...`
*   Submit the PR with a clear summary of your changes.

---

## 📏 Coding Standards & Formatting

To maintain code quality across the codebase, we enforce the following rules:

*   **Go Code Formatting**: All Go code must be formatted using `go fmt`. Run the following before committing:
    ```bash
    go fmt ./...
    ```
*   **Linting**: Run Go vet to inspect common software issues:
    ```bash
    go vet ./...
    ```
*   **Documentation**: Ensure all exported types, variables, functions, and packages have appropriate Go docstrings.
*   **Small, Focused Pull Requests**: Keep pull requests focused on a single issue or feature. Avoid bundling unrelated refactorings or cosmetic changes into a logic change.
