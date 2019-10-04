# Cosmos Modules

This repo is meant to house modules that are created outside of the [Cosmos-SDK](https://github.com/cosmos/cosmos-sdk).

To find out more about the Cosmos-SDK, you can find documentation [here](https://cosmos.network/docs/).

This repo organizes modules into 3 subfolders:

- `stable/`: this folder houses modules that are stable, production-ready, and well-maintained.
- `incubator/`: this folder houses modules that are buildable but makes no guarantees on stability or production-readiness. Once a module meets all requirements specified in [contributing guidelines](./CONTRIBUTING.md), the owners can make a PR to move module into `stable/` folder. Must be approved by at least one `modules` maintainer for the module to be moved.
- `inactive/`: Any stale module from the previous 2 folders may be moved to the `inactive` folder if it is no longer being maintained by its owners. `modules` maintainers reserve the right to move a module into this folder after public discussion in an issue and a specified grace period for module owners to restart work on module.

### Modules maintainers

While each individual module will be owned and maintained by the individual contributors of that module, there will need to be maintainers of the `modules` repo overall to coordinate moving modules between the different folders and enforcing the requirements for inclusion in the `modules` repo.

For now, the maintainers of the `modules` repo will be the SDK team but we intend to eventually expand this responsibility to other members of the Cosmos community.
