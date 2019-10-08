Thank you for considering making contributions to the Cosmos modules repositories!

### Adding a Module:
To get your module added to this repo make sure you meet the following requirements:

1. Your module has a unique name from any other module in the `cosmos/modules` repository
2. Your module contains a `LICENSE` AND `README`
3. Your module has thorough documentation in a `docs/` subfolder
4. Your module has its own `go.mod` and `go.sum` files
5. Your module has unit tests along with a `simapp.go` and `sim_test.go` for fuzz testing
6. Your module contains a `CODEOWNERS.md` file

### Module Specific Issues/PRs

To make an issue or pull request to a specific module (e.g. `modules/poa`), prefix the issue/pr name with the module in question. For example, we can create a PR to update SDK version in poa module like so, `poa: update SDK version`.
