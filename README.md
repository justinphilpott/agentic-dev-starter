# Agentic Dev Starter

A minimal framework for effective AI-assisted development.

## Contents
- [Purpose](#purpose)
- [The Problem](#the-problem)
- [A Solution](#a-solution)
- [The Files](#the-files)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Development Continuity Prompt](#development-continuity-prompt-to-help-keep-agentic-development-workflows-flows-on-track)
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

This framework uses just two files (DEVSTATE.md, PLAN.md), one tool (git), and one prompt (below) to provide in the settings for your agent so that it reads this on every request.

## The files

- **DEVSTATE.md**: The LLM overwrites this file on completion of every task to record the current development status for the next session to pick up.
- **PLAN.md**: This user manages this file which details ideally short-term and clear goals, constraints, requirements, and focused tasks that lead to the goal. The LLM may update this file also, but the user needs to ensure that the goals, tasks and other information are concise, congruent and focused towards their development aims.

## Requirements

- Git
- VSCode or similar dev environment which allows chat based coding agents to be installed.
- Some agentic coding assistant where the agent can be provided with "the prompt" in settings.

## Development continuity prompt for ultra-minimal interactions

Provide the following prompt below to your agent via settings (for example if using Cline in VSCode, click settings -> input prompt to the custom instructions box -> click Done.)

The ultra-minimal prompt style focuses on maximum efficiency with minimal verbosity, while maintaining all essential functionality.

```
⚠️ YOU MUST FOLLOW THESE PROTOCOLS WITHOUT EXCEPTION FOR EVERY TASK - WITH MINIMAL OUTPUT ⚠️

## 🔄 TASK START PROTOCOL

1. **Silent Context Gathering:**
   - Read DEVSTATE.md, check recent git commits, and review PLAN.md
   - DO NOT output detailed summaries of what you've read
   - DO NOT repeat file contents back to the user

2. **Minimal Proposal:**
   - In 1-3 sentences, state what you understand as the next logical task(s)
   - If concerns exist, state them in 1 or 2 sentences
   - Ask for confirmation in a single, direct question

⛔ WAIT FOR CONFIRMATION BEFORE IMPLEMENTATION ⛔

## 📋 DEVELOPMENT GUIDELINES

1. MINIMIZE CODE: Use existing solutions when possible
2. NARROW SCOPE: Build minimum viable solutions
3. SIMPLIFY: Choose straightforward approaches
4. INCREMENTAL: Deliver in small, independent units
5. CHALLENGE: Question if >20 lines of code are needed
6. COMPLETE: Focus on finishing rather than future-proofing

## 🏁 TASK COMPLETION PROTOCOL

1. **Update DEVSTATE.md** with minimal but sufficient information
2. **Commit changes** using conventional commits style
3. **Suggest next tasks** in 1-2 sentences

## ⚡ RESPONSE STYLE PROTOCOL

- NO verbose explanations of your thought process
- NO lengthy summaries of what you've read
- NO repetition of information the user already knows
- NO unnecessary pleasantries or conversational fillers
- YES to direct, information-dense responses
- YES to getting straight to the point
```

## Getting Started

1. **Clone or create a new repository**:
   ```bash
   # Either create a new project
   mkdir my-agentic-project
   cd my-agentic-project
   git init
   
   # Or clone this starter repository
   git clone https://github.com/yourusername/agentic-dev-starter.git
   cd agentic-dev-starter
   ```

2. **Create the required files**:
   - Copy the PLAN.md, DEVSTATE.md, and README.md files to your project
   - Customize PLAN.md with your project goals, constraints, and initial tasks
   - Initialize DEVSTATE.md with your project's starting state

3. **Configure your AI assistant**:
   - Copy the prompt from the section below into your AI assistant's settings
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
