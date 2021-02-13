package webui

import (
	"bytes"
	"context"
	"fmt"
	htmlt "html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	textt "text/template"

	"github.com/OneOfOne/xxhash"
	"github.com/oxtoacart/bpool"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/errors"
)

func NewServer(
	config ServerConfig,
	templateData interface{},
) (*Server, error) {
	if !strings.HasSuffix(config.ServePath, "/") {
		config.ServePath += "/"
	}

	var err error

	config.FilesDir, err = filepath.Abs(config.FilesDir)
	if err != nil {
		return nil, errors.Ent("config.FilesDir", err)
	}

	var filesDirWithSlash string
	if config.FilesDir != "/" {
		filesDirWithSlash = config.FilesDir + "/"
	}
	filesDirNoSlash := config.FilesDir

	fileServer := ETagHandler(
		http.StripPrefix(config.ServePath,
			http.FileServer(
				http.Dir(filesDirNoSlash))))

	srv := &Server{
		config,
		false,
		nil,
		fileServer,
		filesDirWithSlash,
		templateData,
		map[string]*processedFileInfo{},
	}

	err = srv.processFiles()
	if err != nil {
		return nil, err
	}

	return srv, nil
}

type ServerConfig struct {
	ServePort      int                        `env:"SERVE_PORT"`
	ServePath      string                     `env:"SERVE_PATH"`
	FilesDir       string                     `env:"FILES_DIR"`
	FileProcessors map[string][]FileProcessor `env:"-"`
}

type Server struct {
	config         ServerConfig
	shuttingDown   bool
	httpServer     *http.Server
	fileServer     http.Handler
	filesDir       string
	templateData   interface{}
	processedFiles map[string]*processedFileInfo
}

// ProcessedFilenames returns a list of file names which have processed.
func (srv Server) ProcessedFilenames() []string {
	nameList := make([]string, 0, len(srv.processedFiles))
	for k := range srv.processedFiles {
		nameList = append(nameList, k)
	}
	return nameList
}

var serviceInfo = app.ServiceInfo{
	Name:        "Web UI service",
	Description: "A generic service for serving web UI assets",
}

// ServiceInfo conforms app.ServiceServer interface.
func (srv Server) ServiceInfo() app.ServiceInfo { return serviceInfo }

// Serve conforms app.ServiceServer interface.
func (srv *Server) Serve() error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.config.ServePort),
		Handler: srv}
	srv.httpServer = httpServer
	err := srv.httpServer.ListenAndServe()
	if err == nil {
		if !srv.shuttingDown {
			return errors.Msg("server stopped unexpectedly")
		}
		return nil
	}
	if err == http.ErrServerClosed && srv.shuttingDown {
		return nil
	}
	return err
}

// Shutdown conforms app.ServiceServer interface.
func (srv *Server) Shutdown(ctx context.Context) error {
	//TODO: mutex?
	srv.shuttingDown = true
	return srv.httpServer.Shutdown(ctx)
}

// IsAcceptingClients conforms app.ServiceServer interface.
func (srv Server) IsAcceptingClients() bool {
	return !srv.shuttingDown && srv.IsHealthy()
}

// IsHealthy conforms app.ServiceServer interface.
func (srv Server) IsHealthy() bool { return true }

// ServeHTTP conforms Go's HTTP Handler interface.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL != nil && r.Method == http.MethodGet {
		// Note that for now, we use Path instead of RawPath. It works but
		// we haven't tested it with edge cases. If there's issues with
		// the path etc., might want to experiment with RawPath.
		reqPath := r.URL.Path
		if strings.HasPrefix(reqPath, srv.config.ServePath) {
			fileKey := strings.TrimPrefix(reqPath, srv.config.ServePath)
			if fileKey == "" || strings.HasSuffix(fileKey, "/") {
				fileKey += "index.html"
			}
			if fileInfo := srv.processedFiles[fileKey]; fileInfo != nil {
				if fileInfo.contentType != "" {
					w.Header().Set("Content-Type", fileInfo.contentType)
				}
				reqETag := r.Header.Get("If-None-Match")
				if fileInfo.etag == reqETag {
					w.Header().Set("ETag", reqETag)
					w.WriteHeader(http.StatusNotModified)
					return
				}
				w.Header().Set("ETag", fileInfo.etag)
				w.Write(fileInfo.content)
				return
			}
		}
	}

	srv.fileServer.ServeHTTP(w, r)
}

func (srv *Server) processFiles() error {
	return filepath.Walk(srv.filesDir, srv.processFile)
}

func (srv *Server) processFile(path string, info os.FileInfo, err error) error {
	if info == nil {
		return err
	}
	if info.IsDir() {
		return err
	}

	for pattern, processors := range srv.config.FileProcessors {
		//TODO: use proper globbing
		if !strings.HasPrefix(pattern, "*.") {
			panic("Unsupported pattern (for now) " + pattern)
		}
		ext := pattern[1:]
		if !strings.HasSuffix(path, ext) || len(processors) == 0 {
			continue
		}

		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(fileBytes)
		for _, proc := range processors {
			buf, err = proc.ProcessFile(path, buf)
			if err != nil {
				panic(err)
			}
		}

		content := buf.Bytes()

		h := xxhash.Checksum64(content)
		eTagStr := fmt.Sprintf(`W/"%d-%x"`, len(content), h)

		extWithDot := filepath.Ext(path)
		contentType := mime.TypeByExtension(extWithDot)

		fileKey := strings.TrimPrefix(path, srv.filesDir)
		srv.processedFiles[fileKey] = &processedFileInfo{
			content:     content, // GZip?
			etag:        eTagStr,
			contentType: contentType,
		}

		// One processing per pattern?
		return nil
	}

	return nil
}

type processedFileInfo struct {
	content     []byte
	etag        string
	contentType string
}

type FileProcessor interface {
	ProcessFile(
		filename string,
		inputBuffer *bytes.Buffer,
	) (outputBuffer *bytes.Buffer, err error)
}

type HTMLRenderer struct {
	Config       HTMLRendererConfig
	TemplateData interface{}
}

var _ FileProcessor = &HTMLRenderer{}

func (proc *HTMLRenderer) ProcessFile(
	filename string, inputBuffer *bytes.Buffer,
) (outputBuffer *bytes.Buffer, err error) {
	tpl := htmlt.New(filepath.Base(filename)).
		Delims(
			proc.Config.TemplateDelimBegin,
			proc.Config.TemplateDelimEnd)
	tpl, err = tpl.Parse(inputBuffer.String())
	if err != nil {
		return nil, errors.Wrap("parse html template", err)
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, proc.TemplateData)
	if err != nil {
		return nil, errors.Wrap("render html template", err)
	}

	return buf, nil
}

type HTMLRendererConfig struct {
	TemplateDelimBegin string `env:"TEMPLATE_DELIM_BEGIN"`
	TemplateDelimEnd   string `env:"TEMPLATE_DELIM_END"`
}

type JSRenderer struct {
	Config       JSRendererConfig
	TemplateData interface{}
}

var _ FileProcessor = &JSRenderer{}

func (proc *JSRenderer) ProcessFile(
	filename string, inputBuffer *bytes.Buffer,
) (outputBuffer *bytes.Buffer, err error) {
	tpl := textt.New(filepath.Base(filename)).
		Delims(
			proc.Config.TemplateDelimBegin,
			proc.Config.TemplateDelimEnd)
	tpl, err = tpl.Parse(inputBuffer.String())
	if err != nil {
		return nil, errors.Wrap("parse html template", err)
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, proc.TemplateData)
	if err != nil {
		return nil, errors.Wrap("render html template", err)
	}

	return buf, nil
}

type JSRendererConfig struct {
	TemplateDelimBegin string `env:"TEMPLATE_DELIM_BEGIN"`
	TemplateDelimEnd   string `env:"TEMPLATE_DELIM_END"`
}

type StringReplacer struct {
	Old string
	New string
}

var _ FileProcessor = &StringReplacer{}

func (proc *StringReplacer) ProcessFile(
	filename string, inputBuffer *bytes.Buffer,
) (outputBuffer *bytes.Buffer, err error) {
	buf := bytes.ReplaceAll(inputBuffer.Bytes(), []byte(proc.Old), []byte(proc.New))
	return bytes.NewBuffer(buf), nil
}

func ETagHandler(innerHandler http.Handler) http.Handler {
	responseBufferPool := bpool.NewBufferPool(24)
	pathETagsMutex := sync.RWMutex{}
	pathETags := map[string]string{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			innerHandler.ServeHTTP(w, r)
			return
		}

		pathETagsMutex.RLock()
		eTagStr := pathETags[r.URL.String()]
		pathETagsMutex.RUnlock()

		if eTagStr != "" {
			reqETag := r.Header.Get("If-None-Match")
			if eTagStr == reqETag {
				w.Header().Set("ETag", eTagStr)
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		buf := responseBufferPool.Get()
		defer responseBufferPool.Put(buf)

		eTagWriter := &ETagWriter{w.Header(), 0, buf}

		innerHandler.ServeHTTP(eTagWriter, r)

		if eTagWriter.statusCode < 200 || eTagWriter.statusCode >= 300 ||
			eTagWriter.statusCode == http.StatusNoContent || buf.Len() == 0 {
			w.WriteHeader(eTagWriter.statusCode)
			w.Write(buf.Bytes())
			return
		}

		h := xxhash.Checksum64(buf.Bytes())
		// The pattern we use here is <size>-<hash> . We use this pattern
		// because hashes could collide. By including the size, it should
		// significantly reducing collision probability.
		eTagStr = fmt.Sprintf(`W/"%d-%x"`, buf.Len(), h)

		pathETagsMutex.Lock()
		pathETags[r.URL.String()] = eTagStr
		pathETagsMutex.Unlock()

		w.Header().Set("ETag", eTagStr)
		w.WriteHeader(eTagWriter.statusCode)
		w.Write(buf.Bytes())
	})
}

type ETagWriter struct {
	innerHeader http.Header
	statusCode  int
	buffer      *bytes.Buffer
}

var _ http.ResponseWriter = &ETagWriter{}

func (eTagWriter *ETagWriter) Header() http.Header {
	return eTagWriter.innerHeader
}

func (eTagWriter *ETagWriter) WriteHeader(statusCode int) {
	eTagWriter.statusCode = statusCode
}

func (eTagWriter *ETagWriter) Write(p []byte) (int, error) {
	if eTagWriter.statusCode == 0 {
		eTagWriter.statusCode = http.StatusOK
	}
	return eTagWriter.buffer.Write(p)
}
