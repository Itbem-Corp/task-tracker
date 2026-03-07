# Wireframe: Dashboard Layout

**Screen:** Authenticated Professional Dashboard
**Framework:** Next.js 14 (App Router) + shadcn/ui + Tailwind CSS
**Typography:** Inter (all dashboard text)
**Layout:** Sidebar + main content area, responsive

---

## Overall Layout Structure

```
┌──────────────────────────────────────────────────────────────────────┐
│ [TOP BAR — h-16, bg-white, border-b, z-sticky, px-6]               │
│ ┌──────────────────────────────────────────────────────────────────┐ │
│ │ [☰ Toggle]  [Search... 🔍]              [🔔 3] [Avatar ▼]     │ │
│ │  lg:hidden   flex-1 max-w-md             badge   dropdown      │ │
│ │  min-touch   h-10 bg-gray-50 rounded-lg  red dot  w-8 h-8     │ │
│ └──────────────────────────────────────────────────────────────────┘ │
│                                                                      │
│ ┌────────────┬───────────────────────────────────────────────────┐   │
│ │            │                                                   │   │
│ │  SIDEBAR   │   MAIN CONTENT AREA                               │   │
│ │  w-64      │   flex-1, bg-gray-50                              │   │
│ │  bg-white  │   p-6 lg:p-8                                      │   │
│ │  border-r  │   overflow-y-auto                                  │   │
│ │            │                                                   │   │
│ │  (details  │   (varies per page — see below)                   │   │
│ │   below)   │                                                   │   │
│ │            │                                                   │   │
│ │            │                                                   │   │
│ │            │                                                   │   │
│ │            │                                                   │   │
│ └────────────┴───────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────┘

Mobile (< lg): Sidebar hidden, toggleable via hamburger overlay
```

---

## Sidebar Navigation

```
[SIDEBAR — w-64, h-screen, fixed, bg-white, border-r border-gray-200]
┌────────────────────┐
│                    │
│  [Logo/Brand]      │  font-serif, text-xl, bold, text-primary-800
│  px-6 py-5         │  border-b border-gray-100
│                    │
│  MAIN MENU         │  text-xs uppercase text-gray-400 px-4 pt-6 pb-2
│                    │  tracking-wide font-semibold
│  ┌──────────────┐  │
│  │ 📊 Overview  │  │  ACTIVE: bg-primary-50, text-primary-700
│  │              │  │  font-semibold, rounded-lg
│  │ 👥 Clients   │  │  h-11 (44px touch target)
│  │              │  │  px-4, flex items-center gap-3
│  │ 📋 Services  │  │  text-sm, text-gray-700
│  │              │  │  hover: bg-gray-50
│  │ 💰 Payments  │  │  icon: w-5 h-5, text-gray-400
│  │              │  │  (active icon: text-primary-600)
│  │ 📅 Meetings  │  │
│  │              │  │  Badge for pending items:
│  │ 📊 Analytics │  │  ml-auto, badge-sm, badge-info
│  └──────────────┘  │
│                    │
│  SETTINGS          │  text-xs uppercase text-gray-400 px-4 pt-6 pb-2
│  ┌──────────────┐  │
│  │ ⚙ Settings   │  │  Same styling as main items
│  │ 🔗 Public    │  │  "Public Profile" links to /username
│  │    Profile   │  │
│  └──────────────┘  │
│                    │
│  ─────────────     │  border-t border-gray-100, mt-auto
│  ┌──────────────┐  │
│  │ [Avatar]      │  │  w-8 h-8 rounded-full
│  │ Sofia M.      │  │  text-sm font-medium text-gray-900
│  │ Pro Plan      │  │  text-xs text-gray-500
│  │ [Upgrade →]   │  │  text-xs text-primary-600 (if Starter)
│  └──────────────┘  │
│  px-4 py-4         │
└────────────────────┘

Accessibility:
- nav element with aria-label="Dashboard navigation"
- aria-current="page" on active item
- All items are links or buttons with min 44px height
```

---

## Page: Overview (Dashboard Home)

```
[MAIN CONTENT — p-6 lg:p-8]
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  PAGE HEADER:                                                       │
│  "Welcome back, Sofia"              text-2xl font-bold text-gray-900│
│  "Here's your business overview"    text-sm text-gray-500           │
│                                                                     │
│  [METRIC CARDS — grid cols-1 sm:cols-2 lg:cols-4, gap-6, mt-6]    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐              │
│  │Revenue   │ │Bookings  │ │Clients   │ │Avg Value │              │
│  │$4,250    │ │28        │ │15        │ │$151.78   │              │
│  │+12.5% ▲  │ │+8.3% ▲   │ │+3 new    │ │-2.1% ▼   │              │
│  │vs last mo│ │vs last mo│ │this mo   │ │vs last mo│              │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘              │
│  Cards: bg-white, rounded-xl, p-6, border border-gray-200          │
│  Metric: text-3xl font-bold text-gray-900                           │
│  Change: text-sm, success-600 for ▲, error-600 for ▼               │
│  Label: text-sm text-gray-500                                       │
│                                                                     │
│  [TWO-COLUMN LAYOUT — grid cols-1 lg:cols-3, gap-6, mt-6]         │
│                                                                     │
│  ┌───────────────────────────────┐ ┌────────────────────┐          │
│  │ REVENUE CHART (lg:col-span-2) │ │ UPCOMING MEETINGS  │          │
│  │ bg-white, rounded-xl, p-6     │ │ bg-white rounded-xl│          │
│  │ border border-gray-200        │ │ p-6 border         │          │
│  │                               │ │                    │          │
│  │ "Revenue Overview"  text-lg   │ │ "Upcoming" text-lg │          │
│  │ [7d] [30d] [90d] [12m]       │ │                    │          │
│  │  tab buttons, text-sm         │ │ ┌────────────────┐ │          │
│  │                               │ │ │ Mar 8, 10:00am │ │          │
│  │ ┌───────────────────────┐     │ │ │ Sofia ↔ Elena  │ │          │
│  │ │                       │     │ │ │ 1:1 Consulting │ │          │
│  │ │   [Line Chart Area]   │     │ │ │ $150           │ │          │
│  │ │   h-64                │     │ │ │ badge-success  │ │          │
│  │ │   primary-500 line    │     │ │ └────────────────┘ │          │
│  │ │   primary-50 fill     │     │ │ ┌────────────────┐ │          │
│  │ │                       │     │ │ │ Mar 9, 2:00pm  │ │          │
│  │ └───────────────────────┘     │ │ │ ...            │ │          │
│  │                               │ │ └────────────────┘ │          │
│  └───────────────────────────────┘ │                    │          │
│                                    │ [View All →]       │          │
│                                    │ text-sm primary    │          │
│                                    └────────────────────┘          │
│                                                                     │
│  [RECENT ACTIVITY TABLE — mt-6]                                     │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ "Recent Activity"                    [Export CSV ↓]          │    │
│  │ text-lg font-semibold                ghost btn text-sm      │    │
│  │                                                             │    │
│  │ ┌────────┬──────────┬─────────┬────────┬────────┐          │    │
│  │ │ Date   │ Client   │ Service │ Amount │ Status │          │    │
│  │ ├────────┼──────────┼─────────┼────────┼────────┤          │    │
│  │ │ Mar 7  │ Elena R. │ 1:1     │ $150   │ ● Paid │          │    │
│  │ │ Mar 6  │ James L. │ Workshop│ $500   │ ● Paid │          │    │
│  │ │ Mar 5  │ Ana K.   │ 1:1     │ $150   │ ○ Pend │          │    │
│  │ └────────┴──────────┴─────────┴────────┴────────┘          │    │
│  │                                                             │    │
│  │ Status badges: success for Paid, warning for Pending        │    │
│  │ Table: bg-white, rounded-xl, border, overflow-x-auto        │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Page: Clients

```
[MAIN CONTENT]
┌─────────────────────────────────────────────────────────────────────┐
│  PAGE HEADER:                                                       │
│  "Clients"                          text-2xl font-bold              │
│  [+ Add Client]                     primary btn, min-h-touch        │
│                                                                     │
│  [FILTERS — flex gap-4, mt-4]                                       │
│  [Search by name or email...🔍]     input, flex-1, max-w-sm         │
│  [Sort: Recent ▼]                   select dropdown                  │
│                                                                     │
│  [CLIENT TABLE — mt-6]                                              │
│  ┌──────┬──────────┬─────────┬──────────┬──────────┬────────┐      │
│  │      │ Name     │ Email   │ Bookings │ Revenue  │ Last   │      │
│  │      │          │         │          │          │ Booking│      │
│  ├──────┼──────────┼─────────┼──────────┼──────────┼────────┤      │
│  │ [Av] │ Elena R. │ e@...   │ 8        │ $1,200   │ Mar 7  │      │
│  │  ○   │          │         │          │          │        │      │
│  │ [Av] │ James L. │ j@...   │ 3        │ $1,500   │ Mar 6  │      │
│  │ [Av] │ Ana K.   │ a@...   │ 12       │ $1,800   │ Mar 5  │      │
│  └──────┴──────────┴─────────┴──────────┴──────────┴────────┘      │
│                                                                     │
│  Row click → navigates to client detail page                        │
│  Avatars: w-8 h-8 rounded-full bg-primary-100 text-primary-700     │
│  hover: bg-gray-50 on rows                                          │
│                                                                     │
│  [Pagination — flex justify-between, mt-4]                          │
│  "Showing 1-10 of 15"              [← Prev] [1] [2] [Next →]       │
│  text-sm text-gray-500             ghost btns, active: primary bg   │
│                                    min-h-touch                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Page: Meetings

```
[MAIN CONTENT]
┌─────────────────────────────────────────────────────────────────────┐
│  PAGE HEADER:                                                       │
│  "Meetings"                         text-2xl font-bold              │
│  [Upcoming] [Past] [Canceled]       tab buttons, underline active   │
│                                                                     │
│  [MEETING CARDS — grid cols-1 md:cols-2, gap-4, mt-6]             │
│  ┌──────────────────────────────┐ ┌──────────────────────────────┐ │
│  │ 📅 March 8, 2026             │ │ 📅 March 9, 2026             │ │
│  │ 10:00 AM — 11:00 AM (EST)   │ │ 2:00 PM — 3:00 PM (EST)     │ │
│  │                              │ │                              │ │
│  │ Client: Elena Rodriguez      │ │ Client: James Liu            │ │
│  │ Service: 1:1 Consulting      │ │ Service: Strategy Workshop   │ │
│  │ Amount: $150                 │ │ Amount: $500                 │ │
│  │                              │ │                              │ │
│  │ [badge: Confirmed ●]        │ │ [badge: Confirmed ●]        │ │
│  │                              │ │                              │ │
│  │ [Join Meeting]  [Cancel]     │ │ [Join Meeting]  [Cancel]     │ │
│  │  primary btn-sm  ghost-sm   │ │  primary btn-sm  ghost-sm   │ │
│  │  min-h-touch                 │ │                              │ │
│  └──────────────────────────────┘ └──────────────────────────────┘ │
│                                                                     │
│  Cards: bg-white, rounded-xl, p-6, border border-gray-200          │
│  Date: text-lg font-semibold text-gray-900                          │
│  Time: text-sm text-gray-600                                        │
│  Client/Service: text-sm text-gray-700                              │
│  Badge: badge-success                                               │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Page: Analytics

```
[MAIN CONTENT]
┌─────────────────────────────────────────────────────────────────────┐
│  PAGE HEADER:                                                       │
│  "Analytics"                        text-2xl font-bold              │
│  [Date range: Last 30 days ▼]       select, min-h-touch             │
│                                                                     │
│  [METRIC CARDS — same as overview but with more detail]            │
│  4-column grid of KPI cards                                         │
│                                                                     │
│  [CHARTS SECTION — grid cols-1 lg:cols-2, gap-6, mt-6]            │
│  ┌──────────────────────────┐ ┌──────────────────────────┐         │
│  │ Revenue Over Time        │ │ Bookings by Service      │         │
│  │ [Line chart, h-72]       │ │ [Bar chart, h-72]        │         │
│  │ X: dates, Y: revenue     │ │ X: service names         │         │
│  └──────────────────────────┘ └──────────────────────────┘         │
│                                                                     │
│  ┌──────────────────────────┐ ┌──────────────────────────┐         │
│  │ New vs Returning Clients │ │ Top Clients by Revenue   │         │
│  │ [Donut chart, h-72]      │ │ [Horizontal bar, h-72]   │         │
│  │ Legend below              │ │ Top 5 clients             │         │
│  └──────────────────────────┘ └──────────────────────────┘         │
│                                                                     │
│  Charts: bg-white, rounded-xl, p-6, border                          │
│  Starter plan: charts 3 & 4 blurred with upgrade overlay            │
│  Overlay: "Upgrade to Pro for advanced analytics" + btn             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Page: Settings

```
[MAIN CONTENT]
┌─────────────────────────────────────────────────────────────────────┐
│  PAGE HEADER:                                                       │
│  "Settings"                         text-2xl font-bold              │
│                                                                     │
│  [TAB NAV — border-b]                                               │
│  [Profile] [Availability] [Billing] [Notifications]                 │
│   active: border-b-2 border-primary-600 text-primary-600            │
│   min-h-touch per tab                                               │
│                                                                     │
│  [SETTINGS FORM — max-w-2xl, mt-8]                                 │
│  ┌─────────────────────────────────────────────────────┐            │
│  │ Profile Photo                                       │            │
│  │ [Current photo ○]  [Upload New]  [Remove]           │            │
│  │  w-20 h-20          secondary     ghost destructive │            │
│  │  rounded-full       btn-sm        btn-sm            │            │
│  │                                                     │            │
│  │ Full Name           [Sofia Martinez          ]      │            │
│  │ label text-sm       input w-full                    │            │
│  │ font-medium                                         │            │
│  │                                                     │            │
│  │ Email               [sofia@example.com       ]      │            │
│  │                     input disabled                  │            │
│  │                                                     │            │
│  │ Tagline             [Marketing consultant... ]      │            │
│  │                     input, max 100 chars             │            │
│  │                     char counter: text-xs gray-400   │            │
│  │                                                     │            │
│  │ Bio                 [textarea, 4 rows        ]      │            │
│  │                     max 500 chars                    │            │
│  │                                                     │            │
│  │ Specialties         [Tag input field          ]     │            │
│  │                     chips: bg-primary-100 rounded    │            │
│  │                     x to remove                     │            │
│  │                                                     │            │
│  │ Category            [Consulting ▼             ]     │            │
│  │                     select dropdown                  │            │
│  │                                                     │            │
│  │ Timezone            [America/New_York ▼       ]     │            │
│  │                                                     │            │
│  │ Social Links                                        │            │
│  │ Twitter             [@sofiamarketing          ]     │            │
│  │ LinkedIn            [linkedin.com/in/sofia    ]     │            │
│  │ Website             [sofiaconsulting.com       ]     │            │
│  │                                                     │            │
│  │              [Cancel]  [Save Changes]               │            │
│  │               ghost     primary btn                 │            │
│  │               min-h-touch                           │            │
│  └─────────────────────────────────────────────────────┘            │
│                                                                     │
│  All inputs: h-11 (44px), rounded-lg, border border-gray-300       │
│  Focus: ring-2 ring-primary-500                                     │
│  Labels: text-sm font-medium text-gray-700 mb-1.5                  │
│  Error state: border-error-500, error message text-sm text-error-600│
└─────────────────────────────────────────────────────────────────────┘
```

---

## Responsive Behavior

| Breakpoint | Layout |
|------------|--------|
| < 640px | No sidebar (overlay), single-column cards, stacked metrics |
| 640-767px | No sidebar (overlay), 2-column metric cards |
| 768-1023px | No sidebar (overlay), 2-column grids |
| >= 1024px | Fixed sidebar visible, multi-column content |
| >= 1280px | Max-width content area with generous padding |

## Accessibility Notes

- Sidebar nav uses `<nav>` with `aria-label="Dashboard navigation"`
- Active page marked with `aria-current="page"`
- All interactive elements meet 44px min touch target
- Tables use proper `<th scope="col">` and `<caption>`
- Charts include `aria-label` describing the data
- Focus visible on all interactive elements (ring-2 ring-primary-500 ring-offset-2)
- Skip-to-content link targets main content area
- Notification badge: `aria-label="3 unread notifications"`
