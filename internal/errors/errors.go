package errors

import (
	"fmt"
)

// ==================== Error Code Constants ====================

const (
	// Player errors
	ErrorPlayerNotFound      = "PLAYER_NOT_FOUND"
	ErrorPlayerAlreadyExists = "PLAYER_ALREADY_EXISTS"
	ErrorPlayerInactive      = "PLAYER_INACTIVE"
	ErrorPlayerInvalidData   = "PLAYER_INVALID_DATA"

	// Season errors
	ErrorSeasonNotFound      = "SEASON_NOT_FOUND"
	ErrorSeasonNotActive     = "SEASON_NOT_ACTIVE"
	ErrorSeasonAlreadyExists = "SEASON_ALREADY_EXISTS"

	// Team errors
	ErrorTeamNotFound    = "TEAM_NOT_FOUND"
	ErrorTeamInvalidData = "TEAM_INVALID_DATA"

	// Rank errors
	ErrorRankNotFound = "RANK_NOT_FOUND"

	// Player-Season errors
	ErrorPlayerAlreadyInSeason = "PLAYER_ALREADY_IN_SEASON"
	ErrorPlayerSeasonNotFound  = "PLAYER_SEASON_NOT_FOUND"

	// Fixture errors
	ErrorFixtureNotFound  = "FIXTURE_NOT_FOUND"
	ErrorFixtureNotActive = "FIXTURE_NOT_ACTIVE"
	ErrorInvalidTeamMatch = "INVALID_TEAM_MATCH"
	ErrorSameTeamMatch    = "SAME_TEAM_MATCH"

	// Match errors
	ErrorMatchNotFound        = "MATCH_NOT_FOUND"
	ErrorMatchInvalidSets     = "MATCH_INVALID_SETS"
	ErrorInvalidPlayers       = "INVALID_PLAYERS"
	ErrorMatchAlreadyRecorded = "MATCH_ALREADY_RECORDED"

	// Point errors
	ErrorInvalidPointAdjustment = "INVALID_POINT_ADJUSTMENT"
	ErrorNegativePointsResult   = "NEGATIVE_POINTS_RESULT"

	// Database errors
	ErrorDatabaseError   = "DATABASE_ERROR"
	ErrorDatabaseTimeout = "DATABASE_TIMEOUT"

	// Validation errors
	ErrorInvalidInput    = "INVALID_INPUT"
	ErrorMissingRequired = "MISSING_REQUIRED_FIELD"

	// Authorization errors (future)
	ErrorUnauthorized = "UNAUTHORIZED"
	ErrorForbidden    = "FORBIDDEN"
)

// ==================== Custom Error Type ====================

// AppError represents an application error with standardized code and message
type AppError struct {
	Code       string      `json:"error_code"`
	Message    string      `json:"message"`
	StatusCode int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
	Cause      error       `json:"-"` // Internal error for logging
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ==================== Error Constructors ====================

// NewAppError creates a new app error with given code and message
func NewAppError(code string, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithCause adds the underlying error for logging purposes
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// WithDetails adds additional details to the error response
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// ==================== Common Error Generators ====================

func PlayerNotFound() *AppError {
	return NewAppError(ErrorPlayerNotFound, "VĐV không tồn tại", 404)
}

func PlayerAlreadyExists() *AppError {
	return NewAppError(ErrorPlayerAlreadyExists, "VĐV đã tồn tại", 409)
}

func SeasonNotFound() *AppError {
	return NewAppError(ErrorSeasonNotFound, "Mùa giải không tồn tại", 404)
}

func TeamNotFound() *AppError {
	return NewAppError(ErrorTeamNotFound, "Đội bóng không tồn tại", 404)
}

func RankNotFound() *AppError {
	return NewAppError(ErrorRankNotFound, "Hạng trình độ không tồn tại", 404)
}

func PlayerAlreadyInSeason() *AppError {
	return NewAppError(ErrorPlayerAlreadyInSeason, "VĐV đã tồn tại trong mùa giải này", 409)
}

func PlayerSeasonNotFound() *AppError {
	return NewAppError(ErrorPlayerSeasonNotFound, "Không tìm thấy dữ liệu VĐV trong mùa giải", 404)
}

func FixtureNotFound() *AppError {
	return NewAppError(ErrorFixtureNotFound, "Trận đấu CLB không tồn tại", 404)
}

func MatchNotFound() *AppError {
	return NewAppError(ErrorMatchNotFound, "Trận đấu con không tồn tại", 404)
}

func InvalidPointAdjustment() *AppError {
	return NewAppError(ErrorInvalidPointAdjustment, "Điểm điều chỉnh không hợp lệ", 400)
}

func DatabaseError(cause error) *AppError {
	return NewAppError(ErrorDatabaseError, "Lỗi cơ sở dữ liệu", 500).WithCause(cause)
}

func InvalidInput(message string) *AppError {
	return NewAppError(ErrorInvalidInput, message, 400)
}

func SameTeamMatch() *AppError {
	return NewAppError(ErrorSameTeamMatch, "Sự kiện không thể là giữa hai đội giống nhau", 400)
}

func InvalidPlayers() *AppError {
	return NewAppError(ErrorInvalidPlayers, "Danh sách cầu thủ không hợp lệ", 400)
}

func MatchAlreadyRecorded() *AppError {
	return NewAppError(ErrorMatchAlreadyRecorded, "Trận đấu đã được ghi lại kết quả", 409)
}

func NegativePointsResult() *AppError {
	return NewAppError(ErrorNegativePointsResult, "Điểm không thể là số âm", 400)
}
