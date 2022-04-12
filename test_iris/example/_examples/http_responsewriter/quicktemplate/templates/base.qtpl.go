// This file is automatically generated by qtc from "base.qtpl".
// See https://github.com/valyala/quicktemplate for details.

// This is our templates' base implementation.
//

//line base.qtpl:3

package templates

//line base.qtpl:3

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line base.qtpl:3

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line base.qtpl:4

type Partial interface {
	//line base.qtpl:4
	Body() string
	//line base.qtpl:4
	StreamBody(qw422016 *qt422016.Writer)
	//line base.qtpl:4
	WriteBody(qq422016 qtio422016.Writer)
	//line base.qtpl:4

}

// Template writes a template implementing the Partial interface.

//line base.qtpl:11

func StreamTemplate(qw422016 *qt422016.Writer, p Partial) {
	//line base.qtpl:11
	qw422016.N().S(`
<html>
	<head>
		<title>Quicktemplate integration with Iris</title>
	</head>
	<body>
		<div>
			Header contents here...
		</div>

		<div style="margin:10px;">
			`)
	//line base.qtpl:22
	p.StreamBody(qw422016)
	//line base.qtpl:22
	qw422016.N().S(`
		</div>

	</body>
	<footer>
		Footer contents here...
	</footer>
</html>
`)
	//line base.qtpl:30

}

//line base.qtpl:30

func WriteTemplate(qq422016 qtio422016.Writer, p Partial) {
	//line base.qtpl:30
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line base.qtpl:30
	StreamTemplate(qw422016, p)
	//line base.qtpl:30
	qt422016.ReleaseWriter(qw422016)
	//line base.qtpl:30

}

//line base.qtpl:30

func Template(p Partial) string {
	//line base.qtpl:30
	qb422016 := qt422016.AcquireByteBuffer()
	//line base.qtpl:30
	WriteTemplate(qb422016, p)
	//line base.qtpl:30
	qs422016 := string(qb422016.B)
	//line base.qtpl:30
	qt422016.ReleaseByteBuffer(qb422016)
	//line base.qtpl:30
	return qs422016
	//line base.qtpl:30

}

// Base template implementation. Other pages may inherit from it if they need
// overriding only certain Partial methods.

//line base.qtpl:35

type Base struct{}

//line base.qtpl:36

func (b *Base) StreamBody(qw422016 *qt422016.Writer) {
	//line base.qtpl:36

	qw422016.N().S(`This is the base body`)
}

//line base.qtpl:36
//line base.qtpl:36

func (b *Base) WriteBody(qq422016 qtio422016.Writer) {
	//line base.qtpl:36
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line base.qtpl:36
	b.StreamBody(qw422016)
	//line base.qtpl:36
	qt422016.ReleaseWriter(qw422016)
	//line base.qtpl:36

}

//line base.qtpl:36

func (b *Base) Body() string {
	//line base.qtpl:36
	qb422016 := qt422016.AcquireByteBuffer()
	//line base.qtpl:36
	b.WriteBody(qb422016)
	//line base.qtpl:36
	qs422016 := string(qb422016.B)
	//line base.qtpl:36
	qt422016.ReleaseByteBuffer(qb422016)
	//line base.qtpl:36
	return qs422016
	//line base.qtpl:36

}
