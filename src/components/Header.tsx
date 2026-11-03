import React from 'react';
import { Bookmark, Moon, Search, Sun } from 'lucide-react';

interface HeaderProps {
  darkMode: boolean;
  onToggleDarkMode: () => void;
  savedCount: number;
  activeView: 'home' | 'saved' | 'playground' | 'about';
  onNavigate: (view: 'home' | 'saved' | 'playground' | 'about') => void;
  onOpenNewsletter: () => void;
  searchOpen: boolean;
  onToggleSearch: () => void;
}

export const Header: React.FC<HeaderProps> = ({
  darkMode,
  onToggleDarkMode,
  savedCount,
  activeView,
  onNavigate,
  onOpenNewsletter,
  searchOpen,
  onToggleSearch,
}) => {
  return (
    <header className="sticky top-0 z-40 w-full border-b border-stone-200/80 bg-stone-50/90 backdrop-blur-md dark:border-stone-800 dark:bg-stone-950/90 transition-colors">
      <div className="max-w-7xl mx-auto flex items-center justify-between gap-8 px-4 sm:px-6 lg:px-8 h-16">
        {/* Zone 1: Brand Wordmark (Single text element) */}
        <button
          onClick={() => onNavigate('home')}
          className="text-left group cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-stone-500 rounded"
          aria-label="Página inicial do Sintaxe Tech Blog"
        >
          <span className="text-xl sm:text-2xl font-serif font-bold tracking-tight text-stone-900 dark:text-stone-100 group-hover:text-amber-800 dark:group-hover:text-amber-400 transition-colors whitespace-nowrap shrink-0">
            SINTAXE
          </span>
          <span className="sr-only">Jornal de Tecnologia e Sistemas</span>
        </button>

        {/* Zone 2: 4-5 Clean single-line text navigation links */}
        <nav aria-label="Navegação principal" className="hidden md:flex items-center gap-7 text-sm font-medium">
          <button
            onClick={() => onNavigate('home')}
            className={`whitespace-nowrap shrink-0 transition-colors pb-1 border-b-2 cursor-pointer ${
              activeView === 'home'
                ? 'border-stone-900 text-stone-900 dark:border-stone-100 dark:text-stone-100'
                : 'border-transparent text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100'
            }`}
          >
            Artigos
          </button>
          <button
            onClick={() => onNavigate('playground')}
            className={`whitespace-nowrap shrink-0 transition-colors pb-1 border-b-2 cursor-pointer ${
              activeView === 'playground'
                ? 'border-stone-900 text-stone-900 dark:border-stone-100 dark:text-stone-100'
                : 'border-transparent text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100'
            }`}
          >
            Playground & Benchmarks
          </button>
          <button
            onClick={() => onNavigate('saved')}
            className={`whitespace-nowrap shrink-0 flex items-center gap-1.5 transition-colors pb-1 border-b-2 cursor-pointer ${
              activeView === 'saved'
                ? 'border-stone-900 text-stone-900 dark:border-stone-100 dark:text-stone-100'
                : 'border-transparent text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100'
            }`}
          >
            <span>Salvos</span>
            {savedCount > 0 && (
              <span className="text-xs font-mono font-semibold px-1.5 py-0.2 rounded-full bg-stone-200 text-stone-800 dark:bg-stone-800 dark:text-stone-200">
                {savedCount}
              </span>
            )}
          </button>
          <button
            onClick={() => onNavigate('about')}
            className={`whitespace-nowrap shrink-0 transition-colors pb-1 border-b-2 cursor-pointer ${
              activeView === 'about'
                ? 'border-stone-900 text-stone-900 dark:border-stone-100 dark:text-stone-100'
                : 'border-transparent text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100'
            }`}
          >
            Sobre a Redação
          </button>
        </nav>

        {/* Zone 3: 1 Primary action + Utilities */}
        <div className="flex items-center gap-2 sm:gap-3 shrink-0">
          <button
            onClick={onToggleSearch}
            className={`p-2 rounded-lg text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100 hover:bg-stone-200/60 dark:hover:bg-stone-800/60 transition-colors cursor-pointer ${
              searchOpen ? 'bg-stone-200/80 dark:bg-stone-800 text-stone-900 dark:text-stone-100' : ''
            }`}
            aria-label="Abrir busca de artigos"
            title="Buscar artigos"
          >
            <Search className="w-4 h-4" />
          </button>

          <button
            onClick={onToggleDarkMode}
            className="p-2 rounded-lg text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100 hover:bg-stone-200/60 dark:hover:bg-stone-800/60 transition-colors cursor-pointer"
            aria-label={darkMode ? 'Mudar para modo claro' : 'Mudar para modo escuro'}
            title={darkMode ? 'Modo claro' : 'Modo escuro'}
          >
            {darkMode ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
          </button>

          <button
            onClick={() => onNavigate('saved')}
            className="md:hidden relative p-2 rounded-lg text-stone-600 hover:text-stone-900 dark:text-stone-400 dark:hover:text-stone-100 cursor-pointer"
            aria-label="Ver artigos salvos"
          >
            <Bookmark className="w-4 h-4" />
            {savedCount > 0 && (
              <span className="absolute top-1 right-1 w-2 h-2 rounded-full bg-amber-600" />
            )}
          </button>

          <button
            onClick={onOpenNewsletter}
            className="px-3.5 py-2 text-xs sm:text-sm font-medium text-stone-50 bg-stone-900 dark:bg-stone-100 dark:text-stone-900 rounded-lg hover:bg-stone-800 dark:hover:bg-stone-200 transition-colors whitespace-nowrap shrink-0 shadow-sm cursor-pointer"
          >
            Assinar Boletim
          </button>
        </div>
      </div>
    </header>
  );
};
