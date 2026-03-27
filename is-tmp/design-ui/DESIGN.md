# Design System Document: The Precision Curator

## 1. Overview & Creative North Star
The core objective of this design system is to transform raw data management from a utilitarian chore into a "High-End Editorial" experience. We move beyond the standard "SaaS Dashboard" look by adopting the **Creative North Star: The Digital Curator.**

This system treats data points like artifacts in a gallery. It breaks the "template" look through **intentional asymmetry**, where heavy information clusters are balanced by expansive, purposeful white space. We reject the rigid, boxed-in grids of traditional mobile apps in favor of a fluid, layered interface that feels as much like a premium financial broadsheet as it does a data tool.

**The Signature Feel:**
- **Authoritative:** High-contrast typography scales.
- **Atmospheric:** Deep, tonal blues instead of flat neutrals.
- **Weightless:** Using light and glass rather than lines and borders.

---

## 2. Colors & Surface Philosophy
The palette is rooted in a "Trustworthy Deep Blue" (#1A73E8), but its application is sophisticated, utilizing Material Design 3 tonal transitions to create depth.

### The "No-Line" Rule
**Explicit Instruction:** Designers are prohibited from using 1px solid borders to define sections. Boundaries must be defined solely through:
1. **Background Color Shifts:** A `surface-container-low` section sitting on a `surface` background.
2. **Tonal Transitions:** Moving from `surface-dim` to `surface-bright` to guide the eye.

### Surface Hierarchy & Nesting
Treat the mobile screen as a series of physical layers—like stacked sheets of frosted glass. 
- **Base Layer:** `surface` (#f9f9ff).
- **Secondary Content Blocks:** `surface-container-low` (#f2f3fd).
- **High-Priority Data Cards:** `surface-container-lowest` (#ffffff) to create a "pop" against the background.
- **Interactive Modals:** `surface-container-highest` (#e0e2ec).

### The "Glass & Gradient" Rule
To avoid a "flat" feel, use **Glassmorphism** for floating action buttons or sticky headers. Use `surface_variant` at 70% opacity with a `20px` backdrop blur. For primary CTAs, apply a subtle linear gradient from `primary` (#005bbf) to `primary_container` (#1a73e8) at a 135-degree angle to provide visual "soul."

---

## 3. Typography
We utilize a dual-font system to balance editorial authority with functional clarity.

*   **Display & Headlines (Manrope):** Used for data summaries and page titles. The wide aperture of Manrope feels modern and expansive.
    *   *Headline-LG:* 2rem — For primary screen headers.
*   **Body & Labels (Inter):** Used for the data itself. Inter’s tall x-height ensures legibility even at the `label-sm` (0.6875rem) level for dense data tables.
    *   *Title-MD:* 1.125rem — For card titles, creating a clear entry point.
    *   *Body-MD:* 0.875rem — The standard for data entries and descriptions.

**The Editorial Scale:** Do not be afraid of the contrast. A `display-sm` metric sitting next to a `label-sm` caption creates a hierarchy that feels intentional and premium.

---

## 4. Elevation & Depth
In this system, elevation is a property of light, not lines.

### The Layering Principle
Depth is achieved by "stacking" the surface-container tiers. Place a `surface-container-lowest` card on a `surface-container-low` background. This creates a soft, natural lift without the "dirty" look of heavy shadows.

### Ambient Shadows
When a "floating" effect is required (e.g., a Bottom Sheet), shadows must be:
- **Color:** A tinted version of `on-surface` (using #191c23 at 6% opacity).
- **Blur:** Extra-diffused (Blur: 24px, Y-Offset: 8px).
- **The "Ghost Border" Fallback:** If a border is required for accessibility, use `outline-variant` (#c1c6d6) at **15% opacity**. Never use 100% opaque borders.

---

## 5. Components

### Buttons
- **Primary:** Gradient-fill (`primary` to `primary_container`), `xl` (1.5rem) roundedness. No border.
- **Secondary:** `surface-container-high` fill with `on-surface` text.
- **Tertiary:** Ghost style. No background; text-only using `primary` color.

### Data Cards & Lists
- **Rule:** Forbid divider lines.
- **Structure:** Use `spacing-4` (1.4rem) between list items. Separate logical groups by switching the background from `surface` to `surface-container-low`.
- **Roundedness:** All cards must use `DEFAULT` (0.5rem/8px) as per the brand requirement, but use `xl` (1.5rem) for outer containers to create a "nested" aesthetic.

### Input Fields
- **Styling:** Use `surface-container-highest` for the input track. 
- **Active State:** Instead of a thick border, use a `primary` "glow" (a 2px outer shadow using the primary color at 20% opacity).

### Status Chips
- **Success:** `tertiary-fixed` background with `on-tertiary-fixed` text.
- **Error:** `error-container` background with `on-error-container` text.
- **Shape:** Use `full` (9999px) roundedness for chips to distinguish them from square data cards.

---

## 6. Do’s and Don'ts

### Do
*   **Do** use asymmetrical margins. For example, give a header a larger top-margin (`spacing-12`) than bottom-margin (`spacing-4`) to create an editorial feel.
*   **Do** use `surface-tint` sparingly to highlight active navigation states.
*   **Do** leverage `white-space` as a functional element to group data.

### Don't
*   **Don't** use black (#000000) for text. Always use `on-surface` (#191c23) to maintain tonal harmony with the blue palette.
*   **Don't** use 1px dividers to separate list items. Use spacing or subtle background shifts.
*   **Don't** use standard "drop shadows." If it looks like a default shadow, it is too heavy.
*   **Don't** crowd the edges. Respect the `spacing-6` (2rem) minimum gutter for mobile screens.