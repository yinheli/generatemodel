package main

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "struct.text",
		FileModTime: time.Unix(1526461909, 0),
		Content:     string("package model\n\n// {{.TitleCaseName}} table: {{.Name}} {{.Comment}}\ntype {{.TitleCaseName}} struct {\n    {{range .Columns -}}\n        {{- .TitleCaseName}} {{.GoType}} {{.Tag}}\n    {{end -}}\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1526461909, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "struct.text"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`template`, &embedded.EmbeddedBox{
		Name: `template`,
		Time: time.Unix(1526461909, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"struct.text": file2,
		},
	})
}
