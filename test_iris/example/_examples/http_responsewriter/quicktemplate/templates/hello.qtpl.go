// This file is automatically generated by qtc from "hello.qtpl".
// See https://github.com/valyala/quicktemplate for details.

// Hello template, implements the Partial's methods.
//

//line hello.qtpl:3

package templates

//line hello.qtpl:3

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line hello.qtpl:3

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line hello.qtpl:4

type Hello struct {
	Vars map[string]interface{}
}

//line hello.qtpl:9

func (h *Hello) StreamBody(qw422016 *qt422016.Writer) {
	//line hello.qtpl:9
	qw422016.N().S(`
	<h1>`)
	//line hello.qtpl:10
	qw422016.E().V(h.Vars["message"])
	//line hello.qtpl:10
	qw422016.N().S(`</h1>
	<div>
		Hello <b>`)
	//line hello.qtpl:12
	qw422016.E().V(h.Vars["name"])
	//line hello.qtpl:12
	qw422016.N().S(`!</b>
	</div>
`)
	//line hello.qtpl:14

}

//line hello.qtpl:14

func (h *Hello) WriteBody(qq422016 qtio422016.Writer) {
	//line hello.qtpl:14
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line hello.qtpl:14
	h.StreamBody(qw422016)
	//line hello.qtpl:14
	qt422016.ReleaseWriter(qw422016)
	//line hello.qtpl:14

}

//line hello.qtpl:14

func (h *Hello) Body() string {
	//line hello.qtpl:14
	qb422016 := qt422016.AcquireByteBuffer()
	//line hello.qtpl:14
	h.WriteBody(qb422016)
	//line hello.qtpl:14
	qs422016 := string(qb422016.B)
	//line hello.qtpl:14
	qt422016.ReleaseByteBuffer(qb422016)
	//line hello.qtpl:14
	return qs422016
	//line hello.qtpl:14

}
