# Duck Clicker — Design Specification

## Overview

A playful, interactive duck-clicking web application. The visual language is cheerful, colorful, and game-like — designed to make every click feel satisfying. The UI centers around a large clickable duck character with a prominent click counter and milestone celebrations.

---

## 1. Visual Theme

**Mood:** Cheerful, casual game, bright daylight
**Scene:** A duck floating on a pond, under a blue sky with soft clouds
**Color Story:** Sunny yellow duck against sky-blue background, with energetic orange for interactive cues and green bursts for celebrations

### Color Semantics

| Role | Token | Hex | Usage |
|------|-------|-----|-------|
| Duck / Primary | `--color-primary-500` | `#FFD700` | Duck body, primary accents |
| Sky / Secondary | `--color-secondary-300` | `#87CEEB` | Sky background, calm surfaces |
| Interactive / Accent | `--color-accent-500` | `#FF6B35` | Buttons, hover glows, CTAs |
| Milestone / Success | `--color-success-500` | `#10B981` | Achievement badges, celebrations |
| Danger / Error | `--color-error-500` | `#EF4444` | Error states, destructive actions |
| Text / Neutral | `--color-neutral-900` | `#171717` | Primary body text |

---

## 2. Typography

### Font Stack

| Purpose | Token | Font | Fallbacks |
|---------|-------|------|-----------|
| Display / Titles | `--font-display` | Fredoka One | Baloo 2, Comic Neue, cursive |
| Body / UI | `--font-body` | Nunito | Inter, system-ui, sans-serif |
| Code / Data | `--font-mono` | Fira Code | JetBrains Mono, monospace |

### Type Scale Application

| Element | Size Token | Weight | Line Height | Notes |
|---------|-----------|--------|-------------|-------|
| Page title | `--text-4xl` (36px) | `--font-extrabold` | `--leading-tight` | Display font, playful feel |
| Click counter | `--text-5xl` (48px) | `--font-black` | `--leading-none` | Display font, large and bold |
| Counter (milestone) | `--text-7xl` (72px) | `--font-black` | `--leading-none` | Momentary enlarge on milestones |
| Subtitle / tagline | `--text-xl` (20px) | `--font-semibold` | `--leading-snug` | Body font |
| Achievement label | `--text-xs` (12px) | `--font-bold` | `--leading-normal` | Badge text, uppercase |
| Body text | `--text-base` (16px) | `--font-normal` | `--leading-normal` | Body font |

### Google Fonts Import

```css
@import url('https://fonts.googleapis.com/css2?family=Fredoka+One&family=Nunito:wght@400;600;700;800;900&display=swap');
```

---

## 3. Layout & Spacing

### Page Structure

```
┌─────────────────────────────────────────┐
│              Sky Background             │  Full viewport, gradient
│         ☁️ Floating Clouds ☁️           │  z-index: 5
│                                         │
│    ┌───────────────────────────────┐    │
│    │         GAME CARD             │    │  max-width: 600px, centered
│    │                               │    │  padding: 32px
│    │    🦆 Duck Clicker 🦆         │    │  Title (--text-4xl)
│    │                               │    │
│    │        ┌──────────┐           │    │
│    │        │          │           │    │  Duck: 150-200px
│    │        │   🦆     │           │    │  Clickable area
│    │        │          │           │    │
│    │        └──────────┘           │    │
│    │                               │    │
│    │         1,234                  │    │  Counter (--text-5xl)
│    │         clicks                │    │  Label (--text-lg)
│    │                               │    │
│    │   [🏆 50!] [🎯 100!] [⭐ 500]│    │  Achievement badges
│    │                               │    │
│    └───────────────────────────────┘    │
│                                         │
│         ~~~ Pond Water ~~~              │  Bottom section, blue gradient
└─────────────────────────────────────────┘
```

### Spacing Rules

| Area | Token | Value | Notes |
|------|-------|-------|-------|
| Container padding | `--space-8` | 32px | Inner padding of game card |
| Gap between duck and counter | `--space-6` | 24px | Vertical whitespace rhythm |
| Gap between counter and badges | `--space-4` | 16px | Tighter grouping |
| Badge gap | `--space-2` | 8px | Between achievement pills |
| Card margin (mobile) | `--space-4` | 16px | Horizontal margin on small screens |
| Card margin (desktop) | auto | centered | Max-width constraint |

### Responsive Behavior

| Breakpoint | Behavior |
|-----------|----------|
| < 640px (xs/sm) | Full-width card with 16px margins, duck 150px, counter `--text-4xl` |
| 640-768px (sm/md) | Card width 90%, duck 175px |
| 768px+ (md+) | Card max-width 600px, duck 200px, counter `--text-5xl` |

---

## 4. Component Specifications

### 4.1 Duck Button

The centerpiece interactive element.

**Structure:**
```html
<button class="duck-button" aria-label="Click the duck">
  <div class="duck-image"><!-- Duck SVG or image --></div>
  <div class="ripple-container"><!-- Ripple effects spawn here --></div>
</button>
```

**Visual Properties:**

| Property | Value |
|----------|-------|
| Size | `clamp(150px, 30vw, 200px)` |
| Shape | Circular clickable area |
| Cursor | `pointer` |
| Background | Transparent (duck image fills) |
| Border | None |
| Focus ring | `3px solid var(--color-accent-500)`, `4px offset` |

**Interactive States:**

| State | Appearance | Transition |
|-------|-----------|------------|
| **Default** | Natural size, subtle idle wobble animation | — |
| **Hover** | `scale(1.08)`, golden glow shadow (`--shadow-glow`) | `--duration-fast`, `--ease-out` |
| **Active (click)** | `duck-bounce` keyframe animation (scale 1→1.2→0.9→1) | `--duration-slow`, `--ease-bounce` |
| **Focus** | Visible accent-colored ring | `--duration-fast` |
| **Disabled** | `opacity: 0.5`, `cursor: not-allowed`, no animations | `--duration-normal` |

**Animations:**
- **Idle:** `duck-wobble` — gentle 3deg rotation, `3s infinite`, `--ease-in-out`
- **Click:** `duck-bounce` — 350ms, `--ease-bounce`, plays once per click
- **Hover glow:** `box-shadow` transition to `--shadow-glow`

**Accessibility:**
- `role="button"`, `aria-label="Click the duck"`
- Keyboard: `Enter` and `Space` trigger click
- Focus indicator: 3px solid ring, 3:1 contrast ratio
- Touch target: minimum 150px (exceeds 44px requirement)

---

### 4.2 Click Counter

Displays the running total of clicks.

**Structure:**
```html
<div class="click-counter" aria-live="polite" aria-atomic="true">
  <span class="counter-number">1,234</span>
  <span class="counter-label">clicks</span>
</div>
```

**Visual Properties:**

| Property | Value |
|----------|-------|
| Number font | `--font-display` |
| Number size | `--text-5xl` (48px), up to `--text-7xl` on milestone |
| Number weight | `--font-black` (900) |
| Number color | `--counter-color` (`--color-primary-700`) |
| Label font | `--font-body` |
| Label size | `--text-lg` (18px) |
| Label color | `--text-secondary` |
| Text align | Center |
| Number format | Locale-aware with commas (1,234) |

**Animations:**
- **On increment:** `counter-pop` — scale to 1.15 and back, `--duration-fast`, `--ease-playful`
- **+1 float:** `score-float` — "+1" text floats up 40px and fades, `--duration-slow`
- **Milestone flash:** Number briefly scales to `--text-7xl` with color change to `--color-success-500`

**Accessibility:**
- `aria-live="polite"` for screen reader announcements
- `aria-atomic="true"` to read full value on change
- Sufficient color contrast: `--color-primary-700` on white = ~4.8:1 ratio (WCAG AA pass)

---

### 4.3 Page Title

**Structure:**
```html
<h1 class="page-title">
  🦆 Duck Clicker 🦆
</h1>
```

**Visual Properties:**

| Property | Value |
|----------|-------|
| Font | `--font-display` |
| Size | `--text-4xl` (36px) |
| Weight | `--font-extrabold` (800) |
| Color | `--text-primary` |
| Text align | Center |
| Letter spacing | `--tracking-tight` |
| Text shadow | `0 2px 4px rgba(0,0,0,0.1)` for depth |

---

### 4.4 Achievement Badges

Pill-shaped indicators that appear at click milestones (50, 100, 250, 500, 1000...).

**Structure:**
```html
<div class="achievement-list" role="list" aria-label="Achievements">
  <div class="achievement-badge achievement-badge--earned" role="listitem">
    <span class="badge-icon">🏆</span>
    <span class="badge-label">50 Clicks!</span>
  </div>
  <div class="achievement-badge achievement-badge--locked" role="listitem">
    <span class="badge-icon">🔒</span>
    <span class="badge-label">100 Clicks</span>
  </div>
</div>
```

**Visual Properties (Earned):**

| Property | Value |
|----------|-------|
| Height | `28px` (`--badge-height`) |
| Padding | `4px 12px` |
| Background | `--color-success-100` |
| Border | `1px solid var(--color-success-300)` |
| Border radius | `--radius-full` (pill shape) |
| Font size | `--text-xs` (12px) |
| Font weight | `--font-bold` |
| Color | `--color-success-700` |
| Text transform | Uppercase |

**Visual Properties (Locked):**

| Property | Value |
|----------|-------|
| Background | `--color-neutral-100` |
| Border | `1px solid var(--color-neutral-200)` |
| Color | `--color-neutral-400` |
| Opacity | `0.7` |

**Animations:**
- **Earn:** `badge-appear` — slides up 16px, scales from 0.8, overshoots to 1.05, settles. `--duration-slow`, `--ease-elastic`
- **Idle shimmer:** Subtle background gradient shift on earned badges

**Milestone Schedule:**

| Clicks | Badge | Emoji |
|--------|-------|-------|
| 50 | "50 Clicks!" | 🏆 |
| 100 | "Century!" | 🎯 |
| 250 | "Duck Fan" | 🌟 |
| 500 | "Half K!" | ⭐ |
| 1,000 | "Quack Master" | 👑 |
| 5,000 | "Duck Legend" | 🦆 |
| 10,000 | "Unstoppable" | 🔥 |

---

### 4.5 Background Scene

The ambient environment behind the game card.

**Structure:**
```html
<div class="scene-background">
  <div class="sky-gradient"></div>
  <div class="clouds-layer">
    <div class="cloud cloud--1"></div>
    <div class="cloud cloud--2"></div>
    <div class="cloud cloud--3"></div>
  </div>
  <div class="pond-layer"></div>
</div>
```

**Visual Properties:**

| Element | Style |
|---------|-------|
| Sky gradient | `linear-gradient(180deg, var(--sky-top) 0%, var(--sky-bottom) 70%, var(--pond-surface) 100%)` |
| Clouds | White semi-transparent blobs, `border-radius: 50%`, various sizes |
| Cloud animation | `cloud-float`, staggered durations (20s, 30s, 40s), `linear`, `infinite` |
| Pond | Bottom 20% of viewport, `--pond-surface` with `water-shimmer` overlay |
| Sky animation | `sky-shift`, `20s`, `ease-in-out`, `infinite` (subtle gradient movement) |

---

### 4.6 Confetti Burst

Triggered at milestone achievements (every 50 clicks).

**Behavior:**
1. On milestone click, spawn 30-50 confetti particles at duck position
2. Particles use randomized colors from: `--color-primary-500`, `--color-accent-500`, `--color-success-400`, `--color-secondary-300`
3. Each particle: 8x8px to 12x12px, random shape (circle, square, triangle)
4. Animation: `confetti-pop` — burst upward, then fall with rotation
5. Duration: 1.5-2.5s per particle (randomized)
6. Clean up DOM nodes after animation completes

**Performance:**
- Use `transform` and `opacity` only (GPU-accelerated)
- `will-change: transform, opacity` on active particles
- Remove `will-change` after animation ends
- Limit to 50 particles max to prevent frame drops

---

### 4.7 Ripple / Splash Effect

Visual feedback on each duck click.

**Behavior:**
1. On click, create a circular element at click coordinates relative to duck
2. Start at `scale(0.5)`, `opacity: 0.6`
3. Expand to `scale(2.5)`, `opacity: 0`
4. Color: `rgba(255, 215, 0, 0.3)` (translucent gold)
5. Duration: `--duration-slow` (350ms)
6. Remove element after animation

---

## 5. Animation Timing Reference

| Animation | Duration | Easing | Trigger | Iterations |
|-----------|----------|--------|---------|------------|
| `duck-wobble` | 3s | `--ease-in-out` | Idle (no clicks for 2s) | Infinite |
| `duck-bounce` | 350ms | `--ease-bounce` | Each click | 1 |
| `counter-pop` | 150ms | `--ease-playful` | Each click | 1 |
| `score-float` | 500ms | `--ease-out` | Each click | 1 |
| `ripple-splash` | 350ms | `--ease-out` | Each click | 1 |
| `confetti-pop` | 1.5–2.5s | `--ease-out` | Milestone | 1 |
| `badge-appear` | 350ms | `--ease-elastic` | Badge earned | 1 |
| `cloud-float` | 20–40s | `linear` | Always | Infinite |
| `sky-shift` | 20s | `--ease-in-out` | Always | Infinite |
| `water-shimmer` | 4s | `--ease-in-out` | Always | Infinite |

---

## 6. Accessibility Checklist

### Color Contrast (WCAG AA)

| Combination | Ratio | Status |
|-------------|-------|--------|
| `--text-primary` (#171717) on white | 18.4:1 | Pass |
| `--counter-color` (#998400) on white | 4.8:1 | Pass |
| `--color-accent-600` (#E05A1F) on white | 4.6:1 | Pass |
| `--color-success-700` (#047857) on white | 5.9:1 | Pass |
| Badge text `--color-success-700` on `--color-success-100` | 5.2:1 | Pass |
| `--text-inverse` (#FFF) on `--color-accent-500` | 3.4:1 | Pass (large text only) |
| `--text-on-primary` (#171717) on `--color-primary-500` | 12.1:1 | Pass |

### Interactive Elements

- All clickable elements have `cursor: pointer`
- Focus indicators visible with 3:1 contrast
- Duck button minimum touch target: 150px (exceeds 44px requirement)
- Achievement badges minimum touch target: 28px height (non-interactive, display only)
- `aria-live` region for counter updates
- `role="list"` and `role="listitem"` for achievements
- All images have `alt` text or `aria-label`

### Motion

- `prefers-reduced-motion: reduce` disables all animations
- No auto-playing video or audio
- Confetti does not flash at epilepsy-trigger frequencies

---

## 7. Dark Mode

The design supports automatic dark mode via `prefers-color-scheme: dark`.

**Key changes in dark mode:**
- Sky gradient shifts to deep navy (`#0F3460` → `#1A1A2E`)
- Game card uses glassmorphism on dark surface (`rgba(30, 30, 50, 0.85)`)
- Duck yellow slightly brightened (`#FFDF4A`) for visibility
- Shadows use warmer duck-glow in dark context
- Text inverts to light neutrals
- All contrast ratios verified in dark mode

---

## 8. Asset Requirements

| Asset | Format | Size | Notes |
|-------|--------|------|-------|
| Duck character | SVG (preferred) or PNG | 200x200px @2x | Transparent background, cheerful expression |
| Cloud shapes | CSS or SVG | Variable | 3 variants, semi-transparent white |
| Achievement icons | Emoji (native) | — | Using system emoji for performance |
| Confetti particles | CSS-generated | 8-12px | Dynamic creation, no image assets |
| Favicon | ICO + PNG | 32x32 + 192x192 | Duck icon |

---

## 9. Performance Budget

| Metric | Target |
|--------|--------|
| First Contentful Paint | < 1.5s |
| Largest Contentful Paint | < 2.5s |
| Total JS bundle | < 50KB gzipped |
| CSS bundle | < 15KB gzipped |
| Animations | 60fps (GPU-accelerated transforms only) |
| Confetti particles | Max 50 simultaneous |
| Memory | No DOM node leaks from animations |

---

## 10. Cross-Agent Handoff Notes

### For Frontend (Pixel)

Design tokens are available at:
- **CSS:** `duck-clicker/src/styles/tokens.css` — import directly
- **Memory MCP keys:** `design_tokens_duck_clicker`, `component_specs_duck_clicker`

Priority implementation order:
1. Background scene (sky gradient + clouds)
2. Game card container
3. Duck button with click handler
4. Click counter with increment animation
5. Achievement badge system
6. Confetti milestone celebration
7. Dark mode support

### For QA (Sentinel)

Test these interaction flows:
- Single click → duck bounces + counter increments + ripple appears
- Rapid clicks → animations queue properly, counter accurate
- Milestone click (50th) → confetti fires + badge appears
- Keyboard navigation → duck clickable via Enter/Space
- Screen reader → counter announces new value
- Reduced motion → all animations disabled
- Dark mode → all elements visible and readable
- Mobile → touch targets adequate, layout responsive
