---
name: "♻️ Refactor"
about: Improve code structure or infra without changing behavior.
title: "refactor: <title>"
labels: refactor, technical-debt
assignees: ''

---

## 1. Objective
<!-- What code or infra changing? Why now? -->

## 2. Context / Motivation
<!-- What wrong with current code? Performance? Complexity? Debt? -->

## 3. Scope of Changes
<!-- Which packages or files affected? -->

## 4. Expected Benefits
<!-- Faster? Less code? Better DX? -->

## 5. Potential Risks
<!-- Regressions? Breaking changes to internal APIs? -->

## 6. Implementation Plan
- [ ] **Step 1**: Initial changes
- [ ] **Step 2**: Refactor logic
- [ ] **Step 3**: Verify behavior

## 7. Verification
### Test Case Matrix
| Test Case | Scenario | Expected Outcome |
| :--- | :--- | :--- |
| **Existing Logic** | <!-- Run old flows --> | <!-- No change in behavior --> |
| **Edge Cases** | <!-- Boundary values --> | <!-- Handled same as before --> |
| **Performance** | <!-- Load check --> | <!-- Parity or improvement --> |

<!-- How ensure no logic change? -->
- [ ] Existing tests pass.
- [ ] New unit tests for refactored parts.
- [ ] Load/Performance test (if applicable).
