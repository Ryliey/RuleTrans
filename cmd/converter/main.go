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
	changes, err := git.GetDiffFiles(
		os.Getenv("BEFORE_COMMIT"),
		os.Getenv("AFTER_COMMIT"),
	)
	if err != nil {
		log.Fatalf("Failed to get changes: %v", err)
	}

	clashConv := clash.NewConverter()
	singboxConv := singbox.NewConverter()

	proc := processor.NewProcessor(clashConv, singboxConv)
	proc.Process(changes)
}
