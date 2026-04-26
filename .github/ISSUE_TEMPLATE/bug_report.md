---
name: "🐛 Bug Report (Design Doc Style)"
about: Report a bug with a detailed RCA and proposed fix.
title: "bug: <title>"
labels: bug, triage
assignees: ''

---

## 1. Overview
<!-- A concise summary of what is failing. -->

## 2. Problem Statement
### Current Behavior
<!-- What is happening right now? -->

### Expected Behavior
<!-- What should happen? -->

### Steps to Reproduce
1. 
2. 
3. 

## 3. Root Cause Analysis (RCA)
<!-- Describe WHY this is happening. Point to specific files or logic if possible. -->

## 4. Proposed Fix (Technical Design)
<!-- How do you plan to fix this? Describe changes to handlers, usecases, or repositories. -->
```go
// Example of proposed change
```

## 5. Testing & Verification
### Test Case Matrix
| Test Case | Scenario | Expected Outcome |
| :--- | :--- | :--- |
| **Fix Validation** | <!-- e.g. Valid payload --> | <!-- e.g. 200 OK, DB updated --> |
| **Regression** | <!-- e.g. Existing flows --> | <!-- e.g. No side effects --> |
| **Edge Case** | <!-- e.g. Empty values --> | <!-- e.g. 400 Bad Request --> |

- [ ] **Unit Tests**: Description...
- [ ] **E2E Tests**: Description...

## 6. Impact Assessment
<!-- Does this change affect other resources? Breaking changes? -->
