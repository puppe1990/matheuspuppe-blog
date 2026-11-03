import React from 'react';
import { Rss, ArrowUp } from 'lucide-react';

interface FooterProps {
  onSelectCategory: (category: string) => void;
  onNavigate: (view: 'home' | 'saved' | 'playground' | 'about') => void;
  onOpenNewsletter: () => void;
}

export const Footer: React.FC<FooterProps> = ({
  onSelectCategory,
  onNavigate,
  onOpenNewsletter,
}) => {
  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <footer className="w-full border-t border-stone-200 dark:border-stone-800 bg-white dark:bg-stone-950 text-stone-600 dark:text-stone-400 mt-20 transition-colors">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 lg:py-16">
        <div className="grid grid-cols-1 md:grid-cols-12 gap-8 lg:gap-12">
          {/* Brand & Editorial Mission */}
          <div className="md:col-span-5 space-y-4">
            <div className="flex items-center gap-2">
              <span className="text-xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight">
                SINTAXE
              </span>
              <span className="text-xs font-mono text-stone-400">/ TECH BLOG</span>
            </div>
            <p className="text-xs sm:text-sm text-stone-500 dark:text-stone-400 leading-relaxed max-w-sm">
              Publicação independente de tecnologia focada em arquitetura de software, inteligência artificial, engenharia de sistemas e ciência dos computadores.
            </p>
            <div className="pt-2 flex items-center gap-3 text-xs">
              <button
                onClick={onOpenNewsletter}
                className="px-3.5 py-1.5 rounded-lg bg-stone-900 text-stone-50 dark:bg-stone-100 dark:text-stone-900 font-medium hover:opacity-90 transition-opacity cursor-pointer"
              >
                Assinar Boletim Gratuito
              </button>
              <button
                onClick={() => alert('Feed RSS: https://sintaxe.tech/rss.xml (Disponível nos principais agregadores)')}
                className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-stone-200 dark:border-stone-800 hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
              >
                <Rss className="w-3.5 h-3.5 text-amber-600" />
                <span>RSS</span>
              </button>
            </div>
          </div>

          {/* Navigation Sections */}
          <div className="md:col-span-3 space-y-3">
            <h4 className="text-xs font-mono font-bold uppercase tracking-widest text-stone-900 dark:text-stone-100">
              Cadernos & Tópicos
            </h4>
            <ul className="space-y-2 text-xs">
              {['IA & Modelos', 'Engenharia de Sistemas', 'Hardware & Chips', 'Design & UX', 'Segurança & Criptografia'].map((cat) => (
                <li key={cat}>
                  <button
                    onClick={() => {
                      onSelectCategory(cat);
                      onNavigate('home');
                    }}
                    className="hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
                  >
                    {cat}
                  </button>
                </li>
              ))}
            </ul>
          </div>

          {/* Editorial & Tools */}
          <div className="md:col-span-4 space-y-3">
            <h4 className="text-xs font-mono font-bold uppercase tracking-widest text-stone-900 dark:text-stone-100">
              Publicação & Recursos
            </h4>
            <ul className="space-y-2 text-xs">
              <li>
                <button
                  onClick={() => onNavigate('playground')}
                  className="hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
                >
                  Playground de Runtimes & Latência
                </button>
              </li>
              <li>
                <button
                  onClick={() => onNavigate('about')}
                  className="hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
                >
                  Manifesto Editorial & Conselho
                </button>
              </li>
              <li>
                <button
                  onClick={() => onNavigate('saved')}
                  className="hover:text-stone-900 dark:hover:text-stone-100 transition-colors cursor-pointer"
                >
                  Artigos Marcados como Favoritos
                </button>
              </li>
            </ul>

            <div className="pt-3">
              <button
                onClick={scrollToTop}
                className="inline-flex items-center gap-1.5 text-xs text-stone-500 hover:text-stone-900 dark:hover:text-stone-200 transition-colors cursor-pointer"
              >
                <ArrowUp className="w-3.5 h-3.5" />
                <span>Voltar ao topo da página</span>
              </button>
            </div>
          </div>
        </div>

        {/* Bottom Hairline & Legal */}
        <div className="mt-12 pt-6 border-t border-stone-100 dark:border-stone-900 flex flex-col sm:flex-row items-center justify-between gap-4 text-[11px] text-stone-400">
          <p>© {new Date().getFullYear()} Sintaxe Editorial. Todos os direitos reservados.</p>
          <p className="font-mono">Desenvolvido com tipografia e engenharia de alto padrão.</p>
        </div>
      </div>
    </footer>
  );
};
