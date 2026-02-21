# QuantatomAI: The UI/UX Moat Framework & Laws of Interaction

To ensure the frontend of QuantatomAI is not just functional, but an undeniable, addictive "UX Moat" against legacy competitors, we enforce a strict, heavily disciplined methodology. UI/UX is fluid and complex, which means without rigid constraints and principles, entropy inevitably degrades the user experience.

This document serves as the constitutional framework for all frontend development. Every component, view, animation, and flow MUST align with these laws. If a PR violates a law without an explicit, well-documented exemption, it fails review.

---

## Part 1: The Core 20 UX Laws (The EPM Moat)

| # | Law Name | Core Principle (What It Means) | Why It’s Non-Negotiable (Moat Impact) | Enforcement / How to Guard It |
|---|---|---|---|---|
| 1 | **Ruthless Simplicity & Progressive Disclosure** | Show only what the user needs right now. Hide complexity until explicitly requested. Max 5–7 visible top-level nav items. | Overwhelm kills perceived speed and trust. Simplicity is the #1 reason products feel “fast” even when they’re not. | Lint rule: forbid >7 items in nav arrays. Visual regression fails on clutter. |
| 2 | **One Primary Action per Context** | Every screen/modal/card must have exactly one dominant CTA (size/color/position). Others secondary/tertiary. | Reduces decision fatigue → users act faster. Ambiguity feels slow and unprofessional. | Component lint: require exactly one `<Button variant="primary" />` per context. |
| 3 | **Consistent Mental Model & Affordances** | Buttons depress, draggables show grip, expanders rotate 180°. All four states (rest/hover/focus/active/disabled) distinct. | Builds instant trust — users “feel” the UI without thinking. Inconsistency causes hesitation. | Stylelint bans non-token affordances. Storybook tests enforce all four states per component. |
| 4 | **Zero Surprises – Instant Feedback** | Every input gets visual/haptic acknowledgment in ≤100 ms (cursor move, button depress, drag follow, loading spinner). | Eliminates doubt → addictive responsiveness. Most enterprise tools are 300–800 ms. | Playwright/Cypress tests clamp reaction times. Performance budget enforced in CI. |
| 5 | **Undo & Forgiveness Are Mandatory** | Every destructive action has visible undo (toast + timer), soft-delete, or version history in ≤1 click. | Prevents regret → users experiment freely. Forgiveness builds massive adoption lock-in. | Every mutative action wraps in `useUndo` hook. Lint forbids irreversible actions without undo. |
| 6 | **Full Multi-Modal Parity** | Keyboard, mouse, touch, voice are equal citizens. Every flow fully operable in all four. | Serves all users (desk/mobile/voice) → opens markets competitors ignore. | Cypress tests all modes. Accessibility score ≥95 in Lighthouse CI. |
| 7 | **Contextual & Just-In-Time Actions** | Right-click / long-press / hover / selection shows only relevant actions (≤5–7 items max unless asked for more). | Zero overwhelm → speed. Contextual reduces clicks by 40%+. | Lint forbids global menus. Context menu component enforces ≤7 items. |
| 8 | **Motion With Meaning** | Every transition has purpose. Use spring physics or consistent cubic-bezier. No random motion. Tokens enforced. | Natural flow guides eyes → reduces cognitive load. Bad motion makes products feel cheap. | Motion lint bans non-token animations. `framer-motion` configurations strictly enforced via presets. |
| 9 | **Accessibility Is a Feature, Not a Checkbox** | WCAG 2.2 AA minimum (AAA where possible). ARIA live regions for dynamic content. Contrast ≥4.5:1. Keyboard testable. | Opens gov/enterprise deals + ethical edge. True accessibility is rare in EPM. | CI runs Lighthouse ≥95 accessibility. Axe/WAVE tests in PR gate. |
| 10 | **Predictability Over Discoverability** | Once learned, never break a pattern (placement, timing, gesture, shortcut). Primary action bottom-right/top-right. | Muscle memory = addiction = switching cost. Predictability wins power users. | Visual regression fails on layout shifts. Consistent timing enforced via motion tokens. |
| 11 | **Metrics-Driven Iteration** | Instrument every interaction (time-to-task, error rate, rage-clicks, drop-offs). Run A/B tests relentlessly. | Data-driven products evolve faster than competitors can copy. | PostHog/Mixpanel tracking mandatory. Every PR includes metric baseline. |
| 12 | **No Dark Patterns — Ever** | No hidden declines, fake scarcity, misleading toggles, forced upsells. Confirm irreversibles twice. Transparent limits. | Trust compounds. Lost trust is impossible to regain in enterprise. | Design review vetoes any tricks. Audit for hesitations in flows. |
| 13 | **Consistency Is Sacred – No Unauthorized Deviations** | Logo, colors, typography, spacing, radius, shadows, icons, motion, grid must never change without design-system update. Tokens only. No raw hex/magic numbers. | Develops unmistakable “feel” — extremely hard to replicate. Drift makes products feel cheap. | ESLint/Stylelint block raw values. Visual regression fails drift. Token validator in CI. |
| 14 | **Adaptive & Responsive by Default** | Every layout/grid/menu/modal/component auto-adapts to screen size/orientation/device. Fluid typography (clamp), minmax(). | Device-agnostic UX reduces abandonment on mobile/tablet — competitors often fail here. | Responsive lint checks clamp/minmax. Cypress tests multiple viewports. |
| 15 | **Dynamic Menus, Dropdowns & Popovers** | Open with 200ms fade+scale from origin. Auto-flip to stay in viewport. Type-ahead fuzzy search + virtualization for long lists. | Menus always visible/useful — no clipping, no scroll hell. | Floating UI (Popper) positioning enforced. Test edge cases in CI. |
| 16 | **Sidebar Expansion & Collapse** | Hover → expand from icons to labels (300ms); >> / << toggle persists state. Tooltips on collapsed icons. | Saves space while remaining intuitive — adapts to user preference seamlessly. | Gesture lint for touch. Persisted state via localStorage + Zustand API sync. |
| 17 | **Grid & View Controls** | Every grid/panel/modal/frame has restore/maximize/minimize/close controls (hover top-right icons). | Feels like native windows — users manage views intuitively without frustration. | High-order component (HOC) wrappers must include all 4 controls. Lint enforces wrapper presence. |
| 18 | **Dynamic Scrollbars & Overflow** | Scrollbars appear only when needed, auto-hide after 2s inactivity. Thin/custom styling. Momentum on touch. | Clean, modern look — no permanent scrollbars cluttering UI. | Overflow lint checks `auto` vs `scroll`. Touch events tested for momentum natively. |
| 19 | **Multi-Theme Harmony** | ≥5 curated themes remap every visual token instantly (no flicker). User preference persists across sessions/devices. | Personalized polish — users perceive choice without overwhelm. Signature feel hard to copy. | Theme lint ensures no token leaks. Next-Themes integration tested for zero flash-of-unstyled-content (FOUC). |
| 20 | **Containment & Viewport Discipline** | No element/tooltip/popover/modal/overflow may extend beyond visible viewport unless explicit (e.g. full-screen). Auto-flip/prevent-overflow. | No clipping hell — UI always usable, never broken-looking. | Floating UI modifiers enforced. CI tests for overflow/clipping in all viewports via bounding box assertions. |

---

## Part 2: The Antigravity Expansion (AI + WebGPU Specifics)

As your AI Co-Pilot, I have reviewed the core 20 laws and identified that to truly dominate the legacy EPM space (like Anaplan or SAP), we must add rules specific to our technological advantage: **WebGPU, Real-time CRDTs, and Fluxion Generative AI.**

Here are 6 additional Antigravity Laws to complete the framework:

| # | Law Name | Core Principle (What It Means) | Why It’s Non-Negotiable (Moat Impact) | Enforcement / How to Guard It |
|---|---|---|---|---|
| 21 | **WebGPU / DOM Segregation (The Dual-Engine Approach)** | Heavy matrix renderings (10,000+ cells) MUST bypass the DOM entirely and render directly to WebGPU/WebGL canvases. The DOM is strictly for shell UI, forms, and popovers. | The DOM stalls at ~3,000 nodes. WebGPU can render 5,000,000. This is the difference between sluggish legacy grids and instant scrolling. | Architectural lint: No mapping `<tr>`/`<td>` or `<div>` grids over 50 rows. Only `<canvas>` mounts allowed for `Grid` components. |
| 22 | **Zero-Jank Scrolling & 60-120fps Floor** | Scroll events must never drop frames, regardless of dataset size. Pagination is banned in grids; infinite virtualized scrolling is mandatory. | "Jank" destroys the illusion of continuous data. 120fps scrolling feels like native desktop software, completely obliterating web-based competition. | Chrome Performance profiler runs in headless mode during CI to assert >60fps framerates during heavy automated scrolls. |
| 23 | **Optimistic UI Updates (Assume Success)** | When a user edits a grid cell or submits a formula, immediately show the updated state locally BEFORE the network roundtrip (CRDT/WebSocket event) returns. | Network latency (e.g., 50ms to 200ms) creates a sluggish feel. Optimistic UI makes a cloud application feel local and offline-first. | `useMutation` hooks must implement `onMutate` optimistic caches. Redpanda WS confirmations quietly resolve in the background. |
| 24 | **Actionable Intelligence (Never Just Text)** | When Fluxion AI outputs an insight (e.g., "Revenue dropped due to APAC"), it MUST include an executable action button (e.g., `[Drill into APAC]`, `[Create Scenario]`). | Passive AI is quickly ignored. Agentic AI that reduces clicks and operates the software for the user becomes essential workflow infrastructure. | LLM Tool Output schemas must strictly require an `actions[]` array. UI lints require `<FluxionAction>` components. |
| 25 | **Information Density Sovereignty** | Users must have total control over density (Compact, Comfortable, Touch). Power users demand maximum rows/columns. Executives require breathing room. | Legacy systems force one size fits all. Offering density controls respects different personas natively instead of requiring custom CSS hacks. | Root layout applies `--spacing-scale` CSS variable globally. Components must use `calc(var(--spacing-scale) * N)` instead of hardcoded `px` padding. |
| 26 | **Z-Index Governance & Spatial Elevation** | Strict 10-layer Z-index hierarchy. Elements elevate based on transience (Background -> Canvas -> Shell -> Dropdown -> Modal -> Toast/Error -> AI Overlays). Magic Z-Index (`9999`) is banned. | Z-index wars lead to tooltips trapped under headers, destroying trust. A governed system ensures layers stack correctly predictably. | Linter utterly bans `z-index` values in CSS unless they map to a predefined SCSS/CSS variable like `var(--z-modal)`. |

---

## Part 3: Establishing the "Method to the Madness"

Because UI is fluid, we cannot rely on memory. The **"Method to Madness"** means transferring these 26 laws from "human guidelines" into **compiler errors, automated tests, and CI/CD gates.** 

### 1. The Design Token Supremacy
We will build a central `tokens.css` or `theme.ts` file containing every allowable color, spacing unit, elevation shadow, animation bezier curve, and font tier. 
**Rule:** No React component may define a raw value. `padding: '16px'` is illegal. `padding: 'var(--space-md)'` is mandatory. 

### 2. Component Composition Strictness
Instead of building a `<table>` everywhere, we build a single heavily tested `<WebGPUGrid>` primitive that handles virtualization, scrollbars (Law 18), viewport containment (Law 20), keyboard parity (Law 6), and 120fps rendering (Law 22). By forcing all views to use this single primitive, we ensure global conformity to the 26 Laws.

### 3. The "Cypress / Playwright Core Gate"
Before any PR merges to `main`, E2E tests will:
- Tab through the entire app (Checking Law 6: Multi-Modal Parity).
- Assert the menu opens in `< 100ms` (Checking Law 4: Instant Feedback).
- Drag the app to a mobile viewport and assert no horizontal overflow (Checking Law 14: Adaptive by Default).

These 26 laws now form our **QuantatomAI Frontend Bible**.
