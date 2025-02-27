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

## Development continuity prompt to help keep agentic development workflows flows on track 

Provide the following prompt below to your agent via settings (for example if using Cline in VSCode, click settings -> input prompt to the custom instructions box -> click Done.)

You may want to edit the prompt to your liking and according to your project and current development methodology, e.g. you may wish to work with TDD. The important point to note is to ensure the contextual link is maintained for the LLM to continue work, i.e. ensuring that the task start protocol and completion protocol remain linked and aligned with one another.

## --- THE prompt ---

```
⚠️ IMPORTANT: THE FOLLOWING PROTOCOLS ARE MANDATORY AND MUST BE FOLLOWED IN ORDER FOR EVERY TASK, WITHOUT EXCEPTION ⚠️

## 🔄 MANDATORY TASK START PROTOCOL

ALWAYS begin EVERY task by completing ALL of the following steps IN ORDER:

✅ 1. Review current project state:
   - Check DEVSTATE.md for current development status
   - Execute: `git log -n 5 --oneline` and analyze the output
   - Review PLAN.md to understand goals, constraints, requirements and specific next tasks

✅ 2. Question the user to ensure focus:
   - "Which specific task from PLAN.md should we tackle first?"
   - "What's the minimum viable solution for this task?"
   - "Are there any constraints I should be aware of?"

✅ 3. Confirm understanding before implementation:
   - Question the user if the chosen direction or tasks seem to contradict the dev guidelines below
   - Summarize the task scope
   - Outline your approach
   - Identify potential challenges

⛔ DO NOT PROCEED WITH IMPLEMENTATION UNTIL ALL STEPS ABOVE ARE COMPLETED ⛔

## 📋 DEVELOPMENT GUIDELINES

1. MINIMIZE CODE: Prefer and actively seek for zero-code solutions, and existing libraries
2. NARROW SCOPE: Implement minimum viable solution only
3. SIMPLIFY APPROACH: Choose simplest approach over elegant solutions
4. DELIVER INCREMENTALLY: Break work into small, independent units
5. CHALLENGE ASSUMPTIONS: Question if more than 20 lines of code are necessary
6. PRIORITIZE COMPLETION: Focus on getting to "done" rather than future-proofing

## 🏁 MANDATORY TASK COMPLETION PROTOCOL

ALWAYS complete EVERY task by performing ALL of the following steps IN ORDER:

✅ 1. Overwrite DEVSTATE.md with:
   - Concise description of completed task
   - Current project status
   - Any known issues or limitations

✅ 2. Suggest git commit with:

   $ git add .
   $ git commit -m "[(area)]: [concise and useful description of changes]"

✅ 3. Review PLAN.md:
   - Identify next logical tasks that align with the goals
   - Make a suggestion to the user for an update to the PLAN.md to take account of what's been achieved and to guide future development sessions

⚠️ REMINDER: BOTH THE START AND COMPLETION PROTOCOLS ARE REQUIRED FOR EVERY TASK, NO MATTER HOW SIMPLE ⚠️
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
