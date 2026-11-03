import React, { useState } from 'react';
import { Article } from '../types/blog';
import { formatPublicationDate } from '../utils/blogLogic';
import { Bookmark, ArrowRight, Code } from 'lucide-react';

interface ArticleCardProps {
  article: Article;
  onSelectArticle: (article: Article) => void;
  isBookmarked: boolean;
  onToggleBookmark: (articleId: string) => void;
  onSelectAuthor: (authorId: string) => void;
  variant?: 'featured' | 'standard';
}

export const ArticleCard: React.FC<ArticleCardProps> = ({
  article,
  onSelectArticle,
  isBookmarked,
  onToggleBookmark,
  onSelectAuthor,
  variant = 'standard',
}) => {
  const [imgError, setImgError] = useState(false);
  const isFeatured = variant === 'featured';

  return (
    <article className="group flex flex-col justify-between rounded-xl bg-white dark:bg-stone-900/60 border border-stone-200/80 dark:border-stone-800/80 overflow-hidden transition-all duration-300 hover:shadow-md hover:border-stone-300 dark:hover:border-stone-700">
      <div>
        {/* Cover Media */}
        <div
          onClick={() => onSelectArticle(article)}
          className={`relative overflow-hidden bg-stone-100 dark:bg-stone-800 cursor-pointer ${
            isFeatured ? 'aspect-[16/10]' : 'aspect-[16/9]'
          }`}
        >
          {!imgError ? (
            <img
              src={article.coverImage}
              alt={article.title}
              referrerPolicy="no-referrer"
              onError={() => setImgError(true)}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
          ) : (
            <div className="w-full h-full flex flex-col items-center justify-center p-6 text-center bg-stone-100 dark:bg-stone-800 text-stone-600 dark:text-stone-300">
              <Code className="w-6 h-6 mb-2 text-stone-400" />
              <p className="text-xs font-mono">{article.category}</p>
            </div>
          )}
        </div>

        {/* Content Area */}
        <div className="p-5 sm:p-6">
          {/* Unboxed Metadata Header */}
          <div className="flex items-center gap-2 text-[11px] font-mono tracking-wider uppercase text-stone-500 dark:text-stone-400 mb-2.5">
            <span className="text-amber-700 dark:text-amber-400 font-semibold">
              {article.category}
            </span>
            <span aria-hidden="true" className="text-stone-300 dark:text-stone-700">·</span>
            <span>{article.readTimeMinutes} min de leitura</span>
          </div>

          {/* Title */}
          <h3
            onClick={() => onSelectArticle(article)}
            className={`font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight leading-snug group-hover:text-amber-800 dark:group-hover:text-amber-400 transition-colors cursor-pointer text-balance ${
              isFeatured ? 'text-lg sm:text-xl' : 'text-base sm:text-lg'
            }`}
          >
            {article.title}
          </h3>

          {/* Subtitle / Excerpt */}
          <p className="mt-2.5 text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed line-clamp-2">
            {article.subtitle || article.leadParagraph}
          </p>
        </div>
      </div>

      {/* Card Footer */}
      <div className="px-5 pb-5 sm:px-6 sm:pb-6 pt-3 border-t border-stone-100 dark:border-stone-800/80 flex items-center justify-between">
        <button
          onClick={() => onSelectAuthor(article.author.id)}
          className="flex items-center gap-2.5 text-left group/auth cursor-pointer"
        >
          <img
            src={article.author.avatar}
            alt={article.author.name}
            referrerPolicy="no-referrer"
            className="w-7 h-7 rounded-full object-cover ring-1 ring-stone-200 dark:ring-stone-700"
          />
          <div className="leading-tight">
            <p className="text-xs font-semibold text-stone-800 dark:text-stone-200 group-hover/auth:text-amber-700 dark:group-hover/auth:text-amber-400 transition-colors">
              {article.author.name}
            </p>
            <p className="text-[10px] text-stone-400 font-mono">
              {formatPublicationDate(article.publishedAt)}
            </p>
          </div>
        </button>

        <div className="flex items-center gap-1.5">
          <button
            onClick={(e) => {
              e.stopPropagation();
              onToggleBookmark(article.id);
            }}
            className={`p-1.5 rounded-md border transition-colors cursor-pointer ${
              isBookmarked
                ? 'border-amber-600 bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400'
                : 'border-transparent text-stone-400 hover:text-stone-700 dark:hover:text-stone-200 hover:bg-stone-100 dark:hover:bg-stone-800'
            }`}
            aria-label={isBookmarked ? 'Remover dos favoritos' : 'Salvar artigo'}
            title={isBookmarked ? 'Salvo' : 'Salvar'}
          >
            <Bookmark className={`w-3.5 h-3.5 ${isBookmarked ? 'fill-current' : ''}`} />
          </button>

          <button
            onClick={() => onSelectArticle(article)}
            className="p-1.5 rounded-md text-stone-400 hover:text-stone-900 dark:hover:text-stone-100 hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
            aria-label="Abrir artigo"
          >
            <ArrowRight className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </article>
  );
};
