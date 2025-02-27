# Agentic Dev Starter

A minimal framework for effective AI-assisted development.

## Core Files

- **DEVSTATE.md**: Current project state (last completed task)
- **PLAN.md**: Short-term goals, constraints, requirements, and focused tasks
- **Git History**: Source of truth for project evolution

## Agentic Development Guidance

```
## TASK START PROTOCOL

1. Review current project state:
   - Check DEVSTATE.md for current status
   - Review PLAN.md for goals and next tasks
   - Examine git history: `git log -n 5 --oneline`

2. Question the user to ensure focus:
   - "Which specific task from PLAN.md should we tackle first?"
   - "What's the minimum viable solution for this task?"
   - "Are there any constraints I should be aware of?"

3. Confirm understanding before implementation:
   - Summarize the task scope
   - Outline your approach
   - Identify potential challenges

## DEVELOPMENT GUIDELINES

1. MINIMIZE CODE: Prefer zero-code solutions (configs, existing libraries)
2. NARROW SCOPE: Implement minimum viable solution only
3. SIMPLIFY APPROACH: Choose simplest approach over elegant solutions
4. DELIVER INCREMENTALLY: Break work into small, independent units
5. CHALLENGE ASSUMPTIONS: Question if more than 20 lines of code are necessary
6. PRIORITIZE COMPLETION: Focus on getting to "done" rather than future-proofing

## TASK COMPLETION PROTOCOL

1. Update DEVSTATE.md with:
   - Concise description of completed task
   - Current project status
   - Any known issues or limitations

2. Suggest git commit with:
   ```
   git add .
   git commit -m "[area]: [brief description of changes]"
   ```

3. Review PLAN.md:
   - Suggest updates to reflect progress
   - Identify next logical tasks
```

## Philosophy

This framework emphasizes:
- Small, incremental changes over large batches
- Focused delivery over comprehensive planning
- Git as the source of truth for project history
- Minimal documentation that encourages action

## Usage

1. Copy DEVSTATE.md and PLAN.md to your project
2. Add the guidance above to your AI assistant's context
3. Start with a focused task in PLAN.md
4. Let the AI assistant guide the development process
