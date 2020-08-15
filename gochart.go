package chart

import "io"

type GoChart interface {
	Render(rp RendererProvider, w io.Writer) error
}
