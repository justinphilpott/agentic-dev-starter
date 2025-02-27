# Agentic Dev Starter

A minimal framework for effective AI-assisted development.

## TLDR
A simple framework to maintain context between AI assistant sessions using just two files (DEVSTATE.md, PLAN.md) and a prompt (PROMPT.md). Helps keep AI-assisted development focused and efficient by preserving state between sessions while keeping each interaction targeted on specific tasks and goals.

## Contents
- [Purpose](#purpose)
- [The Problem](#the-problem)
- [A Solution](#a-solution)
- [The Files](#the-files)
- [Getting Started](#getting-started)
- [Requirements](#requirements)
- [Development Continuity Prompt](#development-continuity-prompt)
- [Usage Example](#usage-example)
- [Contrib](#contrib)

## Purpose

This micro-framework provides a simple set of files and a prompt template to streamline LLM-augmented development workflows.

(Tested with Cline and Claude 3.x Sonnet/Haiku)

## The problem

When working with chat-based coding agents, the length of the prompt increases with every message in the same session, making each request progressively more expensive.

Long-running chats can reach a point of "diminished returns" when focus, direction, or subject changes within the chat.

It makes sense to work with short sessions focused on atomic tasks. However, when working toward a larger goal, we need to retain state *between* tasks by:
- Prompting the agent on how to start and finish tasks
- Prompting the agent on how to read and write state

## A solution

This framework uses just two files (DEVSTATE.md, PLAN.md), one tool (git), and one prompt (in PROMPT.md) that you add to your agent's settings to be read on every request.

## The files and their purpose

- **DEVSTATE.md**: The LLM overwrites this file after completing each task to record the current development status for the next session.
- **PLAN.md**: User-managed file detailing short-term goals, constraints, requirements, and focused tasks. The LLM may update this file, but the user ensures the content remains concise and focused.
- **PROMPT.md**: Contains the prompt to add to your AI assistant's settings.

## Getting Started

1. **Choose your approach**:
   ```bash
   # Either clone this starter repository
   git clone https://github.com/justinphilpott/agentic-dev-starter.git
   cd agentic-dev-starter
   
   # Or copy the files to your existing project
   # Just copy PLAN.md and DEVSTATE.md to your project
   ```

2. **Set up your project**:
   - Customize PLAN.md with your project goals, constraints, and initial tasks
   - Initialize DEVSTATE.md with your project's starting state

3. **Configure your AI assistant**:
   - Copy the prompt from PROMPT.md into your AI assistant's settings
   - Ensure your assistant has access to read and write files in your project

## Requirements

- Git
- VSCode or similar dev environment that supports chat-based coding agents
- An agentic coding assistant that allows custom prompts in settings

## Usage Example

1. **Initialize your project**:
   - Create PLAN.md with your initial goals and tasks
   - Create an empty DEVSTATE.md file

2. **Start your first development session**:
   - Open your project in VSCode with your AI assistant
   - Ask your assistant to help with a specific task from PLAN.md
   - The assistant will follow the Task Start Protocol to understand the current state

3. **Complete the task**:
   - Work with the assistant to implement the solution
   - The assistant will follow the Task Completion Protocol
   - DEVSTATE.md will be updated with the current status
   - A git commit will be suggested

4. **Continue development**:
   - In your next session, the assistant will read DEVSTATE.md to understand what's been done
   - Select the next task from PLAN.md
   - Repeat the process

This approach maintains context between development sessions while keeping each interaction focused and efficient.

## Contrib

Feedback, ideas, and pull requests are welcome.
