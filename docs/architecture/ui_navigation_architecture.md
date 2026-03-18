> SSOT Derivation Notice
> This document derives from the canonical architecture SSOT: [docs/architecture/quantatomai-single-source-of-truth.md](docs/architecture/quantatomai-single-source-of-truth.md).
> If any conflict exists, the SSOT prevails.

# QuantatomAI Red Team UX Manifesto: The 1500-Point Interaction Moat

This is the ultimate EPM UX fortress: Hyper-granular, AI-adaptive, and infinitely extensibleâ€”obliterating legacy systems with zero-friction flows and predictive intelligence.

---

## 0. The Front Door (Pre-Auth, Identity, & Boot Sequence)
**[0.1] Public Landing (`quantatom.ai`)**
*   **0.1.0** Scrolling background: WebGL lattice of data points resolving into structured cubes.
*   **0.1.1** Button hover `[ Sign In ]`: Button glows with a subtle neon border, hinting at WebGPU styling.
*   **0.1.2** Interactive Benchmark: Drag a slider to simulate row counts. UI shows "Their Speed: 12h" vs "Our Speed: 40ms".
*   **Added: 0.1.3** AI Tour: Hover "Demo" â†’ Predictive popup previews key flows (e.g., "See budget cycle in 30s").

**[0.2] Authentication & Security Gateway**
*   **0.2.0** Identity Routing: User types email. Domain is extracted on `onBlur`. If `@globalcorp.com`, redirect to Azure AD.
*   **0.2.1** SSO Negotiation: SAML Handshake. Display "Securely Handshaking with Microsoft..."
*   **0.2.2** JWT Storage: Store token in secure HttpOnly cookie for CSRF protection.
*   **0.2.3** Multi-Factor Challenge: If IP differs from last login, push Auth0 MFA prompt to phone.
*   **Added: 0.2.4** Biometric Hover: On login, detect device â†’ Suggest fingerprint (WebAuthn) with smooth anim.

**[0.3] "Cold Start" Boot Sequence (The 500ms Window)**
*   **0.3.0** Parallel Boot Loading: Logo pulses based on WebSocket connection status.
*   **0.3.1** RBAC Fetch: GET `/api/context`. Cache user permissions in Redux store.
*   **0.3.2** Governance Ping: Determine if Cortex AI Sparkle button is rendered or hidden based on tenant config.
*   **0.3.3** WebGPU Context Init: Pre-warm the canvas context hidden in the DOM.
*   **0.3.4** Arrow Flight Warmup: Pre-fetch dimension metadata headers.
*   **0.3.5** Cache Hydration: Load "Recently Visited" pointers from IndexedDB.
*   **Added: 0.3.6** Predictive Preload: AI guesses first page based on history â†’ Load in background.

**[0.4] First-Time Administrator Out-of-Box Experience (OOBE)**
*   **0.4.0** "Blank Slate" Prompt: No apps exist. Large "Let's Build" CTA appears.
*   **0.4.1** AI Org-Chart Upload: Drag PDF of corporate org chart onto screen.
*   **0.4.2** Visual Parser: AI draws bounding boxes over the PDF on screen.
*   **0.4.3** Proposed Hierarchy: Cortex generates the LTree. Shows tree diagram with nodes.
*   **0.4.4** Approve Hierarchy: Admin clicks "Looks Good." API generates 5,000 entities in 0.1s.
*   **Added: 0.4.5** Gamified Steps: Badge for each OOBE task; AI suggests "Next Job: Import Data".

---

## 1. Global Shell & Navigational Omni-Layer (The 5000-Point Moat Core)

**[1.1] Profile & Tenant Management (The Personal Command Hub)**
*   **1.1.0** Avatar Click: Dropdown opens with 200ms ease-out CSS animation; shadow depth increases on hover for depth cue.
*   **1.1.1** Theme Switcher: Click "Dark" â†’ Instant CSS var swap; transition all elements with 300ms fade; track preference in localStorage for cross-session persistence.
*   **1.1.2** Tenant Switcher: Dropdown lists tenants with search bar (fuzzy match on name/description); select â†’ Clear Redux state, reload global POV, trigger WebSocket reconnect; show "Switching Realms..." toast with progress bar.
*   **1.1.3** Active Sessions: List devices with IP/location; "Revoke" button â†’ API call to invalidate token, real-time logout push via WebSocket; confirmation modal with "Are you sure?" and undo timer (5s).
*   **Added: 1.1.4** AI Profile Insights: Avatar hover â†’ Tooltip with "Your top action: Export Reports (42x this week)" â€“ AI analyzes usage for tips.
*   **Added: 1.1.5** Voice Logout: Long-press avatar â†’ "Logout?" voice prompt (Web Speech API); confirm verbally.
*   **Added: 1.1.6** Accessibility: ARIA-label="User menu, 5 items"; keyboard nav (arrow keys cycle, Enter selects); high-contrast avatar border.
*   **Added: 1.1.7** Reaction Chain: Theme switch â†’ Auto-adjust chart colors (e.g., invert hex); error if offline â†’ Fallback to default with toast "Re-syncing on reconnect".
*   **Added: 1.1.8** Extension Hook: Devs add custom dropdown items via plugin API (e.g., "Integrate Slack" button).

**[1.2] The "Omni-Switcher" Mega-Menu (Left Sidebar â€“ Infinite Nav Engine)**
*   **1.2.0** Collapse Sidebar: Click `[<]` â†’ Animate shrink to 48px icons (300ms); main content expands with reflow; persist state in localStorage.
*   **1.2.1** App Switcher: Click `Apps v` â†’ Mega-dropdown with grid layout (virtualized for 1000+ apps); thumbnails + search.
*   **1.2.2** Fuzzy Search Box: Type "sal" â†’ Instant filter (debounced 100ms); bold matches; AI-rank results by usage history.
*   **1.2.3** "Open All" Folder: Hover folder â†’ Preview sub-items in tooltip; click â†’ Spawn tabs in background (lazy-load content).
*   **Added: 1.2.4** Infinite Nesting: Sub-menus expand recursively (virtualized list, load on demand); drag to reorder items (persist via API).
*   **Added: 1.2.5** AI-Adaptive Reorder: Sidebar auto-rearranges weekly based on click frequency (e.g., "Planning" bubbles to top); opt-out toggle.
*   **Added: 1.2.6** Gesture Nav: Swipe right on touch â†’ Expand sidebar; pinch to zoom menu font size (persistent).
*   **Added: 1.2.7** Contextual Reactions: Hover item â†’ Show "Pin to Top" button; right-click â†’ "Duplicate Menu" for custom views.
*   **Added: 1.2.8** Plugin Marketplace: Bottom "Add Extension" â†’ Browse/search plugins (e.g., "Add Custom KPI Menu"); install â†’ Dynamic sidebar injection.
*   **Added: 1.2.9** Accessibility: ARIA-tree for nested menus; voice nav ("Go to Reports"); focus trap on open.
*   **Added: 1.2.10** Error Handling: If menu load fails (network) â†’ Fallback to cached version with "Offline Mode" banner; auto-retry on reconnect.

**[1.3] Global Point-of-View (POV) & Dimensional Slice Matrix (TopNav â€“ The Context Engine)**
*   **1.3.0** Pinned Dropdowns: Display globally pinned dimensions (e.g., Region, Time).
*   **1.3.1** List Virtualization: Dropdown opens containing 50k items. No lag, rendered dynamically as user scrolls.
*   **1.3.2** Type-Ahead Filtering: Focus a dropdown, start typing. Instantly filters list.
*   **1.3.3** Global State Update: User clicks France. Redux fires `updateGlobalPOV({ Region: 'France' })`. Every component on screen (grids, charts, text) redraws synchronously.
*   **Added: 1.3.4** Multi-Dim Nesting: Drag dims to nest (e.g., Region under Time); auto-save configs per user.
*   **Added: 1.3.5** POV Preview: Hover option â†’ Mini-grid popup shows slice impact (lazy-fetch data).
*   **Added: 1.3.6** AI POV Suggest: "Based on your role, switch to EMEA?" button; auto-apply on confirm.
*   **Added: 1.3.7** Reaction Chain: POV change â†’ Trigger workflow checks (e.g., "This unlocks Q4 approvals"); haptic feedback on mobile.
*   **Added: 1.3.8** Extension Hook: Plugins add custom POV dims (e.g., "AI Scenario Dim").
*   **Added: 1.3.9** Accessibility: ARIA-combobox for dropdowns; voice select ("Set Region to France").
*   **Added: 1.3.10** Error Handling: Invalid POV (e.g., locked dim) â†’ Shake animation + toast "Permission Denied â€“ Upgrade Role?".

**[1.4] Spatial Git-Flow (The Sandbox Visualizer â€“ Versioning Moat)**
*   **1.4.0** Branch Selector: TopNav says main. User clicks, dropdown shows dev-upside.
*   **1.4.1** Branch Change: Click dev-upside. TopNav turns Neon Orange.
*   **1.4.2** Matrix Overlay: The entire UI gets a 2px glowing orange border. You know you are in a sandbox.
*   **1.4.3** Diff Mode Toggle: Click `[Compare]`. Viewport splits 50/50. Left is Main, Right is Branch.
*   **1.4.4** Commit Action: Click `[Merge to Prod]`. Opens diff summary. Enforces comment before push.
*   **Added: 1.4.5** Infinite Branching: "Fork from Here" button â†’ Auto-create branch with current state; tree view for branch history.
*   **Added: 1.4.6** AI Merge Preview: On merge, show simulated outcomes (e.g., "This affects 500 cellsâ€”variance +3%").
*   **Added: 1.4.7** Gesture Undo: Swipe left â†’ Undo last commit with preview diff modal.
*   **Added: 1.4.8** Reaction Chain: Branch switch â†’ Auto-save open grids; error if conflict â†’ "Resolve Merge" wizard with side-by-side edits.
*   **Added: 1.4.9** Plugin Hook: "Custom Version Hook" for integrations (e.g., GitHub sync).
*   **Added: 1.4.10** Accessibility: ARIA-live announces "Branch changed to devâ€”sandbox mode active".

**[1.5] Command Palette (Ctrl+K â€“ The Action Accelerator)**
*   **1.5.0** Summon Palette: `Ctrl+K`. Translucent blur overlay behind palette.
*   **1.5.1** Command Registry: List displays recent actions.
*   **1.5.2** Prefix Routing: Type `>`. Expects actions (e.g., "Export"). Type `@`. Expects users (e.g., "@Bob"). Type `/`. Expects page navigation.
*   **1.5.3** Direct Action: Hit `Enter` on "Run Consolidation". Job fires in background.
*   **Added: 1.5.4** Infinite Commands: Virtualized list for 1000+ actions; AI ranks by context (e.g., in grid â†’ Suggest "Fill Down").
*   **Added: 1.5.5** Voice Integration: Mic icon â†’ Speak command; confirm with Enter.
*   **Added: 1.5.6** Reaction Preview: Hover command â†’ Tooltip shows outcome (e.g., "This exports 5k rows").
*   **Added: 1.5.7** Custom Commands: Plugin API to add (e.g., "Integrate with Slack").
*   **Added: 1.5.8** Error Handling: Invalid command â†’ AI suggests alternatives ("Did you mean Export PDF?").
*   **Added: 1.5.9** Accessibility: VoiceOver reads options; tab nav cycles.
*   **Added: 1.5.10** Moat Macro: Chain commands (e.g., "Create Page â†’ Add Grid â†’ Apply Formula") in one palette session.

**[1.6] AI Nav Co-Pilot (The Predictive Moat Layer)**
*   **1.6.0** Always-On AI: Floating orb (minimize to dock); click to expand chat.
*   **1.6.1** Predictive Suggestions: On nav hover â†’ "Users who viewed Reports next exported PDF (85% time)â€”Go?".
*   **1.6.2** Auto-Flow Builder: "Guide me through Budget Cycle" â†’ Step-by-step overlay with next-button advances.
*   **1.6.3** Habit Optimizer: Weekly email "Your nav heatmapâ€”try pinning Planning?"; auto-apply on approve.
*   **1.6.4** Reaction Integration: Any error â†’ AI "Fix this?" button with one-click resolution.

**[1.7] Extension Marketplace (The Infinite Moat)**
*   **1.7.0** Bottom Sidebar Button: "Add Extensions" â†’ Marketplace modal with search/browse.
*   **1.7.1** Install Flow: Click "Slack Integration" â†’ API hook injects new menu item dynamically.
*   **1.7.2** Custom Reactions: Extensions add hotkeys/context menus (e.g., "Export to PowerBI" right-click).
*   **1.7.3** Dev SDK: Docs for building (React components + Rust hooks); auto-test on upload.
*   **1.7.4** Security Scan: AI vets extensions for risks before install.

**[1.8] Accessibility & Inclusive Layer (The Universal Moat)**
*   **1.8.0** Auto-ARIA: All elements get labels/roles (e.g., sidebar as ARIA-tree).
*   **1.8.1** Keyboard-Only: Full tab nav; Enter to activate; custom shortcuts configurable.
*   **1.8.2** Screen Reader Flows: Announcements for reactions (e.g., "Cell edited, value 500k").
*   **1.8.3** High-Contrast Toggle: One-click (`Ctrl+H`) â†’ Boost colors/typography.
*   **1.8.4** Mobile Gestures: Swipe for nav; voice for commands; adaptive grids (collapse to cards).

**[1.9] Metrics & Telemetry Dashboard (The Self-Improving Moat)**
*   **1.9.0** Admin View: "UX Health" page with heatmaps (time per menu, drop-offs).
*   **1.9.1** A/B Testing: Randomly test nav variants (e.g., sidebar icons vs text); auto-rollout winner.
*   **1.9.2** User Feedback Loops: Inline thumbs-up/down on reactions â†’ AI refines.

**[1.10] Error & Edge Case Handling (The Resilient Moat)**
*   **1.10.0** Global Error Boundary: Catch crashes â†’ Graceful fallback UI with "Retry" button.
*   **1.10.1** Offline Mode: Detect disconnect â†’ Gray nav; queue actions for sync.
*   **1.10.2** Overload Protection: If 1000+ items â†’ AI paginate/suggest filters.
*   **1.10.3** Chaos Resilience: Simulate lags in dev mode; auto-recover from bad states.
**[2.1] Greetings & Inbox**
*   **2.1.0** Personalization: "Good morning, `Name`." Changes greeting based on local browser time API.
*   **2.1.1** Action Inbox: Red badge counter indicating pending workflow approvals.
*   **2.1.2** Inline Approve/Reject: Buttons for budgets/forecasts.
*   **Added: 2.1.3** Workflow Escalation: If pending >24h, auto-notify manager (Slack integration).  
*   **Added: 2.1.4** AI Job Prioritizer: Sorts inbox by urgency (variance impact).  

**[2.2] "Recently Visited" Carousel**
*   **2.2.0** Intelligent Resumption: Hovering over a recent model card shows a tooltip of the exact POV state when last closed.
*   **2.2.1** Deep Link Launch: Click a card, immediately restore layout and scroll position inside the grid.
*   **Added: 2.2.2** Carousel Auto-Play: Highlights changes since last visit (diff badges).  

**[2.3] Entropy & Alert Feed**
*   **2.3.0** Redpanda Stream: Subscribes to SSE. List updates organically.
*   **2.3.1** Anomaly Iconography: Distinct icons for data loads (green), lock status changes (grey), AI anomaly alerts (red flash).
*   **2.3.2** Actionable Alert: Click on alert "Variance in NA exceeded 10%". Opens side-pane highlighting the exact Grid rows.
*   **Added: 2.3.3** Alert Reactions: Right-click alert â†’ "Snooze" or "Assign to User".  

---

## 3. App Studio: Data Modeling & Architecture (Admin UX)
**[3.1] Dimensions Explorer**
*   **3.1.0** Create Dimension Modal: Input fields for Name, Type (Time, Standard, Version, List).
*   **3.1.1** Time Dimension Generator: Specify Start Year (2020), End Year (2030), Granularity (Months/Quarters). Instantly builds hierarchy.
*   **3.1.2** Orphan Nodes Alert: Banner warns if list items aren't mapped in a hierarchy.
*   **Added: 3.1.3** Dim Import Wizard: Bulk from CSV â†’ AI maps to hierarchy.  

**[3.2] The Interactive Hierarchy Builder (LTree Grid)**
*   **3.2.0** Drag-and-Drop Reparenting: Click node `West Coast`, drag it under `Americas` instead of `USA`.
*   **3.2.1** Real-time Validation: Attempt to drag node under itself -> Flash red, reject with shake animation.
*   **3.2.2** Property Editing Context Pane: Click a node; right sidebar opens showing properties (Alias, Display Format, Weight).
*   **3.2.3** Bulk Edit (Excel Mode): Toggle from Tree to Flat Grid. Select 50 rows, press `Ctrl+C`. Open Excel, paste, edit names, copy, return, `Ctrl+V`. System diffs and upserts.
*   **Added: 3.2.4** Multi-Select Reactions: Shift-click nodes â†’ Bulk color/format (swatch drag).  

**[3.3] The Cube/Model Configurator**
*   **3.3.0** Block Initializer: Click "New Cube".
*   **3.3.1** Axis Assigner: Drag `Time` to Columns. Drag `Line Items` to Rows. Drag `Region`, `Version` to Pages (POV).
*   **3.3.2** Sparsity Matrix Estimator: As axes are added, a live speedometer gauges theoretical cell count (e.g., 10 Trillion) vs estimated actual density. Warns if architecture is fundamentally flawed.
*   **3.3.3** Dependency Mapper: View a visual DAG representing which cubes feed data into this cube.
*   **Added: 3.3.4** Page Builder Moat: Infinite canvas for multi-cube pages; drag grids/charts â†’ AI auto-layout.  

---

## 4. AtomScript Rules Engine & Formula Editor
**[4.1] The Monaco Editor Integration**
*   **4.1.0** Syntax Highlighting: Custom grammar highlights functions (blue), members (orange), keywords (purple).
*   **4.1.1** IntelliSense: Type `LOOKUP(`. Pop-up shows parameter signatures.
*   **4.1.2** Contextual Type-ahead: Type `[Cost Center].`. Editor suggests available node names underneath.
*   **4.1.3** Live Linting: Missing closing bracket flags red squiggly line instantly via WASM parser.
*   **Added: 4.1.4** Bulk Formula Apply: Select grid range â†’ Paste formula â†’ AI validates.  

**[4.2] Formula Debugger UX**
*   **4.2.0** Evaluate Expression: Highlight a sub-section of a formula `([Rev]-[COGS])`. Right-click > "Evaluate". Shows temporary calculated value for current POV.
*   **4.2.1** The AST Visualizer Switch: Click `{}` icon. Editor switches to block-based visual layout of the formula's Abstract Syntax Tree.
*   **Added: 4.2.2** Reaction Preview: Hover formula â†’ Show affected cells in grid overlay.  

---

## 5. The Holographic Grid (Planner UX - The Core Moat)
**[5.1] WebGPU Viewport & Virtualization Operations**
*   **5.1.0** Initial Render: Paint 2,000 cells to HTML Canvas in <16ms via WebGPU bindings.
*   **5.1.1** Smooth Scrolling: Mouse wheel triggers sub-pixel translation of canvas coordinates. No DOM element lag.
*   **5.1.2** Dynamic Chunk Fetching: As scroll approaches viewport edge, Arrow buffers pre-fetch next spatial chunk.
*   **Added: 5.1.3** Infinite Scroll Moat: Auto-load dims as you scroll (no pagination).  

**[5.2] Pivot & Structuring**
*   **5.2.0** Drag Handles: Hover over `Version` label. Cursor changes to hand. Drag to column header.
*   **5.2.1** Matrix Rotation: Dimensions cross-multiply. Nested headers render dynamically.
*   **5.2.2** Asymmetric Expand: Click `+` on Q1 Actuals to show months, but keep Q2 collapsed.
*   **Added: 5.2.3** AI Pivot Suggest: "Optimize Layout" button â†’ Recommends based on data density.  

**[5.3] Basic Keyboard-Driven Editing (Excel Emulation)**
*   **5.3.0** Arrow Navigation: `Up/Down/Left/Right` shifts active cell border.
*   **5.3.1** Direct Typing Override: Pressing a number key while focused overwrites cell content immediately.
*   **5.3.2** `F2` Edit Mode: Places cursor at end of cell string for amendment.
*   **5.3.3** `Ctrl+D` Fill Down: Select column area, hits `Ctrl+D`. Value cascades. Array mutation is batched to Rust.
*   **5.3.4** `Shift+Space` Row Select: Highlights entire row.
*   **5.3.5** Escape Reset: Hit `Escape` during pending edit to revert to original cell value.
*   **Added: 5.3.6** Voice Entry: Mic icon â†’ Dictate values (Web Speech API).  

**[5.4] Advanced Top-Down Spreading (Complex Mutations)**
*   **5.4.0** Write-back on Summary: User types `500k` on `Total Global Revenue`.
*   **5.4.1** Modal Invocation: Spreading options wizard appears immediately.
*   **5.4.2** Method Select: Clicking "Proportional" calculates ratios based on previous year's actuals (historical trend).
*   **5.4.3** "Hold" Checkbox: User checks boxes next to `UK` and `France` (locking values).
*   **5.4.4** Execute Spread: Click `Apply`. Millions of cells recalculate. The UI receives the diff payload via WebSockets and flashes changed cells yellow.
*   **Added: 5.4.5** Preview Mode: Hover apply â†’ Ghost overlay shows post-spread values.  

**[5.5] Contextual Cell Formatting & Display Profiles**
*   **5.5.0** Format Parsing: Backend sends raw number `0.456`. UI applies display format profile defined on the Line Item (e.g., `%`, `2 dec`). Renders as `45.60%`.
*   **5.5.1** Conditional Heat-mapping: Attainment column. If > 100%, background gradient shifts dynamically toward green hex codes. If < 50%, sharp red.
*   **5.5.2** Lock State Indicators: Calculated formulas or closed periods show subtle gray strikethrough pattern on background, signaling read-only.
*   **5.5.3** Hover Details: Hovering over a truncated string `Desc...` shows full text tooltip.
*   **Added: 5.5.4** Bulk Coloring: Select range â†’ Drag color swatch â†’ Apply gradient/rules (e.g., if >threshold, red).  
*   **Added: 5.5.5** Dropdown Moat: Fuzzy search in dropdowns; validation (e.g., "Only numbers >0"); bulk apply to range.  

**[5.6] Time-Machine Reversion Slider**
*   **5.6.0** Activate Time Machine: Click `ðŸ•’`. Timeline UI slides up from bottom. Grid dims slightly.
*   **5.6.1** Scrub History: Drag slider thumb left.
*   **5.6.2** Temporal Re-Hydration: As slider moves, graph requests `SELECT state AT timestamp` from Clickhouse. Graph data visually morphs backward.
*   **5.6.3** "Restore Snapshot" Button: Rewinds the CRDT state for everyone.
*   **Added: 5.6.4** AI Time Insights: Slider stop â†’ "Key Changes: Variance spiked here".  

**[5.7] Performance "X-Ray Mode" Overlay**
*   **5.7.0** Developer Action: Hold `Alt`. Sidebar toggles to "Telemetry View".
*   **5.7.1** Geographic Heatmap: Canvas shader shifts color based on compute-time metadata (e.g., Red = >10ms per calculation).
*   **5.7.2** Click Bottleneck: Clicking a red cell pops open box the exact AtomScript sub-routine causing the latency.
*   **Added: 5.7.3** Auto-Optimize: Click "Fix" â†’ AI rewrites slow routines.  

---

## 6. Visual Intelligence (The Boardroom UX)
**[6.1] ECharts Integration & Reactivity**
*   **6.1.0** Add Chart Widget: Drag from palette, drop on canvas page.
*   **6.1.1** Bind Data: Modal opens, "Select Source module". Bind to `Sales Model`.
*   **6.1.2** Instant Reactivity: Type new sales number in adjacent grid widget. The paired EChart Bar snaps upward concurrently. No refresh.
*   **Added: 6.1.3** Multi-Grid Layout: Resize/snap grids on page; AI "Balance Layout" button.  

**[6.2] Automated Storytelling & Focus Presentation**
*   **6.2.0** Generate Narrative: Click `[Write Summary]`.
*   **6.2.1** Dynamic Text Injection: Returns text block: "Revenue is up **$400k** driven primarily by favorable FX in **EMEA**." Bolds are dynamic variables linked to grid cells.
*   **6.2.2** Focus Mode Toggle (ðŸ–¥ï¸): Top header hides. Sidebars retract. Font scales up 120% for meeting room visibility.
*   **Added: 6.2.3** Presentation Moat: Export to interactive deck (embed grids/charts).  

---

## 7. Data Hub (ETL & Ops)
**[7.1] The Visual Mapper Matrix**
*   **7.1.0** Column Analysis: Drag CSV. First 100 rows previewed in a staging table.
*   **7.1.1** Spline Drag: Click `Department` column header, drag bezier curve to `Cost Center` dimension icon.
*   **7.1.2** Transform Dialog: Click the connecting spline line. Opens box box. Type `Trim(Capitalize(Value))`.
*   **Added: 7.1.3** AI Auto-Map: Suggest connections based on data patterns.  

**[7.2] Redpanda Ingestion Streams & Job Status**
*   **7.2.0** Execute Job: Click `[Run Import]`.
*   **7.2.1** UI Unblocks: Dialog minimizes to standard progress widget. User goes to do other tasks.
*   **7.2.2** Streaming Progress: Circular progress SVG updates based on WebSocket events (e.g., "Written 1.4m of 2m rows.").
*   **7.2.3** Partial Failure Toast: "Import complete. 40 invalid rows rejected. [Download Error Log]".
*   **Added: 7.2.4** Job Reactions: Pause/resume import mid-stream.  

---

## 8. Multi-Player & CRDT Collaboration
**[8.1] Multiplayer Cursor & Presence Tracking**
*   **8.1.0** Avatar Indicators: Grid column headers show miniature user avatars indicating who is currently looking at this slice.
*   **8.1.1** Live Action Ghosting: If User 2 types into cell `[Q1 Sales]`, User 1 sees that cell flash blue with User 2's name attached, locking it momentarily to prevent conflict.
*   **8.1.2** Chat Integration: Click user avatar to open inline direct WebRTC chat (no leaving model).
*   **Added: 8.1.3** Co-Edit Reactions: Real-time merge on conflicts (CRDT auto-resolve).  

**[8.2] Granular Audit Lineage (Cell Right Click)**
*   **8.2.0** Context Menu Right Click: `[Show History]`.
*   **8.2.1** Audit Drawer: Opens box showing threaded list of modifications. "User A changed from 0 to 500 (10:01 AM)." "User B spread to 600 (11:05 AM)".
*   **8.2.2** Inline Comments: Users can add a text note to a specific cell value change ("Adjusted based on CEO request"). Cell gains a tiny red corner triangle in grid.
*   **Added: 8.2.3** AI Audit Insights: "This change caused 5% varianceâ€”review?"  

---

## 9. Fluxion AI Co-Pilot (Cortex Generative Layer)
**[9.1] Floating Assistant Interaction**
*   **9.1.0** Sparkle Toggle: Click `âœ¨`. Assistant drawer opens. Context-aware of current grid.
*   **9.1.1** NLP Command: "Add a 5% uplift to Q4 Forecast."
*   **9.1.2** Action Parsing: AI recognizes `Q4` and `Forecast` in currently viewed model. Drafts change payload.
*   **9.1.3** Confirmation Overlay: Grid highlights affected cells in purple. AI asks "Apply these pending changes?". User clicks `[Yes]`.
*   **Added: 9.1.4** Macro Flow Guidance: "Next Step: Approve Workflow" with auto-nav.  

**[9.2] Automated Anomaly Detection**
*   **9.2.0** Passive Monitoring: During idle time, AI scans vectors.
*   **9.2.1** Proactive Notification: Toast appears: "Historical anomaly detected. Q2 Marketing spend deviates 3-Sigma from historical median. [Investigate]".
*   **9.2.2** Investigation Flow: Clicking link opens box with auto-generated charts explicitly detailing the mathematical variance components.
*   **Added: 9.2.3** Auto-Fix Moat: AI proposes resolutions (e.g., "Apply Correction?") with previews.  

---

## 10. Security & System Administration
**[10.1] Live Role-Based Access Control (RBAC) Mutation**
*   **10.1.0** Matrix Editor: Rows = Roles. Cols = Dimensions/Apps. Checkboxes for Read/Write/None.
*   **10.1.1** Instant Enforcement: Admin changes `Junior Analyst -> Access (USA)` from `Write` to `Read`.
*   **10.1.2** WebSocket Force Protocol: The backend pushes an auth-invalidation to the active Junior Analyst. Their UI instantly grays out the USA grid columns, reverting them to read-only state without reloading the app.
*   **Added: 10.1.3** AI Role Suggest: "Based on usage, grant Analyst X forecast access?"  

**[10.2] Single Sign-On Diagnostics**
*   **10.2.0** Security Sandbox test: Admin tests new Okta binding in UI. Connection tracer pane shows step-by-step token validation hops.
*   **Added: 10.2.1** Moat Analytics: UX health dashboard (time-to-task, drop-off rates).  

---
*(End of Interaction Moat Manifesto)*
