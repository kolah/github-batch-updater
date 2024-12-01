package application

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/kolah/github-batch-updater/internal/domain/github"
	"github.com/kolah/github-batch-updater/internal/domain/transformer"
	"github.com/kolah/github-batch-updater/internal/pkg/slices"
)

type BatchProcessingService struct {
	repoService    github.RepositoryService
	contentService github.ContentService
	prService      github.PullRequestService
	transformer    transformer.Service
	logger         *log.Logger
	count          int
}

func NewBatchProcessingService(
	repoService github.RepositoryService,
	contentService github.ContentService,
	transformer transformer.Service,
	prService github.PullRequestService,
	logger *log.Logger,
) *BatchProcessingService {
	return &BatchProcessingService{
		repoService:    repoService,
		contentService: contentService,
		transformer:    transformer,
		prService:      prService,
		logger:         logger,
	}
}

func (s *BatchProcessingService) Process(ctx context.Context, input Input) {
	for _, repo := range input.Repositories {
		err := s.processRepository(ctx, repo, input)
		if err != nil {
			s.logger.Error("error processing repository", "error", err, "repository", repo)
			continue
		}
	}
	s.logger.Info("processing complete", "count", s.count)
}

func (s *BatchProcessingService) processRepository(ctx context.Context, repo RepositoryInput, input Input) error {
	repoInfo, err := s.repoService.GetRepository(ctx, repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	changes := make([]*github.File, 0)
	for path, rules := range input.Rules {
		file, err := s.processFile(ctx, *repoInfo, path, rules)
		if err != nil {
			s.logger.Error("error processing file", "error", err, "file", path)
			if file == nil {
				continue
			}
		}

		if file.Modified() {
			s.logger.Info("file modified", "file", path, "diff", file.Diff())
			s.count++
		}

		changes = append(changes, file)
	}

	if len(changes) > 0 && !input.DryRun {
		err := s.applyChanges(ctx, repoInfo, changes, input.TargetBranch, input.CreatePR, input.PullRequestTemplate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *BatchProcessingService) processFile(ctx context.Context, repo github.Repository, path string, rules []transformer.Rule) (*github.File, error) {
	file, err := s.contentService.GetFile(ctx, repo, path)
	if err != nil {
		return nil, err //todo error handling
	}

	// if there is an error, we don't apply
	result, err := s.transformer.ApplyRules(file.Content(), rules)
	if err != nil {
		return nil, err // todo: recoverable error
	}

	return file.Modify(result), nil
}

func (s *BatchProcessingService) applyChanges(
	ctx context.Context,
	repo *github.Repository,
	changes []*github.File,
	targetBranch string,
	createPR bool,
	prConfig github.PullRequestTemplate,
) error {
	modified := slices.Filter(changes, func(change *github.File) bool {
		return change.Modified()
	})

	if len(modified) == 0 {
		s.logger.Info("no changes to apply")
		return nil
	}

	err := s.repoService.CreateBranchFromDefaultBranch(ctx, repo, targetBranch)
	if err != nil {
		s.logger.Error("error creating branch", "error", err)
		return err
	}

	for _, change := range changes {
		err := s.contentService.UpdateContent(ctx, *repo, targetBranch, change)
		if err != nil {
			s.logger.Error("error updating content", "error", err)
			return err
		}
	}
	if !createPR {
		s.logger.Info("PR creation disabled")
		return nil
	}

	pr, err := s.prService.Create(ctx, *repo, targetBranch, prConfig)

	if err != nil {
		s.logger.Error("error creating PR", "error", err)
		return err
	}

	err = s.prService.RequestReview(ctx, *repo, pr, prConfig)
	if err != nil {
		s.logger.Error("error requesting review", "error", err)
		return err
	}

	return nil
}
