package main

import (
	"context"
	"log"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// Connect to Dagger (Cloud auto-enabled via env var)
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Writer()))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	src := client.Host().Directory(".")

	goCache := client.CacheVolume("go-build")
	modCache := client.CacheVolume("go-mod")

	container := client.Container().
		From("cgr.dev/chainguard/go:latest").
		WithDirectory("/src", src).
		WithWorkdir("/src").
		WithMountedCache("/root/.cache/go-build", goCache).
		WithMountedCache("/go/pkg/mod", modCache).
		WithExec([]string{"go", "test", "./..."})

	_, err = container.ExitCode(ctx)
	if err != nil {
		panic(err)
	}
}
