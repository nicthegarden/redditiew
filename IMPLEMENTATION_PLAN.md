# RedditView - Mailbox Mode Implementation Plan

## Goal
Create a proper mailbox-style Reddit viewer where:
- Left pane: Thread list (posts)
- Right pane: Thread content (post + comments)
- Navigation: Keyboard-driven, focus-based

## Current Issues
1. ❌ Right pane not displaying content properly
2. ❌ Keyboard navigation not working as intended
3. ❌ Focus management needs improvement
4. ❌ Proxy redirect handling incomplete

## Desired UX Flow

```
1. App loads → Focus on subreddit input
   
2. Type subreddit name [e.g., "linux"] + ENTER
   → Focus shifts to filter/search box
   → Posts load in left pane
   
3. Type search query (optional) + ENTER
   → Filter posts or search Reddit
   → Focus shifts to post list
   
4. Use ↑↓ arrow keys to navigate posts
   → Selected post highlighted in left pane
   
5. Press ENTER on selected post
   → Post opens in right pane
   → Can read full content + comments
   
6. Press TAB to switch back to search
   → Return to step 3
```

## Implementation Steps

### Phase 1: Fix Right Pane Display (Priority: CRITICAL)
**Issue:** Iframe not loading content

**Solution:**
1. Remove iframe approach - use direct HTML/markdown rendering
2. Create a PostViewer component that:
   - Fetches post data from Reddit API
   - Renders post title, body, metadata
   - Shows comments in a scrollable list
   - No CORS/iframe issues

**Files to modify:**
- `src/App.tsx` - Replace iframe with React component
- Create `src/components/PostViewer.tsx` - New component

**Expected result:** Right pane shows post content immediately when clicked

---

### Phase 2: Fix Keyboard Navigation (Priority: HIGH)
**Issue:** Focus management broken, TAB cycling not working

**Solution:**
1. Implement proper focus state machine:
   ```
   'subreddit' → 'search' → 'list' → 'subreddit'
   ```

2. Update handlers:
   - Subreddit input: ENTER → move to search
   - Search input: ENTER → move to list (load posts)
   - Post list: ↑↓ navigate, ENTER to open

3. Remove complex useEffect dependencies

**Files to modify:**
- `src/App.tsx` - Simplify focus/keyboard logic

**Expected result:** TAB/Shift+TAB cycles through inputs smoothly

---

### Phase 3: Polish Right Pane (Priority: MEDIUM)
**Issue:** Content display not optimal

**Solution:**
1. Create `src/components/PostDetail.tsx`:
   - Show post title, author, score, time
   - Show post body (text posts)
   - Show link preview (link posts)
   - Show comments tree

2. Add styling:
   - Readable fonts
   - Good spacing
   - Comment indentation

**Files to modify:**
- Create `src/components/PostDetail.tsx`
- Update `src/index.css` with PostDetail styles

**Expected result:** Professional-looking post viewer

---

### Phase 4: Fix Proxy (Priority: LOW - for optimization later)
**Current:** Redirect handling added but needs testing

**Solution:**
- Test redirect handling after other fixes
- Only if iframe is still used

---

## Quick Checklist

### Must Have (MVP)
- [ ] Right pane shows post content (no iframe)
- [ ] Keyboard navigation works (subreddit → search → list)
- [ ] Click post on left = opens on right
- [ ] Arrow keys navigate posts
- [ ] ENTER opens/closes post

### Nice to Have
- [ ] Comments displayed below post
- [ ] Load more comments
- [ ] Vote buttons
- [ ] Share/save buttons

### Future
- [ ] Multi-subreddit tabs
- [ ] User preferences
- [ ] Theme toggle

---

## Implementation Order

1. **Today:** Fix right pane (remove iframe, use React)
2. **Today:** Fix keyboard navigation (simplify logic)
3. **Tomorrow:** Polish styling
4. **Later:** Add comments, voting, etc.

---

## Files to Create/Modify

**Create:**
- `src/components/PostDetail.tsx` - Displays selected post
- `src/components/PostComments.tsx` - Displays comments

**Modify:**
- `src/App.tsx` - Main component (simplify, remove iframe)
- `src/index.css` - Add PostDetail/PostComments styles
- `proxy.ts` - May need minor adjustments

**Keep:**
- `src/main.tsx` - Entry point (no changes)
- `vite.config.ts` - Build config (no changes)
- `package.json` - Dependencies (no changes)

---

## Expected Timeline

| Phase | Time | Status |
|-------|------|--------|
| Phase 1 (Right Pane) | 30 min | TODO |
| Phase 2 (Navigation) | 20 min | TODO |
| Phase 3 (Polish) | 30 min | TODO |
| Testing | 15 min | TODO |
| **Total** | **~1.5 hrs** | **IN PROGRESS** |

---

## Success Criteria

✅ Posts load in left pane when subreddit entered  
✅ Right pane shows selected post content  
✅ Keyboard navigation works perfectly  
✅ No console errors  
✅ No iframe issues  
✅ Responsive on desktop  
✅ Mailbox UX feels natural  

