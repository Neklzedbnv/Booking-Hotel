package auth

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// Stub method (authentication not required yet)
func (s *Service) Ping() bool {
	return true
}