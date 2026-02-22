import { useEffect, useState } from 'react'
import CommentsList from './CommentsList'

interface PostDetailProps {
  post: {
    data: {
      id: string
      title: string
      subreddit: string
      author: string
      created_utc: number
      score: number
      num_comments: number
      selftext?: string
      url?: string
      permalink: string
    }
  } | null
}

function formatTime(ts: number): string {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

function formatNum(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toString()
}

export default function PostDetail({ post }: PostDetailProps) {
  const [loading, setLoading] = useState(false)

  if (!post) {
    return (
      <div className="post-detail-empty">
        <h2>Select a post</h2>
        <p>Use ↑↓ to navigate, Enter to open</p>
      </div>
    )
  }

  const data = post.data
  const isTextPost = !data.url || data.url.includes('reddit.com')

  return (
    <div className="post-detail">
      <div className="post-detail-header">
        <h1 className="post-detail-title">{data.title}</h1>
        <div className="post-detail-meta">
          <span>r/{data.subreddit}</span>
          <span>u/{data.author}</span>
          <span>{formatTime(data.created_utc)}</span>
          <span>↑ {formatNum(data.score)}</span>
          <span>💬 {formatNum(data.num_comments)}</span>
        </div>
      </div>

      {isTextPost && data.selftext && (
        <div className="post-body">
          {data.selftext}
        </div>
      )}

      {!isTextPost && data.url && (
        <div className="post-link">
          <a href={data.url} target="_blank" rel="noopener noreferrer">
            🔗 Open link: {new URL(data.url).hostname}
          </a>
        </div>
      )}

      <div className="post-actions">
        <button className="action-btn">↑ Upvote</button>
        <button className="action-btn">↓ Downvote</button>
        <button className="action-btn">💾 Save</button>
        <a href={`https://reddit.com${data.permalink}`} target="_blank" rel="noopener noreferrer" className="action-btn">
          🔗 Open on Reddit
        </a>
      </div>

      <div className="comments-section">
        <h3>Comments ({formatNum(data.num_comments)})</h3>
        <CommentsList permalink={data.permalink} />
      </div>
    </div>
  )
}
