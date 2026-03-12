import { useState, useCallback, useRef, useEffect } from 'react'
import PostDetail from './components/PostDetail'
import { loadConfig, getConfig } from './config'
import { APP_VERSION } from './version'

const PROXY = '/api'
const API_BASE = '/api'

// Touch gesture handler
interface TouchStartPos {
  x: number
  y: number
  time: number
}

function useTouchGestures(onSwipeLeft: () => void, onSwipeRight: () => void, onSwipeDown: () => void, onSwipeUp: () => void) {
  const touchStartRef = useRef<TouchStartPos | null>(null)
  
  const handleTouchStart = (e: React.TouchEvent) => {
    if (e.touches.length === 1) {
      const touch = e.touches[0]
      touchStartRef.current = {
        x: touch.clientX,
        y: touch.clientY,
        time: Date.now()
      }
    }
  }

  const handleTouchEnd = (e: React.TouchEvent) => {
    if (!touchStartRef.current) return
    
    const touch = e.changedTouches[0]
    const dx = touch.clientX - touchStartRef.current.x
    const dy = touch.clientY - touchStartRef.current.y
    const dt = Date.now() - touchStartRef.current.time
    
    // Minimum swipe distance and maximum time
    const minDistance = 50
    const maxTime = 500
    
    if (dt > maxTime) return
    
    const absDx = Math.abs(dx)
    const absDy = Math.abs(dy)
    
    // Swipe right (navigate to previous post)
    if (dx > minDistance && absDx > absDy) {
      onSwipeRight()
    }
    // Swipe left (navigate to next post)
    else if (dx < -minDistance && absDx > absDy) {
      onSwipeLeft()
    }
    // Swipe down (scroll up)
    else if (dy > minDistance && absDy > absDx) {
      onSwipeDown()
    }
    // Swipe up (scroll down)
    else if (dy < -minDistance && absDy > absDx) {
      onSwipeUp()
    }
    
    touchStartRef.current = null
  }

  return { handleTouchStart, handleTouchEnd }
}

const SUBREDDIT_SUGGESTIONS = [
  // Core Sysadmin (default start with sysadmin)
  'sysadmin',
  // Virtualization & Infrastructure
  'proxmox', 'selfhosted', 'homelab', 'unraid',
  // Linux Distributions
  'linux', 'archlinux', 'debian', 'ubuntu',
  // Windows & IT
  'windows', 'IT',
  // DevOps & Containers
  'devops', 'docker', 'kubernetes', 'ceph',
  // Networking & Security
  'networking', 'netsec',
  // Cloud & Business
  'msp', 'aws',
  // General Tech
  'technology',
  // Development
  'programming', 'javascript', 'python',
]

interface RedditPost {
  kind: string
  data: {
    id: string
    title: string
    subreddit: string
    author: string
    created_utc: number
    score: number
    num_comments: number
    thumbnail?: string
    preview?: {
      images: Array<{
        source: {
          url: string
        }
      }>
    }
    permalink: string
  }
}

interface CacheEntry {
  posts: RedditPost[]
  time: number
}

function formatTime(ts: number): string {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

function formatNum(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toString()
}

function getThumb(post: RedditPost): string | null {
  const p = post.data
  if (p.thumbnail && p.thumbnail.startsWith('http')) return p.thumbnail
  if (p.preview?.images?.[0]?.source?.url) {
    return p.preview.images[0].source.url.replace(/&amp;/g, '&')
  }
  return null
}

interface PostItemProps {
  post: RedditPost
  active: boolean
  onClick: () => void
}

function PostItem({ post, active, onClick }: PostItemProps) {
   const thumb = getThumb(post)
   const p = post.data
   
   return (
     <div className={`post-item ${active ? 'active' : ''}`} onClick={onClick}>
       {thumb && <span className="post-image-indicator">📷</span>}
       <div className="post-info">
        <div className="post-title">{p.title}</div>
        <div className="post-meta">
          <span className="sub">r/{p.subreddit}</span>
          <span>{formatTime(p.created_utc)}</span>
        </div>
        <div className="post-stats">
          <span>▲ {formatNum(p.score)}</span>
          <span>💬 {formatNum(p.num_comments)}</span>
        </div>
      </div>
    </div>
  )
}

function usePostCache() {
  const [cached, setCached] = useState<Record<string, CacheEntry>>(() => {
    const saved = localStorage.getItem('postCache')
    return saved ? JSON.parse(saved) : {}
  })

  const saveToCache = useCallback((sub: string, posts: RedditPost[]) => {
    setCached(prev => {
      const updated = { ...prev, [sub]: { posts, time: Date.now() } }
      localStorage.setItem('postCache', JSON.stringify(updated))
      return updated
    })
  }, [])

  return { cached, saveToCache }
}

export default function App() {
  const [configLoaded, setConfigLoaded] = useState(false)
  const [defaultSubreddit, setDefaultSubreddit] = useState('sysadmin')
  const [sub, setSub] = useState('sysadmin')
  const [input, setInput] = useState('sysadmin')
  const [search, setSearch] = useState('')
  const [posts, setPosts] = useState<RedditPost[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selected, setSelected] = useState<RedditPost | null>(null)
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [after, setAfter] = useState<string | null>(null)
  const [focused, setFocused] = useState<'search' | 'filter' | 'list'>('search')
  const [suggestions, setSuggestions] = useState<RedditPost[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [suggestionIndex, setSuggestionIndex] = useState(-1)
  const [isRedditSearch, setIsRedditSearch] = useState(false)
  const [theme, setTheme] = useState<'dark' | 'light'>(() => {
    const saved = localStorage.getItem('theme')
    return (saved as 'dark' | 'light') || 'dark'
  })
  const [commentScrollPos, setCommentScrollPos] = useState(0)
  const [commentIsAtBottom, setCommentIsAtBottom] = useState(false)
  const [readerMode, setReaderMode] = useState(false)
  const [twoFingerStartX, setTwoFingerStartX] = useState(0)
  const [fingerCount, setFingerCount] = useState(0)
  const [sort, setSort] = useState<'hot' | 'new' | 'top'>('hot')
  const [sortingButtonsEnabled, setSortingButtonsEnabled] = useState(true)
  
  const listRef = useRef<HTMLDivElement>(null)
  const searchRef = useRef<HTMLInputElement>(null)
  const filterRef = useRef<HTMLInputElement>(null)
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const rightPaneRef = useRef<HTMLDivElement>(null)
  const { cached, saveToCache } = usePostCache()

  // Touch gesture handlers
  const handleSwipeLeft = useCallback(() => {
    // Swipe left = next post
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts
    const nextIdx = Math.min(selectedIndex + 1, filteredPosts.length - 1)
    if (nextIdx !== selectedIndex) {
      setSelectedIndex(nextIdx)
      setSelected(filteredPosts[nextIdx])
    }
  }, [selectedIndex, posts, search])

  const handleSwipeRight = useCallback(() => {
    // Swipe right = previous post
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts
    const prevIdx = Math.max(selectedIndex - 1, 0)
    if (prevIdx !== selectedIndex) {
      setSelectedIndex(prevIdx)
      setSelected(filteredPosts[prevIdx])
    }
  }, [selectedIndex, posts, search])

  const handleSwipeDown = useCallback(() => {
    // Swipe down = scroll up in right pane
    const rightPane = rightPaneRef.current
    if (rightPane) {
      const postDetail = rightPane.querySelector('.post-detail')
      if (postDetail) {
        postDetail.scrollTop -= 150
      }
    }
  }, [])

  const handleSwipeUp = useCallback(() => {
    // Swipe up = scroll down in right pane
    const rightPane = rightPaneRef.current
    if (rightPane) {
      const postDetail = rightPane.querySelector('.post-detail')
      if (postDetail) {
        postDetail.scrollTop += 150
      }
    }
  }, [])

  // Two-finger swipe to change subreddit
  const handleTwoFingerTouchStart = useCallback((e: React.TouchEvent) => {
    if (e.touches.length === 2) {
      setFingerCount(2)
      setTwoFingerStartX(e.touches[0].clientX)
    }
  }, [])

  const handleTwoFingerTouchEnd = useCallback((e: React.TouchEvent) => {
    if (fingerCount === 2 && e.changedTouches.length > 0) {
      const endX = e.changedTouches[0].clientX
      const distance = Math.abs(endX - twoFingerStartX)
      const threshold = 100

      if (distance > threshold) {
        const currentIndex = SUBREDDIT_SUGGESTIONS.indexOf(sub)
        let nextIndex = currentIndex

        if (endX > twoFingerStartX) {
          // Swiped right → previous subreddit
          nextIndex = currentIndex === 0 ? SUBREDDIT_SUGGESTIONS.length - 1 : currentIndex - 1
        } else {
          // Swiped left → next subreddit
          nextIndex = (currentIndex + 1) % SUBREDDIT_SUGGESTIONS.length
        }

        handleSub(SUBREDDIT_SUGGESTIONS[nextIndex])
      }

      setFingerCount(0)
    }
  }, [fingerCount, twoFingerStartX, sub])

  const { handleTouchStart, handleTouchEnd } = useTouchGestures(
    handleSwipeLeft,
    handleSwipeRight,
    handleSwipeDown,
    handleSwipeUp
  )

  // Load config on mount
  useEffect(() => {
    const initConfig = async () => {
      try {
        await loadConfig()
        const config = getConfig()
        const defaultSub = config.web.default_subreddit
        setDefaultSubreddit(defaultSub)
        setSub(defaultSub)
        setInput(defaultSub)
        setSortingButtonsEnabled(config.web.sortingButtonsEnabled ?? true)
      } catch (err) {
        console.error('Failed to load config:', err)
        // Use hardcoded defaults if config fails
        setDefaultSubreddit('sysadmin')
        setSub('sysadmin')
        setInput('sysadmin')
        setSortingButtonsEnabled(true)
      }
      setConfigLoaded(true)
    }
    initConfig()
  }, [])

  const fetchPosts = useCallback(async (subreddit: string, cursor: string | null = null, sortType: 'hot' | 'new' | 'top' = 'hot') => {
    const isMore = !!cursor
    const cachedData = cached[subreddit]
    
    if (!isMore && cachedData && Date.now() - cachedData.time < 3600000) {
      setPosts(cachedData.posts)
      return
    }

    setLoading(true)
    setError(null)
    
    try {
      const limit = 50
      const target = `${API_BASE}/r/${subreddit}.json?limit=${limit}&sort=${sortType}${cursor ? '&after=' + cursor : ''}`
      const res = await fetch(target)
      
      // Handle rate limiting
      if (res.status === 429) {
        const data = await res.json().catch(() => ({}))
        const retryAfter = data.retry_after || 60
        setError(`Rate limited. Please try again in ${retryAfter} seconds.`)
        return
      }

      if (!res.ok) {
        if (res.status === 404) {
          throw new Error(`Subreddit r/${subreddit} not found`)
        }
        if (res.status >= 500) {
          throw new Error('Reddit server error. Please try again later.')
        }
        throw new Error(`Failed to load (${res.status})`)
      }

      const data = await res.json()
      const items = data.data.children as RedditPost[]
      
      if (isMore) {
        setPosts(prev => [...prev, ...items])
      } else {
        setPosts(items)
        saveToCache(subreddit, items)
      }
      setAfter(data.data.after)
      setSub(subreddit)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      console.error('Fetch error:', message)
      setError(message)
    } finally {
      setLoading(false)
    }
  }, [cached, saveToCache])

  // Auto-select first post when posts load (fixes empty reader mode after subreddit change)
  useEffect(() => {
    if (posts.length > 0 && selected === null) {
      setSelected(posts[0])
      setSelectedIndex(0)
    }
  }, [posts, selected])

  const handleSub = (subreddit: string) => {
    setInput(subreddit)
    setSelected(null)
    setSelectedIndex(0)
    fetchPosts(subreddit, null, sort)
  }

  const handleSortChange = (newSort: 'hot' | 'new' | 'top') => {
    setSort(newSort)
    setSelected(null)
    setSelectedIndex(0)
    setAfter(null)
    setPosts([])
    fetchPosts(sub, null, newSort)
  }

  const handleNextPost = () => {
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts
    const nextIdx = Math.min(selectedIndex + 1, filteredPosts.length - 1)
    if (nextIdx !== selectedIndex) {
      setSelectedIndex(nextIdx)
      setSelected(filteredPosts[nextIdx])
    }
  }

  const handlePreviousPost = () => {
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts
    const prevIdx = Math.max(selectedIndex - 1, 0)
    if (prevIdx !== selectedIndex) {
      setSelectedIndex(prevIdx)
      setSelected(filteredPosts[prevIdx])
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const clean = input.replace(/^\/?r\/?/, '').trim()
    if (clean) handleSub(clean)
  }

  const handleSearchSubmit = () => {
    if (search.trim()) {
      // Search scoped to current subreddit, or all Reddit if search is empty
      setIsRedditSearch(true)
      setLoading(true)
      setError(null)
      
      const query = encodeURIComponent(search.trim())
      // Search within the current subreddit
      const searchUrl = `${API_BASE}/r/${sub}/search.json?q=${query}&type=link&limit=50`
      
      fetch(searchUrl)
        .then(res => {
          if (!res.ok) {
            throw new Error(`Search failed (HTTP ${res.status}): ${res.statusText}`)
          }
          return res.json()
        })
        .then(data => {
          const items = data.data.children as RedditPost[]
          setPosts(items)
          setSub(`${sub} search: "${search}"`)
          setLoading(false)
        })
        .catch(err => {
          setError(err instanceof Error ? err.message : 'Search failed')
          setLoading(false)
        })
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts

    if (focused === 'list') {
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        setSelectedIndex(prev => Math.min(prev + 1, filteredPosts.length - 1))
        setSelected(filteredPosts[Math.min(selectedIndex + 1, filteredPosts.length - 1)])
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        setSelectedIndex(prev => Math.max(prev - 1, 0))
        setSelected(filteredPosts[Math.max(selectedIndex - 1, 0)])
      } else if (e.key === 'Enter') {
        e.preventDefault()
        if (filteredPosts[selectedIndex]) {
          setSelected(filteredPosts[selectedIndex])
        }
      }
    }

    if (focused === 'filter') {
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        setSelectedIndex(prev => Math.min(prev + 1, filteredPosts.length - 1))
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        setSelectedIndex(prev => Math.max(prev - 1, 0))
      } else if (e.key === 'Enter') {
        e.preventDefault()
        if (search.trim()) {
          handleSearchSubmit()
        } else {
          setFocused('list')
          listRef.current?.focus()
        }
      }
    }

    if (e.ctrlKey && e.key === 'f') {
      e.preventDefault()
      iframeRef.current?.focus()
    }
  }

  useEffect(() => {
    fetchPosts(sub)
  }, [])

  useEffect(() => {
    const handleTab = (e: KeyboardEvent) => {
      if (e.key === 'Tab') {
        e.preventDefault()
        e.stopPropagation()
        
        if (e.shiftKey) {
          if (focused === 'list') {
            setFocused('filter')
            filterRef.current?.focus()
          } else {
            setFocused('search')
            searchRef.current?.focus()
          }
        } else {
          if (focused === 'search') {
            setFocused('filter')
            filterRef.current?.focus()
          } else if (focused === 'filter') {
            setFocused('list')
            listRef.current?.focus()
          } else {
            setFocused('list')
            listRef.current?.focus()
          }
        }
      }
    }
    document.addEventListener('keydown', handleTab, true)
    return () => document.removeEventListener('keydown', handleTab, true)
  }, [focused])

  // Handle right pane navigation and scrolling
  useEffect(() => {
    const handleGlobalKeys = (e: KeyboardEvent) => {
      // Only handle navigation when a post is selected (right pane is active)
      if (!selected) return

      const rightPane = rightPaneRef.current
      if (!rightPane) return

      // Get the comments list element
      const commentsList = rightPane.querySelector('.comments-list')
      const postDetail = rightPane.querySelector('.post-detail')

      if (e.key === 'PageDown') {
        e.preventDefault()
        e.stopPropagation()
        
        // Scroll right pane down
        if (postDetail) {
          postDetail.scrollTop += 400
          
          // Check if at bottom
          const isAtBottom = postDetail.scrollHeight - postDetail.scrollTop <= postDetail.clientHeight + 50
          if (isAtBottom && commentsList) {
            // Move to next post
            const filteredPosts = search.trim()
              ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
              : posts
            const nextIdx = Math.min(selectedIndex + 1, filteredPosts.length - 1)
            if (nextIdx !== selectedIndex) {
              setSelectedIndex(nextIdx)
              setSelected(filteredPosts[nextIdx])
            }
          }
        }
      } else if (e.key === 'PageUp') {
        e.preventDefault()
        e.stopPropagation()
        if (postDetail) {
          postDetail.scrollTop -= 400
        }
      } else if (e.key === 'ArrowRight') {
        e.preventDefault()
        e.stopPropagation()
        if (postDetail) {
          postDetail.scrollTop += 150
          
          // Check if at bottom and advance to next post
          const isAtBottom = postDetail.scrollHeight - postDetail.scrollTop <= postDetail.clientHeight + 50
          if (isAtBottom && commentsList) {
            const filteredPosts = search.trim()
              ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
              : posts
            const nextIdx = Math.min(selectedIndex + 1, filteredPosts.length - 1)
            if (nextIdx !== selectedIndex) {
              setSelectedIndex(nextIdx)
              setSelected(filteredPosts[nextIdx])
              // Reset scroll for new post
              setTimeout(() => {
                postDetail.scrollTop = 0
              }, 0)
            }
          }
        }
      } else if (e.key === 'ArrowLeft') {
        e.preventDefault()
        e.stopPropagation()
        if (postDetail) {
          postDetail.scrollTop -= 150
        }
      } else if (e.key === ' ') {
        // Spacebar: scroll comments down, auto-advance when at bottom
        e.preventDefault()
        e.stopPropagation()
        
        if (postDetail) {
          postDetail.scrollTop += 400
          
          // Check if at bottom and advance to next post
          const isAtBottom = postDetail.scrollHeight - postDetail.scrollTop <= postDetail.clientHeight + 50
          if (isAtBottom && commentsList) {
            const filteredPosts = search.trim()
              ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
              : posts
            const nextIdx = Math.min(selectedIndex + 1, filteredPosts.length - 1)
            if (nextIdx !== selectedIndex) {
              setSelectedIndex(nextIdx)
              setSelected(filteredPosts[nextIdx])
              // Reset scroll for new post
              setTimeout(() => {
                postDetail.scrollTop = 0
              }, 0)
            }
          }
        }
      } else if (e.key === 'Escape') {
        // Escape: exit reader mode if active
        if (readerMode) {
          e.preventDefault()
          e.stopPropagation()
          setReaderMode(false)
        }
      }
    }
    
    document.addEventListener('keydown', handleGlobalKeys, true)
    return () => document.removeEventListener('keydown', handleGlobalKeys, true)
  }, [selected, selectedIndex, posts, search])

  useEffect(() => {
    const filtered = posts.filter(p => 
      p.data.title.toLowerCase().includes(search.toLowerCase())
    ).slice(0, 3)
    setSuggestions(filtered)
    setShowSuggestions(search.length > 0 && focused === 'filter')
  }, [search, posts, focused])

  const filteredPosts = search.trim()
    ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
    : posts

  useEffect(() => {
    if (listRef.current && selectedIndex >= 0) {
      const items = listRef.current.children
      if (items[selectedIndex]) {
        (items[selectedIndex] as HTMLElement).scrollIntoView({ block: 'nearest' })
      }
    }
  }, [selectedIndex, filteredPosts.length])

  return (
    <div className={`app ${readerMode ? 'reader-mode' : ''}`} onTouchStart={handleTwoFingerTouchStart} onTouchEnd={handleTwoFingerTouchEnd}>
      {readerMode && (
        <button 
          className="reader-mode-close-btn"
          onClick={() => setReaderMode(false)}
          title="Exit Reader Mode (Esc)"
          aria-label="Exit Reader Mode"
        >
          ✕
        </button>
      )}
      <div className="left-pane" onTouchStart={handleTouchStart} onTouchEnd={handleTouchEnd}>
        <form className="header" onSubmit={handleSubmit} tabIndex={-1}>
          <input
            ref={searchRef}
            className="sub-input"
            placeholder="subreddit..."
            value={input}
            onChange={e => setInput(e.target.value)}
            onFocus={() => setFocused('search')}
            autoFocus
          />
          <button className="search-btn" tabIndex={-1}>Go</button>
         </form>
         
         {sortingButtonsEnabled && (
           <div className="sort-buttons" tabIndex={-1}>
             <button 
               className={`sort-btn ${sort === 'hot' ? 'active' : ''}`}
               onClick={() => handleSortChange('hot')}
               tabIndex={-1}
             >
               ◉ Hot
             </button>
             <button 
               className={`sort-btn ${sort === 'new' ? 'active' : ''}`}
               onClick={() => handleSortChange('new')}
               tabIndex={-1}
             >
               ◉ New
             </button>
             <button 
               className={`sort-btn ${sort === 'top' ? 'active' : ''}`}
               onClick={() => handleSortChange('top')}
               tabIndex={-1}
             >
               ◉ Top
             </button>
           </div>
         )}
         
         <div className="quick-links" tabIndex={-1}>
          {SUBREDDIT_SUGGESTIONS.slice(0, 8).map(s => (
            <button key={s} className={`quick-link ${sub === s ? 'active' : ''}`}
              onClick={() => handleSub(s)} tabIndex={-1}>
              r/{s}
            </button>
          ))}
        </div>

        <div className={`filter-bar ${readerMode ? 'hidden' : ''}`}>
          <input
            ref={filterRef}
            className="filter-input"
            placeholder={isRedditSearch ? "Search Reddit..." : "Filter posts... (Enter to search Reddit)"}
            value={search}
            onChange={e => {
              setSearch(e.target.value)
              if (isRedditSearch) setIsRedditSearch(false)
            }}
            onFocus={() => setFocused('filter')}
            onKeyDown={handleKeyDown}
          />
        </div>
        
        <div className={`posts-list ${readerMode ? 'hidden' : ''}`} ref={listRef} onKeyDown={handleKeyDown} tabIndex={0}>
          {loading && posts.length === 0 && <div className="loading">Loading...</div>}
          {error && posts.length === 0 && <div className="error">⚠️ {error}</div>}
          {filteredPosts.map((p, i) => (
            <PostItem key={p.data.id} post={p} active={selected?.data?.id === p.data.id || i === selectedIndex}
              onClick={() => {
                setSelectedIndex(i)
                setSelected(p)
              }} />
          ))}
          {search && !isRedditSearch && filteredPosts.length === 0 && posts.length > 0 && (
            <div className="empty">No local matches - press Enter to search Reddit</div>
          )}
          {isRedditSearch && filteredPosts.length === 0 && posts.length === 0 && (
            <div className="empty">No results found</div>
          )}
          {after && !loading && !search.trim() && (
            <div className="load-more">
              <button onClick={() => fetchPosts(sub, after)} className="load-btn">
                Load More Posts
              </button>
            </div>
          )}
           {loading && posts.length > 0 && <div className="loading">Loading more...</div>}
         </div>
          
          <div className="version-footer">
            <span className="version-text">v{APP_VERSION}</span>
            <button 
              className="reader-mode-btn"
              onClick={() => setReaderMode(!readerMode)}
              title={readerMode ? "Exit Reader Mode" : "Enter Reader Mode"}
            >
              {readerMode ? '📖 Exit Reader' : '📖 Reader Mode'}
            </button>
          </div>
        </div>
        
         <div className="right-pane" ref={rightPaneRef} onTouchStart={handleTouchStart} onTouchEnd={handleTouchEnd}>
        <PostDetail 
          post={selected}
          onNextPost={handleNextPost}
          onPreviousPost={handlePreviousPost}
          canGoNext={selectedIndex < (search.trim() ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase())).length - 1 : posts.length - 1)}
          canGoPrevious={selectedIndex > 0}
        />
      </div>
    </div>
  )
}
