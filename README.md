# Pulumi Native Provider Boilerplate

This repository is a boilerplate showing how to create and locally test a native Pulumi provider **generated from an OpenAPI 3.0 spec**.

## Authoring a Pulumi Native Provider

This boilerplate lets you create a working native Pulumi provider named `xyz`.
You'll need to update `provider/pkg/provider/provider.go` and implement the necessary CRUD functions from Pulumi's `ResourceProviderServer` interface.

### Prerequisites

Prerequisites for this repository are already satisfied by the [Pulumi Devcontainer](https://github.com/pulumi/devcontainer) if you are using Github Codespaces, or VSCode.

If you are not using VSCode, you will need to ensure the following tools are installed and present in your `$PATH`:

-   [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
-   [Go 1.21](https://golang.org/dl/) or 1.latest
-   [NodeJS](https://nodejs.org/en/) 14.x. We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
-   [Yarn](https://yarnpkg.com/)
-   [TypeScript](https://www.typescriptlang.org/)
-   [Python](https://www.python.org/downloads/) (called as `python3`). For recent versions of MacOS, the system-installed version is fine.
-   [.NET](https://dotnet.microsoft.com/download)

### Build & test the boilerplate XYZ provider

1. Create a new Github CodeSpaces environment using this repository.
1. Open a terminal in the CodeSpaces environment.
1. Run `make build install` to build and install the provider.
1. Run `make gen_examples` to generate the example programs in `examples/` off of the source `examples/yaml` example program.
1. Run `make up` to run the example program in `examples/yaml`.
1. Run `make down` to tear down the example program.

### Creating a new provider repository

Pulumi offers this repository as a [GitHub template repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) for convenience. From this repository:

1. Click "Use this template".
1. Set the following options:
    - Owner: pulumi
    - Repository name: pulumi-xyz-native (replace "xyz" with the name of your provider)
    - Description: Pulumi provider for xyz
    - Repository type: Public
1. Clone the generated repository.

From the templated repository:

1. Search-replace `xyz` with the name of your desired provider.

#### Build the provider and install the plugin

```bash
$ make build install
```

This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### A brief repository overview

You now have:

1. A `provider/` folder containing the building and implementation logic
    1. `cmd/pulumi-resource-xyz/main.go` - holds the provider's sample implementation logic.
2. `deployment-templates` - a set of files to help you around deployment and publication
3. `sdk` - holds the generated code libraries created by `pulumi-gen-xyz/main.go`
4. `examples` a folder of Pulumi programs to try locally and/or use in CI.
5. A `Makefile` and this `README`.

#### Additional Details

This repository depends on [`pulschema`](https://github.com/cloudy-sky-software/pulschema) library to handle generating a Pulumi schema from an OpenAPI 3.x spec. For a successful schema generation, you should ensure that your OpenAPI spec is valid and that it conforms to certain expectations. Learn more at https://github.com/cloudy-sky-software/cloud-provider-api-conformance.

### Build Examples

Create an example program using the resources defined in your provider, and place it in the `examples/` folder.

You can now repeat the steps for [build, install, and test](#test-against-the-example).

## Configuring CI and releases

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

## References

Other resources/examples for implementing providers:

-   [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
-   [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
