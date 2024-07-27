package pineline

import "context"

type Context context.Context

type Handler func(ctx *Context) error

type Node struct {
	handler Handler
	next    *Node
}

func (n *Node) Exec(ctx *Context) error {
	if n == nil || n.handler == nil {
		return nil
	}

	if err := n.handler(ctx); err != nil {
		return err
	}

	return n.next.Exec(ctx)
}

type Pineline struct {
	ctx       *Context
	firstNode *Node
}

func (p *Pineline) Exec() error {
	if err := p.firstNode.Exec(p.ctx); err != nil {
		return err
	}

	return nil
}

func NewFirstNode(handlers ...Handler) *Node {
	firstNode := &Node{}
	curNode := firstNode

	for _, handler := range handlers {
		curNode.handler = handler
		curNode.next = &Node{}
		curNode = curNode.next
	}

	return firstNode
}
