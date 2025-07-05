# GoCSS Example Project

This project demonstrates how to use GoCSS with Templ and Chi.

## Setup

1.  **Install Templ CLI:**
    ```bash
    go install github.com/a-h/templ/cmd/templ@latest
    ```

2.  **Build GoCSS CLI:**
    Navigate to the root of the `gocss` project and build the CLI:
    ```bash
    cd ../
    go build -o ./cmd/gocss/main.go ./cmd/gocss
    cd examples
    ```

## Running the Example

1.  **Navigate to the `examples` directory:**
    ```bash
    cd examples
    ```

2.  **Run the build script:**
    ```bash
    ./build.sh
    ```

    This script will:
    *   Generate Go code from Templ templates.
    *   Generate `gocss.css` using the local `gocss` CLI.
    *   Start the Go web server.

3.  **Open your browser:**
    Visit `http://localhost:3000` to see the example in action.

## Development (with Watch Mode)

To enable watch mode for both Templ and GoCSS during development:

1.  **In one terminal, run Templ watch:**
    ```bash
    templ generate -path templates -watch
    ```

2.  **In another terminal, run GoCSS watch:**
    ```bash
    ../cmd/gocss/main.go --input "./templates/*.templ" --output ./static/gocss.css --watch
    ```

3.  **In a third terminal, run the Go application:**
    ```bash
    go run .
    ```

Now, any changes to `*.templ` files will trigger Templ to regenerate Go code, and GoCSS to regenerate CSS, and the Go server will automatically pick up the changes (you might need to refresh your browser).
