import React, { useState, useEffect } from 'react';
import { Article, Comment } from '../types/blog';
import {
  formatPublicationDate,
  calculateRemainingReadingTime,
  calculateScrollProgress,
  buildArticleShareData,
} from '../utils/blogLogic';
import {
  ArrowLeft,
  Bookmark,
  Share2,
  ThumbsUp,
  Volume2,
  VolumeX,
  Copy,
  Check,
  MessageSquare,
  Send,
  Type,
  Sparkles,
  Maximize2,
  Minimize2,
  Sun,
  Moon,
  Coffee,
} from 'lucide-react';

interface ArticleReaderProps {
  article: Article;
  onBack: () => void;
  isBookmarked: boolean;
  onToggleBookmark: (id: string) => void;
  onSelectAuthor: (authorId: string) => void;
  comments: Comment[];
  onAddComment: (comment: { authorName: string; authorRole: string; content: string }) => void;
  onFocusModeChange?: (active: boolean) => void;
}

export type FocusTone = 'default' | 'sepia' | 'dark';

export const ArticleReader: React.FC<ArticleReaderProps> = ({
  article,
  onBack,
  isBookmarked,
  onToggleBookmark,
  onSelectAuthor,
  comments,
  onAddComment,
  onFocusModeChange,
}) => {
  const [isFocusMode, setIsFocusMode] = useState(false);
  const [focusTone, setFocusTone] = useState<FocusTone>('default');
  const [scrollProgress, setScrollProgress] = useState(0);
  const [claps, setClaps] = useState(article.claps);
  const [hasClapped, setHasClapped] = useState(false);
  const [copiedCodeId, setCopiedCodeId] = useState<string | null>(null);
  const [shareToast, setShareToast] = useState(false);
  const [fontSize, setFontSize] = useState<'normal' | 'large' | 'xl'>('normal');

  // Simulated Audio Synthesizer Player state
  const [isPlayingAudio, setIsPlayingAudio] = useState(false);
  const [audioProgress, setAudioProgress] = useState(0);

  // New Comment Form state
  const [commentName, setCommentName] = useState('');
  const [commentRole, setCommentRole] = useState('');
  const [commentContent, setCommentContent] = useState('');
  const [commentSuccess, setCommentSuccess] = useState(false);

  // Notify parent component about Focus Mode state (e.g. to hide site header & footer)
  useEffect(() => {
    onFocusModeChange?.(isFocusMode);
  }, [isFocusMode, onFocusModeChange]);

  // Keyboard shortcut listener: ESC exits Focus Mode
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isFocusMode) {
        setIsFocusMode(false);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isFocusMode]);

  // Track scroll progress for fixed reading depth bar
  useEffect(() => {
    const handleScroll = () => {
      const totalHeight = document.documentElement.scrollHeight - window.innerHeight;
      setScrollProgress(calculateScrollProgress(window.scrollY, totalHeight));
    };
    handleScroll();
    window.addEventListener('scroll', handleScroll, { passive: true });
    window.addEventListener('resize', handleScroll, { passive: true });
    return () => {
      window.removeEventListener('scroll', handleScroll);
      window.removeEventListener('resize', handleScroll);
    };
  }, []);

  // Audio timer simulation
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isPlayingAudio) {
      interval = setInterval(() => {
        setAudioProgress((prev) => {
          if (prev >= 100) {
            setIsPlayingAudio(false);
            return 0;
          }
          return prev + 1;
        });
      }, 500);
    }
    return () => clearInterval(interval);
  }, [isPlayingAudio]);

  const handleClap = () => {
    setClaps((prev) => prev + 1);
    setHasClapped(true);
  };

  const handleCopyCode = (code: string, id: string) => {
    navigator.clipboard.writeText(code);
    setCopiedCodeId(id);
    setTimeout(() => setCopiedCodeId(null), 2000);
  };

  const handleShare = async () => {
    const shareUrl = typeof window !== 'undefined' ? window.location.href : '';
    const shareData = buildArticleShareData(article, shareUrl);

    // Call native Web Share API if supported by the browser
    if (typeof navigator !== 'undefined' && typeof navigator.share === 'function') {
      try {
        await navigator.share(shareData);
        return;
      } catch (err: unknown) {
        // AbortError indicates the user cancelled the share dialog
        if (err instanceof Error && err.name === 'AbortError') {
          return;
        }
      }
    }

    // Fallback: Copy current link to clipboard and show feedback toast
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      try {
        await navigator.clipboard.writeText(shareData.url);
        setShareToast(true);
        setTimeout(() => setShareToast(false), 3000);
      } catch {
        // Ignore clipboard failure
      }
    }
  };

  const handleSubmitComment = (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentName.trim() || !commentContent.trim()) return;

    onAddComment({
      authorName: commentName.trim(),
      authorRole: commentRole.trim() || 'Leitor Sintaxe',
      content: commentContent.trim(),
    });

    setCommentName('');
    setCommentRole('');
    setCommentContent('');
    setCommentSuccess(true);
    setTimeout(() => setCommentSuccess(false), 4000);
  };

  const fontSizeClasses = {
    normal: 'text-base sm:text-lg leading-relaxed sm:leading-loose',
    large: 'text-lg sm:text-xl leading-relaxed sm:leading-loose',
    xl: 'text-xl sm:text-2xl leading-loose',
  };

  // Tone styling for Focus Mode (Default, Sepia, Dark)
  const focusToneContainerClasses = {
    default: 'bg-white dark:bg-stone-950 text-stone-900 dark:text-stone-100',
    sepia: 'bg-[#FBF7EE] text-[#2C2416]',
    dark: 'bg-[#0F0E0D] text-[#ECE7DF]',
  };

  const remainingMinutes = calculateRemainingReadingTime(article.readTimeMinutes, scrollProgress);

  return (
    <article
      className={`min-h-screen transition-colors duration-300 ${
        isFocusMode ? `${focusToneContainerClasses[focusTone]} py-12 sm:py-16 px-4 sm:px-6` : 'pb-24'
      }`}
    >
      {/* Fixed Top Reading Progress Bar (Visible only when ArticleReader is mounted) */}
      <div
        role="progressbar"
        aria-valuenow={Math.round(scrollProgress)}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label="Progresso de leitura do artigo"
        data-testid="fixed-reading-progress-bar"
        className="fixed top-0 left-0 right-0 z-[60] h-1 sm:h-1.5 bg-stone-200/50 dark:bg-stone-800/50 backdrop-blur-xs pointer-events-none"
      >
        <div
          className="h-full bg-gradient-to-r from-amber-600 via-amber-500 to-amber-400 dark:from-amber-500 dark:via-amber-400 dark:to-amber-300 shadow-[0_1px_8px_rgba(217,119,6,0.6)] transition-all duration-150 ease-out"
          style={{ width: `${scrollProgress}%` }}
        />
      </div>

      {/* Share Toast */}
      {shareToast && (
        <div className="fixed bottom-6 right-6 z-50 flex items-center gap-2 bg-stone-900 text-stone-100 dark:bg-stone-100 dark:text-stone-900 px-4 py-2.5 rounded-xl shadow-xl text-xs font-medium animate-fade-in">
          <Check className="w-4 h-4 text-emerald-500" />
          <span>Link do artigo copiado com sucesso!</span>
        </div>
      )}

      {/* ========================================================================= */}
      {/* FOCUS MODE: FLOATING MINIMAL HUD CONTROLS */}
      {/* ========================================================================= */}
      {isFocusMode ? (
        <aside
          aria-label="Controles do Modo de Foco"
          className="fixed top-4 right-4 sm:top-6 sm:right-6 z-50 flex items-center gap-2 p-1.5 rounded-2xl bg-white/95 dark:bg-stone-900/95 border border-stone-200/80 dark:border-stone-800 shadow-xl backdrop-blur-md text-xs"
        >
          {/* Progress & Remaining Time */}
          <div className="hidden md:flex items-center gap-2 px-3 py-1 font-mono text-[11px] text-stone-500 dark:text-stone-400 border-r border-stone-200 dark:border-stone-800">
            <span>{Math.round(scrollProgress)}% lido</span>
            <span>·</span>
            <span>{remainingMinutes}m restantes</span>
          </div>

          {/* Tone switch for focus mode (Default / Sepia / Dark) */}
          <div className="flex items-center gap-0.5 bg-stone-100 dark:bg-stone-800 p-0.5 rounded-lg border border-stone-200/70 dark:border-stone-700/70">
            <button
              onClick={() => setFocusTone('default')}
              className={`p-1.5 rounded-md text-stone-600 dark:text-stone-300 hover:text-stone-900 dark:hover:text-stone-100 cursor-pointer ${
                focusTone === 'default' ? 'bg-white dark:bg-stone-700 shadow-xs font-bold' : ''
              }`}
              title="Tema Padrão"
              aria-label="Tom de Leitura Padrão"
            >
              <Sun className="w-3.5 h-3.5" />
            </button>
            <button
              onClick={() => setFocusTone('sepia')}
              className={`p-1.5 rounded-md text-[#614E32] hover:text-[#2C2416] cursor-pointer ${
                focusTone === 'sepia' ? 'bg-[#F4ECDC] shadow-xs font-bold' : ''
              }`}
              title="Tema Sépia / Papel"
              aria-label="Tom de Leitura Sépia Papel"
            >
              <Coffee className="w-3.5 h-3.5" />
            </button>
            <button
              onClick={() => setFocusTone('dark')}
              className={`p-1.5 rounded-md text-stone-600 dark:text-stone-300 hover:text-stone-900 dark:hover:text-stone-100 cursor-pointer ${
                focusTone === 'dark' ? 'bg-stone-900 text-stone-100 shadow-xs font-bold' : ''
              }`}
              title="Tema Noturno"
              aria-label="Tom de Leitura Noturno"
            >
              <Moon className="w-3.5 h-3.5" />
            </button>
          </div>

          {/* Font Size Selector */}
          <div className="flex items-center gap-0.5 bg-stone-100 dark:bg-stone-800 p-0.5 rounded-lg border border-stone-200/70 dark:border-stone-700/70">
            <button
              onClick={() => setFontSize('normal')}
              className={`px-2 py-1 rounded text-[11px] font-medium cursor-pointer ${
                fontSize === 'normal' ? 'bg-white dark:bg-stone-700 text-stone-900 dark:text-stone-100 shadow-xs' : 'text-stone-500'
              }`}
            >
              A
            </button>
            <button
              onClick={() => setFontSize('large')}
              className={`px-2 py-1 rounded text-[11px] font-medium cursor-pointer ${
                fontSize === 'large' ? 'bg-white dark:bg-stone-700 text-stone-900 dark:text-stone-100 shadow-xs' : 'text-stone-500'
              }`}
            >
              A+
            </button>
            <button
              onClick={() => setFontSize('xl')}
              className={`px-2 py-1 rounded text-[11px] font-medium cursor-pointer ${
                fontSize === 'xl' ? 'bg-white dark:bg-stone-700 text-stone-900 dark:text-stone-100 shadow-xs' : 'text-stone-500'
              }`}
            >
              A++
            </button>
          </div>

          {/* Quick Share in Focus Mode */}
          <button
            onClick={handleShare}
            data-testid="share-article-button-focus"
            className="p-1.5 rounded-lg text-stone-600 dark:text-stone-300 hover:text-stone-900 dark:hover:text-stone-100 bg-stone-100 dark:bg-stone-800 hover:bg-stone-200 dark:hover:bg-stone-700 transition-colors cursor-pointer"
            title="Compartilhar artigo (Web Share API)"
            aria-label="Compartilhar artigo no modo foco"
          >
            <Share2 className="w-3.5 h-3.5" />
          </button>

          {/* Exit Focus Mode Button */}
          <button
            onClick={() => setIsFocusMode(false)}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 font-medium hover:opacity-90 transition-opacity cursor-pointer shadow-xs whitespace-nowrap"
            title="Sair do Modo de Foco (Esc)"
          >
            <Minimize2 className="w-3.5 h-3.5" />
            <span className="hidden sm:inline">Sair do Foco</span>
            <kbd className="hidden lg:inline text-[10px] font-mono opacity-70 bg-stone-800 dark:bg-stone-200 px-1 py-0.2 rounded">
              Esc
            </kbd>
          </button>
        </aside>
      ) : null}

      {/* ========================================================================= */}
      {/* STANDARD NAVIGATION & HEADER (Hidden in Focus Mode) */}
      {/* ========================================================================= */}
      {!isFocusMode && (
        <div className="max-w-4xl mx-auto px-4 sm:px-6 pt-6 sm:pt-10">
          <div className="flex items-center justify-between mb-8">
            <button
              onClick={onBack}
              className="inline-flex items-center gap-2 text-xs sm:text-sm font-medium text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100 transition-colors cursor-pointer group"
            >
              <ArrowLeft className="w-4 h-4 transition-transform group-hover:-translate-x-1" />
              <span>Voltar para todas as publicações</span>
            </button>

            {/* Reading Tools (Focus Mode trigger, Font size, Bookmark, Share) */}
            <div className="flex items-center gap-2">
              {/* TRIGGER FOCUS MODE */}
              <button
                onClick={() => setIsFocusMode(true)}
                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg bg-amber-50 text-amber-900 border border-amber-200 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-900/60 hover:bg-amber-100 dark:hover:bg-amber-900/40 transition-colors cursor-pointer shadow-2xs"
                title="Ativar Modo de Leitura Imersivo (oculta barras e margens)"
                aria-label="Ativar Modo de Foco"
              >
                <Maximize2 className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">Modo de Foco</span>
              </button>

              {/* Font Size Toggle */}
              <div className="flex items-center gap-1 bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 rounded-lg p-0.5 text-xs text-stone-600 dark:text-stone-300">
                <span className="p-1 text-stone-400" title="Tamanho do texto">
                  <Type className="w-3.5 h-3.5" />
                </span>
                <button
                  onClick={() => setFontSize('normal')}
                  className={`px-2 py-0.5 rounded text-[11px] font-medium cursor-pointer ${
                    fontSize === 'normal' ? 'bg-stone-100 dark:bg-stone-800 text-stone-900 dark:text-stone-100' : ''
                  }`}
                >
                  A
                </button>
                <button
                  onClick={() => setFontSize('large')}
                  className={`px-2 py-0.5 rounded text-[11px] font-medium cursor-pointer ${
                    fontSize === 'large' ? 'bg-stone-100 dark:bg-stone-800 text-stone-900 dark:text-stone-100' : ''
                  }`}
                >
                  A+
                </button>
                <button
                  onClick={() => setFontSize('xl')}
                  className={`px-2 py-0.5 rounded text-[11px] font-medium cursor-pointer ${
                    fontSize === 'xl' ? 'bg-stone-100 dark:bg-stone-800 text-stone-900 dark:text-stone-100' : ''
                  }`}
                >
                  A++
                </button>
              </div>

              <button
                onClick={() => onToggleBookmark(article.id)}
                className={`p-2 rounded-lg border transition-colors cursor-pointer ${
                  isBookmarked
                    ? 'border-amber-600 bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400'
                    : 'border-stone-200 dark:border-stone-800 text-stone-500 hover:text-stone-900 dark:hover:text-stone-100 bg-white dark:bg-stone-900'
                }`}
                title={isBookmarked ? 'Salvo nos favoritos' : 'Salvar artigo'}
              >
                <Bookmark className={`w-4 h-4 ${isBookmarked ? 'fill-current' : ''}`} />
              </button>

              <button
                onClick={handleShare}
                data-testid="share-article-button"
                aria-label="Compartilhar artigo via janela nativa do navegador"
                className="p-2 rounded-lg border border-stone-200 dark:border-stone-800 text-stone-500 hover:text-stone-900 dark:hover:text-stone-100 bg-white dark:bg-stone-900 transition-colors cursor-pointer hover:border-amber-500/50"
                title="Compartilhar artigo (Web Share API)"
              >
                <Share2 className="w-4 h-4" />
              </button>
            </div>
          </div>

          {/* Unboxed Metadata (Section 1.A) */}
          <div className="flex flex-wrap items-center gap-2 text-xs font-mono tracking-widest uppercase text-stone-500 dark:text-stone-400 mb-4">
            <span className="text-amber-700 dark:text-amber-400 font-semibold">{article.category}</span>
            <span aria-hidden="true">·</span>
            <span>Publicado em {formatPublicationDate(article.publishedAt)}</span>
            <span aria-hidden="true">·</span>
            <span>{article.readTimeMinutes} min de leitura</span>
          </div>

          {/* Article Headline */}
          <h1 className="text-3xl sm:text-4xl lg:text-5xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight leading-[1.15] text-balance">
            {article.title}
          </h1>

          {/* Subtitle */}
          <p className="mt-4 text-lg sm:text-xl text-stone-600 dark:text-stone-300 font-light leading-relaxed">
            {article.subtitle}
          </p>

          {/* Author Byline & Audio Player Bar */}
          <div className="mt-8 pt-6 pb-6 border-y border-stone-200 dark:border-stone-800 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <button
              onClick={() => onSelectAuthor(article.author.id)}
              className="flex items-center gap-3.5 text-left group/author cursor-pointer"
            >
              <img
                src={article.author.avatar}
                alt={article.author.name}
                referrerPolicy="no-referrer"
                className="w-12 h-12 rounded-full object-cover ring-2 ring-stone-200 dark:ring-stone-700"
              />
              <div>
                <p className="text-sm font-semibold text-stone-900 dark:text-stone-100 group-hover/author:text-amber-700 dark:group-hover/author:text-amber-400 transition-colors">
                  {article.author.name}
                </p>
                <p className="text-xs text-stone-500 dark:text-stone-400">
                  {article.author.role}
                </p>
              </div>
            </button>

            {/* Interactive Audio Player widget */}
            <div className="flex items-center gap-3 bg-stone-100 dark:bg-stone-900 px-3.5 py-2 rounded-xl border border-stone-200 dark:border-stone-800">
              <button
                onClick={() => setIsPlayingAudio(!isPlayingAudio)}
                className="p-1.5 rounded-lg bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 hover:opacity-90 transition-opacity cursor-pointer"
                aria-label={isPlayingAudio ? 'Pausar áudio' : 'Ouvir artigo'}
              >
                {isPlayingAudio ? <VolumeX className="w-3.5 h-3.5" /> : <Volume2 className="w-3.5 h-3.5" />}
              </button>
              <div className="flex flex-col text-left">
                <span className="text-[11px] font-medium text-stone-800 dark:text-stone-200">
                  {isPlayingAudio ? 'Sintetizando Narração...' : 'Ouvir Versão em Áudio'}
                </span>
                <div className="flex items-center gap-1.5 mt-1">
                  <div className="w-24 bg-stone-200 dark:bg-stone-800 h-1.5 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-amber-600 transition-all duration-300"
                      style={{ width: `${audioProgress}%` }}
                    />
                  </div>
                  <span className="text-[10px] font-mono text-stone-500">{article.readTimeMinutes}m</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ========================================================================= */}
      {/* FOCUS MODE: MINIMAL HEADER (Only in Focus Mode) */}
      {/* ========================================================================= */}
      {isFocusMode && (
        <div className="max-w-2xl mx-auto mb-10 text-center">
          <div className="inline-block text-[11px] font-mono tracking-widest uppercase text-amber-700 dark:text-amber-400 mb-3 font-semibold">
            Modo de Leitura Imersivo · {article.category}
          </div>
          <h1 className="text-3xl sm:text-4xl lg:text-5xl font-serif font-bold tracking-tight leading-tight text-balance">
            {article.title}
          </h1>
          <p className="mt-3 text-base sm:text-lg opacity-80 leading-relaxed font-light">
            {article.subtitle}
          </p>
          <div className="mt-4 flex items-center justify-center gap-3 text-xs opacity-70 font-mono">
            <span>{article.author.name}</span>
            <span>·</span>
            <span>{article.readTimeMinutes} min</span>
          </div>
          <div className="mt-8 border-b border-stone-200/40 dark:border-stone-800/60 max-w-xs mx-auto" />
        </div>
      )}

      {/* Featured Cover Image (Shown in Standard Mode only to keep Focus Mode 100% distraction-free) */}
      {!isFocusMode && (
        <div className="max-w-4xl mx-auto px-4 sm:px-6 my-8">
          <div className="rounded-2xl overflow-hidden aspect-[16/9] bg-stone-100 dark:bg-stone-900 border border-stone-200 dark:border-stone-800">
            <img
              src={article.coverImage}
              alt={article.title}
              referrerPolicy="no-referrer"
              className="w-full h-full object-cover"
            />
          </div>
          <p className="text-xs font-serif italic text-stone-500 dark:text-stone-400 mt-2 text-center">
            Figura 1 — Ensaio visual editorial dedicado a &ldquo;{article.title}&rdquo;
          </p>
        </div>
      )}

      {/* ========================================================================= */}
      {/* ARTICLE CONTENT BODY (Constrained measure 65-75ch) */}
      {/* ========================================================================= */}
      <div
        className={`${
          isFocusMode ? 'max-w-2xl' : 'max-w-3xl'
        } mx-auto px-4 sm:px-6 transition-all duration-300`}
      >
        {/* Lead Paragraph with Elegant Initial Drop Cap */}
        <p
          className={`${fontSizeClasses[fontSize]} mb-8 first-letter:text-5xl first-letter:font-serif first-letter:font-bold first-letter:float-left first-letter:mr-3.5 first-letter:mt-1 first-letter:text-amber-800 dark:first-letter:text-amber-400`}
        >
          {article.leadParagraph}
        </p>

        {/* Dynamic Content Blocks */}
        {article.blocks.map((block, idx) => {
          switch (block.type) {
            case 'takeaways':
              return (
                <div
                  key={idx}
                  className={`my-8 p-6 rounded-xl border ${
                    focusTone === 'sepia'
                      ? 'bg-[#F2ECE0] border-[#E0D4C0]'
                      : 'bg-amber-50/70 dark:bg-amber-950/20 border-amber-200/80 dark:border-amber-900/50'
                  }`}
                >
                  <div className="flex items-center gap-2 text-xs font-mono font-bold tracking-wider uppercase text-amber-800 dark:text-amber-400 mb-3">
                    <Sparkles className="w-4 h-4 text-amber-600" />
                    <span>Conclusões Fundamentais & Invariants</span>
                  </div>
                  <ul className="space-y-2.5">
                    {block.items?.map((item, itemIdx) => (
                      <li key={itemIdx} className="flex items-start gap-2.5 text-sm leading-relaxed">
                        <span className="text-amber-700 dark:text-amber-400 font-mono font-bold text-xs mt-0.5">
                          0{itemIdx + 1}.
                        </span>
                        <span>{item}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              );

            case 'heading2':
              return (
                <h2
                  key={idx}
                  className="text-2xl sm:text-3xl font-serif font-bold tracking-tight mt-12 mb-4 pt-6 border-t border-stone-200/60 dark:border-stone-800/60"
                >
                  {block.text}
                </h2>
              );

            case 'heading3':
              return (
                <h3
                  key={idx}
                  className="text-xl sm:text-2xl font-serif font-semibold tracking-tight mt-8 mb-3"
                >
                  {block.text}
                </h3>
              );

            case 'paragraph':
              return (
                <p key={idx} className={`${fontSizeClasses[fontSize]} mb-6`}>
                  {block.text}
                </p>
              );

            case 'quote':
              return (
                <figure
                  key={idx}
                  className="my-10 py-6 px-6 border-y border-stone-200/70 dark:border-stone-800/70 text-center"
                >
                  <blockquote className="text-xl sm:text-2xl font-serif italic leading-snug">
                    &ldquo;{block.text}&rdquo;
                  </blockquote>
                  {block.authorQuote && (
                    <figcaption className="mt-3 text-xs font-mono uppercase tracking-widest opacity-70">
                      — {block.authorQuote}
                    </figcaption>
                  )}
                </figure>
              );

            case 'code':
              return (
                <div key={idx} className="my-8 rounded-xl overflow-hidden border border-stone-200 dark:border-stone-800 shadow-xs">
                  <div className="flex items-center justify-between px-4 py-2 bg-stone-100 dark:bg-stone-900 border-b border-stone-200 dark:border-stone-800 text-xs font-mono text-stone-600 dark:text-stone-400">
                    <span>{block.language || 'code'}</span>
                    <button
                      onClick={() => handleCopyCode(block.code || '', `code-${idx}`)}
                      className="flex items-center gap-1.5 px-2 py-0.5 rounded text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100 hover:bg-stone-200 dark:hover:bg-stone-800 transition-colors cursor-pointer"
                    >
                      {copiedCodeId === `code-${idx}` ? (
                        <>
                          <Check className="w-3.5 h-3.5 text-emerald-500" />
                          <span className="text-emerald-600 dark:text-emerald-400 font-sans">Copiado</span>
                        </>
                      ) : (
                        <>
                          <Copy className="w-3.5 h-3.5" />
                          <span className="font-sans">Copiar</span>
                        </>
                      )}
                    </button>
                  </div>
                  <pre className="p-4 sm:p-5 overflow-x-auto bg-stone-950 text-stone-200 text-xs sm:text-sm font-mono leading-relaxed">
                    <code>{block.code}</code>
                  </pre>
                  {block.caption && (
                    <div className="px-4 py-2 bg-stone-50 dark:bg-stone-900/60 text-[11px] font-serif italic opacity-70 border-t border-stone-200 dark:border-stone-800">
                      {block.caption}
                    </div>
                  )}
                </div>
              );

            case 'interactive-benchmark':
              if (!block.benchmarkData) return null;
              return (
                <div
                  key={idx}
                  className="my-8 p-6 rounded-xl bg-stone-50 dark:bg-stone-900/80 border border-stone-200 dark:border-stone-800"
                >
                  <h4 className="text-sm font-serif font-bold">{block.benchmarkData.title}</h4>
                  <p className="text-xs opacity-70 mt-1 mb-5">{block.benchmarkData.description}</p>
                  <div className="space-y-4">
                    {block.benchmarkData.metrics.map((m, mIdx) => (
                      <div key={mIdx}>
                        <div className="flex justify-between text-xs font-mono mb-1.5">
                          <span className="font-semibold">{m.label}</span>
                          <span className="opacity-70">
                            Zero-Shot: <strong>{m.valueA}{m.unit}</strong> | Raciocínio: <strong className="text-amber-700 dark:text-amber-400">{m.valueB}{m.unit}</strong>
                          </span>
                        </div>
                        <div className="grid grid-cols-2 gap-2 h-2">
                          <div className="bg-stone-200 dark:bg-stone-800 rounded-full overflow-hidden">
                            <div className="h-full bg-stone-500" style={{ width: `${Math.min(100, m.valueA)}%` }} />
                          </div>
                          <div className="bg-stone-200 dark:bg-stone-800 rounded-full overflow-hidden">
                            <div className="h-full bg-amber-600" style={{ width: `${Math.min(100, m.valueB)}%` }} />
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              );

            default:
              return null;
          }
        })}

        {/* End of Focus Reading Mode footer */}
        {isFocusMode ? (
          <div className="mt-16 pt-8 border-t border-stone-200/50 dark:border-stone-800/60 text-center space-y-4">
            <p className="font-serif italic text-sm opacity-70">
              Você concluiu a leitura deste ensaio no Modo de Foco.
            </p>
            <div className="flex items-center justify-center gap-3">
              <button
                onClick={() => setIsFocusMode(false)}
                className="px-4 py-2 rounded-xl bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 text-xs font-semibold hover:opacity-90 transition-opacity cursor-pointer"
              >
                Retornar à Visualização Completa
              </button>
              <button
                onClick={onBack}
                className="px-4 py-2 rounded-xl border border-stone-300 dark:border-stone-700 text-xs font-semibold hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
              >
                Voltar aos Artigos
              </button>
            </div>
          </div>
        ) : (
          /* ========================================================================= */
          /* STANDARD MODE: ENGAGEMENT, AUTHOR BIO & COMMENTS */
          /* ========================================================================= */
          <>
            {/* Tags */}
            <div className="mt-12 pt-6 border-t border-stone-200 dark:border-stone-800 flex flex-wrap items-center gap-2">
              <span className="text-xs font-mono text-stone-500 uppercase">Tópicos:</span>
              {article.tags.map((tag) => (
                <span
                  key={tag}
                  className="text-xs font-medium px-2.5 py-1 rounded-md bg-stone-100 text-stone-700 dark:bg-stone-800 dark:text-stone-300"
                >
                  #{tag}
                </span>
              ))}
            </div>

            {/* Reaction Bar (Claps & Share) */}
            <div className="mt-8 p-6 rounded-2xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 flex items-center justify-between gap-4">
              <div className="flex items-center gap-4">
                <button
                  onClick={handleClap}
                  className={`flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-semibold transition-all cursor-pointer ${
                    hasClapped
                      ? 'bg-amber-100 text-amber-900 dark:bg-amber-950/60 dark:text-amber-300 scale-105'
                      : 'bg-stone-100 hover:bg-stone-200 dark:bg-stone-800 dark:hover:bg-stone-700 text-stone-800 dark:text-stone-200'
                  }`}
                >
                  <ThumbsUp className={`w-4 h-4 ${hasClapped ? 'fill-current' : ''}`} />
                  <span className="font-mono tabular-nums">{claps} aplausos</span>
                </button>
                <span className="text-xs text-stone-500 font-mono">
                  {article.views.toLocaleString('pt-BR')} visualizações
                </span>
              </div>

              <div className="flex items-center gap-2">
                <button
                  onClick={handleShare}
                  data-testid="share-article-button-bottom"
                  aria-label="Compartilhar artigo via janela nativa do navegador"
                  className="p-2.5 rounded-lg border border-stone-200 dark:border-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-50 dark:hover:bg-stone-800 hover:border-amber-500/50 transition-colors cursor-pointer"
                  title="Compartilhar artigo (Web Share API)"
                >
                  <Share2 className="w-4 h-4" />
                </button>
                <button
                  onClick={() => onToggleBookmark(article.id)}
                  className={`p-2.5 rounded-lg border transition-colors cursor-pointer ${
                    isBookmarked
                      ? 'border-amber-600 bg-amber-50 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400'
                      : 'border-stone-200 dark:border-stone-800 text-stone-600 dark:text-stone-300 hover:bg-stone-50 dark:hover:bg-stone-800'
                  }`}
                  title={isBookmarked ? 'Salvo' : 'Salvar'}
                >
                  <Bookmark className={`w-4 h-4 ${isBookmarked ? 'fill-current' : ''}`} />
                </button>
              </div>
            </div>

            {/* Author Expanded Bio Box */}
            <div className="mt-8 p-6 rounded-2xl bg-stone-100/80 dark:bg-stone-900/60 border border-stone-200 dark:border-stone-800 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-6">
              <div className="flex items-start gap-4">
                <img
                  src={article.author.avatar}
                  alt={article.author.name}
                  referrerPolicy="no-referrer"
                  className="w-14 h-14 rounded-full object-cover ring-2 ring-stone-300 dark:ring-stone-700 shrink-0"
                />
                <div>
                  <h4 className="text-base font-serif font-bold text-stone-900 dark:text-stone-100">
                    {article.author.name}
                  </h4>
                  <p className="text-xs text-stone-500 dark:text-stone-400 mb-2">
                    {article.author.role}
                  </p>
                  <p className="text-xs text-stone-600 dark:text-stone-300 leading-relaxed max-w-lg">
                    {article.author.bio}
                  </p>
                </div>
              </div>
              <button
                onClick={() => onSelectAuthor(article.author.id)}
                className="px-4 py-2 text-xs font-semibold text-stone-900 dark:text-stone-100 border border-stone-300 dark:border-stone-700 rounded-lg hover:bg-stone-200/50 dark:hover:bg-stone-800 transition-colors whitespace-nowrap shrink-0 cursor-pointer"
              >
                Ver Outros Artigos
              </button>
            </div>

            {/* Comments Section */}
            <section className="mt-14 pt-10 border-t border-stone-200 dark:border-stone-800">
              <div className="flex items-center justify-between mb-8">
                <div className="flex items-center gap-2">
                  <MessageSquare className="w-5 h-5 text-amber-700 dark:text-amber-400" />
                  <h3 className="text-xl font-serif font-bold text-stone-900 dark:text-stone-100">
                    Discussão da Comunidade Técnica
                  </h3>
                </div>
                <span className="text-xs font-mono text-stone-500">
                  {comments.length} {comments.length === 1 ? 'comentário' : 'comentários'}
                </span>
              </div>

              {/* New Comment Form */}
              <form onSubmit={handleSubmitComment} className="mb-10 p-5 rounded-xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 shadow-xs">
                <h4 className="text-xs font-mono font-semibold uppercase tracking-wider text-stone-500 dark:text-stone-400 mb-3">
                  Deixe sua reflexão técnica
                </h4>

                {commentSuccess && (
                  <div className="mb-4 p-3 rounded-lg bg-emerald-50 text-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-300 text-xs font-medium flex items-center gap-2">
                    <Check className="w-4 h-4" />
                    <span>Seu comentário foi publicado com sucesso!</span>
                  </div>
                )}

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-3">
                  <input
                    type="text"
                    required
                    value={commentName}
                    onChange={(e) => setCommentName(e.target.value)}
                    placeholder="Seu Nome *"
                    className="px-3.5 py-2 text-xs rounded-lg border border-stone-200 dark:border-stone-700 bg-stone-50 dark:bg-stone-950 text-stone-900 dark:text-stone-100 placeholder-stone-400 focus:outline-none focus:ring-2 focus:ring-amber-700/30"
                  />
                  <input
                    type="text"
                    value={commentRole}
                    onChange={(e) => setCommentRole(e.target.value)}
                    placeholder="Cargo / Organização (ex: Engenheiro na TechCorp)"
                    className="px-3.5 py-2 text-xs rounded-lg border border-stone-200 dark:border-stone-700 bg-stone-50 dark:bg-stone-950 text-stone-900 dark:text-stone-100 placeholder-stone-400 focus:outline-none focus:ring-2 focus:ring-amber-700/30"
                  />
                </div>

                <textarea
                  required
                  rows={3}
                  value={commentContent}
                  onChange={(e) => setCommentContent(e.target.value)}
                  placeholder="Contribua com suas observações, benchmarks ou contrapontos..."
                  className="w-full px-3.5 py-2 text-xs rounded-lg border border-stone-200 dark:border-stone-700 bg-stone-50 dark:bg-stone-950 text-stone-900 dark:text-stone-100 placeholder-stone-400 focus:outline-none focus:ring-2 focus:ring-amber-700/30 mb-3"
                />

                <div className="flex justify-end">
                  <button
                    type="submit"
                    className="inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-stone-50 bg-stone-900 dark:bg-stone-100 dark:text-stone-900 rounded-lg hover:bg-stone-800 dark:hover:bg-stone-200 transition-colors cursor-pointer"
                  >
                    <Send className="w-3.5 h-3.5" />
                    <span>Publicar Comentário</span>
                  </button>
                </div>
              </form>

              {/* Comment List */}
              <div className="space-y-4">
                {comments.map((comment) => (
                  <div
                    key={comment.id}
                    className="p-5 rounded-xl bg-white dark:bg-stone-900 border border-stone-200/80 dark:border-stone-800/80"
                  >
                    <div className="flex items-center justify-between mb-2">
                      <div>
                        <span className="text-xs font-bold text-stone-900 dark:text-stone-100">
                          {comment.authorName}
                        </span>
                        <span className="mx-2 text-stone-300 dark:text-stone-700">·</span>
                        <span className="text-[11px] text-stone-500 dark:text-stone-400">
                          {comment.authorRole}
                        </span>
                      </div>
                      <span className="text-[11px] font-mono text-stone-400">
                        {formatPublicationDate(comment.createdAt)}
                      </span>
                    </div>
                    <p className="text-xs sm:text-sm text-stone-700 dark:text-stone-300 leading-relaxed">
                      {comment.content}
                    </p>
                  </div>
                ))}
              </div>
            </section>
          </>
        )}
      </div>
    </article>
  );
};
