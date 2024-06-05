package types

import (
	"errors"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ActorSuite struct{}

var _ = Suite(&ActorSuite{})

func (s *ActorSuite) TestInit(c *C) {
	root := NewActor("root")
	child1 := NewActor("child1")
	child2 := NewActor("child2")
	grandchild1 := NewActor("grandchild1")
	grandchild2 := NewActor("grandchild2")

	root.Children[child1] = true
	root.Children[child2] = true
	child1.Children[grandchild1] = true
	child2.Children[grandchild2] = true

	root.InitRoot()

	// init the atomics
	c.Assert(child1.started, NotNil)
	c.Assert(child1.finished, NotNil)
	c.Assert(grandchild2.started, NotNil)
	c.Assert(grandchild2.finished, NotNil)

	// ensure parents are set
	c.Assert(child1.parents, DeepEquals, map[*Actor]bool{root: true})
	c.Assert(grandchild2.parents, DeepEquals, map[*Actor]bool{child2: true})
}

func (s *ActorSuite) TestWalkDepthFirst(c *C) {
	root := NewActor("root")
	child1 := NewActor("child1")
	child2 := NewActor("child2")
	grandchild1 := NewActor("grandchild1")
	grandchild2 := NewActor("grandchild2")

	root.Children[child1] = true
	root.Children[child2] = true
	child1.Children[grandchild1] = true
	child2.Children[grandchild2] = true

	root.InitRoot()

	visited := map[string]bool{}
	root.WalkDepthFirst(func(a *Actor) bool {
		visited[a.Name] = true
		return a.Execute(nil) == nil
	})

	expected := map[string]bool{
		"root":        true,
		"child1":      true,
		"grandchild1": true,
		"child2":      true,
		"grandchild2": true,
	}
	c.Assert(visited, DeepEquals, expected)
}

func (s *ActorSuite) TestWalkDepthFirstFail(c *C) {
	opFail := func(config *OpConfig) OpResult {
		return OpResult{Continue: false, Finish: true, Error: errors.New("foo")}
	}

	root := NewActor("root")
	child1 := NewActor("child1")
	child2 := NewActor("child2")
	child2.Ops = []Op{opFail}
	child3 := NewActor("child3")
	grandchild1 := NewActor("grandchild1")
	grandchild2 := NewActor("grandchild2")
	grandchild3 := NewActor("grandchild3")

	root.Children[child1] = true
	root.Children[child2] = true
	root.Children[child3] = true
	child1.Children[grandchild1] = true
	child2.Children[grandchild2] = true
	child3.Children[grandchild3] = true

	root.InitRoot()

	visited := map[string]bool{}
	root.WalkDepthFirst(func(a *Actor) bool {
		visited[a.Name] = true
		return a.Execute(nil) == nil
	})

	expected := map[string]bool{
		"root":        true,
		"child1":      true,
		"child2":      true,
		"child3":      true,
		"grandchild1": true,
		"grandchild3": true,
	}
	c.Assert(visited, DeepEquals, expected)
}
