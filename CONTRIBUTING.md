# Contributing to Adaptive OS Mesh

First off, thank you for considering contributing! This project is an experimental framework, and we welcome any and all contributions.

## Development Workflow

1.  **Fork & Clone**: Fork the repository to your own GitHub account and then clone it to your local machine.
2.  **Create a Branch**: Create a new branch for your feature or bug fix.
    ```bash
    git checkout -b your-feature-name
    ```
3.  **Develop**: Make your changes. Ensure you follow the existing code style and architectural patterns.
4.  **Test**: Run the test suites to ensure your changes haven't broken anything.
    ```bash
    # Run backend tests
    go test ./...

    # Run frontend/integrity tests
    npm test
    ```
5.  **Commit**: Write a clear, concise commit message.
6.  **Push**: Push your branch to your fork.
7.  **Open a Pull Request**: Open a PR against the `main` branch of the original repository.

## Code Style

- **Go**: We follow standard `gofmt` and `golint` conventions.
- **TypeScript/Node.js**: We use `prettier` for code formatting.

## Submitting a Pull Request

- Provide a clear title and description for your pull request.
- Link to any relevant issues.
- Explain the "why" and "what" of your changes.

Thank you for your contribution!
