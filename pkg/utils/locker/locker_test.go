package locker

import (
	"testing"
)

func TestLocker(t *testing.T) {
	pathfile := "/var/lock/test.lock"

	if IsLocked(pathfile) == true {
		t.Errorf("Fail - lock file already exist")
	}

	if err := Lock(pathfile); err != nil {
		t.Errorf("Error: %s", err)
		t.FailNow()
	}

	if IsLocked(pathfile) == false {
		t.Errorf("Fail - lock file not exist after creating")
	}

	if err := Unlock(); err != nil {
		t.Errorf("Error: %s", err)
		t.FailNow()
	}

	if IsLocked(pathfile) == true {
		t.Errorf("Fail - lock file exists after unlocking")
	}
}