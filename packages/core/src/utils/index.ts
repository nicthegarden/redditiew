/**
 * Utility functions for formatting and data manipulation
 * Shared across all platforms
 */

export function formatTime(ts: number): string {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s`
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

export function formatTimeAgo(ts: number): string {
  const diff = Math.floor((Date.now() / 1000) - ts)
  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

export function formatNum(n: number): string {
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return n.toString()
}

export function getSubredditFromInput(input: string): string {
  return input.replace(/^\/?r\/?/, '').trim()
}

export function isTextPost(post: { url?: string }): boolean {
  return !post.url || post.url.includes('reddit.com')
}

export function getPostThumbnail(post: { thumbnail?: string; preview?: { images: Array<{ source: { url: string } }> } }): string | null {
  if (post.thumbnail && post.thumbnail.startsWith('http')) {
    return post.thumbnail
  }
  if (post.preview?.images?.[0]?.source?.url) {
    return post.preview.images[0].source.url.replace(/&amp;/g, '&')
  }
  return null
}
