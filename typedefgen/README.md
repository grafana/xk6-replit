# tsjson

**tsjson** is a tool that recursively parses TypeScript definition files (`.d.ts`) and generates a hierarchical JSON schema that includes type information along with JSDoc documentation. This output can then be used by other tools (for example, a Go REPL) to display documentation.

## Prerequisites

- [Node.js](https://nodejs.org/) (v14 or later is recommended)
- npm

## Installation

1. **Initialize an npm Project**

   If you haven't already, run:

   ```bash
   npm init -y
   ```

1. **Install Required Dependencies**

   Install the necessary packages:

   ```bash
   npm install ts-morph comment-parser
   npm install --save-dev @types/node
   npm install -g ts-node
   ```

## Usage

To generate the JSON documentation, run:

```bash
ts-node tsjson.ts [k6-type-definitions-directory] [output-directory-json-directory]
```

This command will:
- Process the `.d.ts` files in the `k6` directory.
- Generate a JSON file for each `.d.ts` file that preserves the hierarchical structure and includes JSDoc documentation.

## Output

The tool produces JSON files that mirror the structure of your TypeScript definitions. Each JSON file includes:

- **Type Kind** (e.g. class, interface, variable)
- **Signature**
- **Documentation** (extracted from JSDoc comments)
- **Members/Properties** (for classes, interfaces, or object literal types)
