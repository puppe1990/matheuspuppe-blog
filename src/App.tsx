/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState, useEffect, useMemo } from 'react';
import { Article, Comment, Author } from './types/blog';
import { INITIAL_ARTICLES, INITIAL_COMMENTS, AUTHORS } from './data/articles';
import { filterArticles, getCategoryCounts, toggleBookmark } from './utils/blogLogic';
import { Header } from './components/Header';
import { CategoryNav } from './components/CategoryNav';
import { HeroLeadArticle } from './components/HeroLeadArticle';
import { ArticleCard } from './components/ArticleCard';
import { ArticleReader } from './components/ArticleReader';
import { BenchmarkPlayground } from './components/BenchmarkPlayground';
import { AboutView } from './components/AboutView';
import { AuthorModal } from './components/AuthorModal';
import { NewsletterModal } from './components/NewsletterModal';
import { Footer } from './components/Footer';
import { Bookmark, Sparkles, FilterX, Newspaper, ArrowRight } from 'lucide-react';

export default function App() {
  // Theme state
  const [darkMode, setDarkMode] = useState<boolean>(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('sintaxe_dark_mode');
      if (saved !== null) return saved === 'true';
      return window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    return false;
  });

  // Navigation & View state
  const [activeView, setActiveView] = useState<'home' | 'saved' | 'playground' | 'about'>('home');
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [selectedAuthorId, setSelectedAuthorId] = useState<string | null>(null);
  const [newsletterModalOpen, setNewsletterModalOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [isFocusMode, setIsFocusMode] = useState(false);

  // Filter state
  const [selectedCategory, setSelectedCategory] = useState<string>('Todos');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [maxReadTime, setMaxReadTime] = useState<number | null>(null);

  // Persistent storage state (Bookmarks & Comments)
  const [bookmarks, setBookmarks] = useState<string[]>(() => {
    if (typeof window !== 'undefined') {
      try {
        const saved = localStorage.getItem('sintaxe_bookmarks');
        return saved ? JSON.parse(saved) : ['art-1', 'art-2'];
      } catch {
        return ['art-1', 'art-2'];
      }
    }
    return ['art-1', 'art-2'];
  });

  const [articles, setArticles] = useState<Article[]>(INITIAL_ARTICLES);
  const [comments, setComments] = useState<Comment[]>(() => {
    if (typeof window !== 'undefined') {
      try {
        const saved = localStorage.getItem('sintaxe_comments');
        return saved ? JSON.parse(saved) : INITIAL_COMMENTS;
      } catch {
        return INITIAL_COMMENTS;
      }
    }
    return INITIAL_COMMENTS;
  });

  // Synchronize dark mode class on <html>
  useEffect(() => {
    const root = document.documentElement;
    if (darkMode) {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
    try {
      localStorage.setItem('sintaxe_dark_mode', String(darkMode));
    } catch {
      // ignore
    }
  }, [darkMode]);

  // Synchronize bookmarks to localStorage
  useEffect(() => {
    try {
      localStorage.setItem('sintaxe_bookmarks', JSON.stringify(bookmarks));
    } catch {
      // ignore
    }
  }, [bookmarks]);

  // Synchronize comments to localStorage
  useEffect(() => {
    try {
      localStorage.setItem('sintaxe_comments', JSON.stringify(comments));
    } catch {
      // ignore
    }
  }, [comments]);

  const toggleDarkMode = () => setDarkMode((prev) => !prev);

  const handleToggleBookmark = (articleId: string) => {
    setBookmarks((prev) => toggleBookmark(prev, articleId));
  };

  const handleSelectArticle = (article: Article) => {
    setSelectedArticle(article);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleBackToArticles = () => {
    setSelectedArticle(null);
    setIsFocusMode(false);
  };

  const handleAddComment = (commentData: { authorName: string; authorRole: string; content: string }) => {
    if (!selectedArticle) return;
    const newComment: Comment = {
      id: `c-${Date.now()}`,
      articleId: selectedArticle.id,
      authorName: commentData.authorName,
      authorRole: commentData.authorRole,
      content: commentData.content,
      createdAt: new Date().toISOString(),
      claps: 0,
    };
    setComments((prev) => [newComment, ...prev]);
  };

  // Author details lookup
  const currentAuthor: Author | null = useMemo(() => {
    if (!selectedAuthorId) return null;
    return Object.values(AUTHORS).find((a) => a.id === selectedAuthorId) || null;
  }, [selectedAuthorId]);

  // Categories list & counts
  const categoriesList = useMemo(() => {
    return [
      'Todos',
      'IA & Modelos',
      'Engenharia de Sistemas',
      'Hardware & Chips',
      'Design & UX',
      'Segurança & Criptografia',
      'Opinião & Futuro',
    ];
  }, []);

  const categoryCounts = useMemo(() => {
    return getCategoryCounts(articles);
  }, [articles]);

  // Filtered articles list
  const filteredArticles = useMemo(() => {
    const isSavedView = activeView === 'saved';
    return filterArticles(
      articles,
      searchQuery,
      selectedCategory,
      maxReadTime,
      isSavedView,
      bookmarks
    );
  }, [articles, searchQuery, selectedCategory, maxReadTime, activeView, bookmarks]);

  // 3-Tier editorial segmentation when on home view without active search
  const isFiltering = searchQuery.trim() !== '' || selectedCategory !== 'Todos' || maxReadTime !== null;
  const leadArticle = useMemo(() => {
    if (isFiltering) return filteredArticles[0] || null;
    return filteredArticles.find((a) => a.tier === 1) || filteredArticles[0] || null;
  }, [filteredArticles, isFiltering]);

  const secondaryArticles = useMemo(() => {
    if (isFiltering) return filteredArticles.slice(1, 3);
    return filteredArticles.filter((a) => a.tier === 2 && a.id !== leadArticle?.id).slice(0, 2);
  }, [filteredArticles, isFiltering, leadArticle]);

  const gridArticles = useMemo(() => {
    if (isFiltering) return filteredArticles.slice(3);
    const featuredIds = new Set([leadArticle?.id, ...secondaryArticles.map((a) => a.id)].filter(Boolean));
    return filteredArticles.filter((a) => !featuredIds.has(a.id));
  }, [filteredArticles, isFiltering, leadArticle, secondaryArticles]);

  // Article comments for the reader view
  const currentArticleComments = useMemo(() => {
    if (!selectedArticle) return [];
    return comments.filter((c) => c.articleId === selectedArticle.id);
  }, [comments, selectedArticle]);

  return (
    <div className="min-h-screen flex flex-col bg-stone-50 text-stone-900 dark:bg-stone-950 dark:text-stone-100 transition-colors">
      {/* Header with Top Bar Contract (Hidden when Focus Mode is active) */}
      {!(selectedArticle && isFocusMode) && (
        <Header
          darkMode={darkMode}
          onToggleDarkMode={toggleDarkMode}
          savedCount={bookmarks.length}
          activeView={activeView}
          onNavigate={(view) => {
            setActiveView(view);
            setSelectedArticle(null);
            setIsFocusMode(false);
            window.scrollTo({ top: 0, behavior: 'smooth' });
          }}
          onOpenNewsletter={() => setNewsletterModalOpen(true)}
          searchOpen={searchOpen}
          onToggleSearch={() => setSearchOpen((prev) => !prev)}
        />
      )}

      {/* Main Content Area */}
      <main
        className={`flex-1 w-full ${
          selectedArticle && isFocusMode
            ? 'p-0 max-w-full'
            : 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-10'
        }`}
      >
        {/* VIEW 1: Reader Mode */}
        {selectedArticle ? (
          <ArticleReader
            article={selectedArticle}
            onBack={handleBackToArticles}
            isBookmarked={bookmarks.includes(selectedArticle.id)}
            onToggleBookmark={handleToggleBookmark}
            onSelectAuthor={(authorId) => setSelectedAuthorId(authorId)}
            comments={currentArticleComments}
            onAddComment={handleAddComment}
            onFocusModeChange={setIsFocusMode}
          />
        ) : activeView === 'playground' ? (
          /* VIEW 2: Interactive Engineering Playground */
          <div className="space-y-6">
            <div className="mb-2">
              <button
                onClick={() => setActiveView('home')}
                className="text-xs font-semibold text-stone-500 hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
              >
                &larr; Voltar para a Capa
              </button>
            </div>
            <BenchmarkPlayground />
          </div>
        ) : activeView === 'about' ? (
          /* VIEW 3: About Editorial Manifesto */
          <AboutView
            onSelectAuthor={(authorId) => setSelectedAuthorId(authorId)}
            onNavigateHome={() => setActiveView('home')}
          />
        ) : (
          /* VIEW 4: Frontpage (Home & Saved Articles) */
          <div className="space-y-10 sm:space-y-12">
            {/* Header Hero Title for Saved View */}
            {activeView === 'saved' ? (
              <div className="border-b border-stone-200 dark:border-stone-800 pb-6">
                <div className="flex items-center gap-2 text-xs font-mono uppercase text-amber-700 dark:text-amber-400 mb-2">
                  <Bookmark className="w-4 h-4 fill-current" />
                  <span>Sua Biblioteca Pessoal</span>
                </div>
                <h1 className="text-3xl sm:text-4xl font-serif font-bold text-stone-900 dark:text-stone-100">
                  Artigos Salvos ({filteredArticles.length})
                </h1>
                <p className="mt-2 text-sm text-stone-600 dark:text-stone-400">
                  Publicações marcadas para leitura aprofundada posterior. Sincronizadas no seu navegador.
                </p>
              </div>
            ) : (
              /* Editorial Marquee Sub-Header */
              <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 border-b border-stone-200 dark:border-stone-800 pb-5">
                <div>
                  <div className="text-[11px] font-mono tracking-widest uppercase text-stone-400 dark:text-stone-500 mb-1">
                    EDIÇÃO VOL. IV · PRIMAVERA 2026 · SINTAXE.TECH
                  </div>
                  <h1 className="text-2xl sm:text-3xl lg:text-4xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight">
                    O Estado da Arte em Computação e Sistemas
                  </h1>
                </div>

                <div className="flex items-center gap-3 text-xs text-stone-500 font-mono">
                  <span>Curadoria Independente</span>
                  <span aria-hidden="true">·</span>
                  <span>6 Ensaios Ativos</span>
                </div>
              </div>
            )}

            {/* Category Navigation & Search Bar */}
            <CategoryNav
              categories={categoriesList}
              categoryCounts={categoryCounts}
              selectedCategory={selectedCategory}
              onSelectCategory={(cat) => setSelectedCategory(cat)}
              searchQuery={searchQuery}
              onSearchChange={(q) => setSearchQuery(q)}
              maxReadTime={maxReadTime}
              onMaxReadTimeChange={(t) => setMaxReadTime(t)}
              showSearchInput={searchOpen}
            />

            {/* Empty State when search or bookmarks yield nothing */}
            {filteredArticles.length === 0 ? (
              <div className="py-20 text-center rounded-2xl bg-white dark:bg-stone-900/60 border border-stone-200 dark:border-stone-800 p-8">
                <div className="w-12 h-12 rounded-full bg-stone-100 dark:bg-stone-800 flex items-center justify-center mx-auto text-stone-400 mb-4">
                  <FilterX className="w-6 h-6" />
                </div>
                <h3 className="text-lg font-serif font-bold text-stone-900 dark:text-stone-100">
                  {activeView === 'saved' ? 'Nenhum artigo salvo ainda' : 'Nenhuma publicação encontrada'}
                </h3>
                <p className="mt-2 text-xs sm:text-sm text-stone-500 dark:text-stone-400 max-w-sm mx-auto">
                  {activeView === 'saved'
                    ? 'Explore nossa edição e clique no ícone de marcador para salvar ensaios aqui.'
                    : `Não encontramos resultados para "${searchQuery}" com os filtros selecionados.`}
                </p>
                <div className="mt-6">
                  {activeView === 'saved' ? (
                    <button
                      onClick={() => setActiveView('home')}
                      className="px-4 py-2 text-xs font-semibold rounded-lg bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 cursor-pointer"
                    >
                      Explorar Todos os Artigos
                    </button>
                  ) : (
                    <button
                      onClick={() => {
                        setSearchQuery('');
                        setSelectedCategory('Todos');
                        setMaxReadTime(null);
                      }}
                      className="px-4 py-2 text-xs font-semibold rounded-lg border border-stone-300 dark:border-stone-700 hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
                    >
                      Limpar Filtros de Busca
                    </button>
                  )}
                </div>
              </div>
            ) : (
              /* Articles Feed */
              <div className="space-y-12">
                {/* TIER 1: Lead Story */}
                {leadArticle && (
                  <section aria-label="História Principal">
                    <HeroLeadArticle
                      article={leadArticle}
                      onSelectArticle={handleSelectArticle}
                      isBookmarked={bookmarks.includes(leadArticle.id)}
                      onToggleBookmark={handleToggleBookmark}
                      onSelectAuthor={(id) => setSelectedAuthorId(id)}
                    />
                  </section>
                )}

                {/* TIER 2: Secondary Featured Stories (2 Column Split) */}
                {secondaryArticles.length > 0 && (
                  <section aria-label="Destaques Secundários">
                    <div className="flex items-center gap-2 text-xs font-mono font-semibold uppercase tracking-wider text-stone-500 mb-4">
                      <Sparkles className="w-3.5 h-3.5 text-amber-600" />
                      <span>Análises em Destaque</span>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      {secondaryArticles.map((article) => (
                        <ArticleCard
                          key={article.id}
                          article={article}
                          variant="featured"
                          onSelectArticle={handleSelectArticle}
                          isBookmarked={bookmarks.includes(article.id)}
                          onToggleBookmark={handleToggleBookmark}
                          onSelectAuthor={(id) => setSelectedAuthorId(id)}
                        />
                      ))}
                    </div>
                  </section>
                )}

                {/* Mid-Page Interactive Promotion Banner for Playground */}
                {activeView === 'home' && !isFiltering && (
                  <div className="rounded-2xl bg-gradient-to-r from-stone-900 to-stone-800 text-stone-100 p-6 sm:p-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-6 shadow-md border border-stone-800">
                    <div className="space-y-2 max-w-xl">
                      <div className="flex items-center gap-2 text-xs font-mono uppercase text-amber-400">
                        <Newspaper className="w-4 h-4" />
                        <span>Laboratório Técnico Sintaxe</span>
                      </div>
                      <h3 className="text-xl sm:text-2xl font-serif font-bold text-white tracking-tight">
                        Experimente nosso Simulador de Runtimes em Alta Concorrência
                      </h3>
                      <p className="text-xs sm:text-sm text-stone-300 leading-relaxed">
                        Compare a curva de latência p99 e uso de RAM entre Rust, Go, Node e Python sob cargas dinâmicas de 1k a 100k requisições por segundo.
                      </p>
                    </div>

                    <button
                      onClick={() => {
                        setActiveView('playground');
                        window.scrollTo({ top: 0, behavior: 'smooth' });
                      }}
                      className="px-4 py-2.5 rounded-xl bg-amber-500 text-stone-950 font-semibold text-xs sm:text-sm hover:bg-amber-400 transition-colors whitespace-nowrap shrink-0 flex items-center gap-1.5 cursor-pointer shadow-xs"
                    >
                      <span>Abrir Simulador</span>
                      <ArrowRight className="w-4 h-4" />
                    </button>
                  </div>
                )}

                {/* TIER 3: Standard Editorial Grid */}
                {gridArticles.length > 0 && (
                  <section aria-label="Caderno de Publicações">
                    <div className="flex items-center justify-between mb-4 border-b border-stone-200 dark:border-stone-800 pb-2">
                      <h3 className="text-sm font-mono font-bold uppercase tracking-wider text-stone-500">
                        Demais Ensaios & Despachos ({gridArticles.length})
                      </h3>
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                      {gridArticles.map((article) => (
                        <ArticleCard
                          key={article.id}
                          article={article}
                          variant="standard"
                          onSelectArticle={handleSelectArticle}
                          isBookmarked={bookmarks.includes(article.id)}
                          onToggleBookmark={handleToggleBookmark}
                          onSelectAuthor={(id) => setSelectedAuthorId(id)}
                        />
                      ))}
                    </div>
                  </section>
                )}
              </div>
            )}
          </div>
        )}
      </main>

      {/* Author Details Modal */}
      <AuthorModal
        author={currentAuthor}
        onClose={() => setSelectedAuthorId(null)}
        articles={articles}
        onSelectArticle={handleSelectArticle}
      />

      {/* Newsletter Subscribe Dialog */}
      <NewsletterModal
        isOpen={newsletterModalOpen}
        onClose={() => setNewsletterModalOpen(false)}
      />

      {/* Footer (Hidden when Focus Mode is active) */}
      {!(selectedArticle && isFocusMode) && (
        <Footer
          onSelectCategory={(cat) => {
            setSelectedCategory(cat);
            setActiveView('home');
            setSelectedArticle(null);
            setIsFocusMode(false);
            window.scrollTo({ top: 0, behavior: 'smooth' });
          }}
          onNavigate={(view) => {
            setActiveView(view);
            setSelectedArticle(null);
            setIsFocusMode(false);
            window.scrollTo({ top: 0, behavior: 'smooth' });
          }}
          onOpenNewsletter={() => setNewsletterModalOpen(true)}
        />
      )}
    </div>
  );
}
