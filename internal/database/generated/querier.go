// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateProblem(ctx context.Context, arg CreateProblemParams) (Problem, error)
	CreateSubmission(ctx context.Context, arg CreateSubmissionParams) (Submission, error)
	CreateTestCase(ctx context.Context, arg CreateTestCaseParams) (TestCase, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProblem(ctx context.Context, id int32) error
	DeleteTestCase(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	GetProblemById(ctx context.Context, id int32) (Problem, error)
	// TODO: optimize this query
	GetProblemsCount(ctx context.Context) (int64, error)
	// TODO: optimize this query
	GetPublishedProblemsCount(ctx context.Context) (int64, error)
	GetSubmissionById(ctx context.Context, id int32) (Submission, error)
	GetTestCaseById(ctx context.Context, id int32) (TestCase, error)
	GetUserById(ctx context.Context, id int32) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	// TODO: optimize this query
	GetUserProblemsCount(ctx context.Context, ownerID int32) (int64, error)
	GetUserRankingById(ctx context.Context, limit int32) ([]GetUserRankingByIdRow, error)
	GetUserStatsById(ctx context.Context, userID int32) (GetUserStatsByIdRow, error)
	GetUserSubmissionsCount(ctx context.Context, userID pgtype.Int4) (int64, error)
	ListProblemSubmissions(ctx context.Context, problemID pgtype.Int4) ([]Submission, error)
	ListProblems(ctx context.Context, arg ListProblemsParams) ([]Problem, error)
	ListPublishedProblems(ctx context.Context, arg ListPublishedProblemsParams) ([]Problem, error)
	ListTestCases(ctx context.Context, problemID pgtype.Int4) ([]TestCase, error)
	ListUserProblems(ctx context.Context, arg ListUserProblemsParams) ([]Problem, error)
	ListUserSubmissions(ctx context.Context, arg ListUserSubmissionsParams) ([]Submission, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateProblem(ctx context.Context, arg UpdateProblemParams) (Problem, error)
	UpdateSubmissionStatusTimeMemory(ctx context.Context, arg UpdateSubmissionStatusTimeMemoryParams) (Submission, error)
	UpdateTestCase(ctx context.Context, arg UpdateTestCaseParams) (TestCase, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UserHasSolvedProblem(ctx context.Context, arg UserHasSolvedProblemParams) (bool, error)
}

var _ Querier = (*Queries)(nil)
