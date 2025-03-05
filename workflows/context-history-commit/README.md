# Basic flow - "context-history-commit"

Agentic directive overview:

- Read the project history from DEVLOG.md
- Read the project context from CONTEXT.md 
- Proceed according to the users instructions
- Ensure tests pass prior to completion.
- Summarise in detail any thought processes, conclusions and specific actions taken in this session, in the DEVLOG.
- Attempt to commit work.
- Complete.

Files:

- DEVLOG.md - goes in your root folder
- CONTEXT.md - goes in your root folder (fill this out with your project context as needed)
- prompt.md - to be supplied to your agent as custom instructions (see agent settings)

Suggestions:
- Consider setting a pre-commit hook to run tests.
- If your test suite is huge, ensure tests can run in parallel.

Disclaimer:
- Please ensure that the prompt is worded to your requirements, needs and project context. Perform experiments in an inconsequential context before using this for important real work.
