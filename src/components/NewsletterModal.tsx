import React, { useState } from 'react';
import { X, Mail, CheckCircle2, ShieldCheck } from 'lucide-react';

interface NewsletterModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const NewsletterModal: React.FC<NewsletterModalProps> = ({ isOpen, onClose }) => {
  const [email, setEmail] = useState('');
  const [topics, setTopics] = useState<string[]>(['IA & Modelos', 'Engenharia de Sistemas']);
  const [subscribed, setSubscribed] = useState(false);

  if (!isOpen) return null;

  const handleToggleTopic = (topic: string) => {
    if (topics.includes(topic)) {
      setTopics(topics.filter((t) => t !== topic));
    } else {
      setTopics([...topics, topic]);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !email.includes('@')) return;
    setSubscribed(true);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs animate-fade-in">
      <div className="relative w-full max-w-lg rounded-2xl bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800 p-6 sm:p-8 shadow-2xl overflow-hidden">
        <button
          onClick={onClose}
          className="absolute top-5 right-5 p-1.5 rounded-lg text-stone-400 hover:text-stone-900 dark:hover:text-stone-100 hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
          aria-label="Fechar janela de inscrição"
        >
          <X className="w-5 h-5" />
        </button>

        {!subscribed ? (
          <div>
            <div className="w-10 h-10 rounded-xl bg-amber-100 dark:bg-amber-950/60 flex items-center justify-center text-amber-800 dark:text-amber-400 mb-4">
              <Mail className="w-5 h-5" />
            </div>

            <h3 className="text-xl sm:text-2xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight">
              Boletim Semanal Sintaxe
            </h3>
            <p className="mt-2 text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed">
              Toda terça-feira pela manhã, uma seleção curada de ensaios aprofundados sobre arquitetura de software, avanços de IA e semicondutores. Zero clickbait, 100% rigor técnico.
            </p>

            <form onSubmit={handleSubmit} className="mt-6 space-y-4">
              <div>
                <label className="block text-xs font-semibold text-stone-700 dark:text-stone-300 mb-1.5">
                  Endereço de e-mail corporativo ou pessoal
                </label>
                <input
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="voce@empresa.com.br"
                  className="w-full px-4 py-2.5 text-sm rounded-xl border border-stone-200 dark:border-stone-700 bg-stone-50 dark:bg-stone-950 text-stone-900 dark:text-stone-100 placeholder-stone-400 focus:outline-none focus:ring-2 focus:ring-amber-700/40"
                />
              </div>

              <div>
                <span className="block text-xs font-semibold text-stone-700 dark:text-stone-300 mb-2">
                  Tópicos de preferência:
                </span>
                <div className="grid grid-cols-2 gap-2 text-xs">
                  {['IA & Modelos', 'Engenharia de Sistemas', 'Hardware & Chips', 'Segurança & Cripto'].map((t) => (
                    <label
                      key={t}
                      className="flex items-center gap-2 p-2 rounded-lg border border-stone-200/80 dark:border-stone-800 bg-stone-50/50 dark:bg-stone-950/50 text-stone-700 dark:text-stone-300 cursor-pointer"
                    >
                      <input
                        type="checkbox"
                        checked={topics.includes(t)}
                        onChange={() => handleToggleTopic(t)}
                        className="rounded text-amber-600 focus:ring-amber-500"
                      />
                      <span className="text-[11px] truncate">{t}</span>
                    </label>
                  ))}
                </div>
              </div>

              <div className="pt-2">
                <button
                  type="submit"
                  className="w-full py-2.5 px-4 text-xs sm:text-sm font-semibold text-stone-50 bg-stone-900 dark:bg-stone-100 dark:text-stone-900 rounded-xl hover:bg-stone-800 dark:hover:bg-stone-200 transition-colors shadow-xs cursor-pointer"
                >
                  Confirmar Inscrição Gratuita
                </button>
              </div>

              <div className="flex items-center gap-2 text-[11px] text-stone-400 justify-center">
                <ShieldCheck className="w-3.5 h-3.5" />
                <span>Sem spam. Cancele com um clique a qualquer momento.</span>
              </div>
            </form>
          </div>
        ) : (
          <div className="py-6 text-center">
            <div className="w-12 h-12 rounded-full bg-emerald-100 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 flex items-center justify-center mx-auto mb-4">
              <CheckCircle2 className="w-7 h-7" />
            </div>
            <h3 className="text-xl font-serif font-bold text-stone-900 dark:text-stone-100">
              Inscrição Confirmada!
            </h3>
            <p className="mt-2 text-xs sm:text-sm text-stone-600 dark:text-stone-300 leading-relaxed max-w-sm mx-auto">
              Enviamos uma mensagem de boas-vindas para <strong className="text-stone-900 dark:text-stone-100">{email}</strong> com nossa edição mais recente.
            </p>
            <button
              onClick={onClose}
              className="mt-6 px-5 py-2 text-xs font-semibold text-stone-900 dark:text-stone-100 border border-stone-300 dark:border-stone-700 rounded-lg hover:bg-stone-100 dark:hover:bg-stone-800 transition-colors cursor-pointer"
            >
              Concluir
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
