# Design System Documentation

## Overview

Design system for the Professional Services Marketplace & Client Management Platform. Built with a token-based architecture for consistency across all UI surfaces.

## File Structure

| File | Purpose |
|------|---------|
| `design-tokens.json` | Source of truth — W3C Design Tokens Community Group format |
| `styles/tokens.css` | CSS custom properties generated from tokens (light + dark mode) |
| `tailwind.config.js` | Tailwind CSS v3 config extending defaults with token values |
| `designs/wireframes/landing-page.md` | ASCII wireframe for the public landing page |
| `designs/wireframes/dashboard.md` | ASCII wireframe for the authenticated dashboard |

## Color Palette

### Brand Colors
- **Primary** (Blue): `#3B82F6` (500) — used for CTAs, links, active states
- **Gray**: Neutral scale for text, borders, backgrounds

### Semantic Colors
- **Success** (Green): `#22C55E` (500) — confirmations, positive status
- **Warning** (Amber): `#F59E0B` (500) — caution states, pending status
- **Error** (Red): `#EF4444` (500) — errors, destructive actions

### Dark Mode
Toggle via `[data-theme="dark"]`, `.dark` class, or `prefers-color-scheme: dark`. Semantic surface/text/border tokens swap automatically via CSS custom properties.

## Typography

| Role | Family | Usage |
|------|--------|-------|
| `--font-sans` | Inter | Body text, UI elements |
| `--font-serif` | Playfair Display | Marketing headings |
| `--font-mono` | JetBrains Mono | Code, data values |

Font sizes use a modular scale from `xs` (0.75rem/12px) to `6xl` (3.75rem/60px).

## Component Tokens

Defined in `design-tokens.json` under `components`:

| Component | Variants | Sizes |
|-----------|----------|-------|
| **Button** | primary, secondary, outline, ghost, destructive | sm, md, lg |
| **Input** | default, error, disabled | sm, md, lg |
| **Card** | default, elevated, outlined, interactive | - |
| **Modal** | - | sm, md, lg, xl, full |
| **Badge** | default, success, warning, error, info | sm, md, lg |
| **Alert** | info, success, warning, error | - |
| **Table** | - (header, row, striped variants via props) | - |
| **Navigation** | sidebar (collapsible), topbar | - |

All interactive components include accessibility tokens (ARIA roles, focus indicators, 44px min touch targets).

## Spacing & Layout

- **Base unit**: 4px (0.25rem)
- **Scale**: 0, 0.5, 1, 1.5, 2 ... 32 (0 to 8rem)
- **Breakpoints**: sm (640px), md (768px), lg (1024px), xl (1280px)
- **Approach**: Mobile-first with min-width media queries

## Accessibility

- WCAG AA contrast ratios for all text/background pairs
- 44x44px minimum touch targets (`--touch-target-min`)
- Focus-visible rings on all interactive elements
- `prefers-reduced-motion` respected (use `transition` tokens, not hardcoded)
- Semantic ARIA roles documented per component

## Integration

### CSS Custom Properties
```html
<link rel="stylesheet" href="styles/tokens.css" />
```
Use variables directly: `color: var(--color-primary-600);`

### Tailwind CSS
```js
// Already configured in tailwind.config.js
// Use standard Tailwind classes: bg-primary-600, text-gray-900, etc.
```

### Design Tokens JSON
Import `design-tokens.json` for tooling integration (Style Dictionary, Figma plugins, etc.).
