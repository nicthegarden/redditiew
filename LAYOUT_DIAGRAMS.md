# RedditView Layout Diagrams and Structure

## 1. DOM TREE - CURRENT STRUCTURE

```
#root
└── <div className="app">
    ├── <div className="left-pane">
    │   ├── <form className="header">
    │   │   ├── <input className="sub-input" />
    │   │   └── <button className="search-btn" />
    │   │
    │   ├── <div className="quick-links">
    │   │   └── (multiple quick-link buttons)
    │   │
    │   ├── <div className="filter-bar">
    │   │   └── <input className="filter-input" />
    │   │
    │   └── <div className="posts-list">
    │       └── (multiple post-item elements)
    │
    ├── <div className="version-footer">  ← CURRENTLY HERE
    │   ├── <span className="version-text" />
    │   └── <button className="reader-mode-btn" />
    │
    └── <div className="right-pane">
        └── <PostDetail post={selected} />
            └── (post content, comments, etc.)
```

## 2. DESKTOP LAYOUT (> 768px) - VISUAL

### Current State
```
┌─────────────────────────────────────────────────────────────────┐
│  VIEWPORT (100vw × 100vh)                                       │
├──────────────────┬──────────────────────────────────────────────┤
│                  │                                               │
│  LEFT-PANE       │       RIGHT-PANE                              │
│  (40% width)     │       (Flex: 1 / remaining space)             │
│                  │                                               │
│ ┌──────────────┐ │ ┌────────────────────────────────────────┐   │
│ │   HEADER     │ │ │  Post Detail Header                    │   │
│ │ (search)     │ │ │  • Title                               │   │
│ └──────────────┘ │ │  • Author, Score, Comments            │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │ Post Content                            │   │
│ │ QUICK LINKS  │ │ │  • Body or Media                        │   │
│ │ r/sysadmin.. │ │ │  • Links                                │   │
│ └──────────────┘ │ │  • Actions (upvote, etc.)              │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │ Comments Section                        │   │
│ │ FILTER BAR   │ │ │  • Sorted by score                      │   │
│ │ (search)     │ │ │  • Collapsible threads                  │   │
│ └──────────────┘ │ │  (scrollable)                           │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │                                         │   │
│ │  POSTS LIST  │ │ │                                         │   │
│ │ (scrollable) │ │ │                                         │   │
│ │              │ │ │                                         │   │
│ │ • Post 1     │ │ │                                         │   │
│ │ • Post 2     │ │ │                                         │   │
│ │ • Post 3 ✓   │ │ │                                         │   │
│ │ • Post 4     │ │ │                                         │   │
│ │ • Post 5     │ │ │                                         │   │
│ │ • [Load More]│ │ │                                         │   │
│ └──────────────┘ │ │                                         │   │
│                  │ └────────────────────────────────────────┘   │
└──────────────────┴──────────────────────────────────────────────┘
└──────────────── VERSION-FOOTER ──────────────────────────────────┘
    (display: flex; align-items: center; justify-content: space-between)
    v2.2.0                                            [Reader Mode btn]
```

### If Moved to Bottom (Recommended)
```
┌─────────────────────────────────────────────────────────────────┐
│  VIEWPORT (100vw × 100vh)                                       │
├──────────────────┬──────────────────────────────────────────────┤
│                  │                                               │
│  LEFT-PANE       │       RIGHT-PANE                              │
│  (40% width)     │       (Flex: 1 / remaining space)             │
│                  │                                               │
│ ┌──────────────┐ │ ┌────────────────────────────────────────┐   │
│ │   HEADER     │ │ │  Post Detail Header                    │   │
│ │ (search)     │ │ │  • Title                               │   │
│ └──────────────┘ │ │  • Author, Score, Comments            │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │ Post Content                            │   │
│ │ QUICK LINKS  │ │ │  • Body or Media                        │   │
│ │ r/sysadmin.. │ │ │  • Links                                │   │
│ └──────────────┘ │ │  • Actions (upvote, etc.)              │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │ Comments Section                        │   │
│ │ FILTER BAR   │ │ │  • Sorted by score                      │   │
│ │ (search)     │ │ │  • Collapsible threads                  │   │
│ └──────────────┘ │ │  (scrollable)                           │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │                                         │   │
│ │  POSTS LIST  │ │ │                                         │   │
│ │ (flex: 1)    │ │ │                                         │   │
│ │ (scrollable) │ │ │                                         │   │
│ │              │ │ │                                         │   │
│ │ • Post 1     │ │ │                                         │   │
│ │ • Post 2     │ │ │                                         │   │
│ │ • Post 3 ✓   │ │ │                                         │   │
│ │ • Post 4     │ │ │                                         │   │
│ │ • Post 5     │ │ │                                         │   │
│ │ • [Load More]│ │ │                                         │   │
│ └──────────────┘ │ │                                         │   │
│                  │ │                                         │   │
│ ┌──────────────┐ │ │                                         │   │
│ │VERSION-FOOTER│ │ │                                         │   │
│ │v2.2.0 [ReBtn]│ │ │                                         │   │
│ └──────────────┘ │ └────────────────────────────────────────┘   │
└──────────────────┴──────────────────────────────────────────────┘
```

## 3. MOBILE/TABLET LAYOUT (≤ 768px) - VISUAL

### Current State with order: -1
```
┌────────────────────────────────────────┐
│  VIEWPORT (100vw × 100vh)              │
│  (flex-direction: column)               │
├────────────────────────────────────────┤
│ ┌────────────────────────────────────┐ │
│ │   VERSION-FOOTER (order: -1)       │ │  ← Moved to TOP
│ │   v2.2.0          [Reader Mode btn] │ │     using order
│ │                                     │ │
│ └────────────────────────────────────┘ │
├────────────────────────────────────────┤
│ ┌────────────────────────────────────┐ │
│ │     LEFT-PANE (50% height)         │ │
│ │  (height: 50%; border-bottom)      │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │      HEADER (search)         │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │  QUICK LINKS (scrollable)    │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │ FILTER BAR (search)          │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │ POSTS LIST (scrollable)      │   │ │
│ │ │ • Post 1                     │   │ │
│ │ │ • Post 2                     │   │ │
│ │ │ • Post 3 ✓                   │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ └────────────────────────────────────┘ │
├────────────────────────────────────────┤
│ ┌────────────────────────────────────┐ │
│ │    RIGHT-PANE (50% height)         │ │
│ │                                     │ │
│ │  Post Detail Content               │ │
│ │  (scrollable)                       │ │
│ │                                     │ │
│ └────────────────────────────────────┘ │
└────────────────────────────────────────┘
```

### If Moved to Bottom of Left-Pane
```
┌────────────────────────────────────────┐
│  VIEWPORT (100vw × 100vh)              │
│  (flex-direction: column)               │
├────────────────────────────────────────┤
│ ┌────────────────────────────────────┐ │
│ │     LEFT-PANE (50% height)         │ │
│ │  (height: 50%; border-bottom)      │ │
│ │  (display: flex; flex-direction: column)
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │      HEADER (search)         │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │  QUICK LINKS (scrollable)    │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │ FILTER BAR (search)          │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │ POSTS LIST (flex: 1)         │   │ │
│ │ │ (scrollable)                 │   │ │
│ │ │                              │   │ │
│ │ │ • Post 1                     │   │ │
│ │ │ • Post 2                     │   │ │
│ │ │ • Post 3 ✓                   │   │ │
│ │ │ • ...                        │   │ │
│ │ │                              │   │ │
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ │ ┌──────────────────────────────┐   │ │
│ │ │VERSION-FOOTER (flex-shrink:0)│   │ │  ← NOW AT BOTTOM
│ │ │v2.2.0  [Reader Mode btn]     │   │ │     naturally!
│ │ └──────────────────────────────┘   │ │
│ │                                     │ │
│ └────────────────────────────────────┘ │
├────────────────────────────────────────┤
│ ┌────────────────────────────────────┐ │
│ │    RIGHT-PANE (50% height)         │ │
│ │                                     │ │
│ │  Post Detail Content               │ │
│ │  (scrollable)                       │ │
│ │                                     │ │
│ └────────────────────────────────────┘ │
└────────────────────────────────────────┘
```

## 4. FLEXBOX PROPERTIES EXPLAINED

### `.app` Container
```
Initial State (Desktop):
  display: flex
  flex-direction: row         (default - horizontal)
  height: 100vh
  width: 100vw
  
  Children Layout:
  [LEFT (40%)] [VERSION-FOOTER] [RIGHT (flex:1)]

Mobile (≤768px):
  flex-direction: column      (switch to vertical)
  
  Children Layout:
  [VERSION-FOOTER order:-1]
  [LEFT (50% height)]
  [RIGHT (50% height)]
```

### `.left-pane` Container
```
  display: flex
  flex-direction: column      (vertical stacking)
  width: 40%                  (desktop)
  width: 100%                 (mobile)
  
  Children:
  [HEADER - fixed height]
  [QUICK-LINKS - fixed height]
  [FILTER-BAR - fixed height]
  [POSTS-LIST - flex:1 GROWS]
  [VERSION-FOOTER - flex-shrink:0 FIXED]
```

### `.right-pane` Container
```
  flex: 1                     (grows to fill space)
  display: flex
  flex-direction: column
  overflow: hidden
  
  Content:
  [POST-DETAIL - flex:1 with overflow-y:auto]
```

## 5. SCROLLING BEHAVIOR

### Current Scrolling Boundaries
```
┌─────────────────────────┐
│  LEFT-PANE              │
├─────────────────────────┤
│ HEADER (fixed)          │ overflow: hidden
├─────────────────────────┤
│ QUICK-LINKS (fixed)     │ overflow: hidden
├─────────────────────────┤
│ FILTER-BAR (fixed)      │ overflow: hidden
├─────────────────────────┤
│ POSTS-LIST (scrollable) │ ← overflow-y: auto
│                         │
│ Can scroll post items   │
│ independently           │
├─────────────────────────┤
│ VERSION-FOOTER (fixed)  │ overflow: hidden
└─────────────────────────┘

┌─────────────────────────┐
│  RIGHT-PANE             │
├─────────────────────────┤
│ POST-DETAIL (scrollable)│ ← overflow-y: auto
│                         │
│ Can scroll comments,    │
│ post content, etc.      │
│                         │
└─────────────────────────┘
```

## 6. CSS SPECIFICITY & CASCADE

### Version Footer Declarations (Lines 769-776 and 964-1027)
```css
/* First declaration (line 769-776) */
.version-footer {
  padding: 8px 16px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 12px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-primary);
}

/* Second declaration (line 964-971) - OVERRIDES ABOVE */
.version-footer {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
/* Result: Later declaration takes precedence */
```

### Media Query Override
```css
@media (max-width: 768px) {
  .version-footer {
    order: -1;              /* Reorder to top */
    border-bottom: 1px solid var(--border-color);
    border-top: none;       /* No top border on mobile */
  }
}

@media (max-width: 768px) {
  .version-footer {
    padding: 8px 12px;
    gap: 8px;
  }
}
```

## 7. COLOR THEME VARIABLES

```css
:root {
  /* Light backgrounds */
  --bg-primary: #0e1113;      /* Main background (dark) */
  --bg-secondary: #1a1d21;    /* Secondary (slightly lighter) */
  --bg-tertiary: #242729;     /* Tertiary (even lighter) */
  
  /* Text colors */
  --text-primary: #d7dadc;    /* Main text (light) */
  --text-secondary: #818384;  /* Secondary text (gray) */
  
  /* Accent colors */
  --accent-orange: #ff4500;   /* Reddit orange */
  --accent-blue: #7193ff;     /* Link blue */
  
  /* Borders & misc */
  --border-color: #343536;
  --success: #46d160;
  --danger: #ff4500;
}

body.light {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --text-primary: #1a1a1a;
  /* ... etc */
}
```

## 8. TRANSITION EFFECTS

```css
.reader-mode-btn {
  transition: background-color 0.2s;  /* Smooth color change */
}

.quick-link {
  transition: all 0.2s;               /* All properties animate */
}

.post-item {
  transition: background 0.15s;       /* Faster background change */
}

.search-btn {
  transition: opacity 0.2s;           /* Fade effect */
}
```

---

## SUMMARY TABLE

| Element | Desktop Width | Desktop Height | Mobile Width | Mobile Height | Overflow | Flex-Grow |
|---------|--------------|----------------|--------------|---------------|----------|-----------|
| .app | 100vw | 100vh | 100vw | 100vh | hidden | N/A |
| .left-pane | 40% | 100% | 100% | 50% | hidden | 0 |
| .header | 100% | auto | 100% | auto | hidden | 0 |
| .quick-links | 100% | auto | 100% | auto | x-auto | 0 |
| .filter-bar | 100% | auto | 100% | auto | hidden | 0 |
| .posts-list | 100% | flex:1 | 100% | flex:1 | y-auto | 1 |
| .version-footer | 100% | auto | 100% | auto | hidden | 0 |
| .right-pane | flex:1 | 100% | 100% | 50% | hidden | 1 |
| .post-detail | 100% | flex:1 | 100% | 100% | y-auto | 1 |

