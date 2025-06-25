# process management

We use process compose for running and controlling binaries that we need.

https://github.com/F1bonacc1/process-compose


These could be anything, allowing use to use XTask as a foundation for any projects running on a server or a laptops.

On Fly.io, we use it.

## WASM

We need to allow Devs to bring their Logic to the system, and so we use WASM for that.  

Our WASM runner is Wazero. https://github.com/tetratelabs/wazero

We run wasm in the browser and on the Servers.

1. Rendering off the Decksh as wasm to the DOM as SVG

2. Logic layer built in WASM, and brought in via nats





