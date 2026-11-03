package db

import (
	"time"

	"github.com/puppe1990/matheuspuppe-blog/internal/models"
)

// primeiroEnsaio é o primeiro texto publicado do Sintaxe: um relato de
// engenharia sobre o tradutor de EPUBs que roda capítulo a capítulo.
func primeiroEnsaio() models.Article {
	return models.Article{
		Slug:            "traducao-epubs-llms-fila-por-capitulo",
		Title:           "Tradução de EPUBs com LLMs: fila resiliente, retry por capítulo e modelos baratos",
		Subtitle:        "O deploy chega no meio da tradução, o modelo esquece uma tag e a conta sobe. Foi operando um tradutor de livros que aprendi a desenhar a fila, o cancelamento e o custo por capítulo.",
		LeadParagraph:   "Um livro não é uma requisição. São centenas de documentos XHTML, cada um levando minutos de inferência, e um servidor que reinicia no meio do trabalho. A primeira versão do meu tradutor tratava o livro como um job único: um deploy no meio da tarde era suficiente para perder horas de tradução. O desenho que sobrou depois disso vive num aplicativo pequeno, com SQLite e um arquivo por livro, e resolve os três problemas que um serviço de LLM sempre tem: interrupção, validação e custo.",
		Category:        "Engenharia de Sistemas",
		ReadTimeMinutes: 8,
		PublishedAt:     time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC),
		CoverImage:      "/static/img/tech_rust_systems.jpg",
		Author:          authorRef("matheus-puppe"),
		Tags:            []string{"LLM", "EPUB", "Filas", "Confiabilidade", "Custo"},
		Tier:            1,
		Claps:           0,
		Views:           0,
		Blocks: []models.ContentBlock{
			{Type: "takeaways", Items: []string{
				"O capítulo é a unidade atômica: retry, pausa, custo e retomada depois do restart vivem nele, não no livro.",
				"Em produção o deploy chega no meio: o SIGTERM precisa de um drain que devolve o capítulo em execução à fila e deixa o livro marcado como traduzindo para o boot retomar.",
				"O modelo é a parte menos confiável do pipeline: valide a resposta por contrato — tags, proporção, idioma — antes de gastar a próxima tentativa.",
			}},
			{Type: "heading2", Text: "O livro não é uma requisição"},
			{Type: "paragraph", Text: "Um EPUB é um zip com uma lista ordenada de documentos XHTML. O tradutor abre o pacote, lê o spine e transforma cada documento num capítulo: uma linha na tabela chapters, com caminho único dentro do livro, texto de origem, contagem de tokens cl100k e o estado da execução. Essa decisão muda tudo o que vem depois. O retry deixa de ser \"tentar o livro de novo\" e passa a ser \"tentar o capítulo 14\"; a falha de um documento não derruba os outros 199; e o custo de cada chamada fica observável por unidade, porque os tokens estão medidos antes de qualquer requisição."},
			{Type: "paragraph", Text: "A tentação de mandar o livro inteiro numa chamada só é grande: menos código, um streaming único, um único resultado. Na prática, o contexto do modelo vira o limite do produto, o retry fica caro, e qualquer falha no meio derruba a tradução inteira. Traduzir por documento é mais chato de implementar e infinitamente mais barato de operar."},
			{Type: "heading2", Text: "Três estados, uma fila e uma pausa que não é falha"},
			{Type: "paragraph", Text: "O capítulo anda por quatro estados: pendente, traduzindo, traduzido e falhou. A fila é simples de propósito: pendentes e falhos, em ordem do documento. O runner mantém um goroutine por livro, então dois capítulos do mesmo livro nunca competem pelo mesmo arquivo nem pela mesma linha do banco. O botão de retry de uma linha usa outra entrada: o capítulo escolhido fura a fila e vira o próximo depois do atual, sem reiniciar o livro inteiro."},
			{Type: "paragraph", Text: "Falha é um estado de primeira classe. O capítulo que estourou o contexto, recebeu uma resposta com as tags erradas ou bateu num provedor fora do ar termina como falhou, com o erro gravado na linha. O livro só é marcado como falho se algum capítulo ficou assim. Isso mantém a interface honesta: a tabela diz exatamente o que aconteceu, e o retry aponta para o problema certo."},
			{Type: "paragraph", Text: "Pausar é outra história, e a distinção importa. Pausar cancela a chamada em andamento, espera o término da execução e devolve à fila apenas o capítulo que estava traduzindo. O livro fica pausado, os concluídos ficam intactos e o capítulo interrompido não conta como falha. Quando o usuário retoma, a fila continua de onde parou."},
			{Type: "heading2", Text: "O deploy no meio do capítulo"},
			{Type: "paragraph", Text: "O primeiro besouro apareceu num deploy de rotina. O systemd manda SIGTERM e espera quinze segundos; o servidor HTTP, por sua vez, dava cinco segundos para as conexões fecharem. Só que a página do livro mantém um stream SSE aberto o tempo inteiro para mostrar o progresso: esse stream nunca termina sozinho, o shutdown estourava o prazo, o processo saía com erro e matava a chamada de tradução em andamento. O capítulo ficava marcado como traduzindo para sempre, e a interface mentia até alguém reiniciar de novo. Na semana do incidente, uma pausa clicada depois do deploy deixou o livro parado por horas, porque o estado de pausa foi gravado por cima de um trabalho que ninguém tinha pausado."},
			{Type: "paragraph", Text: "A correção tem duas partes. Primeiro, o runner precisa de um drain que reconheça o desligamento: cancelar o contexto de cada livro, esperar até dois segundos pelo término e devolver os capítulos travados em traduzindo para a fila. Segundo, a ordem do shutdown passa a importar: primeiro o drain, depois o servidor HTTP, por último o SQLite. Sem essa ordem, o drain grava num banco que já fechou."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "A ordem do desligamento: drenar a fila antes de fechar o banco",
				Code: `case <-ctx.Done():
    // SIGTERM: devolve o capítulo em execução à fila.
    if a.runner != nil {
        _ = a.runner.Drain(2 * time.Second)
    }
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := a.server.Shutdown(shutdownCtx); err != nil {
        _ = a.server.Close()
    }
    _ = a.store.Close()
    return nil`,
			},
			{Type: "paragraph", Text: "O detalhe que faz o drain funcionar é o que ele não faz: não marca o livro como pausado. Drenar é sobrevivência do processo, não decisão do usuário. O livro continua marcado como traduzindo, e o próximo boot entende que deve retomá-lo."},
			{Type: "quote", Text: "O deploy não é um evento raro: é o clima em que o serviço vive. O código precisa ser escrito para ser interrompido.", AuthorQuote: "Matheus Puppe"},
			{Type: "heading2", Text: "O boot que reconcilia"},
			{Type: "paragraph", Text: "Na subida, antes de qualquer requisição, o runner arruma a bagunça que o processo anterior deixou. Capítulos travados em traduzindo voltam para a fila; livros pausados mantêm a pausa, e o que estava pendente neles estaciona em vez de reaparecer como trabalho em aberto. Só os livros marcados como traduzindo voltam a rodar. A regra é simples e vale a pena repetir: o boot retoma o que o usuário pediu, nunca o que ele pausou."},
			{
				Type:     "code",
				Language: "go",
				Caption:  "ResumeInterrupted: o boot só retoma o que estava traduzindo",
				Code: `func (r *Runner) ResumeInterrupted() error {
    if err := r.store.ResetRunningChapters(); err != nil {
        return err
    }
    if err := r.store.ParkPausedBookQueue(); err != nil {
        return err
    }
    books, err := r.store.ListBooks()
    if err != nil {
        return err
    }
    for _, book := range books {
        if book.Status == "translating" {
            r.Start(book.ID)
        }
    }
    return nil
}`,
			},
			{Type: "paragraph", Text: "É esse par — drain no SIGTERM e reconciliação no boot — que transforma um deploy em algo invisível para quem está traduzindo um livro. O máximo que o usuário vê é a linha do capítulo em andamento piscar e o progresso continuar do próximo pendente."},
			{Type: "heading2", Text: "Validação antes da próxima tentativa"},
			{Type: "paragraph", Text: "Traduzir XHTML não é trocar texto por texto: o modelo devolve markup. Se ele come uma tag, o capítulo inteiro fica inválido. A defesa é mascarar as tags antes da chamada, deixando marcadores no lugar, e conferir a resposta por contrato: a mesma sequência normalizada de tags da origem, texto realmente traduzido, proporção de palavras entre 0,3 e 2,5, nenhum preâmbulo e nenhum capítulo misturado. Uma proporção fora dessa faixa quase sempre significa resposta truncada ou o modelo devolvendo o original."},
			{Type: "paragraph", Text: "Quando a validação falha, o serviço tenta de novo — três tentativas, com um lembrete anexado à conversa dizendo o que quebrou. Mas nem todo erro merece retentativa. Estouro de contexto e um 503 de \"no upstreams\" do OpenRouter não são traduções inválidas: são o provedor sem capacidade naquele minuto, e insistir só gasta tempo e dinheiro. Esses erros sobem para o capítulo como falha imediata."},
			{Type: "quote", Text: "O XHTML original é a fonte da verdade; o modelo só traduz o texto. Se ele omitiu um fechamento de parágrafo, o correto é recolocá-lo — não aceitar a página quebrada.", AuthorQuote: "Comentário no código do keepIfDropped"},
			{Type: "paragraph", Text: "É por isso que o mascaramento tem uma segunda metade. Qualquer tag cujo marcador o modelo omitiu é recolocada na posição certa, na ordem original. Um índice perdeu um fechamento de parágrafo — 401 tags na fonte, 400 na resposta — e o capítulo foi salvo com a página íntegra porque a fonte manda."},
			{Type: "heading2", Text: "Modelos baratos, custo por capítulo"},
			{Type: "paragraph", Text: "O tradutor fala o protocolo OpenAI de chat completions e nada mais. Cada provedor é um perfil com base fixa — OpenAI, Gemini, Groq, NVIDIA, OpenRouter e um personalizado — e só um fica ativo por vez, com a própria chave salva. A configuração é resolvida a cada chamada, então trocar de provedor no meio de um livro vale no próximo capítulo, sem reiniciar nada. Cada linha de capítulo registra provedor, modelo, início, fim e duração da última execução; a configuração é congelada antes da requisição para a tabela não mentir depois."},
			{Type: "paragraph", Text: "O orçamento é por capítulo por um motivo prático: é a única unidade em que dá para comparar modelos sem enganar a si mesmo. Os padrões são deliberadamente baratos — um flash lite do Gemini, um gpt-oss via Groq — porque a maior parte do trabalho é tradução de prosa longa, não raciocínio difícil. Quando o capítulo importa mais, troca-se o modelo e a linha seguinte já sai com o novo; o que ficou traduzido não é retraduzido. Não existe troca automática de provedor quando o selecionado falha. É uma escolha explícita: um fallback silencioso esconde a falha e muda a conta sem ninguém perceber."},
			{Type: "paragraph", Text: "O tamanho do pedaço também é um orçamento. Cada provedor tem um teto de caracteres por chamada: o Gemini aceita 12 mil e devolve até 64 mil tokens; o DeepSeek da NVIDIA fica em 6 mil, porque o gateway devolve 504 perto de quatro minutos de completion bufferizada; o Riva trabalha com 6 mil para caber nos 8 mil de contexto. O cliente tem dez minutos de timeout por chamada, o suficiente para provedores que só mandam cabeçalhos no fim da geração. E os catálogos de modelo de Groq e OpenRouter são cacheados por quinze minutos: escolher modelo é uma consulta, não um evento."},
			{Type: "heading2", Text: "O que eu levaria para o próximo serviço"},
			{Type: "paragraph", Text: "Cinco regras que ficaram deste projeto. Primeira: grave o estado antes de chamar o provedor — a linha traduzindo precisa existir antes do primeiro byte, senão o restart não tem o que reconciliar. Segunda: faça da unidade de trabalho a menor coisa que faz sentido repetir; retry caro é retry que ninguém usa. Terceira: cancelamento é cooperativo e precisa de orçamento — um drain de dois segundos dentro de um SIGTERM de quinze resolve mais que qualquer retry infinito. Quarta: o boot é parte do desenho; todo serviço que mantém estado precisa de uma função que arruma o que o processo anterior deixou. Quinta: meça o custo na unidade que você retenta; sem isso, a fatura é a única telemetria."},
			{Type: "paragraph", Text: "O tradutor hoje não é um sistema distribuído, não tem fila externa e roda num SQLite. Ainda assim, cada uma dessas decisões veio de um incidente real: um deploy, um capítulo travado, uma tag omitida, uma conta inesperada. É o tipo de engenharia que não aparece na demo e aparece inteira no dia em que o serviço precisa continuar sozinho."},
		},
	}
}
