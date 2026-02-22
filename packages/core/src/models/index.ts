/**
 * Data models and interfaces for RedditView
 * Shared between web and CLI/TUI applications
 */

export interface RedditPostData {
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
  thumbnail?: string
  preview?: {
    images: Array<{
      source: {
        url: string
      }
    }>
  }
}

export interface RedditPost {
  kind: string
  data: RedditPostData
}

export interface Comment {
  id: string
  author: string
  body: string
  score: number
  created_utc: number
  replies: Comment[]
  depth: number
}

export interface CacheEntry {
  posts: RedditPost[]
  time: number
}

export interface FetchPostsResult {
  posts: RedditPost[]
  after: string | null
  error?: string
}

export interface FetchCommentsResult {
  comments: Comment[]
  error?: string
}

export interface SearchResult {
  posts: RedditPost[]
  error?: string
}
