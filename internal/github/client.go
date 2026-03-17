package github

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"

	"github.com/lichenglife/issue2md/internal/models"
)

// Client GitHub API 客户端
type Client struct {
	client *github.Client
}

// NewClient 创建新的客户端
// token 为空时使用匿名访问（仅限公开仓库）
func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &Client{client: client}
}

// FetchIssue 获取 Issue 或 PR 数据
// resourceType: "issue" 或 "pull"
// 返回的 Comments 已按 CreatedAt 升序排列
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int, resourceType string) (*models.IssueData, error) {
	var (
		issue *github.Issue
		resp  *github.Response
		err   error
	)

	switch resourceType {
	case "issue":
		issue, resp, err = c.client.Issues.Get(ctx, owner, repo, number)
	case "pull":
		var pull *github.PullRequest
		pull, resp, err = c.client.PullRequests.Get(ctx, owner, repo, number)
		if pull != nil {
			// 将 PR 转换为 Issue 格式
			issue = &github.Issue{
				Title:     pull.Title,
				Body:      pull.Body,
				User:      pull.User,
				CreatedAt: pull.CreatedAt,
				State:     pull.State,
				ClosedAt:  pull.ClosedAt,
				Labels:    pull.Labels,
				Assignees: pull.Assignees,
				Milestone: pull.Milestone,
			}
		}
	default:
		return nil, fmt.Errorf("不支持的资源类型：%s", resourceType)
	}

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		if resp != nil && resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			return nil, ErrRateLimited
		}
		return nil, fmt.Errorf("获取 Issue 失败：%w", err)
	}

	if issue == nil {
		return nil, ErrNotFound
	}

	// 获取所有评论
	comments, err := c.fetchAllComments(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("获取评论失败：%w", err)
	}

	// 转换 Labels
	var labels []models.Label
	for _, label := range issue.Labels {
		if label != nil {
			labels = append(labels, models.Label{
				Name:        label.GetName(),
				Color:       label.GetColor(),
				Description: label.GetDescription(),
			})
		}
	}

	// 转换 Assignees
	var assignees []models.User
	for _, assignee := range issue.Assignees {
		if assignee != nil {
			assignees = append(assignees, models.User{
				Login:   assignee.GetLogin(),
				HTMLURL: assignee.GetHTMLURL(),
			})
		}
	}

	// 转换 Milestone
	var milestone *models.Milestone
	if issue.Milestone != nil {
		milestone = &models.Milestone{
			Title:       issue.Milestone.GetTitle(),
			Number:      issue.Milestone.GetNumber(),
			State:       issue.Milestone.GetState(),
			Description: issue.Milestone.GetDescription(),
			DueOn:       milestoneTime(issue.Milestone.DueOn),
		}
	}

	// 转换 User
	user := models.User{
		Login:   issue.User.GetLogin(),
		HTMLURL: issue.User.GetHTMLURL(),
	}

	return &models.IssueData{
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		HTMLURL:   issue.GetHTMLURL(),
		Number:    issue.GetNumber(),
		Repo:      repo,
		User:      user,
		CreatedAt: timestampTime(issue.GetCreatedAt()),
		State:     issue.GetState(),
		ClosedAt:  timestampTimePtr(issue.ClosedAt),
		Labels:    labels,
		Assignees: assignees,
		Milestone: milestone,
		Comments:  comments,
	}, nil
}

// fetchAllComments 获取所有评论并按 CreatedAt 升序排序
func (c *Client) fetchAllComments(ctx context.Context, owner, repo string, number int) ([]models.Comment, error) {
	var allComments []models.Comment

	opt := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Direction:   github.String("asc"), // 按时间升序
	}

	for {
		comments, resp, err := c.client.Issues.ListComments(ctx, owner, repo, number, opt)
		if err != nil {
			return nil, err
		}

		for _, comment := range comments {
			if comment != nil {
				allComments = append(allComments, models.Comment{
					ID:        comment.GetID(),
					Body:      comment.GetBody(),
					User: models.User{
						Login:   comment.User.GetLogin(),
						HTMLURL: comment.User.GetHTMLURL(),
					},
					CreatedAt: timestampTime(comment.GetCreatedAt()),
					UpdatedAt: timestampTime(comment.GetUpdatedAt()),
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// 确保按 CreatedAt 升序排序
	sort.Slice(allComments, func(i, j int) bool {
		return allComments[i].CreatedAt.Before(allComments[j].CreatedAt)
	})

	return allComments, nil
}

// FetchDiscussion 获取 Discussion 数据（使用 GraphQL API）
// TODO: 实现 GraphQL API 调用
func (c *Client) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*models.IssueData, error) {
	return nil, fmt.Errorf("FetchDiscussion 尚未实现")
}

// timestampTime 将 github.Timestamp 转换为 time.Time
func timestampTime(ts github.Timestamp) time.Time {
	if ts.IsZero() {
		return time.Time{}
	}
	return ts.Time
}

// timestampTimePtr 将 *github.Timestamp 转换为 *time.Time
func timestampTimePtr(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.Time
	return &t
}

// milestoneTime 将 *github.Timestamp 转换为 *time.Time
func milestoneTime(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.Time
	return &t
}
