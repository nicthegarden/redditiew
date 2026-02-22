import { useState, useCallback, useRef, useEffect } from 'react'

const PROXY = 'http://localhost:3001'
const API_BASE = 'http://localhost:3001/api'

const SUBREDDIT_SUGGESTIONS = [
  'sysadmin', 'IT', 'linux', 'homelab', 'networking', 'devops', 'technology',
  'programming', 'javascript', 'python', 'docker', 'kubernetes', 'aws',
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
  
  const listRef = useRef<HTMLDivElement>(null)
  const searchRef = useRef<HTMLInputElement>(null)
  const filterRef = useRef<HTMLInputElement>(null)
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const { cached, saveToCache } = usePostCache()

  const fetchPosts = useCallback(async (subreddit: string, cursor: string | null = null) => {
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
      const target = `${API_BASE}/r/${subreddit}.json?limit=${limit}${cursor ? '&after=' + cursor : ''}`
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

  const handleSub = (subreddit: string) => {
    setInput(subreddit)
    setSelected(null)
    setSelectedIndex(0)
    fetchPosts(subreddit)
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const clean = input.replace(/^\/?r\/?/, '').trim()
    if (clean) handleSub(clean)
  }

  const handleSearchSubmit = () => {
    if (search.trim()) {
      // Search Reddit across all subreddits
      setIsRedditSearch(true)
      setLoading(true)
      setError(null)
      
      const query = encodeURIComponent(search.trim())
      const searchUrl = `${API_BASE}/search.json?q=${query}&type=link&limit=50`
      
      fetch(searchUrl)
        .then(res => res.json())
        .then(data => {
          const items = data.data.children as RedditPost[]
          setPosts(items)
          setSub(`search: "${search}"`)
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

  useEffect(() => {
    const handleGlobalKeys = (e: KeyboardEvent) => {
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
            placeholder={isRedditSearch ? "Search Reddit..." : "Filter posts... (Enter to search Reddit)"}
            value={search}
            onChange={e => setSearch(e.target.value)}
            onFocus={() => setFocused('filter')}
          />
        </div>
        
        <div className="posts-list" ref={listRef} onKeyDown={handleKeyDown} tabIndex={0}>
          {loading && posts.length === 0 && <div className="loading">Loading...</div>}
          {error && posts.length === 0 && <div className="error">⚠️ {error}</div>}
          {filteredPosts.map((p, i) => (
            <PostItem key={p.data.id} post={p} active={selected?.data?.id === p.data.id || i === selectedIndex}
              onClick={() => setSelected(p)} />
          ))}
          {search && !isRedditSearch && filteredPosts.length === 0 && posts.length > 0 && (
            <div className="empty">No local matches - press Enter to search Reddit</div>
          )}
          {isRedditSearch && filteredPosts.length === 0 && posts.length === 0 && (
            <div className="empty">No results found</div>
          )}
          {after && !loading && (
            <div className="load-more">
              <button onClick={() => fetchPosts(sub, after)} className="load-btn">
                Load More Posts
              </button>
            </div>
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
