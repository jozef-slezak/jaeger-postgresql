package release

import (
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/client"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/strategy"
	"github.com/pkg/errors"
)

// Pipe for github release
type Pipe struct{}

func (Pipe) String() string {
	return "GitHub Releases"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Release.NameTemplate == "" {
		ctx.Config.Release.NameTemplate = "{{.Tag}}"
	}
	if ctx.Config.Release.GitHub.Name == "" {
		repo, err := remoteRepo()
		if err != nil && !ctx.Snapshot {
			return err
		}
		ctx.Config.Release.GitHub = repo
	}

	// Check if we have to check the git tag for an indicator to mark as pre release
	switch ctx.Config.Release.Prerelease {
	case "auto":
		sv, err := semver.NewVersion(ctx.Git.CurrentTag)
		if err != nil {
			return errors.Wrapf(err, "failed to parse tag %s as semver", ctx.Git.CurrentTag)
		}

		if sv.Prerelease() != "" {
			ctx.PreRelease = true
		}
		log.Debugf("pre-release was detected for tag %s: %v", ctx.Git.CurrentTag, ctx.PreRelease)
	case "true":
		ctx.PreRelease = true
	}
	log.Debugf("pre-release for tag %s set to %v", ctx.Git.CurrentTag, ctx.PreRelease)

	return nil
}

// Publish github release
func (Pipe) Publish(ctx *context.Context) error {
	c, err := client.NewGitHub(ctx)
	if err != nil {
		return err
	}
	return doPublish(ctx, c)
}

func doPublish(ctx *context.Context, c client.Client) error {
	if ctx.Config.Release.Disable {
		return pipe.Skip("release pipe is disabled")
	}
	log.WithField("tag", ctx.Git.CurrentTag).
		WithField("repo", ctx.Config.Release.GitHub.String()).
		Info("creating or updating release")
	body, err := describeBody(ctx)
	if err != nil {
		return err
	}
	releaseID, err := c.CreateRelease(ctx, body.String())
	if err != nil {
		return err
	}
	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range ctx.Artifacts.Filter(
		artifact.Or(
			artifact.ByType(artifact.UploadableArchive),
			artifact.ByType(artifact.UploadableBinary),
			artifact.ByType(artifact.Checksum),
			artifact.ByType(artifact.Signature),
			artifact.ByType(artifact.LinuxPackage),
		),
	).List() {
		artifact := artifact
		g.Go(func() error {
			var repeats uint
			action := func(try uint) error {
				repeats = try + 1
				if uploadErr := upload(ctx, c, releaseID, artifact); uploadErr != nil {
					log.WithFields(log.Fields{
						"try":      try,
						"artifact": artifact.Name,
					}).Warnf("failed to upload artifact, will retry")
					return uploadErr
				}
				return nil
			}
			strategies := []strategy.Strategy{
				strategy.Limit(10),
				strategy.Backoff(backoff.Linear(50 * time.Millisecond)),
			}
			if retryErr := retry.Retry(ctx.Done(), action, strategies...); retryErr != nil {
				return errors.Wrapf(retryErr, "failed to upload %s after %d retries", artifact.Name, repeats)
			}
			return nil
		})
	}
	return g.Wait()
}

func upload(ctx *context.Context, c client.Client, releaseID int64, artifact artifact.Artifact) error {
	file, err := os.Open(artifact.Path)
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck
	log.WithField("file", file.Name()).WithField("name", artifact.Name).Info("uploading to release")
	return c.Upload(ctx, releaseID, artifact.Name, file)
}
