import { useState, useCallback, useRef, useEffect } from 'react'

const PROXY = 'http://localhost:3001'
const API_BASE = 'http://localhost:3001/api'

const DEFAULT_SHORTCUTS = {
  '1': 'sysadmin',
  '2': 'golang',
  '3': 'programming',
  '4': 'linux',
  '5': 'devops',
  '6': 'webdev',
  '7': 'learnprogramming',
  '8': '100DaysOfCode',
  '9': 'codereview'
}

const SUBREDDIT_SUGGESTIONS = [
  'sysadmin', 'IT', 'linux', 'homelab', 'networking', 'devops', 'technology',
  'programming', 'javascript', 'python', 'docker', 'kubernetes', 'aws',
]

function formatTime(ts) {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

function formatNum(n) {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n
}

function getThumb(post) {
  const p = post.data
  if (p.thumbnail && p.thumbnail.startsWith('http')) return p.thumbnail
  if (p.preview?.images?.[0]?.source?.url) {
    return p.preview.images[0].source.url.replace(/&amp;/g, '&')
  }
  return null
}

function PostItem({ post, active, onClick }) {
  const thumb = getThumb(post)
  const p = post.data
  
  return (
    <div className={`post-item ${active ? 'active' : ''}`} onClick={onClick}>
      {thumb && <img src={thumb} alt="" className="post-thumb" loading="lazy" />}
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
  const [cached, setCached] = useState(() => {
    const saved = localStorage.getItem('postCache')
    return saved ? JSON.parse(saved) : {}
  })

  const saveToCache = useCallback((sub, posts) => {
    setCached(prev => {
      const updated = { ...prev, [sub]: { posts, time: Date.now() } }
      localStorage.setItem('postCache', JSON.stringify(updated))
      return updated
    })
  }, [])

  return { cached, saveToCache }
}

export default function App() {
  const [sub, setSub] = useState('sysadmin')
  const [input, setInput] = useState('sysadmin')
  const [search, setSearch] = useState('')
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [selected, setSelected] = useState(null)
  const [selectedIndex, setSelectedIndex] = useState(0)
  const [after, setAfter] = useState(null)
  const [focused, setFocused] = useState('search')
  const [suggestions, setSuggestions] = useState([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [suggestionIndex, setSuggestionIndex] = useState(-1)
  const [sort, setSort] = useState('hot')
  const [shortcuts, setShortcuts] = useState(DEFAULT_SHORTCUTS)
  
  const listRef = useRef(null)
  const searchRef = useRef(null)
  const filterRef = useRef(null)
  const iframeRef = useRef(null)
  const { cached, saveToCache } = usePostCache()

  const fetchPosts = useCallback(async (subreddit, cursor = null, sortType = 'hot') => {
    const isMore = !!cursor
    const cacheKey = `${subreddit}:${sortType}`
    const cachedData = cached[cacheKey]
    
    if (!isMore && cachedData && Date.now() - cachedData.time < 3600000) {
      setPosts(cachedData.posts)
      return
    }

    setLoading(true)
    setError(null)
    
    try {
      const limit = 50
      const target = `${API_BASE}/r/${subreddit}/${sortType}.json?limit=${limit}${cursor ? '&after=' + cursor : ''}`
      const res = await fetch(target)
      if (!res.ok) throw new Error('Failed to load')
      const data = await res.json()
      const items = data.data.children
      
      if (isMore) {
        setPosts(prev => [...prev, ...items])
      } else {
        setPosts(items)
        saveToCache(cacheKey, items)
      }
      setAfter(data.data.after)
      setSub(subreddit)
      setSort(sortType)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }, [cached, saveToCache])

  const handleSub = (subreddit, sortType = sort) => {
    setInput(subreddit)
    setSelected(null)
    setSelectedIndex(0)
    fetchPosts(subreddit, null, sortType)
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const clean = input.replace(/^\/?r\/?/, '').trim()
    if (clean) handleSub(clean)
  }

  const handleSearchSubmit = () => {
    if (search.trim()) {
      setSelected({ data: { permalink: `/search?q=${encodeURIComponent(search)}` } })
    }
  }

  const handleKeyDown = (e) => {
    const filteredPosts = search.trim()
      ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
      : posts

    if (focused === 'list') {
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        setSelectedIndex(prev => Math.min(prev + 1, filteredPosts.length - 1))
      } else if (e.key === 'ArrowUp') {
        e.preventDefault()
        setSelectedIndex(prev => Math.max(prev - 1, 0))
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
    fetchPosts(sub, null, sort)
  }, [])

  useEffect(() => {
    // Load config from backend
    const loadConfig = async () => {
      try {
        const res = await fetch(`${API_BASE}/config`)
        if (res.ok) {
          const config = await res.json()
          if (config.tui?.subreddit_shortcuts) {
            setShortcuts(config.tui.subreddit_shortcuts)
          }
        }
      } catch (err) {
        console.warn('Failed to load config from backend:', err.message)
        // Fall back to DEFAULT_SHORTCUTS which is already set
      }
    }
    loadConfig()
  }, [])

  useEffect(() => {
    const handleShortcuts = (e) => {
      // 't' key - toggle sort
      if (e.key === 't' || e.key === 'T') {
        e.preventDefault()
        e.stopPropagation()
        const newSort = sort === 'hot' ? 'new' : 'hot'
        fetchPosts(sub, null, newSort)
        return
      }

      // '1'-'9' keys - subreddit shortcuts
      if (e.key >= '1' && e.key <= '9') {
        const shortcutSub = shortcuts[e.key]
        if (shortcutSub) {
          e.preventDefault()
          e.stopPropagation()
          handleSub(shortcutSub, sort)
        }
      }
    }
    document.addEventListener('keydown', handleShortcuts, true)
    return () => document.removeEventListener('keydown', handleShortcuts, true)
  }, [sub, sort, handleSub, fetchPosts, shortcuts])

  useEffect(() => {
    const handleTab = (e) => {
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

  useEffect(() => {
    const handleGlobalKeys = (e) => {
      if (e.key === 'PageDown' || e.key === 'PageUp') {
        e.preventDefault()
        e.stopPropagation()
        if (iframeRef.current) {
          try {
            if (iframeRef.current.contentWindow) {
              iframeRef.current.contentWindow.scrollBy(0, e.key === 'PageDown' ? 400 : -400)
            }
          } catch {}
          iframeRef.current.scrollTop += e.key === 'PageDown' ? 400 : -400
        }
      }
    }
    document.addEventListener('keydown', handleGlobalKeys, true)
    return () => document.removeEventListener('keydown', handleGlobalKeys, true)
  }, [])

  useEffect(() => {
    const filtered = posts.filter(p => 
      p.data.title.toLowerCase().includes(search.toLowerCase())
    ).slice(0, 3)
    setSuggestions(filtered)
    setShowSuggestions(search.length > 0 && focused === 'filter')
  }, [search, posts, focused])

  useEffect(() => {
    if (listRef.current) {
      const items = listRef.current.children
      if (items[selectedIndex]) {
        items[selectedIndex].scrollIntoView({ block: 'nearest' })
      }
    }
  }, [selectedIndex])

  const filteredPosts = search.trim()
    ? posts.filter(p => p.data.title.toLowerCase().includes(search.toLowerCase()))
    : posts

  return (
    <div className="app">
      <div className="left-pane">
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
          <button 
            type="button"
            className="sort-btn" 
            onClick={() => {
              const newSort = sort === 'hot' ? 'new' : 'hot'
              fetchPosts(sub, null, newSort)
            }}
            tabIndex={-1}
            title={`Toggle sort (current: ${sort}). Press 't' to toggle`}
          >
            {sort === 'hot' ? '📊 Hot' : '🆕 New'}
          </button>
        </form>
        
        <div className="quick-links" tabIndex={-1}>
          {SUBREDDIT_SUGGESTIONS.slice(0, 8).map(s => (
            <button key={s} className={`quick-link ${sub === s ? 'active' : ''}`}
              onClick={() => handleSub(s)} tabIndex={-1}>
              r/{s}
            </button>
          ))}
        </div>

        <div className="filter-bar">
          <input
            ref={filterRef}
            className="filter-input"
            placeholder="Filter posts... (Enter to search)"
            value={search}
            onChange={e => setSearch(e.target.value)}
            onFocus={() => setFocused('filter')}
          />
        </div>
        
        <div className="posts-list" ref={listRef} onScroll={() => {}} onKeyDown={handleKeyDown} tabIndex={0}>
          {loading && posts.length === 0 && <div className="loading">Loading...</div>}
          {error && posts.length === 0 && <div className="error">{error}</div>}
          {filteredPosts.map((p, i) => (
            <PostItem key={p.data.id} post={p} active={selected?.data?.id === p.data.id || i === selectedIndex}
              onClick={() => setSelected(p)} />
          ))}
          {search && filteredPosts.length === 0 && posts.length > 0 && (
            <div className="empty">No matches - Enter to search Reddit</div>
          )}
          {loading && posts.length > 0 && <div className="loading">Loading more...</div>}
        </div>
      </div>
      
      <div className="right-pane">
        {selected ? (
          <iframe 
            ref={iframeRef}
            src={`${PROXY}${selected.data.permalink}`} 
            title="Reddit" 
            sandbox="allow-scripts allow-same-origin allow-popups allow-forms" 
          />
        ) : (
          <div className="empty-state">
            <h3>Select a post</h3>
            <p>Use ↑↓ to navigate, Enter to open</p>
          </div>
        )}
      </div>
    </div>
  )
}
