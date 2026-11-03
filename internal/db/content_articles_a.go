package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// seedArticles lists every published essay, newest first. Bodies are verbatim
// ports of the shipped editorial content.
func seedArticles() []models.Article {
	return []models.Article{
		articleAgentesAutonomos(),
		articleGoParaRust(),
		articleFimLeiDeMoore(),
		articleDesignSystems(),
		articlePosQuantica(),
		articleSeniorIa(),
	}
}

func articleAgentesAutonomos() models.Article {
	return models.Article{
		Slug:            "fronteira-modelos-raciocinio-agentes-autonomos",
		Title:           "A Nova Fronteira dos Agentes Autônomos: Por Que Raciocínio Superou a Escala Pura",
		Subtitle:        "Uma dissecação arquitetural do test-time compute, busca em árvore e orquestração determinística.",
		LeadParagraph:   "Durante anos, o dogma da inteligência artificial foi simples: dobrar a quantidade de parâmetros e tokens pré-treinados resultava em inteligência proporcional. Hoje, as fronteiras da engenharia inverteram essa premissa. O salto substancial em precisão não vem mais da memorização estatística na fase de pré-treino, mas do orçamento dinâmico de computação em tempo de inferência.",
		Category:        "IA & Modelos",
		ReadTimeMinutes: 8,
		PublishedAt:     time.Date(2026, 3, 28, 9, 0, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_ai_agents_lead.jpg",
		Author:          authorRef("camila-siqueira"),
		Tags:            []string{"Inteligência Artificial", "Modelos de Raciocínio", "Agentes", "Arquitetura"},
		Tier:            1,
		Claps:           342,
		Views:           4890,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"O trade-off tradicional mudou: gastar 10x mais ciclos em tempo de inferência supera modelos 5x maiores operando em modo zero-shot.",
				"Modelos de raciocínio combinam busca em árvore (Tree of Thoughts) com auto-avaliação e retrocesso (backtracking) quando encontram becos sem saída.",
				"A verdadeira autonomia reside na capacidade de invocar ferramentas com verificadores formais e validação estrita de esquemas em tempo real.",
			}},
			{Type: "heading2", Text: "O Esgotamento do Pré-Treino Passivo"},
			{Type: "paragraph", Text: "A curva de retornos decrescentes do pré-treinamento não é surpresa para quem monitora o custo energético de clusters com dezenas de milhares de GPUs. Alimentar o modelo com a totalidade do texto aberto na web criou máquinas com invejável fluência sintática, porém propensas a lapsos catastróficos em cadeias de raciocínio superiores a cinco derivações lógicas."},
			{Type: "paragraph", Text: "A virada de chave aconteceu quando laboratórios de ponta adotaram o que a literatura chama de \"test-time compute scaling\". Em vez de forçar o modelo a emitir o primeiro token em menos de 20 milissegundos, o sistema recebe um orçamento computacional de 5 a 30 segundos. Ele gera múltiplos rascunhos mentais, pontua cada nó da árvore de decisão e retrocede quando identifica inconsistências matemáticas."},
			{Type: "quote", Text: "Não estamos mais construindo papagaios estocásticos de resposta instantânea. Estamos construindo compiladores iterativos de hipóteses guiados por restrições formais.", AuthorQuote: "Dra. Camila Siqueira"},
			{Type: "heading2", Text: "Orquestração de Ferramentas com Verificação em Loop Fechado"},
			{Type: "paragraph", Text: "Considere a diferença entre um agente ingênuo e um agente de raciocínio avançado diante de uma tarefa crítica de migração de banco de dados. O agente ingênuo gera SQL imediatamente e tenta executá-lo. O agente com orquestração em loop fechado emite um plano abstrato, submete o schema para um validador de tipos, executa em uma réplica efêmera com snapshot transacional e só conclui após checar invariants."},
			{
				Type:     "code",
				Language: "typescript",
				Caption:  "Exemplo de dispatch determinístico com validador estrito de invariants em TypeScript",
				Code: "interface AgentStep<TInput, TOutput> {\n" +
					"  stepId: string;\n" +
					"  hypothesis: string;\n" +
					"  action: () => Promise<TOutput>;\n" +
					"  invariantValidator: (result: TOutput) => boolean;\n" +
					"  rollback: () => Promise<void>;\n" +
					"}\n" +
					"\n" +
					"async function executeReasoningStep<TIn, TOut>(\n" +
					"  step: AgentStep<TIn, TOut>,\n" +
					"  maxRetries = 3\n" +
					"): Promise<TOut> {\n" +
					"  for (let attempt = 1; attempt <= maxRetries; attempt++) {\n" +
					"    const outcome = await step.action();\n" +
					"    const passed = step.invariantValidator(outcome);\n" +
					"\n" +
					"    if (passed) {\n" +
					"      return outcome;\n" +
					"    }\n" +
					"\n" +
					"    console.warn(`Invariant violation on ${step.stepId}, backtracking...`);\n" +
					"    await step.rollback();\n" +
					"  }\n" +
					"  throw new Error(`Failed to verify invariant for step: ${step.stepId}`);\n" +
					"}",
			},
			{Type: "heading2", Text: "Benchmark de Precisão: Zero-Shot vs. Multi-Passo com Raciocínio"},
			{
				Type: "interactive-benchmark",
				Benchmark: &models.BenchmarkData{
					Title:       "Resolução de Tarefas de Engenharia Complexas (SWE-Bench Pro)",
					Description: "Comparativo entre inferência imediata zero-shot e raciocínio com busca e auto-correção.",
					Metrics: []models.BenchmarkMetric{
						{Label: "Precisão Lógica", ValueA: 42, ValueB: 88, Unit: "%"},
						{Label: "Detecção de Regressão", ValueA: 31, ValueB: 94, Unit: "%"},
						{Label: "Alucinação em Código", ValueA: 28, ValueB: 3, Unit: "%"},
						{Label: "Tempo Médio de Resposta", ValueA: 1.2, ValueB: 14.5, Unit: "s"},
					},
				},
			},
			{Type: "heading2", Text: "O Futuro Próximo: Agentes Especializados Heterogêneos"},
			{Type: "paragraph", Text: "A era do modelo monolítico único está chegando ao fim. Em ambientes industriais de 2026 e adiante, o padrão dominante é um conselho de subagentes especializados: um modelo rápido de triagem sintática, um modelo formal de prova lógica, e um auditor adversarial que atua como crítico antes de qualquer modificação de código ir para o repositório principal."},
		},
	}
}

func articleGoParaRust() models.Article {
	return models.Article{
		Slug:            "go-para-rust-sistemas-criticos-2026",
		Title:           "Por Que Migramos Nosso Core de Go para Rust em 2026: Lições de Produção",
		Subtitle:        "Spikes de p99 no Garbage Collector e a busca por previsibilidade mecânica absoluta em alta vazão.",
		LeadParagraph:   "Go é uma das linguagens mais produtivas já concebidas para microsserviços de negócio. No entanto, quando um gateway precisa rotear 450 mil requisições HTTP/3 por segundo sob SLA de latência máxima abaixo de 2 milissegundos, as pausas imprevisíveis do runtime exigem concessões insustentáveis.",
		Category:        "Engenharia de Sistemas",
		ReadTimeMinutes: 6,
		PublishedAt:     time.Date(2026, 3, 24, 14, 15, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_rust_systems.jpg",
		Author:          authorRef("gabriel-alencar"),
		Tags:            []string{"Rust", "Go", "Sistemas Distribuídos", "Performance", "Memória"},
		Tier:            2,
		Claps:           278,
		Views:           3410,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"A redução do percentil p99 de 14.8ms para 0.9ms justificou integralmente o custo de curva de aprendizado da equipe.",
				"O modelo de ownership de Rust eliminou vazamentos de concorrência e reduziu o consumo de memória em 78% nos nós Kubernetes.",
				"A migração não foi total: serviços de CRUD e back-office continuam operando em Go pela rapidez de iteração.",
			}},
			{Type: "heading2", Text: "O Problema Silencioso do p99 em Alta Concorrência"},
			{Type: "paragraph", Text: "Em médias e medianas, nosso cluster de Go apresentava números invejáveis: 3ms de tempo médio. Mas em sistemas de pagamentos financeiros distribuídos, médias são mentiras estatísticas. Um usuário que aguarda na ponta do percentil p99 sofria com pausas de 15ms a 22ms enquanto o runtime realizava varreduras de mark-sweep concorrente sobre centenas de milhares de buffers na memória heap."},
			{
				Type:     "code",
				Language: "rust",
				Caption:  "Manipulação de buffer zero-copy com lifetimes e bytes slice em Rust",
				Code: "use bytes::Bytes;\n" +
					"use tokio::io::AsyncReadExt;\n" +
					"\n" +
					"pub struct StreamBuffer<'a> {\n" +
					"    payload: &'a [u8],\n" +
					"    correlation_id: u64,\n" +
					"}\n" +
					"\n" +
					"impl<'a> StreamBuffer<'a> {\n" +
					"    #[inline(always)]\n" +
					"    pub fn parse(raw: &'a [u8]) -> Result<Self, ProtocolError> {\n" +
					"        if raw.len() < 12 {\n" +
					"            return Err(ProtocolError::IncompleteHeader);\n" +
					"        }\n" +
					"        let correlation_id = u64::from_be_bytes(raw[0..8].try_into().unwrap());\n" +
					"        Ok(Self {\n" +
					"            payload: &raw[8..],\n" +
					"            correlation_id,\n" +
					"        })\n" +
					"    }\n" +
					"}",
			},
			{Type: "heading2", Text: "A Simbiose Entre Rust e o Runtime Tokio"},
			{Type: "paragraph", Text: "Com Rust e o ecossistema Tokio, conseguimos associar tarefas diretamente a núcleos dedicados (work-stealing threadpool afinado) com alocação na stack para a vasta maioria dos frames de rede transitórios. Se um pacote precisa ser descartado, o destrutor (Drop) libera o recurso imediatamente, sem enfileirar trabalho para um coletor assíncrono posterior."},
			{Type: "quote", Text: "Previsibilidade é o santo graal da engenharia de infraestrutura. Não queremos apenas velocidade; queremos a certeza de que a milésima requisição terá a mesma latência que a primeira.", AuthorQuote: "Gabriel Alencar"},
		},
	}
}

func articleFimLeiDeMoore() models.Article {
	return models.Article{
		Slug:            "fim-lei-de-moore-aceleradores-silicio",
		Title:           "O Fim da Lei de Moore e o Renascimento dos Aceleradores de Silício Especializados",
		Subtitle:        "Como a física dos nós de 2nm impulsionou a era dos chiplets 3D e da fotônica integrada.",
		LeadParagraph:   "A redução de tamanho dos transistores atingiu o limite onde efeitos de tunelamento quântico deixaram de ser curiosidade acadêmica e viraram obstáculo térmico instransponível. O resultado não é a estagnação da capacidade computacional, mas a reinvenção radical da arquitetura de hardware.",
		Category:        "Hardware & Chips",
		ReadTimeMinutes: 7,
		PublishedAt:     time.Date(2026, 3, 20, 11, 45, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_chips_quantum.jpg",
		Author:          authorRef("renato-vasconcelos"),
		Tags:            []string{"Hardware", "Semicondutores", "Silício", "Chiplets", "GPU"},
		Tier:            2,
		Claps:           195,
		Views:           2980,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Dies monolíticos gigantescos estão sendo abandonados em prol de chiplets modulares interconectados por interfaces dielétricas 2.5D/3D.",
				"O gargalo agora é a memória e a largura de banda de barramento (Memory Wall), não o número bruto de núcleos lógicos.",
				"Interconexões ópticas on-package prometem reduzir o gasto de energia por bit transmitido em até 80%.",
			}},
			{Type: "heading2", Text: "A Barreira Térmica e a Parada do Dimensionamento de Dennard"},
			{Type: "paragraph", Text: "Mesmo que litografias ultravioleta extrema (High-NA EUV) consigam desenhar trilhas sub-2 nanômetros, a densidade de dissipação de calor em watts por milímetro quadrado ultrapassou a capacidade de refrigeradores a líquido convencionais. É impossível manter todos os transistores de uma matriz ativados simultaneamente — o fenômeno conhecido como Dark Silicon."},
			{Type: "paragraph", Text: "A resposta da indústria foi a desagregação: transformar um único processador central em uma constelação de pastilhas de silício especializadas. Módulos dedicados para multiplicação de tensores matriciais, blocos específicos para descompressão de dados e memórias HBM empilhadas verticalmente com milhares de vias condutoras (Through-Silicon Vias)."},
		},
	}
}
