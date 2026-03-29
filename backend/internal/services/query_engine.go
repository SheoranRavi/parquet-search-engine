package services

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/SheoranRavi/parquet-search-engine/internal/logger"
	"github.com/SheoranRavi/parquet-search-engine/internal/model"
	"github.com/SheoranRavi/parquet-search-engine/internal/store"
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
	q.logger.Info().Msgf("input query: %s", input)
	tree := Parse(input)
	q.logger.Info().Msgf("tree expression: %s", printNode(tree))
	ids := q.eval(tree)
	messages := q.store.GetMessages(ids)

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

func (q *QueryEngine) eval(node Node) map[string]struct{} {
	switch n := node.(type) {
	case *WordNode:
		return q.store.Lookup(n.Value) // returns set of message IDs
	case *BinOp:
		left := q.eval(n.Left)
		right := q.eval(n.Right)
		q.logger.Info().Msgf("left=%d, right=%d, op=%s", len(left), len(right), n.Op)
		if n.Op == "and" {
			result := intersect(left, right)
			q.logger.Info().Msgf("after intersect=%d", len(result))
			return result
		}
		result := union(left, right)
		q.logger.Info().Msgf("after union=%d", len(result))
		return result
	}
	return nil
}

func intersect(a, b map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

func union(a, b map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for k := range a {
		result[k] = struct{}{}
	}
	for k := range b {
		result[k] = struct{}{}
	}
	return result
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

func printNode(node Node) string {
	switch n := node.(type) {
	case *WordNode:
		return fmt.Sprintf("Word(%s)", n.Value)
	case *BinOp:
		return fmt.Sprintf("(%s %s %s)", printNode(n.Left), n.Op, printNode(n.Right))
	}
	return "unknown"
}

func Parse(query string) Node {
	p := &Parser{tokens: tokenize(query)}
	return p.parseExpr()
}

func (p *Parser) parseExpr() Node {
	left := p.parseAnd()
	for {
		tok := p.peek()
		if tok == "or" {
			p.next()
			right := p.parseAnd()
			left = &BinOp{Op: "or", Left: left, Right: right}
		} else if tok != "" && tok != "and" && tok != ")" {
			// implicit OR: "due snapshot" treated as "due or snapshot"
			right := p.parseAnd()
			left = &BinOp{Op: "or", Left: left, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseAnd() Node {
	left := p.parseTerm()
	for p.peek() == "and" {
		p.next()
		right := p.parseTerm()
		left = &BinOp{Op: "and", Left: left, Right: right}
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
