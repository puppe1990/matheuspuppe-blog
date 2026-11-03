import { Article, ContentBlock } from '../types/blog';

/**
 * Filters articles based on search text, category, maximum reading time, and bookmark state.
 */
export function filterArticles(
  articles: Article[],
  query: string = '',
  category: string = 'Todos',
  maxReadTime: number | null = null,
  onlyBookmarked: boolean = false,
  bookmarkIds: string[] = []
): Article[] {
  const normalizedQuery = query.trim().toLowerCase();

  return articles.filter((article) => {
    // 1. Bookmark filter
    if (onlyBookmarked && !bookmarkIds.includes(article.id)) {
      return false;
    }

    // 2. Category filter
    if (category !== 'Todos' && article.category !== category) {
      return false;
    }

    // 3. Max read time filter
    if (maxReadTime !== null && article.readTimeMinutes > maxReadTime) {
      return false;
    }

    // 4. Search query filter
    if (normalizedQuery) {
      const titleMatch = article.title.toLowerCase().includes(normalizedQuery);
      const subtitleMatch = article.subtitle.toLowerCase().includes(normalizedQuery);
      const authorMatch = article.author.name.toLowerCase().includes(normalizedQuery);
      const tagMatch = article.tags.some((t) => t.toLowerCase().includes(normalizedQuery));
      const leadMatch = article.leadParagraph.toLowerCase().includes(normalizedQuery);

      if (!titleMatch && !subtitleMatch && !authorMatch && !tagMatch && !leadMatch) {
        return false;
      }
    }

    return true;
  });
}

/**
 * Calculates estimated reading time in minutes based on total words.
 */
export function calculateReadingTime(wordsCount: number, wordsPerMinute: number = 200): number {
  if (wordsCount <= 0) return 1;
  const minutes = Math.ceil(wordsCount / wordsPerMinute);
  return Math.max(1, minutes);
}

/**
 * Estimates words count across all text blocks in an article.
 */
export function estimateWordCount(blocks: ContentBlock[]): number {
  let count = 0;
  for (const block of blocks) {
    if (block.text) {
      count += block.text.trim().split(/\s+/).filter(Boolean).length;
    }
    if (block.items) {
      for (const item of block.items) {
        count += item.trim().split(/\s+/).filter(Boolean).length;
      }
    }
  }
  return count;
}

/**
 * Toggles an article ID in the bookmarks list.
 */
export function toggleBookmark(bookmarkIds: string[], articleId: string): string[] {
  if (bookmarkIds.includes(articleId)) {
    return bookmarkIds.filter((id) => id !== articleId);
  }
  return [...bookmarkIds, articleId];
}

/**
 * Increments claps with a sensible safety cap per session.
 */
export function calculateClap(current: number, increment: number = 1): number {
  return current + Math.max(1, increment);
}

/**
 * Formats an ISO date into a refined Portuguese editorial date representation.
 */
export function formatPublicationDate(isoDate: string): string {
  try {
    const date = new Date(isoDate);
    if (isNaN(date.getTime())) return isoDate;
    return new Intl.DateTimeFormat('pt-BR', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    }).format(date);
  } catch {
    return isoDate;
  }
}

/**
 * Calculates scroll percentage based on current scroll position and total scrollable height.
 */
export function calculateScrollProgress(scrollY: number, totalHeight: number): number {
  if (totalHeight <= 0) return 0;
  const progress = (scrollY / totalHeight) * 100;
  return Math.min(100, Math.max(0, Number(progress.toFixed(1))));
}

/**
 * Calculates remaining reading time in minutes based on current scroll progress percentage.
 */
export function calculateRemainingReadingTime(readTimeMinutes: number, scrollProgressPercent: number): number {
  const remaining = Math.ceil(readTimeMinutes * (1 - Math.min(100, Math.max(0, scrollProgressPercent)) / 100));
  return Math.max(0, remaining);
}

/**
 * Groups article counts per category.
 */
export function getCategoryCounts(articles: Article[]): Record<string, number> {
  const counts: Record<string, number> = { Todos: articles.length };
  for (const article of articles) {
    counts[article.category] = (counts[article.category] || 0) + 1;
  }
  return counts;
}

/**
 * Builds standard metadata payload for the Web Share API.
 */
export function buildArticleShareData(article: { title: string; subtitle?: string }, url: string): { title: string; text: string; url: string } {
  return {
    title: article.title,
    text: article.subtitle ? `${article.title} — ${article.subtitle}` : article.title,
    url: url || '',
  };
}
