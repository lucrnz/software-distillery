# node-nvm: switch builder Python to 3.13.11 (pyenv 2.6.16)

- Updated `node-nvm/Dockerfile.alpine` and `node-nvm/Dockerfile.debian` builder stages to download Python 3.13.11 built with pyenv 2.6.16 from release `released-20251209-1917`.
- Adjusted release/tag variables and tarball URL pattern to the new artifact naming for both Alpine and Debian/Ubuntu variants while keeping architecture detection (`amd64`/`aarch64`) intact.
- Goal: resolve NVM build issues by aligning the builder’s Python with the latest published pyenv artifacts; downloads remain via ripvex for hash-verified fetches.

