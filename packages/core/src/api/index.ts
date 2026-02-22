/**
 * Reddit API client
 * Handles data fetching, parsing, and caching
 * Works across web and CLI applications
 */

import type {
  RedditPost,
  Comment,
  FetchPostsResult,
  FetchCommentsResult,
  SearchResult,
} from '../models/index.js'
import { PostCache } from '../cache/index.js'

export interface ApiClientConfig {
  baseUrl?: string
  cacheConfig?: any
  timeout?: number
}

export class RedditApiClient {
  private baseUrl: string
  private cache: PostCache
  private timeout: number

  constructor(config: ApiClientConfig = {}) {
    this.baseUrl = config.baseUrl ?? '/api'
    this.cache = new PostCache(config.cacheConfig)
    this.timeout = config.timeout ?? 15000
  }

  /**
   * Fetch posts from a subreddit
   * Uses cache to avoid redundant requests
   */
  async fetchPosts(
    subreddit: string,
    options?: {
      cursor?: string | null
      limit?: number
      useCache?: boolean
    }
  ): Promise<FetchPostsResult> {
    const { cursor = null, limit = 50, useCache = true } = options ?? {}
    const isMore = !!cursor

    // Check cache for non-pagination requests
    if (!isMore && useCache) {
      const cached = this.cache.get(subreddit)
      if (cached) {
        return { posts: cached, after: null }
      }
    }

    try {
      const url = new URL(`${this.baseUrl}/r/${subreddit}.json`, this.getBaseUrl())
      url.searchParams.set('limit', limit.toString())
      if (cursor) {
        url.searchParams.set('after', cursor)
      }

      const response = await this.fetchWithTimeout(url.toString())

      if (response.status === 429) {
        const data = await response.json().catch(() => ({}))
        const retryAfter = data.retry_after || 60
        return {
          posts: [],
          after: null,
          error: `Rate limited. Please try again in ${retryAfter} seconds.`,
        }
      }

      if (!response.ok) {
        let error = `Failed to load (${response.status})`
        if (response.status === 404) {
          error = `Subreddit r/${subreddit} not found`
        } else if (response.status >= 500) {
          error = 'Reddit server error. Please try again later.'
        }
        return { posts: [], after: null, error }
      }

      const data = await response.json()
      const posts = data.data.children as RedditPost[]

      // Cache the result
      if (!isMore && useCache) {
        this.cache.set(subreddit, posts)
      }

      return {
        posts,
        after: data.data.after,
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      return { posts: [], after: null, error: message }
    }
  }

  /**
   * Fetch comments for a post
   */
  async fetchComments(permalink: string): Promise<FetchCommentsResult> {
    try {
      const url = new URL(`${this.baseUrl}${permalink}.json`, this.getBaseUrl())
      const response = await this.fetchWithTimeout(url.toString())

      if (!response.ok) {
        return {
          comments: [],
          error: `Failed to load comments (${response.status})`,
        }
      }

      const data = await response.json()
      if (Array.isArray(data) && data.length > 1) {
        const comments = this.parseComments(data[1])
        return { comments }
      }

      return { comments: [] }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      return { comments: [], error: message }
    }
  }

  /**
   * Search Reddit across all subreddits
   */
  async search(query: string, limit = 50): Promise<SearchResult> {
    try {
      const url = new URL(`${this.baseUrl}/search.json`, this.getBaseUrl())
      url.searchParams.set('q', query)
      url.searchParams.set('type', 'link')
      url.searchParams.set('limit', limit.toString())

      const response = await this.fetchWithTimeout(url.toString())

      if (!response.ok) {
        return {
          posts: [],
          error: `Search failed (${response.status})`,
        }
      }

      const data = await response.json()
      const posts = data.data.children as RedditPost[]

      return { posts }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error'
      return { posts: [], error: message }
    }
  }

  /**
   * Clear the cache
   */
  clearCache(): void {
    this.cache.clear()
  }

  /**
   * Get cache stats for debugging
   */
  getCacheStats(): { size: number } {
    return { size: this.cache.size() }
  }

  /**
   * Parse comments from Reddit API response
   */
  private parseComments(data: any, depth = 0): Comment[] {
    if (!data || !data.data || !data.data.children) {
      return []
    }

    return data.data.children
      .filter((child: any) => child.kind === 't1') // Comment kind
      .map((child: any) => ({
        id: child.data.id,
        author: child.data.author || '[deleted]',
        body: child.data.body || '',
        score: child.data.score || 0,
        created_utc: child.data.created_utc || 0,
        depth: depth,
        replies: this.parseComments(child.data.replies, depth + 1),
      }))
  }

  /**
   * Fetch with timeout support
   */
  private fetchWithTimeout(url: string, init?: RequestInit): Promise<Response> {
    return Promise.race([
      fetch(url, init),
      new Promise<Response>((_, reject) =>
        setTimeout(() => reject(new Error('Request timeout')), this.timeout)
      ),
    ])
  }

  /**
   * Get base URL for relative path resolution
   */
  private getBaseUrl(): string {
    if (typeof window !== 'undefined') {
      return window.location.origin
    }
    // For Node.js/CLI usage
    return 'http://localhost:3001'
  }
}

// Export singleton instance
export const apiClient = new RedditApiClient()
