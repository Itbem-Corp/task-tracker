# Wireframe: Landing Page

**Screen:** Public Marketing Landing Page
**Framework:** Astro.js + Tailwind CSS
**Typography:** Playfair Display (headings) + Inter (body)
**Layout:** Single-column, full-width sections, max-content 1280px centered

---

## Global Navigation Bar

```
[STICKY TOP BAR — h-16, bg-white, border-b border-gray-200, z-sticky]
┌─────────────────────────────────────────────────────────────────────────┐
│  [Logo/Brand]          [Features] [Pricing] [Testimonials]   [Login] [Get Started →] │
│  font-serif bold 2xl   font-sans text-sm text-gray-600       ghost    primary btn    │
│  text-primary-800      gap-8, hover:text-primary-600         btn-sm   btn-sm         │
└─────────────────────────────────────────────────────────────────────────┘

Mobile (< md):
┌─────────────────────────────────────────────────────────────────────────┐
│  [Logo]                                              [Hamburger ☰]     │
│                                                      min-touch 44x44   │
└─────────────────────────────────────────────────────────────────────────┘
  → Opens full-screen overlay with stacked nav links + CTA buttons
```

---

## Section 1: Hero

```
[SECTION — py-24 lg:py-32, bg-white]
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│  [Two-column layout: lg:grid-cols-2, gap-16, items-center]             │
│                                                                         │
│  LEFT COLUMN (text):                                                    │
│  ┌───────────────────────────────────────┐                              │
│  │ Badge: "For Independent Professionals"│  badge-info, text-xs         │
│  │                                       │                              │
│  │ HEADLINE (font-serif):                │                              │
│  │ "Sell Your Services.                  │  text-5xl lg:text-6xl        │
│  │  Manage Clients.                      │  font-extrabold              │
│  │  Grow Your Business."                 │  text-gray-900               │
│  │                                       │  tracking-tight              │
│  │ SUBHEADLINE:                          │                              │
│  │ "The all-in-one platform for          │  text-xl, text-gray-600      │
│  │  consultants, coaches, and freelancers│  leading-relaxed             │
│  │  to list services, book clients, and  │  max-w-lg                    │
│  │  get paid — in minutes."              │                              │
│  │                                       │                              │
│  │ [Get Started Free →]  [Watch Demo ▶]  │  primary btn-lg + ghost btn  │
│  │                                       │  min-h-touch (44px)          │
│  │ "No credit card required. Free 14-day │  text-sm, text-gray-400      │
│  │  trial."                              │                              │
│  └───────────────────────────────────────┘                              │
│                                                                         │
│  RIGHT COLUMN (visual):                                                 │
│  ┌───────────────────────────────────────┐                              │
│  │                                       │                              │
│  │  [Dashboard Preview Image/Mockup]     │  rounded-2xl, shadow-2xl     │
│  │  Showing the professional dashboard   │  border border-gray-200      │
│  │  with booking calendar and revenue    │  aspect-[4/3]                │
│  │  chart visible                        │  bg-gray-100 (placeholder)   │
│  │                                       │                              │
│  └───────────────────────────────────────┘                              │
│                                                                         │
│  Mobile: stacks vertically (text first, image below)                    │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 2: Social Proof Bar

```
[SECTION — py-12, bg-gray-50, border-y border-gray-100]
┌─────────────────────────────────────────────────────────────────────────┐
│  "Trusted by 500+ professionals worldwide"  text-sm text-gray-400      │
│                                              text-center                │
│  [Logo 1]  [Logo 2]  [Logo 3]  [Logo 4]  [Logo 5]                     │
│  grayscale, opacity-50, hover:opacity-100                               │
│  flex justify-center gap-12, items-center                               │
│                                                                         │
│  Mobile: 2-row grid, 3+2 logos                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 3: Features Grid

```
[SECTION — py-20 lg:py-24, bg-white]
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│  SECTION HEADER (text-center, max-w-2xl mx-auto):                      │
│  "Everything you need to                     font-serif, text-4xl       │
│   run your service business"                 font-bold, text-gray-900   │
│  "Stop juggling multiple tools. One          text-lg, text-gray-600     │
│   platform for listings, bookings,           mt-4                       │
│   payments, and client management."                                     │
│                                                                         │
│  [FEATURES GRID — grid cols-1 md:cols-2 lg:cols-3, gap-8, mt-16]      │
│                                                                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐          │
│  │ [Icon: 🌐]      │ │ [Icon: 📅]      │ │ [Icon: 💳]      │          │
│  │ w-12 h-12       │ │                  │ │                  │          │
│  │ bg-primary-100   │ │ bg-success-100   │ │ bg-warning-100   │          │
│  │ text-primary-600 │ │ text-success-600 │ │ text-warning-600 │          │
│  │ rounded-xl       │ │                  │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ "Professional    │ │ "Smart           │ │ "Instant         │          │
│  │  Storefront"     │ │  Scheduling"     │ │  Payments"       │          │
│  │ text-lg bold     │ │ text-lg bold     │ │ text-lg bold     │          │
│  │ text-gray-900    │ │                  │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ "Set up your     │ │ "Clients book    │ │ "Accept payments │          │
│  │  professional    │ │  based on your   │ │  via Stripe.     │          │
│  │  page in under   │ │  real-time       │ │  Automatic       │          │
│  │  15 minutes."    │ │  availability."  │ │  invoicing."     │          │
│  │ text-sm          │ │                  │ │                  │          │
│  │ text-gray-500    │ │                  │ │                  │          │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘          │
│                                                                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐          │
│  │ [Icon: 👥]      │ │ [Icon: 📊]      │ │ [Icon: 🔔]      │          │
│  │ bg-error-100     │ │ bg-primary-100   │ │ bg-success-100   │          │
│  │                  │ │                  │ │                  │          │
│  │ "Client          │ │ "Business        │ │ "Automated       │          │
│  │  Management"     │ │  Analytics"      │ │  Notifications"  │          │
│  │                  │ │                  │ │                  │          │
│  │ "Track clients,  │ │ "Revenue charts, │ │ "Email confirms, │          │
│  │  notes, and      │ │  booking stats,  │ │  reminders, and  │          │
│  │  booking history  │ │  and retention   │ │  calendar invites │          │
│  │  in one place."  │ │  insights."      │ │  — automatic."   │          │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘          │
│                                                                         │
│  Each feature card: p-6, rounded-xl, bg-white, border border-gray-200  │
│  hover: shadow-md, transition-base                                      │
│  Mobile: single column stack                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 4: How It Works

```
[SECTION — py-20, bg-gray-50]
┌─────────────────────────────────────────────────────────────────────────┐
│  SECTION HEADER:                                                        │
│  "Get started in 3 simple steps"          font-serif, text-4xl, bold    │
│                                                                         │
│  [3-COLUMN LAYOUT — grid cols-1 md:cols-3, gap-12, mt-16]             │
│                                                                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐          │
│  │  ① (Step 1)     │ │  ② (Step 2)     │ │  ③ (Step 3)     │          │
│  │  w-10 h-10      │ │  w-10 h-10      │ │  w-10 h-10      │          │
│  │  bg-primary-600  │ │  bg-primary-600  │ │  bg-primary-600  │          │
│  │  text-white      │ │                  │ │                  │          │
│  │  rounded-full    │ │                  │ │                  │          │
│  │  font-bold       │ │                  │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ "Create Your     │ │ "Set Your        │ │ "Start Earning"  │          │
│  │  Profile"        │ │  Services"       │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ "Sign up, add    │ │ "List services   │ │ "Share your page │          │
│  │  your photo,     │ │  with pricing,   │ │  and receive     │          │
│  │  bio, and        │ │  duration, and   │ │  instant bookings │          │
│  │  specialties."   │ │  availability."  │ │  and payments."  │          │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘          │
│                                                                         │
│  Connector lines between steps (hidden on mobile):                      │
│  Dashed line, border-gray-300, absolute positioned                      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 5: Pricing Table

```
[SECTION — py-20 lg:py-24, bg-white, id="pricing"]
┌─────────────────────────────────────────────────────────────────────────┐
│  SECTION HEADER:                                                        │
│  "Simple, transparent pricing"            font-serif, text-4xl          │
│  "Start free, upgrade as you grow."       text-lg, text-gray-600       │
│                                                                         │
│  [Toggle: Monthly / Annual (save 20%)]    flex, bg-gray-100, rounded-full│
│                                           p-1, text-sm                  │
│                                                                         │
│  [3-COLUMN PRICING CARDS — grid cols-1 lg:cols-3, gap-8, mt-12]       │
│                                                                         │
│  ┌─────────────────┐ ┌══════════════════┐ ┌─────────────────┐          │
│  │ STARTER          │ ║ PRO ⭐ Popular   ║ │ ENTERPRISE       │          │
│  │                  │ ║                  ║ │                  │          │
│  │ "$29"            │ ║ "$99"            ║ │ "$299"           │          │
│  │ text-4xl bold    │ ║ text-4xl bold    ║ │ text-4xl bold    │          │
│  │ "/month"         │ ║ "/month"         ║ │ "/month"         │          │
│  │ text-gray-500    │ ║                  ║ │                  │          │
│  │                  │ ║                  ║ │                  │          │
│  │ "For getting     │ ║ "For growing     ║ │ "For established │          │
│  │  started"        │ ║  professionals"  ║ │  businesses"     │          │
│  │ text-gray-600    │ ║                  ║ │                  │          │
│  │                  │ ║                  ║ │                  │          │
│  │ ✓ 3 services     │ ║ ✓ Unlimited      ║ │ ✓ Everything in  │          │
│  │ ✓ Basic analytics│ ║   services       ║ │   Pro            │          │
│  │ ✓ 3% txn fee     │ ║ ✓ Advanced       ║ │ ✓ White-label    │          │
│  │ ✓ Email support  │ ║   analytics      ║ │ ✓ API access     │          │
│  │                  │ ║ ✓ 2% txn fee     ║ │ ✓ 1% txn fee     │          │
│  │                  │ ║ ✓ Priority       ║ │ ✓ Dedicated      │          │
│  │                  │ ║   support        ║ │   support        │          │
│  │                  │ ║                  ║ │ ✓ Custom domain   │          │
│  │                  │ ║                  ║ │                  │          │
│  │ [Get Started]    │ ║ [Get Started]    ║ │ [Contact Sales]  │          │
│  │  outline btn     │ ║  primary btn     ║ │  outline btn     │          │
│  │  w-full          │ ║  w-full          ║ │  w-full          │          │
│  │  min-h-touch     │ ║  min-h-touch     ║ │  min-h-touch     │          │
│  └─────────────────┘ └══════════════════┘ └─────────────────┘          │
│                                                                         │
│  PRO card: ring-2 ring-primary-600, relative                            │
│  "Most Popular" badge: absolute -top-4, bg-primary-600, text-white     │
│  All cards: rounded-2xl, p-8, bg-white, border border-gray-200         │
│  ✓ marks: text-success-600                                              │
│                                                                         │
│  Mobile: single column, PRO card first (order-first)                    │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 6: Testimonials

```
[SECTION — py-20, bg-gray-50, id="testimonials"]
┌─────────────────────────────────────────────────────────────────────────┐
│  SECTION HEADER:                                                        │
│  "Loved by professionals worldwide"       font-serif, text-4xl          │
│                                                                         │
│  [TESTIMONIAL CARDS — grid cols-1 md:cols-2 lg:cols-3, gap-8, mt-16]  │
│                                                                         │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐          │
│  │ ★★★★★           │ │ ★★★★★           │ │ ★★★★★           │          │
│  │ text-warning-400 │ │                  │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ "I replaced 4    │ │ "My booking rate │ │ "Setting up took │          │
│  │  different tools │ │  went up 300%    │ │  10 minutes. Got │          │
│  │  with this one   │ │  in the first    │ │  my first client │          │
│  │  platform..."    │ │  month."         │ │  the same week." │          │
│  │ text-gray-700    │ │                  │ │                  │          │
│  │ italic           │ │                  │ │                  │          │
│  │                  │ │                  │ │                  │          │
│  │ ┌──┐            │ │ ┌──┐            │ │ ┌──┐            │          │
│  │ │🧑│ Sofia M.    │ │ │🧑│ Marcus T.   │ │ │🧑│ Elena R.    │          │
│  │ └──┘             │ │ └──┘             │ │ └──┘             │          │
│  │ w-10 h-10        │ │                  │ │                  │          │
│  │ rounded-full     │ │                  │ │                  │          │
│  │ "Marketing       │ │ "Fitness Coach"  │ │ "UX Designer"   │          │
│  │  Consultant"     │ │                  │ │                  │          │
│  │ text-sm gray-500 │ │                  │ │                  │          │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘          │
│                                                                         │
│  Cards: bg-white, rounded-xl, p-6, shadow-sm                            │
│  Mobile: horizontal scroll or stacked                                   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Section 7: Final CTA

```
[SECTION — py-24, bg-primary-900, text-white]
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│  text-center, max-w-3xl mx-auto                                         │
│                                                                         │
│  "Ready to grow your business?"           font-serif, text-4xl          │
│                                           font-bold, text-white         │
│                                                                         │
│  "Join thousands of professionals who     text-xl, text-primary-200     │
│   have simplified their workflow and      leading-relaxed               │
│   increased their revenue."                                             │
│                                                                         │
│  [Start Your Free Trial →]                bg-white, text-primary-900    │
│                                           btn-lg, font-semibold         │
│                                           hover:bg-primary-50           │
│                                           min-h-touch                   │
│                                                                         │
│  "14-day free trial. No credit card."     text-sm, text-primary-300    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Footer

```
[FOOTER — py-16, bg-gray-900, text-gray-400]
┌─────────────────────────────────────────────────────────────────────────┐
│  [4-COLUMN LAYOUT — grid cols-2 md:cols-4, gap-8]                      │
│                                                                         │
│  COLUMN 1: Brand                                                        │
│  [Logo] text-white                                                      │
│  "The all-in-one platform                                               │
│   for service professionals."                                           │
│  [Twitter] [LinkedIn] [GitHub]     icon links, min-touch, gap-4         │
│                                                                         │
│  COLUMN 2: Product                                                      │
│  Features, Pricing, Integrations,                                       │
│  Changelog                                                              │
│                                                                         │
│  COLUMN 3: Company                                                      │
│  About, Blog, Careers, Contact                                          │
│                                                                         │
│  COLUMN 4: Legal                                                        │
│  Privacy, Terms, Cookie Policy                                          │
│                                                                         │
│  ─────────────────────────────────────────────                          │
│  "2026 Professional Services SaaS. All rights reserved."               │
│  text-sm, text-gray-500, pt-8, border-t border-gray-800                │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Responsive Behavior Summary

| Breakpoint | Behavior |
|------------|----------|
| < 640px (mobile) | Single column, hamburger nav, stacked cards, full-width buttons |
| 640-767px (sm) | 2-column grids where applicable |
| 768-1023px (md) | 2-column features, side-by-side hero |
| 1024-1279px (lg) | 3-column features/pricing, full nav |
| >= 1280px (xl) | Max-width container (1280px), generous spacing |

## Accessibility Notes

- All interactive elements have visible focus indicators (ring-2 ring-primary-500)
- Color contrast ratios meet WCAG AA (text on all backgrounds >= 4.5:1)
- Skip-to-content link as first focusable element
- Pricing table uses semantic headings (h3 for plan names)
- All images have descriptive alt text
- CTA buttons have clear, action-oriented labels
- Mobile nav: focus trap when open, Escape to close
