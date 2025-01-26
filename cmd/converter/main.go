package main

import (
	"log"
	"os"

	"github.com/Ryliey/RuleTrans/internal/converter/clash"
	"github.com/Ryliey/RuleTrans/internal/converter/singbox"
	"github.com/Ryliey/RuleTrans/internal/git"
	"github.com/Ryliey/RuleTrans/internal/processor"
)

func main() {
	before := os.Getenv("BEFORE_COMMIT")
	after := os.Getenv("AFTER_COMMIT")
	log.Printf("Comparing commits: %s..%s", before, after)

	changes, err := git.GetDiffFiles(before, after)
	if err != nil {
		log.Fatalf("Failed to get changes: %v", err)
	}
	log.Printf("Found %d file changes", len(changes))

	for i, fc := range changes {
		log.Printf("[%d] %s %s", i+1, fc.Status, fc.Path)
	}

	clashConv := clash.NewConverter()
	singboxConv := singbox.NewConverter()

	proc := processor.NewProcessor(clashConv, singboxConv)
	proc.Process(changes)
	log.Println("Processing completed")
}
