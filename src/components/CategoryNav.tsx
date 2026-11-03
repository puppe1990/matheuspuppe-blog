import React from 'react';
import { Search, X, Clock } from 'lucide-react';

interface CategoryNavProps {
  categories: string[];
  categoryCounts: Record<string, number>;
  selectedCategory: string;
  onSelectCategory: (cat: string) => void;
  searchQuery: string;
  onSearchChange: (q: string) => void;
  maxReadTime: number | null;
  onMaxReadTimeChange: (time: number | null) => void;
  showSearchInput: boolean;
}

export const CategoryNav: React.FC<CategoryNavProps> = ({
  categories,
  categoryCounts,
  selectedCategory,
  onSelectCategory,
  searchQuery,
  onSearchChange,
  maxReadTime,
  onMaxReadTimeChange,
  showSearchInput,
}) => {
  return (
    <div className="w-full space-y-4">
      {/* Search Input Bar (smoothly revealed or toggled) */}
      {showSearchInput && (
        <div className="relative max-w-2xl mx-auto transition-all">
          <div className="relative flex items-center">
            <Search className="absolute left-3.5 w-4 h-4 text-stone-400 dark:text-stone-500 pointer-events-none" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              placeholder="Buscar por títulos, tópicos, autores ou palavras-chave..."
              className="w-full pl-10 pr-10 py-2.5 text-sm rounded-xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 text-stone-900 dark:text-stone-100 placeholder-stone-400 focus:outline-none focus:ring-2 focus:ring-amber-800/40 dark:focus:ring-amber-500/40 shadow-xs"
              autoFocus
            />
            {searchQuery && (
              <button
                onClick={() => onSearchChange('')}
                className="absolute right-3 p-1 rounded-md text-stone-400 hover:text-stone-700 dark:hover:text-stone-200 cursor-pointer"
                aria-label="Limpar busca"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>
      )}

      {/* Segmented Filter Bar & Reading Time selector */}
      <div className="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-3 border-b border-stone-200 dark:border-stone-800 pb-3">
        {/* Categories segmented list */}
        <div className="flex items-center gap-1.5 overflow-x-auto w-full lg:w-auto scrollbar-none py-1">
          {categories.map((cat) => {
            const count = categoryCounts[cat] ?? 0;
            const isSelected = selectedCategory === cat;
            return (
              <button
                key={cat}
                onClick={() => onSelectCategory(cat)}
                className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors whitespace-nowrap shrink-0 flex items-center gap-1.5 cursor-pointer ${
                  isSelected
                    ? 'bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 shadow-xs font-semibold'
                    : 'text-stone-600 hover:text-stone-900 hover:bg-stone-200/50 dark:text-stone-400 dark:hover:text-stone-200 dark:hover:bg-stone-900'
                }`}
              >
                <span>{cat}</span>
                <span className={`text-[11px] font-mono tabular-nums ${isSelected ? 'opacity-80' : 'text-stone-400 dark:text-stone-500'}`}>
                  ({count})
                </span>
              </button>
            );
          })}
        </div>

        {/* Reading Time filter */}
        <div className="flex items-center gap-2 self-end lg:self-auto text-xs text-stone-600 dark:text-stone-400">
          <Clock className="w-3.5 h-3.5 text-stone-400" />
          <span className="hidden sm:inline">Tempo:</span>
          <div className="flex items-center gap-1 bg-stone-100 dark:bg-stone-900 p-0.5 rounded-md border border-stone-200 dark:border-stone-800">
            <button
              onClick={() => onMaxReadTimeChange(null)}
              className={`px-2 py-1 rounded text-[11px] font-medium transition-colors cursor-pointer ${
                maxReadTime === null
                  ? 'bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 shadow-xs'
                  : 'text-stone-500 hover:text-stone-900 dark:hover:text-stone-200'
              }`}
            >
              Todos
            </button>
            <button
              onClick={() => onMaxReadTimeChange(5)}
              className={`px-2 py-1 rounded text-[11px] font-medium transition-colors cursor-pointer ${
                maxReadTime === 5
                  ? 'bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 shadow-xs'
                  : 'text-stone-500 hover:text-stone-900 dark:hover:text-stone-200'
              }`}
            >
              &le; 5 min
            </button>
            <button
              onClick={() => onMaxReadTimeChange(8)}
              className={`px-2 py-1 rounded text-[11px] font-medium transition-colors cursor-pointer ${
                maxReadTime === 8
                  ? 'bg-white dark:bg-stone-800 text-stone-900 dark:text-stone-100 shadow-xs'
                  : 'text-stone-500 hover:text-stone-900 dark:hover:text-stone-200'
              }`}
            >
              &le; 8 min
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
