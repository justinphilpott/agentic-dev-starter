# Agentic Dev Starter

An minimal and experimental micro-framework for custom prompting Cline+Claude (or similar).

## Contents
- [Purpose](#purpose)
- [The Files](#the-files)
- [Getting Started](#getting-started)
- [Requirements](#requirements)
- [Development Continuity Prompt](#development-continuity-prompt)
- [Contrib](#contrib)

## Purpose

This micro-framework provides a simple set of files and a selection of prompts for assisting with different kinds of coding workflows. 

## The files

- **DEVLOG.md**: LLM authored: Tracks development progress and changes between sessions. The LLM updates this file after each task (prior to commit) to maintain a more detailed (if necessary) history than a commit log. This helps retain reasoning and thought process where necessary.
- **CONTEXT.md**: HUMAN authored: Provides project-specific context, constraints and requirements to guide development, this should change rarely.

## Prompts

### context-history-commit.md
This is a prompt to get your agent to obtain context and project history from the two file 
- Session start: read context (CONTEXT.md) and history (DEVLOG.md)
- Dev guidelines,
- Session finish: update DEVLOG.md, commit, complete.

## Getting Started

1. **Choose a prompt**:
   - Browse the prompts/ directory to select a workflow
   - Each prompt comes with its own set of files (CONTEXT.md, DEVLOG.md, etc.)

2. **Copy files to your project**:
   ```bash
   # Copy the selected prompt's files to your project root
   cp -r prompts/your-chosen-prompt/* .
   ```

3. **Set up your project**:
   - Customize CONTEXT.md with your project goals and constraints
   - Initialize DEVLOG.md with your project's starting state
   - Update your AI assistant's settings with the chosen prompt

4. **Start working**:
   - Open your project in VSCode with your AI assistant
   - The assistant will read CONTEXT.md for project requirements
   - The assistant will read DEVLOG.md to understand current state
   - Work with the assistant to implement solutions
   - The assistant will follow the Session Completion Protocol:
     * Update DEVLOG.md with progress
     * Commit changes with conventional commit messages

## Requirements

- Git
- VSCode or similar dev environment that supports chat-based coding agents
- An agentic coding assistant that allows custom prompts in settings, for example Cline

## Contrib

Feedback, ideas, and pull requests are welcome.
