package intepreter

import (
	"github.com/DGTV11/weh-script/nodes"
	"github.com/DGTV11/weh-script/tokens"
)

type Interpreter struct {
	//TODO
}

func (i *Interpreter) Eval(node Node) any { //TODO
	switch n := node.(type) {
	case *nodes.NumberNode: //TODO: (need to look at yt which I can't atm)
		//TODO
	default:
		//TODO: error case
	}
}
