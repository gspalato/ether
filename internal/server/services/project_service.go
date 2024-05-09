package services

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/google/go-github/v61/github"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v3"

	"unreal.sh/ether/internal/structures"
	"unreal.sh/ether/internal/utils"
)

type ProjectMetadata struct {
	Name string `yaml:"name"`
	Description string `yaml:"description"`
}

type ProjectService struct {
	cron *cron.Cron
	client *github.Client

	projects []structures.Project
}

func (ps *ProjectService) GetProjects(ctx context.Context) []structures.Project {
	return ps.projects
}

func (ps *ProjectService) AddProject(ctx context.Context, project structures.Project) {
	ps.projects = append(ps.projects, project)
}

func (ps *ProjectService) UpdateProjects(ctx context.Context) {
	username := utils.GetenvOr("GITHUB_USERNAME", "gspalato");

	opt := &github.RepositoryListByUserOptions{Type: "public"}
	repos, _, err := ps.client.Repositories.ListByUser(ctx, username, opt)

	if err != nil {
		fmt.Printf("Failed to fetch repositories from GitHub. %s\n", err.Error());
		return;
	}

	fmt.Printf("Found %d repositories belonging to %s.\n", len(repos), username);

	for _, repo := range repos {
		fmt.Printf("Checking repository %s...\n", repo.GetName());

		_, dirContent, _, err := ps.client.Repositories.GetContents(ctx, username, repo.GetName(), ".project", nil);
		
		if err != nil {
			fmt.Printf("Repository %s doesn't contain a .project folder. Skipping...\n", repo.GetName());
			continue;
		}

		project := structures.Project{}

		metadataFileIndex := slices.IndexFunc(dirContent, func(c *github.RepositoryContent) bool {
			return c.GetName() == "metadata.yml" || c.GetName() == "metadata.yaml"
		})
		fmt.Printf("Found metadata file index: %d\n", metadataFileIndex)

		if metadataFileIndex == -1 {
			fmt.Printf("No metadata found for project %s. Skipping...", repo.GetName());
			continue;
		}

		bannerFileIndex := slices.IndexFunc(dirContent, func(c *github.RepositoryContent) bool {
			return c.GetName() == "banner.jpg"
		})
		fmt.Printf("Found banner file index: %d\n", bannerFileIndex)

		if bannerFileIndex == -1 {
			fmt.Printf("No banner found for project %s. Skipping...", repo.GetName());
			continue;
		}

		metadataFile := dirContent[metadataFileIndex]
		bannerFile := dirContent[bannerFileIndex]

		fmt.Printf("Found metadata and banner for project %s.\nContent:\n%s\n", repo.GetName(), *metadataFile.Content);

		metadataFileContent, err := metadataFile.GetContent();
		if err != nil {
			fmt.Printf("Failed to get content for metadata file for project %s. %s\n", repo.GetName(), err.Error());
			continue;
		}

		metadata, err := ps.ParseMetadata(ctx, &metadataFileContent)

		if err != nil {
			fmt.Printf("Failed to parse metadata for project %s. %s\n", repo.GetName(), err.Error());
			continue;
		}

		project.Name = metadata.Name;
		project.Description = metadata.Description;
		project.Url = repo.GetHTMLURL();

		// Set banner URL only, since the portfolio website currently only uses the banners.
		// Other images are stored in .project for use in README or other stuff.
		project.BannerUrl = bannerFile.GetDownloadURL();

		ps.AddProject(ctx, project);

		fmt.Printf("Reloaded project %s.\n", repo.GetName());
	}
}

func (ps *ProjectService) ParseMetadata(ctx context.Context, str *string) (ProjectMetadata, error) {
	metadata := ProjectMetadata{}

	err := yaml.Unmarshal([]byte(*str), &metadata)

	if err != nil {
		return metadata, err;
	}

	return metadata, nil;
}

func (ps *ProjectService) Init(ctx context.Context) {
	ps.client = github.NewClient(nil);

	// Authenticate with GitHub API.
	authToken, foundAuthToken := os.LookupEnv("GITHUB_AUTH_TOKEN");
	if foundAuthToken {
		ps.client.WithAuthToken(authToken);
	}

	ps.projects = make([]structures.Project, 0);

	ps.cron = cron.New();

	ps.cron.AddFunc("@every 10m", func() {
		ps.UpdateProjects(ctx);
		fmt.Println("Updated projects.");
	});

	ps.cron.Start();

	ps.UpdateProjects(ctx);
}