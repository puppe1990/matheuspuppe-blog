package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func articleDesignSystems() models.Article {
	return models.Article{
		Slug:            "anatomia-design-systems-que-nao-quebram",
		Title:           "A Anatomia dos Design Systems que Não Quebram: De Tokens a Web Components Nativos",
		Subtitle:        "Como desvincular seu sistema de design das flutuações de frameworks JavaScript e garantir longevidade.",
		LeadParagraph:   "Quantas vezes sua empresa reescreveu a biblioteca de componentes porque a camada de framework migrou do React para Vue, depois para Next.js e agora para Svelte ou Qwik? Interfaces digitais construídas em 2026 estão finalmente adotando o que a própria plataforma web sempre ofereceu: contratos nativos imutáveis.",
		Category:        "Design & UX",
		ReadTimeMinutes: 5,
		PublishedAt:     time.Date(2026, 3, 16, 16, 20, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_interface_design.jpg",
		Author:          authorRef("sofia-mendes"),
		Tags:            []string{"Design Systems", "CSS", "Web Components", "Acessibilidade", "Frontend"},
		Tier:            3,
		Claps:           224,
		Views:           2850,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Tokens de design em JSON agnóstico compilados para CSS Custom Properties criam uma fundação à prova de futuro.",
				"Cascade Layers (@layer) em CSS moderno resolvem para sempre guerras de especificidade e conflitos de estilo.",
				"Acessibilidade não é uma camada adicionada no final; é a própria estrutura semântica dos elementos base.",
			}},
			{Type: "heading2", Text: "A Armadilha do Acoplamento ao Framework"},
			{Type: "paragraph", Text: "Um botão não deveria exigir 50 kilobaites de runtime JavaScript apenas para renderizar um estado desabilitado e disparar um evento. Quando estruturamos design systems com foco em contratos semânticos nativos e CSS composable, garantimos que o mesmo botão funcione num microfrontend legado, numa aplicação moderna em SSR e num painel administrativo estático."},
		},
	}
}

func articlePosQuantica() models.Article {
	return models.Article{
		Slug:            "criptografia-pos-quantica-na-pratica",
		Title:           "Criptografia Pós-Quântica na Prática: O Que Engenheiros Devem Migrar Hoje",
		Subtitle:        "Preparando a infraestrutura de rede para o fim do RSA e da criptografia de curva elíptica.",
		LeadParagraph:   "O perigo da \"colheita antecipada\" (Harvest Now, Decrypt Later) não é uma especulação de ficção científica. Agências estatais e corporações já estão armazenando tráfego cifrado em massa na expectativa de que computadores quânticos decifrem chaves assimétricas na próxima década. O NIST padronizou os primeiros algoritmos; a corrida agora é da engenharia de software.",
		Category:        "Segurança & Criptografia",
		ReadTimeMinutes: 6,
		PublishedAt:     time.Date(2026, 3, 11, 13, 0, 0, 0, time.UTC),
		CoverImage:      "/static/img/cover-post-quantica.jpg",
		Author:          authorRef("marcelo-brandao"),
		Tags:            []string{"Segurança", "Criptografia", "Redes", "TLS", "Computação Quântica"},
		Tier:            3,
		Claps:           167,
		Views:           2130,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Algoritmos baseados em reticulados (Lattice-based cryptography), como o ML-KEM (Kyber), são os novos padrões do NIST.",
				"O principal desafio prático é o tamanho das chaves e das assinaturas: chaves de 1.5KB a 3KB sobrecarregam handshakes TLS antigos.",
				"Estratégias híbridas combinam Curve25519 clássico com ML-KEM para garantir segurança mesmo diante de falhas de implementação inicial.",
			}},
			{Type: "heading2", Text: "Por Que o Handshake TLS 1.3 Precisa Mudar"},
			{Type: "paragraph", Text: "Com chaves públicas significativamente maiores do que as enxutas chaves de 32 bytes do Curve25519, pacotes de handshake frequentemente excedem a unidade máxima de transmissão (MTU) de 1500 bytes da internet padrão. Isso obriga engenheiros a lidar com fragmentação de pacotes TCP e latência extra de round-trip."},
		},
	}
}

func articleSeniorIa() models.Article {
	return models.Article{
		Slug:            "dilema-desenvolvedor-senior-sintese-codigo",
		Title:           "O Dilema do Desenvolvedor Sênior na Era das Ferramentas de Síntese de Código",
		Subtitle:        "Por que a habilidade mais valiosa deixou de ser digitar lógica e passou a ser o discernimento crítico arquitetural.",
		LeadParagraph:   "Assistentes de geração produzem milhares de linhas de código em minutos. Mas código é passivo, não ativo; cada linha adicionada a um repositório é uma responsabilidade futura de manutenção. Na medida em que o custo marginal de escrever código cai para zero, o custo de validar a coerência conceitual do sistema cresce exponencialmente.",
		Category:        "Opinião & Futuro",
		ReadTimeMinutes: 5,
		PublishedAt:     time.Date(2026, 3, 5, 8, 30, 0, 0, time.UTC),
		CoverImage:      "/static/img/cover-senior-ia.jpg",
		Author:          authorRef("helena-vance"),
		Tags:            []string{"Carreira", "Engenharia de Software", "Filosofia", "Liderança"},
		Tier:            3,
		Claps:           412,
		Views:           5210,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"Escrever código tornou-se a parte mais rápida do desenvolvimento; o gargalo agora é a especificação correta do domínio.",
				"Desenvolvedores sêniores operam como editores-chefe e curadores de sistemas, não meros digitadores de rotinas.",
				"A proliferação de código não revisado corre o risco de criar dívidas técnicas invisíveis e arquiteturas incoerentes.",
			}},
			{Type: "heading2", Text: "De Artesão a Curador Crítico"},
			{Type: "paragraph", Text: "Historicamente, a senioridade era medida pela maestria na sintaxe da linguagem, no conhecimento enciclopédico de bibliotecas e na capacidade de memorizar truques de implementação. Esse valor transmutou-se. Hoje, o engenheiro indispensável é aquele capaz de perguntar: \"Esse problema realmente precisa dessa complexidade? Quais são as falhas em cascata que este modelo não previu?\""},
			{Type: "quote", Text: "A inteligência artificial pode gerar a resposta em segundos. Mas é a sabedoria humana que ainda sabe formular a pergunta certa.", AuthorQuote: "Helena Vance"},
		},
	}
}
