import React, { useState } from 'react';
import { Article } from '../types/blog';
import { formatPublicationDate } from '../utils/blogLogic';
import { Bookmark, ArrowRight, Sparkles } from 'lucide-react';

interface HeroLeadArticleProps {
  article: Article;
  onSelectArticle: (article: Article) => void;
  isBookmarked: boolean;
  onToggleBookmark: (articleId: string) => void;
  onSelectAuthor: (authorId: string) => void;
}

export const HeroLeadArticle: React.FC<HeroLeadArticleProps> = ({
  article,
  onSelectArticle,
  isBookmarked,
  onToggleBookmark,
  onSelectAuthor,
}) => {
  const [imgError, setImgError] = useState(false);

  return (
    <article className="group relative rounded-2xl bg-white dark:bg-stone-900/70 border border-stone-200/80 dark:border-stone-800/80 overflow-hidden transition-all duration-300 hover:shadow-lg hover:border-stone-300 dark:hover:border-stone-700">
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-0 items-stretch">
        {/* Left Column: Editorial Content */}
        <div className="lg:col-span-7 p-6 sm:p-8 lg:p-10 flex flex-col justify-between">
          <div>
            {/* Unboxed Metadata Header (No pills, quiet inline text with dot separator) */}
            <div className="flex flex-wrap items-center gap-2 text-xs font-medium text-stone-500 dark:text-stone-400 mb-4 tracking-wide uppercase font-mono">
              <span className="text-amber-700 dark:text-amber-400 font-semibold">
                {article.category}
              </span>
              <span aria-hidden="true" className="text-stone-300 dark:text-stone-700">·</span>
              <span>{formatPublicationDate(article.publishedAt)}</span>
              <span aria-hidden="true" className="text-stone-300 dark:text-stone-700">·</span>
              <span>{article.readTimeMinutes} min de leitura</span>
            </div>

            {/* Dominant Headline */}
            <h2
              onClick={() => onSelectArticle(article)}
              className="text-2xl sm:text-3xl lg:text-4xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight leading-tight group-hover:text-amber-800 dark:group-hover:text-amber-400 transition-colors cursor-pointer text-balance"
            >
              {article.title}
            </h2>

            {/* Subtitle / Excerpt */}
            <p className="mt-4 text-stone-600 dark:text-stone-300 text-sm sm:text-base leading-relaxed line-clamp-3">
              {article.leadParagraph}
            </p>
          </div>

          {/* Author Footer & Actions */}
          <div className="mt-8 pt-6 border-t border-stone-100 dark:border-stone-800 flex items-center justify-between gap-4">
            <button
              onClick={() => onSelectAuthor(article.author.id)}
              className="flex items-center gap-3 text-left group/author cursor-pointer"
            >
              <img
                src={article.author.avatar}
                alt={article.author.name}
                referrerPolicy="no-referrer"
                className="w-10 h-10 rounded-full object-cover ring-1 ring-stone-200 dark:ring-stone-700"
              />
              <div>
                <p className="text-xs sm:text-sm font-semibold text-stone-900 dark:text-stone-100 group-hover/author:text-amber-700 dark:group-hover/author:text-amber-400 transition-colors">
                  {article.author.name}
                </p>
                <p className="text-[11px] text-stone-500 dark:text-stone-400 truncate max-w-[180px] sm:max-w-xs">
                  {article.author.role}
                </p>
              </div>
            </button>

            <div className="flex items-center gap-2">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onToggleBookmark(article.id);
                }}
                className={`p-2 rounded-lg border transition-colors cursor-pointer ${
                  isBookmarked
                    ? 'border-amber-600 bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400'
                    : 'border-stone-200 dark:border-stone-800 text-stone-500 hover:text-stone-900 dark:hover:text-stone-100 hover:bg-stone-50 dark:hover:bg-stone-800'
                }`}
                aria-label={isBookmarked ? 'Remover dos favoritos' : 'Salvar artigo'}
                title={isBookmarked ? 'Salvo' : 'Salvar para depois'}
              >
                <Bookmark className={`w-4 h-4 ${isBookmarked ? 'fill-current' : ''}`} />
              </button>

              <button
                onClick={() => onSelectArticle(article)}
                className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-stone-900 dark:bg-stone-100 dark:text-stone-900 rounded-lg hover:bg-stone-800 dark:hover:bg-stone-200 transition-colors cursor-pointer"
              >
                <span>Ler Artigo</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

        {/* Right Column: Hero Visual Asset */}
        <div
          onClick={() => onSelectArticle(article)}
          className="lg:col-span-5 relative min-h-[260px] lg:min-h-full overflow-hidden bg-stone-100 dark:bg-stone-800 cursor-pointer"
        >
          {!imgError ? (
            <img
              src={article.coverImage}
              alt={article.title}
              referrerPolicy="no-referrer"
              onError={() => setImgError(true)}
              className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105"
            />
          ) : (
            <div className="w-full h-full min-h-[240px] flex flex-col items-center justify-center p-8 text-center bg-gradient-to-br from-stone-800 to-stone-950 text-stone-200">
              <Sparkles className="w-8 h-8 text-amber-400 mb-2" />
              <p className="text-sm font-medium">{article.title}</p>
            </div>
          )}
          {/* Subtle contrast scrim */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/40 via-transparent to-transparent opacity-60 pointer-events-none" />
        </div>
      </div>
    </article>
  );
};
