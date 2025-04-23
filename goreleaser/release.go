package main

import (
	"dagger/goreleaser/internal/dagger"
	"strconv"
)

// Release represents the `goreleaser release` command.
type Release struct {
	// +private
	gr *Goreleaser

	// build flags
	// +private
	flags []string
}

// Release represents the `goreleaser release` command.
func (gr *Goreleaser) Release() *Release {
	return &Release{
		gr:    gr,
		flags: []string{"goreleaser", "release"},
	}
}

// Generate an unversioned snapshot release, skipping all validations and without publishing any artifacts.
//
// e.g. `goreleaser release --snapshot`.
func (r *Release) WithSnapshot() *Release {
	r.flags = append(r.flags, "--snapshot")
	return r
}

// Automatically sets WithSnapshot if the repository is dirty.
//
// e.g. `goreleaser build --auto-snapshot`.
func (r *Release) WithAutoSnapshot() *Release {
	r.flags = append(r.flags, "--auto-snapshot")
	return r
}

// Removes the 'dist' directory before building.
//
// e.g. `goreleaser release --clean`.
func (r *Release) WithClean() *Release {
	r.flags = append(r.flags, "--clean")
	return r
}

// Timeout to the entire release process.
//
// e.g. `goreleaser build --timeout <duration>`.
func (r *Release) WithTimeout(
	// Timeout duration, e.g. 10m, 10m30s. Default is 30m.
	duration string,
) *Release {
	r.flags = append(r.flags, "--timeout", duration)
	return r
}

// Abort the release publishing on the first error.
//
// e.g. `goreleaser release --fail-fast`.
func (r *Release) WithFailFast() *Release {
	r.flags = append(r.flags, "--fail-fast")
	return r
}

// Tasks to run concurrently (default: number of CPUs).
//
// e.g. `goreleaser release --parallelism <n>`.
func (r *Release) WithParallelism(
	// concurrent tasks
	n int,
) *Release {
	r.flags = append(r.flags, "--parallelism", strconv.Itoa(n))
	return r
}

// Load custom release notes from a markdown file, skips changelong generation.
//
// e.g. `goreleaser release --release-notes <notes>`.
func (r *Release) WithNotes(
	// release notes markdown file
	notes *dagger.File,
) *Release {
	notesPath := "/work/notes.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(notesPath, notes)
	r.flags = append(r.flags, "--release-notes", notesPath)
	return r
}

// Load custom release notes from a templated markdown file. Overrides WithNotes.
//
// e.g. `goreleaser release --release-notes-tmpl <notesTmpl>`.
func (r *Release) WithNotesTmpl(
	// release notes templated markdown file
	notesTmpl *dagger.File,
) *Release {
	notesPath := "/work/notes-tmpl.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(notesPath, notesTmpl)
	r.flags = append(r.flags, "--release-notes-tmpl", notesPath)
	return r
}

// Load custom release notes header from a markdown file.
//
// e.g. `goreleaser release --release-header <header>`.
func (r *Release) WithNotesHeader(header *dagger.File) *Release {
	headerPath := "/work/header.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(headerPath, header)
	r.flags = append(r.flags, "--release-header", headerPath)
	return r
}

// Load custom release notes header from a templated markdown file. Overrides WithNotesHeader.
//
// e.g. `goreleaser release --release-header-tmpl <headerTmpl>`.
func (r *Release) WithNotesHeaderTmpl(
	// release notes header templated markdown file
	headerTmpl *dagger.File,
) *Release {
	headerPath := "/work/header-tmpl.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(headerPath, headerTmpl)
	r.flags = append(r.flags, "release-header-tmpl", headerPath)
	return r
}

// Load custom release notes footer from a markdown file.
//
// e.g. `goreleaser release --release-footer <footer>`.
func (r *Release) WithNotesFooter(footer *dagger.File) *Release {
	footerPath := "/work/header.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(footerPath, footer)
	r.flags = append(r.flags, "--release-footer", footerPath)
	return r
}

// Load custom release notes footer from a templated markdown file. Overrides WithNotesFooter.
//
// e.g. `goreleaser release --release-footer-tmpl <footerTmpl>`.
func (r *Release) WithNotesFooterTmpl(
	// release notes footer templated markdown file
	footerTmpl *dagger.File,
) *Release {
	footerPath := "/work/footer-tmpl.md"
	r.gr.Ctr = r.gr.Ctr.WithMountedFile(footerPath, footerTmpl)
	r.flags = append(r.flags, "--release-footer-tmpl", footerPath)
	return r
}
