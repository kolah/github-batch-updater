package github

import (
	"context"
	"fmt"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type ContentService interface {
	GetFile(ctx context.Context, repo Repository, path string) (*File, error)
	UpdateContent(ctx context.Context, repo Repository, branch string, file *File) error
}

type File struct {
	path            string
	originalContent string
	modifiedContent *string
	diff            string
	sourceSHA       string
}

func LoadGitHubFile(path, content, sourceSHA string) *File {
	return &File{
		path:            path,
		originalContent: content,
		sourceSHA:       sourceSHA,
	}
}

func (f *File) Modify(content string) *File {
	modified := f.originalContent != content
	if !modified {
		return f
	}

	return &File{
		path:            f.path,
		originalContent: f.originalContent,
		modifiedContent: &content,
		sourceSHA:       f.sourceSHA,
	}
}

func (f *File) Diff() string {
	if f.modifiedContent == nil {
		return ""
	}

	edits := myers.ComputeEdits(span.URIFromPath(f.path), f.originalContent, *f.modifiedContent)
	difference := gotextdiff.ToUnified(f.path, f.path, f.originalContent, edits)

	return fmt.Sprint(difference)
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Content() string {
	if f.modifiedContent != nil {
		return *f.modifiedContent
	}

	return f.originalContent
}

func (f *File) SourceSHA() string {
	return f.sourceSHA
}

func (f *File) Modified() bool {
	return f.modifiedContent != nil
}
