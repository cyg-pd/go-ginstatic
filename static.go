package static

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

type FilePathFunc func(ctx *gin.Context) string
type TemplateValuesFunc func(req *http.Request) any

func DefaultFilePathFunc(ctx *gin.Context) string {
	return strings.TrimPrefix(path.Clean(ctx.Request.URL.Path), "/")
}

type optionFunc func(*Static)
type option interface{ apply(*Static) }

func (fn optionFunc) apply(cfg *Static) { fn(cfg) }

func WithFilePathFunc(f FilePathFunc) option {
	return optionFunc(func(cfg *Static) {
		cfg.filePathFunc = f
	})
}

func WithRenderHeader(h http.Header) option {
	return optionFunc(func(cfg *Static) {
		cfg.renderHeader = h
	})
}

func WithTemplateValuesFunc(f TemplateValuesFunc) option {
	return optionFunc(func(cfg *Static) {
		cfg.templateValuesFunc = f
	})
}

type Static struct {
	fs fs.FS

	indexContent []byte
	tmpl         *template.Template
	renderHeader http.Header

	filePathFunc       FilePathFunc
	templateValuesFunc TemplateValuesFunc
}

func (s *Static) Handler() gin.HandlerFunc { return s.handler }
func (s *Static) handler(ctx *gin.Context) {
	r := ctx.Request
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		ctx.Next()
		return
	}

	upath := s.filePathFunc(ctx)
	if _, err := s.fs.Open(upath); err == nil {
		s.serveFile(ctx, upath)
		return
	}

	if filepath.Ext(upath) != "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	s.serveIndex(ctx)
	ctx.Abort()
}

func (s *Static) writeCustomHeader(ctx *gin.Context) {
	if len(s.renderHeader) == 0 {
		return
	}

	for k := range s.renderHeader {
		ctx.Header(k, s.renderHeader.Get(k))
	}
}

func (s *Static) serveIndex(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
	s.writeCustomHeader(ctx)
	ctx.Header("Content-Type", "text/html")

	if s.templateValuesFunc == nil {
		ctx.Writer.Write(s.indexContent)
		return
	}

	ctx.Header("Cache-Control", "no-store")
	if err := s.tmpl.Execute(ctx.Writer, s.templateValuesFunc(ctx.Request)); err != nil {
		panic(err)
	}
}

func (s *Static) serveFile(ctx *gin.Context, upath string) {
	if ctx.Request.Method == http.MethodHead {
		ctx.AbortWithStatus(http.StatusOK)
		return
	}

	s.writeCustomHeader(ctx)
	http.ServeFileFS(ctx.Writer, ctx.Request, s.fs, upath)
}

func New(fs fs.FS, opts ...option) *Static {
	tmpl, indexContent := mustIndex(fs)
	s := &Static{
		fs:           fs,
		indexContent: indexContent,
		tmpl:         tmpl,
		filePathFunc: DefaultFilePathFunc,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func mustIndex(f fs.FS) (*template.Template, []byte) {
	b := must(io.ReadAll(must(f.Open("index.html"))))
	tmpl := must(template.New("static").Parse(string(b)))
	return tmpl, b
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
