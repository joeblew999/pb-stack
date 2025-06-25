# Open Cloud + Guigui Integration

This project integrates [`opencloud`](https://github.com/opencloud-eu/opencloud) (a Golang collaboration server) with [`Guigui`](https://github.com/hajimehoshi/guigui) (a Golang GUI framework).

## Goals

- Run an Open cloud server.
- Run an Open cloud cli.
- Run an Open cloud gui that 
 Integrate `opencloud` with `Guigui`.
- Add a Golang-based indexing solution as an alternative to Apache Tika.
- Add support for indexing and searching Markdown files.
- Add support for Decksh presentations.

## Taskfile

This project uses [Task](https://taskfile.dev/) for task automation. The `Taskfile.yml` is in the root of the workspace.

Dont forget that ":" inside echo statements break the taskfile.

## Binaries

This project produces three distinct binaries, each serving a specific purpose:

*   **`opencloud-gui`**: Runs the graphical user interface.
*   **`opencloud-server`**: Runs the OpenCloud collaboration server.
*   **`opencloud-cli`**: Provides command-line utilities for indexing and searching.



### Building the Binaries

You can build each binary individually using standard Go commands or use the provided Taskfile for automated builds.

### Running the Binaries

#### `opencloud-gui`
Runs the graphical user interface.

#### `opencloud-server` (Server)
Runs the collaboration server with configurable port and index directory.

#### `opencloud-cli`
Provides command-line interface for indexing and searching with support for various modes and options.

### Usage

Use the Task runner to see available commands and build the project. The project is set up to be developed within a mono repo structure.

## Search

`opencloud`'s search functionality is a key part of this integration.

### Current Implementation

- **Indexing:** By default, `opencloud` uses [Apache Tika](https://docs.opencloud.eu/docs/dev/server/Services/search/Search-info/) for indexing various file types.
- **Search Service:** The search service is built using [bleve](https://github.com/opencloud-eu/opencloud/tree/main/services/search).
- **Query Language:** It uses [KQL (Keyword Query Language)](https://github.com/opencloud-eu/opencloud/tree/main/pkg/kql).

### Integration Goals

The main goal is to enhance the search capabilities within the `Guigui` client by:

1.  **Adding a native Golang indexing solution:** This will provide an alternative to the Java-based Apache Tika, simplifying the stack.
2.  **Adding Markdown support:** The new indexer should support `.md` files.
3.  **Adding Decksh support:** Once Markdown support is working, we will add support for `Decksh`.
