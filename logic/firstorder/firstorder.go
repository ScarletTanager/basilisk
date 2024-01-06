package firstorder

//
// First order logic constructs
//

// Term is a constant, a variable, a function of terms...
type Term interface{}

// Predicate is a variadic function taking a list of Terms and returning a boolean value.
type Predicate func(Term...) bool

type AtomicSentence Predicate
type Atom AtomicSentence

// Not negates the passed value - semantically equivalent to ~ in FOPL
func Not(v bool) bool {
	return !v
}

// Tell registers a fact with the knowledge base.
// func Tell()