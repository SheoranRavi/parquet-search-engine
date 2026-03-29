package services

import (
	"slices"
	"strings"
	"time"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
	"github.com/SheoranRavi/parquet-search-engine/internal/util"
	"github.com/rs/zerolog"
)

type QueryEngine struct {
	logger zerolog.Logger
	store  *store.InMemoryStore
}

func NewQueryEngine(store *store.InMemoryStore) *QueryEngine {
	return &QueryEngine{
		logger: logger.NewServiceLogger("QueryEngine"),
		store:  store,
	}
}

func (q *QueryEngine) Query(input string) ([]model.Message, time.Duration) {
	t := time.Now()
	// tokenize the query
	tokens := util.Tokenize(input)
	tokens = util.FilterStopWords(tokens)
	messages, _ := q.store.GetIntersection(tokens)

	// order messages by timestamp
	slices.SortFunc(messages, func(a, b model.Message) int {
		if a.NanoTimeStamp > b.NanoTimeStamp {
			return 1
		} else if a.NanoTimeStamp < b.NanoTimeStamp {
			return -1
		}
		return 0
	})
	elapsed := time.Since(t)
	q.logger.Info().Msgf("Fetched %d messages in %d ms", len(messages), elapsed.Milliseconds())
	return messages, elapsed
}

// expr     = andExpr ("OR" andExpr)*
// andExpr  = term ("AND" term)*
// term     = WORD | "(" expr ")"

type Node interface{}

type WordNode struct {
	Value string
}

type BinOp struct {
	Op    string // "AND" | "OR"
	Left  Node
	Right Node
}

type Parser struct {
	tokens []string
	pos    int
}

func Parse(query string) Node {
	p := &Parser{tokens: tokenize(query)}
	return p.parseExpr()
}

func (p *Parser) parseExpr() Node {
	left := p.parseAnd()
	for p.peek() == "OR" {
		p.next()
		right := p.parseAnd()
		left = &BinOp{Op: "OR", Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseAnd() Node {
	left := p.parseTerm()
	for p.peek() == "AND" {
		p.next()
		right := p.parseTerm()
		left = &BinOp{Op: "AND", Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseTerm() Node {
	tok := p.next()
	if tok == "(" {
		node := p.parseExpr()
		p.next() // consume ")"
		return node
	}
	return &WordNode{Value: tok}
}

func tokenize(query string) []string {
	// splits on whitespace, keeps AND/OR/( /) as separate tokens
	query = strings.ToLower(query)
	return strings.Fields(query)
}

func (p *Parser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() string {
	tok := p.peek()
	p.pos++
	return tok
}
