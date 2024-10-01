# Version Control Activity Analyzer

<p align="center">
    <img src="assets/logo.png" alt="Logo" />
</p>

A command-line tool for analyzing version control activity in Git repositories. This tool provides insights into commit history, identifies bottlenecks, and detects anti-patterns in a project's version control workflow.

## Features

- **Commit History Analysis**: Get detailed statistics about the commit history, including the total number of commits and contributions by each developer.
- **Bottleneck Identification**: Identify potential bottlenecks in the workflow, such as long-lived branches and infrequent commits.
- **Anti-Pattern Detection**: Detect common anti-patterns like large commits, force pushes, and other practices that may hinder collaboration and code quality.

## Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/doc/install) (version 1.18 or higher)

## Installation

Clone the repository
```
git clone https://github.com/adigulalkari/VC-Analyzer.git
cd VC-Analyzer
```
Run main
```
chmod +x build.sh
./build.sh
```

## Contributing
Contributions are welcome! Please follow these steps to contribute to the project:

- Fork the repository.
- Create a new branch: ```git checkout -b feature-branch```
- Make your changes and commit them: ```git commit -m "Add new feature"```
- Push to the branch: ```git push origin feature-branch```
- Open a pull request.

Refer to [CONTRIBUTING.md](https://github.com/adigulalkari/VC-Analyzer/blob/main/CONTRIBUTING.md) for more guidelines!

## LICENSE
See the [LICENSE](https://github.com/adigulalkari/VC-Analyzer/blob/main/LICENSE) file for license rights and limitations (MIT).

