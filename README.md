# Agentic Dev Starter

An minimal and experimental micro-framework for prompting Cline+Claude (or similar) to create custom workflows.

## Contents
- [Purpose](#purpose)
- [The Files](#the-files)
- [Getting Started](#getting-started)
- [Requirements](#requirements)
- [Development Continuity Prompt](#development-continuity-prompt)
- [Contrib](#contrib)

## Workflows

This micro-framework provides a set of workflow definitions that each comprise the following:

- A prompt, to be supplied as the "custom instructions" for the agentic extension in use (for example Cline)
- A set of supporting files that are intended to help the LLM to maintain state and direction between sessions, and to augment and provide context for general agentic functionaltiy.

A workflow can be very simple, for example the "context-history-commit" workflow (see below), or potentially more complex involving MCP tools to create some very useful agentic flows.

This is an experimental WIP, please test and use accordinly.

### 1. "context-history-commit"

Named to indicate the basic flow of: read context and history, perform user directed tasks, update the history, commit and finish.

See the [workflow README](workflows/context-history-commit/README.md) for detailed instructions.

Key files:
- **DEVLOG.md**: LLM authored: Tracks development progress and changes between sessions. The LLM updates this file after each task (prior to commit) to maintain a more detailed (if necessary) history than a commit log. This helps retain reasoning chains and associated conclusions and decisions where necessary.
- **CONTEXT.md**: HUMAN authored: Provides high-level project-specific context, goals, constraints and requirements to guide development, this should change rarely.
- **prompt.md**: The specific prompt to be supplied to the custom instructions in your agentic extension.

## Requirements

- Git
- VSCode or similar dev environment that supports chat-based coding agents
- An agentic coding assistant that allows custom prompts in settings, for example Cline
- An API key for Claude, Deepseek or whatever model you like to work with.

## Contrib

Feedback, ideas, and pull requests are welcome.
