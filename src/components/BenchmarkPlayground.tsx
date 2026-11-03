import React, { useState } from 'react';
import { Cpu, Zap, Activity, Server, RefreshCw } from 'lucide-react';

interface RuntimeStats {
  name: string;
  lang: string;
  color: string;
  accentBg: string;
  p99Latency: number;
  memoryMb: number;
  cpuLoad: number;
  throughputPerWatt: number;
}

export const BenchmarkPlayground: React.FC = () => {
  const [concurrency, setConcurrency] = useState<number>(25000); // 1,000 to 100,000 req/s
  const [payloadKb, setPayloadKb] = useState<number>(4); // 1 to 32 KB

  // Calculate realistic simulated metrics based on architectural profiles
  // Rust: Zero GC, minimal memory overhead, lowest latency
  // Go: Low latency with GC spikes under high heap pressure
  // Node.js: Event loop single-threaded with V8 heap
  // Python (FastAPI/Uvicorn): High flexibility, higher latency and memory
  const rustP99 = Math.max(0.4, Number((0.3 + (concurrency / 100000) * 0.8 + (payloadKb / 32) * 0.5).toFixed(2)));
  const rustMem = Math.round(45 + (concurrency / 1000) * 0.6 + payloadKb * 3);
  const rustCpu = Math.min(95, Math.round(15 + (concurrency / 100000) * 45));

  const goP99 = Math.max(1.8, Number((1.5 + Math.pow(concurrency / 40000, 1.8) * 4.2 + (payloadKb / 32) * 2.1).toFixed(2)));
  const goMem = Math.round(120 + (concurrency / 1000) * 2.8 + payloadKb * 9);
  const goCpu = Math.min(98, Math.round(25 + (concurrency / 100000) * 60));

  const nodeP99 = Math.max(5.5, Number((4.5 + Math.pow(concurrency / 30000, 2.1) * 9.5 + (payloadKb / 32) * 5.0).toFixed(2)));
  const nodeMem = Math.round(240 + (concurrency / 1000) * 5.2 + payloadKb * 16);
  const nodeCpu = Math.min(99, Math.round(40 + (concurrency / 100000) * 58));

  const pyP99 = Math.max(12.0, Number((9.0 + Math.pow(concurrency / 20000, 2.3) * 22.0 + (payloadKb / 32) * 9.5).toFixed(2)));
  const pyMem = Math.round(380 + (concurrency / 1000) * 8.5 + payloadKb * 24);
  const pyCpu = Math.min(100, Math.round(55 + (concurrency / 100000) * 45));

  const runtimes: RuntimeStats[] = [
    {
      name: 'Rust (Tokio + Axum)',
      lang: 'Rust',
      color: 'bg-amber-600 dark:bg-amber-500',
      accentBg: 'border-amber-500/40',
      p99Latency: rustP99,
      memoryMb: rustMem,
      cpuLoad: rustCpu,
      throughputPerWatt: Math.round((concurrency / (rustCpu * 0.8 + 20)) * 10),
    },
    {
      name: 'Go (Gin / net/http)',
      lang: 'Go',
      color: 'bg-cyan-600 dark:bg-cyan-500',
      accentBg: 'border-cyan-500/40',
      p99Latency: goP99,
      memoryMb: goMem,
      cpuLoad: goCpu,
      throughputPerWatt: Math.round((concurrency / (goCpu * 0.9 + 30)) * 10),
    },
    {
      name: 'Node.js (Fastify)',
      lang: 'JavaScript / Node',
      color: 'bg-emerald-600 dark:bg-emerald-500',
      accentBg: 'border-emerald-500/40',
      p99Latency: nodeP99,
      memoryMb: nodeMem,
      cpuLoad: nodeCpu,
      throughputPerWatt: Math.round((concurrency / (nodeCpu * 1.1 + 45)) * 10),
    },
    {
      name: 'Python (FastAPI + Uvicorn)',
      lang: 'Python',
      color: 'bg-indigo-600 dark:bg-indigo-500',
      accentBg: 'border-indigo-500/40',
      p99Latency: pyP99,
      memoryMb: pyMem,
      cpuLoad: pyCpu,
      throughputPerWatt: Math.round((concurrency / (pyCpu * 1.3 + 60)) * 10),
    },
  ];

  return (
    <section className="rounded-2xl bg-white dark:bg-stone-900/80 border border-stone-200 dark:border-stone-800 p-6 sm:p-8 lg:p-10 shadow-xs">
      <div className="max-w-3xl mb-8">
        <div className="flex items-center gap-2 text-xs font-mono uppercase tracking-wider text-amber-700 dark:text-amber-400 mb-2">
          <Activity className="w-4 h-4" />
          <span>Laboratório Interativo de Engenharia</span>
        </div>
        <h2 className="text-2xl sm:text-3xl font-serif font-bold text-stone-900 dark:text-stone-100 tracking-tight">
          Simulador de Runtimes em Alta Concorrência
        </h2>
        <p className="mt-2 text-stone-600 dark:text-stone-300 text-sm leading-relaxed">
          Ajuste os parâmetros de tráfego de rede abaixo para inspecionar a resposta mecânica de quatro pilhas de backend. Observe como a pressão sobre o heap afeta o percentil p99 e o custo de infraestrutura em nuvem.
        </p>
      </div>

      {/* Control Sliders */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-5 rounded-xl bg-stone-50 dark:bg-stone-950/60 border border-stone-200/80 dark:border-stone-800/80 mb-8">
        <div>
          <div className="flex items-center justify-between text-xs font-semibold text-stone-800 dark:text-stone-200 mb-2">
            <span className="flex items-center gap-1.5">
              <Server className="w-3.5 h-3.5 text-stone-500" />
              Volume de Requisições
            </span>
            <span className="font-mono tabular-nums text-amber-700 dark:text-amber-400">
              {concurrency.toLocaleString('pt-BR')} req/s
            </span>
          </div>
          <input
            type="range"
            min="2000"
            max="100000"
            step="2000"
            value={concurrency}
            onChange={(e) => setConcurrency(Number(e.target.value))}
            className="w-full accent-amber-600 cursor-pointer h-2 bg-stone-200 dark:bg-stone-800 rounded-lg"
          />
          <div className="flex justify-between text-[11px] text-stone-400 font-mono mt-1">
            <span>2k req/s</span>
            <span>50k req/s</span>
            <span>100k req/s</span>
          </div>
        </div>

        <div>
          <div className="flex items-center justify-between text-xs font-semibold text-stone-800 dark:text-stone-200 mb-2">
            <span className="flex items-center gap-1.5">
              <Zap className="w-3.5 h-3.5 text-stone-500" />
              Tamanho do Payload JSON
            </span>
            <span className="font-mono tabular-nums text-amber-700 dark:text-amber-400">
              {payloadKb} KB / requisição
            </span>
          </div>
          <input
            type="range"
            min="1"
            max="32"
            step="1"
            value={payloadKb}
            onChange={(e) => setPayloadKb(Number(e.target.value))}
            className="w-full accent-amber-600 cursor-pointer h-2 bg-stone-200 dark:bg-stone-800 rounded-lg"
          />
          <div className="flex justify-between text-[11px] text-stone-400 font-mono mt-1">
            <span>1 KB (leve)</span>
            <span>16 KB</span>
            <span>32 KB (pesado)</span>
          </div>
        </div>
      </div>

      {/* Runtimes Metric Comparison Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {runtimes.map((rt) => (
          <div
            key={rt.name}
            className={`p-5 rounded-xl bg-white dark:bg-stone-900 border ${rt.accentBg} flex flex-col justify-between shadow-xs transition-transform hover:-translate-y-0.5`}
          >
            <div>
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-mono font-semibold uppercase text-stone-500 dark:text-stone-400">
                  {rt.lang}
                </span>
                <span className={`w-2.5 h-2.5 rounded-full ${rt.color}`} />
              </div>
              <h4 className="text-sm font-semibold text-stone-900 dark:text-stone-100">
                {rt.name}
              </h4>

              {/* Metric 1: p99 Latency */}
              <div className="mt-4 pt-3 border-t border-stone-100 dark:border-stone-800/80">
                <div className="flex items-center justify-between text-xs text-stone-500 dark:text-stone-400">
                  <span>Latência p99</span>
                  <span className="font-mono font-bold tabular-nums text-stone-900 dark:text-stone-100">
                    {rt.p99Latency} ms
                  </span>
                </div>
                <div className="w-full bg-stone-100 dark:bg-stone-800 h-1.5 rounded-full mt-1.5 overflow-hidden">
                  <div
                    className={`h-full ${rt.color} transition-all duration-300`}
                    style={{ width: `${Math.min(100, (rt.p99Latency / 30) * 100)}%` }}
                  />
                </div>
              </div>

              {/* Metric 2: Memory */}
              <div className="mt-3">
                <div className="flex items-center justify-between text-xs text-stone-500 dark:text-stone-400">
                  <span>RAM por Pod</span>
                  <span className="font-mono font-semibold tabular-nums text-stone-800 dark:text-stone-200">
                    {rt.memoryMb} MB
                  </span>
                </div>
                <div className="w-full bg-stone-100 dark:bg-stone-800 h-1.5 rounded-full mt-1.5 overflow-hidden">
                  <div
                    className="h-full bg-stone-500 dark:bg-stone-400 transition-all duration-300"
                    style={{ width: `${Math.min(100, (rt.memoryMb / 800) * 100)}%` }}
                  />
                </div>
              </div>

              {/* Metric 3: CPU Load */}
              <div className="mt-3">
                <div className="flex items-center justify-between text-xs text-stone-500 dark:text-stone-400">
                  <span>Carga de CPU</span>
                  <span className="font-mono font-semibold tabular-nums text-stone-800 dark:text-stone-200">
                    {rt.cpuLoad}%
                  </span>
                </div>
              </div>
            </div>

            <div className="mt-4 pt-3 border-t border-stone-100 dark:border-stone-800/80 text-[11px] text-stone-500 dark:text-stone-400">
              <span className="font-mono tabular-nums text-amber-700 dark:text-amber-400 font-semibold">
                {rt.throughputPerWatt.toLocaleString('pt-BR')}
              </span>{' '}
              req/s por kW estimado
            </div>
          </div>
        ))}
      </div>

      {/* Engineering Takeaway Footer */}
      <div className="mt-6 p-4 rounded-xl bg-stone-50 dark:bg-stone-950/60 border border-stone-200/80 dark:border-stone-800/80 text-xs text-stone-600 dark:text-stone-300 flex items-start gap-3">
        <Cpu className="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
        <p>
          <strong className="font-semibold text-stone-900 dark:text-stone-100">Conclusão Mecânica:</strong> Em cargas acima de 50.000 req/s, linguagens com coleta de lixo estocástica gastam até 35% de ciclos úteis de CPU gerenciando heap. Para infraestrutura de roteamento de borda e gateways API, o determinismo de desalocação imediata em Rust viabiliza economias de custo de até 4x em clusters de nuvem.
        </p>
      </div>
    </section>
  );
};
