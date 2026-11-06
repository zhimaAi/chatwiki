# Design System: ChatWiki Default

This document describes the existing ChatWiki default theme (`default`). It is extracted from project configuration rather than an external brand specification. Sources: `src/config/themes/default.js`, `src/config/theme.js`, and the existing styles listed in `src/config/themes/README.md`.

## 1. Visual Theme & Atmosphere

The default theme uses white content surfaces, a light gray workspace, and blue interactive accents. Typography follows the project's existing Chinese sans-serif stack. Text tones, surface separation, and restrained popover shadows establish hierarchy.

**Key Characteristics:**

- Application Blue (`#2475fc`) for primary actions.
- White (`#ffffff`) business containers on a light gray (`#f5f6f8`) workspace.
- Three business text levels: `#262626`, `#595959`, and `#8c8c8c`.
- Existing sidebar radii of 10px for items and 8px for child items.
- A 12px menu popover radius with a soft black shadow.
- Ant Design Vue defaults and derived tokens for component interaction states.

## 2. Color Palette & Roles

### Primary

- **Application Blue** (`#2475fc`): Primary actions and the configured `colorPrimary`.
- **Business Links** (`#2475fc`): Existing ordinary link, hover, and active colors preserved by `legacyBusiness`. Ant component link states follow Provider tokens.

### Surface & Background

- **Page White** (`#ffffff`): Existing business container background.
- **Workspace Gray** (`#f5f6f8`): Customer service workspace and application layout variable.
- **Statistics Gray** (`#f2f4f7`): Extracted statistics card background; not a global container override.

### Neutrals & Text

- **Primary Text** (`#262626`): Main business text.
- **Secondary Text** (`#595959`): Descriptions and supporting text.
- **Muted Text** (`#8c8c8c`): Metadata and de-emphasized text.

### Borders & Semantic Colors

- **Header Border** (`#f2f4f7`): Existing header separator reference; not a global border Token override.
- **Delete Text** (`#fb363f`): Existing business deletion color reference; not a global error Token override.
- General borders, success, warning, error, disabled, hover, and selected states use final Provider tokens through the existing application variables.

### Chart Series

| Role | Color |
| --- | --- |
| Total messages | `#0079fe` |
| Hit messages | `#00c292` |
| Missed messages | `#ff6b6b` |
| Hit rate | `#333333` |
| Data point border | `#ffffff` |

Preserve each series' meaning rather than replacing every series with the application primary color.

## 3. Typography Rules

### Font Family

- **Body / UI**: `PingFang SC, Noto Sans SC, sans-serif`.
- No separate branded headline or code font is configured by the default theme.

### Hierarchy

| Role | Font | Size | Weight | Line Height | Letter Spacing | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| Existing body baseline | UI font stack | 15px | Existing stylesheet | 1.6 | Not specified | From `src/assets/base.css` |
| Sidebar item | UI font stack | 14px | Existing stylesheet | Existing component | Not specified | Extracted sidebar baseline |
| Selected sidebar item | UI font stack | 14px | 600 | Existing component | Not specified | Selected-state emphasis |
| Ant Design Vue controls | UI font stack | Provider token | Provider token | Provider token | Existing component | Not globally overridden by the body baseline |

### Principles

- Preserve the existing Chinese font fallback chain.
- Treat body and sidebar values as scoped references rather than universal typography overrides.
- Retain existing heading, code, and component typography during theme adaptation.

## 4. Component Stylings

### Buttons

- Ant Design Vue buttons consume the existing dynamic `ConfigProvider.theme`.
- Primary color: Application Blue (`#2475fc`).
- Hover, pressed, disabled, borders, and contrast text come from final Provider tokens.
- The default theme does not explicitly override the generic control radius.

### Cards & Containers

- Standard business container: Page White (`#ffffff`).
- Existing statistics card: Statistics Gray (`#f2f4f7`), 6px radius.
- Statistics-specific values do not redefine all application cards.

### Inputs & Forms

- Follow Ant Design Vue component tokens under the current Provider.
- No default-theme-specific Input radius or focus-color override is configured.

### Navigation

- Sidebar item radius: 10px; child item radius: 8px.
- Selected sidebar item weight: 600.
- Background, text, borders, hover, and selected states consume shared `--app-*` variables.
- Selected emphasis follows the application primary color.

### Popovers

- Existing custom menu popover: 12px radius and `0 8px 24px rgba(0, 0, 0, 0.12)` shadow.
- Custom popovers consume application popover variables; Ant popovers follow their component tokens.

### Icons

- Monochrome icons inherit semantic text or primary colors through `currentColor`.
- Preserve multicolor artwork, brand assets, and business status meanings.

## 5. Layout Principles

### Spacing System

- No unified spacing scale is declared by the extracted default theme.
- Preserve existing component padding, gaps, dimensions, and alignment.

### Grid & Container

- Keep existing application layouts, sidebar widths, and content grids.
- Theme switching updates appearance without remounting route components.

### Border Radius Scale

| Role | Radius |
| --- | --- |
| Statistics card | 6px |
| Sidebar child item | 8px |
| Sidebar item | 10px |
| Menu popover | 12px |
| Generic control | Final Provider `borderRadius` token |

Component-specific values are extracted references, not a universal radius scale for every component.

## 6. Depth & Elevation

| Treatment | Value | Use |
| --- | --- | --- |
| Workspace surface | `#f5f6f8` | Business layout separation |
| Content surface | `#ffffff` | Business containers |
| Menu popover shadow | `0 8px 24px rgba(0, 0, 0, 0.12)` | Existing custom menu popover |
| Other component elevation | Final Provider tokens | Ant Design Vue components |

No additional global shadow hierarchy is defined. Do not apply the menu shadow to every control or card.

## 7. Do's and Don'ts

### Do

- Use existing `--app-*` semantic variables for custom component styles.
- Let Ant Design Vue components consume the current Provider theme.
- Preserve chart series meanings and business status colors.
- Keep layout, interaction, permissions, and business logic unchanged.
- Add genuinely missing shared variables at the existing token-to-variable output point.

### Don't

- Do not copy palette hex values into components that need theme switching.
- Do not introduce `--customer-*` or other module-specific theme variable layers.
- Do not apply statistics card styling to every container.
- Do not treat header or deletion reference colors as global border or error tokens.
- Do not create fixed interaction palettes that bypass Provider-derived states.

## 8. Responsive Behavior

### Breakpoints

- The default theme does not define application-wide breakpoints.
- Preserve existing page and component media queries and responsive behavior.

### Collapsing Strategy

- Retain current sidebar, navigation, and content resizing behavior.
- Theme adaptation does not add mobile layouts or change component dimensions.

### Color Scheme

- No additional dark theme is configured by the three application theme definitions.
- Existing system dark-mode styles remain in place; this document does not redefine them.

## 9. Agent Prompt Guide

### Quick Color Reference

- Primary: Application Blue (`#2475fc`) → `--app-primary`.
- Business layout: Workspace Gray (`#f5f6f8`) → `--app-bg-layout`.
- Business container: Page White (`#ffffff`) → `--app-bg-container`.
- Primary text: `#262626` → `--app-text-primary`.
- Secondary text: `#595959` → `--app-text-secondary`.
- Tertiary text: `#8c8c8c` → `--app-text-tertiary`.
- Selected background: final `colorPrimaryBg` → `--app-selected`.
- Hover background: final `colorFillTertiary` → `--app-bg-hover`.
- Generic radius: final `borderRadius` → `--app-radius-control`.

### Example Component Prompt

“Adapt this component to the existing application theme. Use `--app-bg-container` for its content surface, semantic text variables for its text hierarchy, and `--app-selected` with `--app-primary` for its selected state. Preserve layout and business behavior. Follow the reusable adaptation prompt in `design-system/README.md` for the full requirements.”

### Implementation Notes

- `antTheme` references `src/config/theme.js`, which explicitly configures only the primary color and UI font.
- `legacyBusiness` preserves existing business backgrounds, text, and ordinary links separately from Ant tokens. It is not a storage-key migration mechanism.
- Runtime variables are emitted by `src/components/global-config-provider/theme-variables.vue` from final Provider tokens and explicit default business values.
- Detailed application mapping and known limitations are documented in `src/config/themes/README.md`.
