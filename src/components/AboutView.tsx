import React from 'react';
import { AUTHORS } from '../data/articles';
import { Terminal, Shield, Feather, BookOpen, ArrowRight } from 'lucide-react';
import { Author } from '../types/blog';

interface AboutViewProps {
  onSelectAuthor: (authorId: string) => void;
  onNavigateHome: () => void;
}

export const AboutView: React.FC<AboutViewProps> = ({ onSelectAuthor, onNavigateHome }) => {
  const authorList: Author[] = Object.values(AUTHORS);

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-10 animate-fade-in">
      {/* Manifesto Header */}
      <div className="mb-12 border-b border-stone-200 dark:border-stone-800 pb-10">
        <div className="flex items-center gap-2 text-xs font-mono font-bold tracking-widest uppercase text-amber-700 dark:text-amber-400 mb-3">
          <Feather className="w-4 h-4" />
          <span>Manifesto Editorial</span>
        </div>
        <h1 className="text-3xl sm:text-4xl lg:text-5xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight leading-tight">
          Engenharia com profundidade. Ideias que resistem ao ciclo de hype.
        </h1>
        <p className="mt-4 text-stone-600 dark:text-stone-300 text-base sm:text-lg leading-relaxed">
          O <strong>Sintaxe</strong> nasceu em resposta à superficialidade da cobertura contemporânea de tecnologia. Enquanto a web foi inundada por listas geradas em segundos e manchetes sensacionalistas sobre inteligência artificial, nós acreditamos no valor do ensaio técnico meticuloso, fundamentado em código real, benchmarks reproduzíveis e reflexão arquitetural.
        </p>
      </div>

      {/* 3 Core Editorial Pillars */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
        <div className="p-6 rounded-2xl bg-white dark:bg-stone-900/60 border border-stone-200 dark:border-stone-800">
          <div className="w-9 h-9 rounded-lg bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-400 flex items-center justify-center mb-4">
            <Terminal className="w-5 h-5" />
          </div>
          <h3 className="text-base font-serif font-bold text-stone-900 dark:text-stone-100 mb-2">
            Primeiros Princípios
          </h3>
          <p className="text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed">
            Investigamos como as coisas realmente funcionam: do circuito integrado ao compilador, do handshake de rede à inferência em matrizes neurais.
          </p>
        </div>

        <div className="p-6 rounded-2xl bg-white dark:bg-stone-900/60 border border-stone-200 dark:border-stone-800">
          <div className="w-9 h-9 rounded-lg bg-stone-200 dark:bg-stone-800 text-stone-800 dark:text-stone-200 flex items-center justify-center mb-4">
            <Shield className="w-5 h-5" />
          </div>
          <h3 className="text-base font-serif font-bold text-stone-900 dark:text-stone-100 mb-2">
            Rigor Mecânico
          </h3>
          <p className="text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed">
            Sem afirmações vazias. Nossas comparações de performance trazem métricas de percentil p99, análise de concorrência e trade-offs documentados.
          </p>
        </div>

        <div className="p-6 rounded-2xl bg-white dark:bg-stone-900/60 border border-stone-200 dark:border-stone-800">
          <div className="w-9 h-9 rounded-lg bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-400 flex items-center justify-center mb-4">
            <BookOpen className="w-5 h-5" />
          </div>
          <h3 className="text-base font-serif font-bold text-stone-900 dark:text-stone-100 mb-2">
            Tipografia & Leitura
          </h3>
          <p className="text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed">
            Uma interface pensada para quem lê ensaios de 10 minutos com foco ininterrupto: contraste calibrado, zero poluição visual e modo escuro nativo.
          </p>
        </div>
      </div>

      {/* Editorial Team */}
      <div>
        <h2 className="text-2xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight mb-6">
          Corpo de Autores & Pesquisadores
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {authorList.map((author) => (
            <div
              key={author.id}
              onClick={() => onSelectAuthor(author.id)}
              className="p-5 rounded-xl bg-white dark:bg-stone-900 border border-stone-200/80 dark:border-stone-800/80 hover:border-stone-400 dark:hover:border-stone-600 transition-colors cursor-pointer flex items-center justify-between group"
            >
              <div className="flex items-center gap-3.5">
                <img
                  src={author.avatar}
                  alt={author.name}
                  referrerPolicy="no-referrer"
                  className="w-12 h-12 rounded-full object-cover ring-1 ring-stone-200 dark:ring-stone-700"
                />
                <div>
                  <h4 className="text-sm font-bold text-stone-900 dark:text-stone-100 group-hover:text-amber-700 dark:group-hover:text-amber-400 transition-colors">
                    {author.name}
                  </h4>
                  <p className="text-xs text-stone-500 dark:text-stone-400">
                    {author.role}
                  </p>
                </div>
              </div>
              <ArrowRight className="w-4 h-4 text-stone-400 group-hover:text-stone-800 dark:group-hover:text-stone-200" />
            </div>
          ))}
        </div>
      </div>

      <div className="mt-12 text-center">
        <button
          onClick={onNavigateHome}
          className="px-5 py-2.5 text-xs font-semibold text-white bg-stone-900 dark:bg-stone-100 dark:text-stone-900 rounded-xl hover:bg-stone-800 dark:hover:bg-stone-200 transition-colors cursor-pointer"
        >
          Explorar Edição Atual de Artigos
        </button>
      </div>
    </div>
  );
};
