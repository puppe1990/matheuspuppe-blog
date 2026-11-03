// Package db carries the Sintaxe editorial catalog: the published authors,
// essays and seed comments that RunSeeds upserts into SQLite.
package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

func seedAuthors() []models.Author {
	return []models.Author{
		{
			Slug:      "camila-siqueira",
			Name:      "Dra. Camila Siqueira",
			Role:      "Pesquisadora em IA & Modelos Cognitivos",
			Bio:       "Doutora pelo Inria e ex-pesquisadora em arquiteturas neurais de raciocínio. Estuda sistemas multi-agentes e verificação formal aplicada a LLMs.",
			AvatarURL: "/static/img/avatar-camila.jpg",
			Twitter:   "@camilasiqueira_ai",
			GitHub:    "camilasiqueira",
		},
		{
			Slug:      "gabriel-alencar",
			Name:      "Gabriel Alencar",
			Role:      "Staff Systems Engineer na InfraCore",
			Bio:       "Especialista em compiladores, Rust, e sistemas de baixa latência em larga escala. Mantenedor de projetos open-source de rede assíncrona.",
			AvatarURL: "/static/img/avatar-gabriel.jpg",
			Twitter:   "@galencar_rust",
			GitHub:    "galencar",
		},
		{
			Slug:      "renato-vasconcelos",
			Name:      "Renato F. Vasconcelos",
			Role:      "Arquiteto de Semicondutores & Silício",
			Bio:       "15 anos projetando aceleradores neurais e ASICs. Entusiasta de interconexões ópticas e empacotamento 3D chiplet.",
			AvatarURL: "/static/img/avatar-renato.jpg",
			Twitter:   "@renatovasc_chips",
			GitHub:    "renatovasc",
		},
		{
			Slug:      "sofia-mendes",
			Name:      "Sofia Mendes",
			Role:      "Principal Design Technologist",
			Bio:       "Pioneira em design systems interoperáveis, tipografia generativa e ergonomia de interfaces modernas.",
			AvatarURL: "/static/img/avatar-sofia.jpg",
			Twitter:   "@sofiamendes_ux",
			GitHub:    "sofiamendes",
		},
		{
			Slug:      "marcelo-brandao",
			Name:      "Marcelo Brandão",
			Role:      "Criptógrafo & Consultor de Segurança",
			Bio:       "Focado em criptografia pós-quântica (PQC), protocolos de chave pública seguros contra ataques de Shor e arquiteturas zero-trust.",
			AvatarURL: "/static/img/avatar-marcelo.jpg",
			Twitter:   "@brandao_crypto",
			GitHub:    "mbrandao",
		},
		{
			Slug:      "helena-vance",
			Name:      "Helena Vance",
			Role:      "Diretora de Engenharia & Ensaísta",
			Bio:       "Escreve sobre sociologia do código, liderança técnica, evolução de ecossistemas open-source e o impacto humano da tecnologia.",
			AvatarURL: "/static/img/avatar-helena.jpg",
			Twitter:   "@helenavance_tech",
			GitHub:    "hvance",
		},
	}
}

// seedComment pairs a starter comment with its article slug so RunSeeds can
// resolve the article id without hardcoding AUTOINCREMENT values.
type seedComment struct {
	ArticleSlug string
	Comment     models.Comment
}

func seedComments() []seedComment {
	return []seedComment{
		{
			ArticleSlug: "fronteira-modelos-raciocinio-agentes-autonomos",
			Comment: models.Comment{
				AuthorName: "Thiago Nogueira",
				AuthorRole: "Engenheiro de ML na Turing Lab",
				Content:    "Brilhante artigo da Camila. Estamos aplicando exatamente essa técnica de test-time compute scaling em nossos agentes de teste. A redução de alucinações ao forçar backtracking foi de mais de 70%.",
				CreatedAt:  time.Date(2026, 3, 28, 14, 20, 0, 0, time.UTC),
				Claps:      18,
			},
		},
		{
			ArticleSlug: "fronteira-modelos-raciocinio-agentes-autonomos",
			Comment: models.Comment{
				AuthorName: "Mariana Duarte",
				AuthorRole: "Tech Lead na CloudStream",
				Content:    "O ponto chave para nós foi a verificação formal de esquemas. Não adianta o LLM \"achar\" que o SQL é seguro; você precisa de uma sandbox com invariants rígidos.",
				CreatedAt:  time.Date(2026, 3, 29, 10, 5, 0, 0, time.UTC),
				Claps:      9,
			},
		},
		{
			ArticleSlug: "go-para-rust-sistemas-criticos-2026",
			Comment: models.Comment{
				AuthorName: "Lucas Ferraz",
				AuthorRole: "Backend Developer",
				Content:    "Essa discussão de p99 em Go é real. Nós sofremos meses tentando tunar o GOGC antes de admitir que para determinados micro-serviços de rede crítica o Rust é insubstituível.",
				CreatedAt:  time.Date(2026, 3, 25, 11, 15, 0, 0, time.UTC),
				Claps:      14,
			},
		},
	}
}

// authorRef points an article at its columnist; RunSeeds resolves the slug to
// the row id before the upsert.
func authorRef(slug string) models.Author {
	return models.Author{Slug: slug}
}
