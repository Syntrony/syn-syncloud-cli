package system

import "synctl/internal/application/interfaces"

type SudoSession struct {
	authorized bool
}

func (s *SudoSession) Ensure(runner interfaces.CommandRunner) error {
	if s.authorized {
		return nil
	}

	_, err := runner.Run("sudo", "-v")

	if err != nil {
		return err
	}

	s.authorized = true

	return nil
}
