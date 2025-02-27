# Agentic Dev Starter

A minimal framework for effective AI-assisted development.

## TLDR
A simple framework to maintain context between AI assistant sessions using just two files (DEVSTATE.md, PLAN.md) and a prompt. Helps keep AI-assisted development focused and efficient by preserving state between sessions while keeping each interaction targeted on specific tasks.

## Contents
- [Purpose](#purpose)
- [The Problem](#the-problem)
- [A Solution](#a-solution)
- [The Files](#the-files)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Development Continuity Prompt](#development-continuity-prompt)
- [Usage Example](#usage-example)
- [Contrib](#contrib)

## Purpose

The purpose of this micro-framework is to provide a simple set of files and a single prompt template that can help to streamline an LLM augmented development workflow.

(N.B. tested with Cline and Claude 3.x Sonnet/Haiku)

## The problem

When working with chat based coding agents such as Cline, the length of the prompt being sent to the API on each request increases with every chat message in the same session, and therefore each request "message -> response" cycle get progressively more expensive as a chat ages. 

Long running chats can also reach a point of "diminished returns" if there have been changes of focus, direction, subject etc within the chat. 

It therefore makes sense to work with short running sessions that focus on completing atomic tasks. However, it may be the case that we are working with an agent to complete many tasks in service of a larger goal or delivery. 

In this context, we need to retain state *between* tasks. This requires two things:
- Prompting the agent on how to start and how to finish tasks.
- Prompting the agent on how to read and write state.

## A solution

This framework uses just two files (DEVSTATE.md, PLAN.md), one tool (git), and one prompt (in PROMPT.md) to provide in the settings for your agent so that it reads this on every request.

## The files

- **DEVSTATE.md**: The LLM overwrites this file on completion of every task to record the current development status for the next session to pick up.
- **PLAN.md**: This user manages this file which details ideally short-term and clear goals, constraints, requirements, and focused tasks that lead to the goal. The LLM may update this file also, but the user needs to ensure that the goals, tasks and other information are concise, congruent and focused towards their development aims.
- **PROMPT.md**: Contains the prompt to be added to your AI assistant's settings.

## Requirements

- Git
- VSCode or similar dev environment which allows chat based coding agents to be installed.
- Some agentic coding assistant where the agent can be provided with "the prompt" in settings.

## Development Continuity Prompt

The development continuity prompt has been moved to its own file. See [PROMPT.md](PROMPT.md) for the prompt to add to your AI assistant's settings.

## Getting Started

1. **Choose your approach**:
   ```bash
   # Either clone this starter repository
   git clone https://github.com/yourusername/agentic-dev-starter.git
   cd agentic-dev-starter
   
   # Or copy the files to your existing project
   # Just copy PLAN.md, DEVSTATE.md, PROMPT.md, and README.md to your project
   ```

2. **Set up your project**:
   - Customize PLAN.md with your project goals, constraints, and initial tasks
   - Initialize DEVSTATE.md with your project's starting state

3. **Configure your AI assistant**:
   - Copy the prompt from PROMPT.md into your AI assistant's settings
   - Ensure your assistant has access to read and write files in your project

## Usage Example

Here's a simple example of how to use this framework:

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

This workflow helps maintain context between development sessions while keeping each interaction focused and efficient.

## Contrib

Any feedback, ideas or PR's are most welcome... I suspect this project may well age faster than a box of blueberries!
