# Contribution Guidelines
Please read this guide if you plan to contribute to the Capten. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with the Capten, please [open an issue on GitHub](https://github.com/intelops/capten/issues/new?labels=kind%2Fbug&template=bug-report.md&title=Bug:) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/intelops/capten/issues/new?labels=kind%2Ffeature&template=feature-request.md&title=Feature%20Request:).

## Contributing Code
This project is written in Golang 

To contribute code.
1. Ensure you are running golang version 1.21 or greater for go module support
2. Set the following environment variables:
    ```
    GO111MODULE=on
    GOFLAGS=-mod=vendor
    ```
3. Fork the project.
4. Clone the project: `git clone https://github.com/[YOUR_USERNAME]/capten && cd capten`
5. kindly refer capten.md file to know the structure of the project.
6. Commit changes *([Please refer the commit message conventions](https://www.conventionalcommits.org/en/v1.0.0/))*
7. Push commits.
8. Open pull request.

