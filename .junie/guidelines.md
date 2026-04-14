# Project Guidelines

## Development Workflow
1.  **TDD & Test Coverage**: Cover all new functionality with tests immediately. Test-Driven Development (TDD) is preferred.
2.  **Verification**: After implementing new functionality, run `make test` and `make build` to verify the state.
3.  **Git Commit**: If tests and build pass, commit the changes to the git repository.
4.  **Read Documentation**: Before implementing new functionality or making changes, read the existing documentation.

## Documentation and Knowledge Base
1.  **Knowledge Base (KB)**: Maintain a structured Knowledge Base in `.junie/kb`. 
    - Store useful information discovered during development or debugging: nuances, features, errors, solutions, and successful patterns.
    - Organize the KB so it's easy to search.
2.  **Project Documentation**: All project documentation (except KB) must be stored in the `docs/` directory.
3.  **Index**: The root `README.md` must contain an index for the `docs/` directory.

## Language and Communication
1.  **Code & Documentation**: All code comments and documentation must be in **English**.
2.  **User Communication**: All communication with the user must be in the **user's language** (Russian in this case).
