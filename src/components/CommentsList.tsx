import { useEffect, useState } from 'react'

interface Comment {
  id: string
  author: string
  body: string
  score: number
  created_utc: number
  replies: Comment[]
  depth: number
}

interface CommentsListProps {
  permalink: string
}

function formatTime(ts: number): string {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

function CommentTree({ comment, depth = 0 }: { comment: Comment; depth?: number }) {
  const [collapsed, setCollapsed] = useState(false)

  if (collapsed) {
    return (
      <div className="comment-collapsed" style={{ marginLeft: `${depth * 16}px` }}>
        <button onClick={() => setCollapsed(false)} className="collapse-btn">
          [{comment.replies?.length || 0} replies]
        </button>
      </div>
    )
  }

  return (
    <div className="comment" style={{ marginLeft: `${depth * 16}px` }}>
      <div className="comment-header">
        <span className="comment-author">{comment.author}</span>
        <span className="comment-time">{formatTime(comment.created_utc)}</span>
        <span className="comment-score">↑ {comment.score}</span>
        {comment.replies?.length > 0 && (
          <button onClick={() => setCollapsed(true)} className="collapse-btn">
            −
          </button>
        )}
      </div>
      <div className="comment-body">{comment.body}</div>
      {comment.replies?.map((reply) => (
        <CommentTree key={reply.id} comment={reply} depth={depth + 1} />
      ))}
    </div>
  )
}

function parseComments(data: any, depth = 0): Comment[] {
  if (!data || !data.data || !data.data.children) return []

  return data.data.children
    .filter((child: any) => child.kind === 't1') // Comment kind
    .map((child: any) => ({
      id: child.data.id,
      author: child.data.author || '[deleted]',
      body: child.data.body || '',
      score: child.data.score || 0,
      created_utc: child.data.created_utc || 0,
      depth: depth,
      replies: parseComments(child.data.replies, depth + 1),
    }))
}

export default function CommentsList({ permalink }: CommentsListProps) {
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    setError(null)
    setComments([])

    // Fetch comments with better error handling
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 15000) // 15s timeout

    fetch(`http://localhost:3001/api${permalink}.json`, { signal: controller.signal })
      .then(async (res) => {
        clearTimeout(timeoutId)
        if (!res.ok) {
          console.error(`Proxy returned HTTP ${res.status}:`, res.statusText)
          throw new Error(`HTTP ${res.status}: ${res.statusText}`)
        }
        return res.json()
      })
      .then((data) => {
        if (Array.isArray(data) && data.length > 1) {
          const commentsList = parseComments(data[1])
          setComments(commentsList)
        }
      })
      .catch((err) => {
        if (err.name === 'AbortError') {
          setError('Request timeout - comments took too long to load')
        } else if (err instanceof TypeError) {
          setError('Network error - check if proxy is running on port 3001')
        } else if (err instanceof Error) {
          setError(err.message)
        } else {
          setError('Unknown error loading comments')
        }
        console.error('Comments fetch error:', err)
      })
      .finally(() => {
        setLoading(false)
      })

    return () => {
      clearTimeout(timeoutId)
      controller.abort()
    }
  }, [permalink])

  if (loading) return <div className="comments-loading">Loading comments...</div>
  if (error) return <div className="comments-error">⚠️ {error}</div>
  if (comments.length === 0) return <div className="comments-empty">No comments yet</div>

  return (
    <div className="comments-list">
      {comments.map((comment) => (
        <CommentTree key={comment.id} comment={comment} />
      ))}
    </div>
  )
}
