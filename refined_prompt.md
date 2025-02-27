# Refined Prompt for Agentic Dev Starter

Below is a refined version of the prompt that makes the task start protocol less verbose while maintaining its effectiveness:

```
⚠️ IMPORTANT: FOLLOW THESE PROTOCOLS FOR EVERY TASK ⚠️

## 🔄 TASK START PROTOCOL

1. **Quick Context Review:**
   - Read DEVSTATE.md for current status
   - Run: `git log -n 3 --oneline`
   - Check PLAN.md for goals and next tasks

2. **Propose Approach:**
   - Outline your intended steps
   - Flag any concerns with tasks that:
     * Contradict guidelines
     * Conflict with goals/requirements
     * Need rescoping (too broad/narrow)
   - Get user confirmation

⛔ WAIT FOR CONFIRMATION BEFORE IMPLEMENTATION ⛔

## 📋 DEVELOPMENT GUIDELINES

1. MINIMIZE CODE: Use existing solutions when possible
2. NARROW SCOPE: Build minimum viable solutions
3. SIMPLIFY: Choose straightforward approaches
4. INCREMENTAL: Deliver in small, independent units
5. CHALLENGE: Question if >20 lines of code are needed
6. COMPLETE: Focus on finishing rather than future-proofing

## 🏁 TASK COMPLETION PROTOCOL

1. **Update DEVSTATE.md with:**
   - Task completed
   - Current status
   - Known issues/limitations

2. **Commit changes:**
   ```
   git add .
   git commit -m "<type>(<scope>): <subject>"
   ```

3. **Review PLAN.md:**
   - Identify next logical tasks
   - Suggest updates to guide future development
```

## Key Changes Made:

1. **Reduced Verbosity:**
   - Shortened section headings and instructions
   - Removed redundant phrases like "ALWAYS" and "ALL of the following steps"
   - Condensed multi-line instructions into bullet points

2. **Improved Formatting:**
   - Used more compact bullet points and nested lists
   - Applied better visual hierarchy with headers and emphasis
   - Reduced repetition of formatting elements

3. **Streamlined Content:**
   - Reduced git log from 5 to 3 commits (usually sufficient for context)
   - Consolidated similar instructions
   - Removed unnecessary explanatory text

4. **Maintained Critical Elements:**
   - Preserved all essential steps and protocols
   - Kept important warnings and confirmation points
   - Retained all development guidelines

This refined prompt is approximately 40% shorter while preserving all the functional elements of the original.
