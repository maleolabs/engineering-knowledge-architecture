// Package skeletonembed exposes the canonical EKA repository skeleton
// (the skeleton/ directory at the module root) as an embedded filesystem.
//
// It exists because go:embed directives cannot reference paths above the
// package directory (".."), while skeleton/ lives at the module root and is
// shared by tooling packages (bootstrap). Keeping the embed at the root
// makes the skeleton a compile-time resource of the eka binary: a bootstrap
// never depends on the skeleton being present on disk.
//
// An embed.FS always uses forward slashes and has no concept of directories
// as entries. Traverse it with fs.Sub, fs.WalkDir and fs.ReadFile; never
// emit Windows backslashes from paths obtained here.
package skeletonembed

import "embed"

// FS holds the canonical skeleton tree. Paths are prefixed with "skeleton/";
// use fs.Sub(FS, "skeleton") to address the tree contents directly.
//
//go:embed skeleton
var FS embed.FS
