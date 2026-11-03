import { Article, Comment } from '../types/blog';

export const AUTHORS = {
  camila: {
    id: 'author-camila',
    name: 'Dra. Camila Siqueira',
    role: 'Pesquisadora em IA & Modelos Cognitivos',
    bio: 'Doutora pelo Inria e ex-pesquisadora em arquiteturas neurais de raciocínio. Estuda sistemas multi-agentes e verificação formal aplicada a LLMs.',
    avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80',
    twitter: '@camilasiqueira_ai',
    github: 'camilasiqueira',
    articleCount: 14,
  },
  gabriel: {
    id: 'author-gabriel',
    name: 'Gabriel Alencar',
    role: 'Staff Systems Engineer na InfraCore',
    bio: 'Especialista em compiladores, Rust, e sistemas de baixa latência em larga escala. Mantenedor de projetos open-source de rede assíncrona.',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80',
    twitter: '@galencar_rust',
    github: 'galencar',
    articleCount: 9,
  },
  renato: {
    id: 'author-renato',
    name: 'Renato F. Vasconcelos',
    role: 'Arquiteto de Semicondutores & Silício',
    bio: '15 anos projetando aceleradores neurais e ASICs. Entusiasta de interconexões ópticas e empacotamento 3D chiplet.',
    avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&auto=format&fit=crop&q=80',
    twitter: '@renatovasc_chips',
    github: 'renatovasc',
    articleCount: 11,
  },
  sofia: {
    id: 'author-sofia',
    name: 'Sofia Mendes',
    role: 'Principal Design Technologist',
    bio: 'Pioneira em design systems interoperáveis, tipografia generativa e ergonomia de interfaces modernas.',
    avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=150&auto=format&fit=crop&q=80',
    twitter: '@sofiamendes_ux',
    github: 'sofiamendes',
    articleCount: 8,
  },
  marcelo: {
    id: 'author-marcelo',
    name: 'Marcelo Brandão',
    role: 'Criptógrafo & Consultor de Segurança',
    bio: 'Focado em criptografia pós-quântica (PQC), protocolos de chave pública seguros contra ataques de Shor e arquiteturas zero-trust.',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&auto=format&fit=crop&q=80',
    twitter: '@brandao_crypto',
    github: 'mbrandao',
    articleCount: 6,
  },
  helena: {
    id: 'author-helena',
    name: 'Helena Vance',
    role: 'Diretora de Engenharia & Ensaísta',
    bio: 'Escreve sobre sociologia do código, liderança técnica, evolução de ecossistemas open-source e o impacto humano da tecnologia.',
    avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=150&auto=format&fit=crop&q=80',
    twitter: '@helenavance_tech',
    github: 'hvance',
    articleCount: 22,
  },
};

export const INITIAL_ARTICLES: Article[] = [
  {
    id: 'art-1',
    slug: 'fronteira-modelos-raciocinio-agentes-autonomos',
    title: 'A Nova Fronteira dos Agentes Autônomos: Por Que Raciocínio Superou a Escala Pura',
    subtitle: 'Uma dissecação arquitetural do test-time compute, busca em árvore e orquestração determinística.',
    leadParagraph: 'Durante anos, o dogma da inteligência artificial foi simples: dobrar a quantidade de parâmetros e tokens pré-treinados resultava em inteligência proporcional. Hoje, as fronteiras da engenharia inverteram essa premissa. O salto substancial em precisão não vem mais da memorização estatística na fase de pré-treino, mas do orçamento dinâmico de computação em tempo de inferência.',
    category: 'IA & Modelos',
    readTimeMinutes: 8,
    publishedAt: '2026-03-28T09:00:00Z',
    coverImage: '/src/assets/images/tech_ai_agents_lead_1793715354649.jpg',
    author: AUTHORS.camila,
    tags: ['Inteligência Artificial', 'Modelos de Raciocínio', 'Agentes', 'Arquitetura'],
    tier: 1,
    featured: true,
    claps: 342,
    views: 4890,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'O trade-off tradicional mudou: gastar 10x mais ciclos em tempo de inferência supera modelos 5x maiores operando em modo zero-shot.',
          'Modelos de raciocínio combinam busca em árvore (Tree of Thoughts) com auto-avaliação e retrocesso (backtracking) quando encontram becos sem saída.',
          'A verdadeira autonomia reside na capacidade de invocar ferramentas com verificadores formais e validação estrita de esquemas em tempo real.',
        ],
      },
      {
        type: 'heading2',
        text: 'O Esgotamento do Pré-Treino Passivo',
      },
      {
        type: 'paragraph',
        text: 'A curva de retornos decrescentes do pré-treinamento não é surpresa para quem monitora o custo energético de clusters com dezenas de milhares de GPUs. Alimentar o modelo com a totalidade do texto aberto na web criou máquinas com invejável fluência sintática, porém propensas a lapsos catastróficos em cadeias de raciocínio superiores a cinco derivações lógicas.',
      },
      {
        type: 'paragraph',
        text: 'A virada de chave aconteceu quando laboratórios de ponta adotaram o que a literatura chama de "test-time compute scaling". Em vez de forçar o modelo a emitir o primeiro token em menos de 20 milissegundos, o sistema recebe um orçamento computacional de 5 a 30 segundos. Ele gera múltiplos rascunhos mentais, pontua cada nó da árvore de decisão e retrocede quando identifica inconsistências matemáticas.',
      },
      {
        type: 'quote',
        text: 'Não estamos mais construindo papagaios estocásticos de resposta instantânea. Estamos construindo compiladores iterativos de hipóteses guiados por restrições formais.',
        authorQuote: 'Dra. Camila Siqueira',
      },
      {
        type: 'heading2',
        text: 'Orquestração de Ferramentas com Verificação em Loop Fechado',
      },
      {
        type: 'paragraph',
        text: 'Considere a diferença entre um agente ingênuo e um agente de raciocínio avançado diante de uma tarefa crítica de migração de banco de dados. O agente ingênuo gera SQL imediatamente e tenta executá-lo. O agente com orquestração em loop fechado emite um plano abstrato, submete o schema para um validador de tipos, executa em uma réplica efêmera com snapshot transacional e só conclui após checar invariants.',
      },
      {
        type: 'code',
        language: 'typescript',
        caption: 'Exemplo de dispatch determinístico com validador estrito de invariants em TypeScript',
        code: `interface AgentStep<TInput, TOutput> {
  stepId: string;
  hypothesis: string;
  action: () => Promise<TOutput>;
  invariantValidator: (result: TOutput) => boolean;
  rollback: () => Promise<void>;
}

async function executeReasoningStep<TIn, TOut>(
  step: AgentStep<TIn, TOut>,
  maxRetries = 3
): Promise<TOut> {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    const outcome = await step.action();
    const passed = step.invariantValidator(outcome);
    
    if (passed) {
      return outcome;
    }
    
    console.warn(\`Invariant violation on \${step.stepId}, backtracking...\`);
    await step.rollback();
  }
  throw new Error(\`Failed to verify invariant for step: \${step.stepId}\`);
}`,
      },
      {
        type: 'heading2',
        text: 'Benchmark de Precisão: Zero-Shot vs. Multi-Passo com Raciocínio',
      },
      {
        type: 'interactive-benchmark',
        benchmarkData: {
          title: 'Resolução de Tarefas de Engenharia Complexas (SWE-Bench Pro)',
          description: 'Comparativo entre inferência imediata zero-shot e raciocínio com busca e auto-correção.',
          metrics: [
            { label: 'Precisão Lógica', valueA: 42, valueB: 88, unit: '%' },
            { label: 'Detecção de Regressão', valueA: 31, valueB: 94, unit: '%' },
            { label: 'Alucinação em Código', valueA: 28, valueB: 3, unit: '%' },
            { label: 'Tempo Médio de Resposta', valueA: 1.2, valueB: 14.5, unit: 's' },
          ],
        },
      },
      {
        type: 'heading2',
        text: 'O Futuro Próximo: Agentes Especializados Heterogêneos',
      },
      {
        type: 'paragraph',
        text: 'A era do modelo monolítico único está chegando ao fim. Em ambientes industriais de 2026 e adiante, o padrão dominante é um conselho de subagentes especializados: um modelo rápido de triagem sintática, um modelo formal de prova lógica, e um auditor adversarial que atua como crítico antes de qualquer modificação de código ir para o repositório principal.',
      },
    ],
  },
  {
    id: 'art-2',
    slug: 'go-para-rust-sistemas-criticos-2026',
    title: 'Por Que Migramos Nosso Core de Go para Rust em 2026: Lições de Produção',
    subtitle: 'Spikes de p99 no Garbage Collector e a busca por previsibilidade mecânica absoluta em alta vazão.',
    leadParagraph: 'Go é uma das linguagens mais produtivas já concebidas para microsserviços de negócio. No entanto, quando um gateway precisa rotear 450 mil requisições HTTP/3 por segundo sob SLA de latência máxima abaixo de 2 milissegundos, as pausas imprevisíveis do runtime exigem concessões insustentáveis.',
    category: 'Engenharia de Sistemas',
    readTimeMinutes: 6,
    publishedAt: '2026-03-24T14:15:00Z',
    coverImage: '/src/assets/images/tech_rust_systems_1793715354876.jpg',
    author: AUTHORS.gabriel,
    tags: ['Rust', 'Go', 'Sistemas Distribuídos', 'Performance', 'Memória'],
    tier: 2,
    featured: true,
    claps: 278,
    views: 3410,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'A redução do percentil p99 de 14.8ms para 0.9ms justificou integralmente o custo de curva de aprendizado da equipe.',
          'O modelo de ownership de Rust eliminou vazamentos de concorrência e reduziu o consumo de memória em 78% nos nós Kubernetes.',
          'A migração não foi total: serviços de CRUD e back-office continuam operando em Go pela rapidez de iteração.',
        ],
      },
      {
        type: 'heading2',
        text: 'O Problema Silencioso do p99 em Alta Concorrência',
      },
      {
        type: 'paragraph',
        text: 'Em médias e medianas, nosso cluster de Go apresentava números invejáveis: 3ms de tempo médio. Mas em sistemas de pagamentos financeiros distribuídos, médias são mentiras estatísticas. Um usuário que aguarda na ponta do percentil p99 sofria com pausas de 15ms a 22ms enquanto o runtime realizava varreduras de mark-sweep concorrente sobre centenas de milhares de buffers na memória heap.',
      },
      {
        type: 'code',
        language: 'rust',
        caption: 'Manipulação de buffer zero-copy com lifetimes e bytes slice em Rust',
        code: `use bytes::Bytes;
use tokio::io::AsyncReadExt;

pub struct StreamBuffer<'a> {
    payload: &'a [u8],
    correlation_id: u64,
}

impl<'a> StreamBuffer<'a> {
    #[inline(always)]
    pub fn parse(raw: &'a [u8]) -> Result<Self, ProtocolError> {
        if raw.len() < 12 {
            return Err(ProtocolError::IncompleteHeader);
        }
        let correlation_id = u64::from_be_bytes(raw[0..8].try_into().unwrap());
        Ok(Self {
            payload: &raw[8..],
            correlation_id,
        })
    }
}`,
      },
      {
        type: 'heading2',
        text: 'A Simbiose Entre Rust e o Runtime Tokio',
      },
      {
        type: 'paragraph',
        text: 'Com Rust e o ecossistema Tokio, conseguimos associar tarefas diretamente a núcleos dedicados (work-stealing threadpool afinado) com alocação na stack para a vasta maioria dos frames de rede transitórios. Se um pacote precisa ser descartado, o destrutor (Drop) libera o recurso imediatamente, sem enfileirar trabalho para um coletor assíncrono posterior.',
      },
      {
        type: 'quote',
        text: 'Previsibilidade é o santo graal da engenharia de infraestrutura. Não queremos apenas velocidade; queremos a certeza de que a milésima requisição terá a mesma latência que a primeira.',
        authorQuote: 'Gabriel Alencar',
      },
    ],
  },
  {
    id: 'art-3',
    slug: 'fim-lei-de-moore-aceleradores-silicio',
    title: 'O Fim da Lei de Moore e o Renascimento dos Aceleradores de Silício Especializados',
    subtitle: 'Como a física dos nós de 2nm impulsionou a era dos chiplets 3D e da fotônica integrada.',
    leadParagraph: 'A redução de tamanho dos transistores atingiu o limite onde efeitos de tunelamento quântico deixaram de ser curiosidade acadêmica e viraram obstáculo térmico instransponível. O resultado não é a estagnação da capacidade computacional, mas a reinvenção radical da arquitetura de hardware.',
    category: 'Hardware & Chips',
    readTimeMinutes: 7,
    publishedAt: '2026-03-20T11:45:00Z',
    coverImage: '/src/assets/images/tech_chips_quantum_1793715355825.jpg',
    author: AUTHORS.renato,
    tags: ['Hardware', 'Semicondutores', 'Silício', 'Chiplets', 'GPU'],
    tier: 2,
    featured: true,
    claps: 195,
    views: 2980,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'Dies monolíticos gigantescos estão sendo abandonados em prol de chiplets modulares interconectados por interfaces dielétricas 2.5D/3D.',
          'O gargalo agora é a memória e a largura de banda de barramento (Memory Wall), não o número bruto de núcleos lógicos.',
          'Interconexões ópticas on-package prometem reduzir o gasto de energia por bit transmitido em até 80%.',
        ],
      },
      {
        type: 'heading2',
        text: 'A Barreira Térmica e a Parada do Dimensionamento de Dennard',
      },
      {
        type: 'paragraph',
        text: 'Mesmo que litografias ultravioleta extrema (High-NA EUV) consigam desenhar trilhas sub-2 nanômetros, a densidade de dissipação de calor em watts por milímetro quadrado ultrapassou a capacidade de refrigeradores a líquido convencionais. É impossível manter todos os transistores de uma matriz ativados simultaneamente — o fenômeno conhecido como Dark Silicon.',
      },
      {
        type: 'paragraph',
        text: 'A resposta da indústria foi a desagregação: transformar um único processador central em uma constelação de pastilhas de silício especializadas. Módulos dedicados para multiplicação de tensores matriciais, blocos específicos para descompressão de dados e memórias HBM empilhadas verticalmente com milhares de vias condutoras (Through-Silicon Vias).',
      },
    ],
  },
  {
    id: 'art-4',
    slug: 'anatomia-design-systems-que-nao-quebram',
    title: 'A Anatomia dos Design Systems que Não Quebram: De Tokens a Web Components Nativos',
    subtitle: 'Como desvincular seu sistema de design das flutuações de frameworks JavaScript e garantir longevidade.',
    leadParagraph: 'Quantas vezes sua empresa reescreveu a biblioteca de componentes porque a camada de framework migrou do React para Vue, depois para Next.js e agora para Svelte ou Qwik? Interfaces digitais construídas em 2026 estão finalmente adotando o que a própria plataforma web sempre ofereceu: contratos nativos imutáveis.',
    category: 'Design & UX',
    readTimeMinutes: 5,
    publishedAt: '2026-03-16T16:20:00Z',
    coverImage: '/src/assets/images/tech_interface_design_1793715355582.jpg',
    author: AUTHORS.sofia,
    tags: ['Design Systems', 'CSS', 'Web Components', 'Acessibilidade', 'Frontend'],
    tier: 3,
    claps: 224,
    views: 2850,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'Tokens de design em JSON agnóstico compilados para CSS Custom Properties criam uma fundação à prova de futuro.',
          'Cascade Layers (@layer) em CSS moderno resolvem para sempre guerras de especificidade e conflitos de estilo.',
          'Acessibilidade não é uma camada adicionada no final; é a própria estrutura semântica dos elementos base.',
        ],
      },
      {
        type: 'heading2',
        text: 'A Armadilha do Acoplamento ao Framework',
      },
      {
        type: 'paragraph',
        text: 'Um botão não deveria exigir 50 kilobaites de runtime JavaScript apenas para renderizar um estado desabilitado e disparar um evento. Quando estruturamos design systems com foco em contratos semânticos nativos e CSS composable, garantimos que o mesmo botão funcione num microfrontend legado, numa aplicação moderna em SSR e num painel administrativo estático.',
      },
    ],
  },
  {
    id: 'art-5',
    slug: 'criptografia-pos-quantica-na-pratica',
    title: 'Criptografia Pós-Quântica na Prática: O Que Engenheiros Devem Migrar Hoje',
    subtitle: 'Preparando a infraestrutura de rede para o fim do RSA e da criptografia de curva elíptica.',
    leadParagraph: 'O perigo da "colheita antecipada" (Harvest Now, Decrypt Later) não é uma especulação de ficção científica. Agências estatais e corporações já estão armazenando tráfego cifrado em massa na expectativa de que computadores quânticos decifrem chaves assimétricas na próxima década. O NIST padronizou os primeiros algoritmos; a corrida agora é da engenharia de software.',
    category: 'Segurança & Criptografia',
    readTimeMinutes: 6,
    publishedAt: '2026-03-11T13:00:00Z',
    coverImage: 'https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=1200&auto=format&fit=crop&q=80',
    author: AUTHORS.marcelo,
    tags: ['Segurança', 'Criptografia', 'Redes', 'TLS', 'Computação Quântica'],
    tier: 3,
    claps: 167,
    views: 2130,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'Algoritmos baseados em reticulados (Lattice-based cryptography), como o ML-KEM (Kyber), são os novos padrões do NIST.',
          'O principal desafio prático é o tamanho das chaves e das assinaturas: chaves de 1.5KB a 3KB sobrecarregam handshakes TLS antigos.',
          'Estratégias híbridas combinam Curve25519 clássico com ML-KEM para garantir segurança mesmo diante de falhas de implementação inicial.',
        ],
      },
      {
        type: 'heading2',
        text: 'Por Que o Handshake TLS 1.3 Precisa Mudar',
      },
      {
        type: 'paragraph',
        text: 'Com chaves públicas significativamente maiores do que as enxutas chaves de 32 bytes do Curve25519, pacotes de handshake frequentemente excedem a unidade máxima de transmissão (MTU) de 1500 bytes da internet padrão. Isso obriga engenheiros a lidar com fragmentação de pacotes TCP e latência extra de round-trip.',
      },
    ],
  },
  {
    id: 'art-6',
    slug: 'dilema-desenvolvedor-senior-sintese-codigo',
    title: 'O Dilema do Desenvolvedor Sênior na Era das Ferramentas de Síntese de Código',
    subtitle: 'Por que a habilidade mais valiosa deixou de ser digitar lógica e passou a ser o discernimento crítico arquitetural.',
    leadParagraph: 'Assistentes de geração produzem milhares de linhas de código em minutos. Mas código é passivo, não ativo; cada linha adicionada a um repositório é uma responsabilidade futura de manutenção. Na medida em que o custo marginal de escrever código cai para zero, o custo de validar a coerência conceitual do sistema cresce exponencialmente.',
    category: 'Opinião & Futuro',
    readTimeMinutes: 5,
    publishedAt: '2026-03-05T08:30:00Z',
    coverImage: 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?w=1200&auto=format&fit=crop&q=80',
    author: AUTHORS.helena,
    tags: ['Carreira', 'Engenharia de Software', 'Filosofia', 'Liderança'],
    tier: 3,
    claps: 412,
    views: 5210,
    blocks: [
      {
        type: 'takeaways',
        items: [
          'Escrever código tornou-se a parte mais rápida do desenvolvimento; o gargalo agora é a especificação correta do domínio.',
          'Desenvolvedores sêniores operam como editores-chefe e curadores de sistemas, não meros digitadores de rotinas.',
          'A proliferação de código não revisado corre o risco de criar dívidas técnicas invisíveis e arquiteturas incoerentes.',
        ],
      },
      {
        type: 'heading2',
        text: 'De Artesão a Curador Crítico',
      },
      {
        type: 'paragraph',
        text: 'Historicamente, a senioridade era medida pela maestria na sintaxe da linguagem, no conhecimento enciclopédico de bibliotecas e na capacidade de memorizar truques de implementação. Esse valor transmutou-se. Hoje, o engenheiro indispensável é aquele capaz de perguntar: "Esse problema realmente precisa dessa complexidade? Quais são as falhas em cascata que este modelo não previu?"',
      },
      {
        type: 'quote',
        text: 'A inteligência artificial pode gerar a resposta em segundos. Mas é a sabedoria humana que ainda sabe formular a pergunta certa.',
        authorQuote: 'Helena Vance',
      },
    ],
  },
];

export const INITIAL_COMMENTS: Comment[] = [
  {
    id: 'c-1',
    articleId: 'art-1',
    authorName: 'Thiago Nogueira',
    authorRole: 'Engenheiro de ML na Turing Lab',
    content: 'Brilhante artigo da Camila. Estamos aplicando exatamente essa técnica de test-time compute scaling em nossos agentes de teste. A redução de alucinações ao forçar backtracking foi de mais de 70%.',
    createdAt: '2026-03-28T14:20:00Z',
    claps: 18,
  },
  {
    id: 'c-2',
    articleId: 'art-1',
    authorName: 'Mariana Duarte',
    authorRole: 'Tech Lead na CloudStream',
    content: 'O ponto chave para nós foi a verificação formal de esquemas. Não adianta o LLM "achar" que o SQL é seguro; você precisa de uma sandbox com invariants rígidos.',
    createdAt: '2026-03-29T10:05:00Z',
    claps: 9,
  },
  {
    id: 'c-3',
    articleId: 'art-2',
    authorName: 'Lucas Ferraz',
    authorRole: 'Backend Developer',
    content: 'Essa discussão de p99 em Go é real. Nós sofremos meses tentando tunar o GOGC antes de admitir que para determinados micro-serviços de rede crítica o Rust é insubstituível.',
    createdAt: '2026-03-25T11:15:00Z',
    claps: 14,
  },
];
