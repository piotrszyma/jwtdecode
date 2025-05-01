# jwtdecode

A simple command-line tool to decode the header and claims of a JSON Web Token (JWT) without performing signature verification. This is useful for inspecting the contents of a JWT.

## Overview

This utility takes a JWT as input and prints its header and payload (claims) in a human-readable JSON format with syntax highlighting. It also attempts to interpret and display timestamp claims (like `iat`, `exp`, `nbf`) in a more informative way, showing the absolute time in ISO 8601 format and a relative time difference from the current time.

```bash
jwtdecode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
```

**Example Output:**

```sh
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload:
{
  "iat": 1516239022, # 2018-01-18T01:30:22Z [2660 day(s) ago]
  "name": "John Doe",
  "sub": "1234567890"
}
```

## Installation

To install `jwtdecode`, you need to have Go installed on your system. Open your terminal and run the following command:

```bash
go install github.com/piotrszyma/jwtdecode@latest
```

This command will download and install the `jwtdecode` executable in your `$GOPATH/bin` directory (or `$HOME/go/bin` if you are using Go 1.18 or later with default settings). Make sure this directory is included in your system's `PATH` environment variable so you can run the `jwtdecode` command directly.

## Usage

Once installed, you can use the `jwtdecode` command followed by the JWT you want to inspect.

```bash
jwtdecode <your_jwt_token>
```

Replace `<your_jwt_token>` with the actual JWT string. The tool will then output the decoded header and payload to your terminal.

## Important Notes

- **No Signature Verification:** This tool does **not** verify the signature of the JWT. It simply decodes the base64-encoded header and payload sections. Therefore, the information displayed might be tampered with or invalid.
- **Timestamp Interpretation:** The tool attempts to recognize and format standard JWT timestamp claims (like `iat`, `exp`, `nbf`). The relative time difference is calculated based on the current system time when the command is executed.
- **Error Handling:** The tool provides basic error handling for invalid JWT formats or decoding issues.

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests on the GitHub repository: [github.com/piotrszyma/jwtdecode](https://www.google.com/search?q=https://github.com/piotrszyma/jwtdecode).
