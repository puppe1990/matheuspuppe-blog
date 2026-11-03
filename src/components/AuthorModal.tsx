import React from 'react';
import { Author, Article } from '../types/blog';
import { X, BookOpen, ExternalLink, ArrowRight } from 'lucide-react';

interface AuthorModalProps {
  author: Author | null;
  onClose: () => void;
  articles: Article[];
  onSelectArticle: (article: Article) => void;
}

export const AuthorModal: React.FC<AuthorModalProps> = ({
  author,
  onClose,
  articles,
  onSelectArticle,
}) => {
  if (!author) return null;

  const authorArticles = articles.filter((a) => a.author.id === author.id);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-fade-in">
      <div className="relative w-full max-w-xl rounded-2xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 p-6 sm:p-8 shadow-2xl overflow-hidden max-h-[90vh] flex flex-col">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-5 right-5 p-1.5 rounded-lg text-stone-400 hover:text-stone-900 dark:hover:text-stone-100 hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
          aria-label="Fechar perfil do autor"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Author Bio Header */}
        <div className="flex items-start gap-4 mb-6">
          <img
            src={author.avatar}
            alt={author.name}
            referrerPolicy="no-referrer"
            className="w-16 h-16 rounded-full object-cover ring-2 ring-stone-200 dark:ring-stone-700 shrink-0"
          />
          <div>
            <h3 className="text-xl font-serif font-bold text-stone-900 dark:text-stone-100">
              {author.name}
            </h3>
            <p className="text-xs font-medium text-amber-700 dark:text-amber-400 mb-1">
              {author.role}
            </p>
            {author.twitter && (
              <span className="text-xs text-stone-500 font-mono">
                {author.twitter}
              </span>
            )}
          </div>
        </div>

        <p className="text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed mb-6">
          {author.bio}
        </p>

        {/* Written Articles List */}
        <div className="flex-1 overflow-y-auto pt-4 border-t border-stone-100 dark:border-stone-800">
          <div className="flex items-center gap-1.5 text-xs font-mono font-semibold uppercase text-stone-500 mb-3">
            <BookOpen className="w-3.5 h-3.5" />
            <span>Artigos Publicados ({authorArticles.length})</span>
          </div>

          <div className="space-y-3">
            {authorArticles.map((art) => (
              <div
                key={art.id}
                onClick={() => {
                  onSelectArticle(art);
                  onClose();
                }}
                className="group p-3.5 rounded-xl border border-stone-200/80 dark:border-stone-800/80 hover:bg-stone-50 dark:hover:bg-stone-800/60 transition-colors cursor-pointer flex items-center justify-between gap-3"
              >
                <div>
                  <h4 className="text-xs sm:text-sm font-semibold text-stone-900 dark:text-stone-100 group-hover:text-amber-700 dark:group-hover:text-amber-400 transition-colors line-clamp-1">
                    {art.title}
                  </h4>
                  <div className="flex items-center gap-2 text-[11px] text-stone-500 font-mono mt-1">
                    <span>{art.category}</span>
                    <span>·</span>
                    <span>{art.readTimeMinutes} min de leitura</span>
                  </div>
                </div>
                <ArrowRight className="w-4 h-4 text-stone-400 group-hover:text-stone-800 dark:group-hover:text-stone-200 shrink-0" />
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
